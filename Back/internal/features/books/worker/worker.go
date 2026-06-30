package book_worker

import (
	"archive/zip"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/Eternity8c/FreeLib/internal/core/domain"
	"go.uber.org/zap"
)

type BookRepository interface {
	GetS3URLFromBook(ctx context.Context, id int) (string, error)
	SaveChapters(ctx context.Context, bookID int, chapters []domain.Chapter) error
}

type BookS3Repository interface {
	GetBookFile(ctx context.Context, fileName string) (io.ReadCloser, error)
}

type Worker struct {
	queue            chan int
	bookRepository   BookRepository
	bookS3Repository BookS3Repository
	logger           *zap.Logger
}

func NewWorker(
	bookRepository BookRepository,
	bookS3Repository BookS3Repository,
	logger *zap.Logger,
) *Worker {
	return &Worker{
		queue:            make(chan int, 100),
		bookRepository:   bookRepository,
		bookS3Repository: bookS3Repository,
		logger:           logger,
	}
}

func (w *Worker) Enqueue(bookID int) {
	w.queue <- bookID
}

func (w *Worker) Run(ctx context.Context) {
	for {
		select {
		case bookID := <-w.queue:
			if err := w.processBook(ctx, bookID); err != nil {
				w.logger.Error("failed to parse book",
					zap.Int("bookID", bookID),
					zap.Error(err),
				)
			} else {
				w.logger.Info("book parsed successfully", zap.Int("bookID", bookID))
			}
		case <-ctx.Done():
			return
		}
	}
}

func (w *Worker) processBook(ctx context.Context, bookID int) error {
	s3URL, err := w.bookRepository.GetS3URLFromBook(ctx, bookID)
	if err != nil {
		return fmt.Errorf("get s3 url: %w", err)
	}

	fileName := strings.TrimPrefix(s3URL, "https://storage.yandexcloud.net/tes-freelib-server/")

	rc, err := w.bookS3Repository.GetBookFile(ctx, fileName)
	if err != nil {
		return fmt.Errorf("get book file from s3: %w", err)
	}
	defer rc.Close()

	tmpFile, err := os.CreateTemp("", "epub-worker-*.tmp")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	size, err := io.Copy(tmpFile, rc)
	if err != nil {
		return fmt.Errorf("copy to temp file: %w", err)
	}

	if _, err := tmpFile.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("seek temp file: %w", err)
	}

	zipReader, err := zip.NewReader(tmpFile, size)
	if err != nil {
		return fmt.Errorf("create zip reader: %w", err)
	}

	chapters, err := parseEpub(zipReader)
	if err != nil {
		return fmt.Errorf("parse epub: %w", err)
	}

	w.logger.Debug("epub parsed", zap.Int("bookID", bookID), zap.Int("chapters", len(chapters)))

	if len(chapters) == 0 {
		return fmt.Errorf("epub has no readable chapters")
	}

	if err := w.bookRepository.SaveChapters(ctx, bookID, chapters); err != nil {
		return fmt.Errorf("save chapters: %w", err)
	}

	return nil
}

// --- epub parsing (epub 2.0 / 3.0, без сторонних библиотек) ---

type epubContainer struct {
	Rootfile struct {
		Path string `xml:"full-path,attr"`
	} `xml:"rootfiles>rootfile"`
}

type opfPackage struct {
	XMLName  xml.Name    `xml:"package"`
	Manifest opfManifest `xml:"manifest"`
	Spine    opfSpine    `xml:"spine"`
}

type opfManifest struct {
	Items []opfItem `xml:"item"`
}

type opfItem struct {
	ID        string `xml:"id,attr"`
	Href      string `xml:"href,attr"`
	MediaType string `xml:"media-type,attr"`
}

type opfSpine struct {
	Toc      string       `xml:"toc,attr,omitempty"`
	ItemRefs []opfItemRef `xml:"itemref"`
}

// ncxDoc — таблица содержания epub 2.0 (toc.ncx).
type ncxDoc struct {
	NavPoints []struct {
		NavLabel struct {
			Text string `xml:"text"`
		} `xml:"navLabel"`
		Content struct {
			Src string `xml:"src,attr"`
		} `xml:"content"`
	} `xml:"navMap>navPoint"`
}

type opfItemRef struct {
	IDRef  string `xml:"idref,attr"`
	Linear string `xml:"linear,attr,omitempty"`
}

func openFromZip(zr *zip.Reader, name string) (io.ReadCloser, error) {
	for _, f := range zr.File {
		if f.Name == name {
			return f.Open()
		}
	}
	return nil, fmt.Errorf("file %q not found in epub", name)
}

func decodeXML(zr *zip.Reader, name string, v any) error {
	rc, err := openFromZip(zr, name)
	if err != nil {
		return err
	}
	defer rc.Close()
	return xml.NewDecoder(rc).Decode(v)
}

// buildTitleMap строит карту filename → title из NCX (epub 2.0).
func buildTitleMap(zr *zip.Reader, itemByID map[string]opfItem, spine opfSpine, opfDir string) map[string]string {
	titles := make(map[string]string)

	if spine.Toc == "" {
		return titles
	}
	ncxItem, ok := itemByID[spine.Toc]
	if !ok {
		return titles
	}

	var ncx ncxDoc
	ncxPath := path.Join(opfDir, ncxItem.Href)
	if err := decodeXML(zr, ncxPath, &ncx); err != nil {
		return titles
	}

	for _, np := range ncx.NavPoints {
		src := np.Content.Src
		if idx := strings.IndexByte(src, '#'); idx >= 0 {
			src = src[:idx]
		}
		if src != "" && np.NavLabel.Text != "" {
			titles[src] = strings.TrimSpace(np.NavLabel.Text)
		}
	}
	return titles
}

func parseEpub(zr *zip.Reader) ([]domain.Chapter, error) {
	var container epubContainer
	if err := decodeXML(zr, "META-INF/container.xml", &container); err != nil {
		return nil, fmt.Errorf("read container.xml: %w", err)
	}

	opfPath := container.Rootfile.Path
	if opfPath == "" {
		return nil, fmt.Errorf("container.xml: rootfile path is empty")
	}

	var opf opfPackage
	if err := decodeXML(zr, opfPath, &opf); err != nil {
		return nil, fmt.Errorf("read opf %s: %w", opfPath, err)
	}

	itemByID := make(map[string]opfItem, len(opf.Manifest.Items))
	for _, item := range opf.Manifest.Items {
		itemByID[item.ID] = item
	}

	opfDir := path.Dir(opfPath)
	titleByHref := buildTitleMap(zr, itemByID, opf.Spine, opfDir)

	// Spine задаёт порядок чтения; фолбэк — весь manifest по порядку.
	readOrder := opf.Spine.ItemRefs
	if len(readOrder) == 0 {
		for _, item := range opf.Manifest.Items {
			readOrder = append(readOrder, opfItemRef{IDRef: item.ID})
		}
	}

	var chapters []domain.Chapter
	chapterNumber := 1

	for _, ref := range readOrder {
		if ref.Linear == "no" {
			continue
		}

		item, ok := itemByID[ref.IDRef]
		if !ok {
			continue
		}
		if item.MediaType != "application/xhtml+xml" && item.MediaType != "text/html" {
			continue
		}

		// Href может быть относительным, резолвим относительно директории OPF.
		itemPath := item.Href
		if !strings.HasPrefix(itemPath, "/") {
			itemPath = path.Join(opfDir, item.Href)
		}

		rc, err := openFromZip(zr, itemPath)
		if err != nil {
			return nil, fmt.Errorf("open chapter file %s: %w", itemPath, err)
		}

		_, content, err := extractChapterContent(rc)
		rc.Close()
		if err != nil {
			return nil, fmt.Errorf("extract content from %s: %w", itemPath, err)
		}

		if strings.TrimSpace(content) == "" {
			continue
		}

		chapters = append(chapters, domain.Chapter{
			Number:  chapterNumber,
			Title:   titleByHref[item.Href],
			Content: content,
		})
		chapterNumber++
	}

	return chapters, nil
}

// extractChapterContent парсит xhtml через xml.Decoder и возвращает заголовок и plain text.
// html.Parse (HTML5) не подходит — трактует <title/> как незакрытый тег и поглощает весь body.
func extractChapterContent(r io.Reader) (string, string, error) {
	decoder := xml.NewDecoder(r)
	decoder.Strict = false
	decoder.Entity = xml.HTMLEntity

	skipTags := map[string]bool{"style": true, "script": true}

	var sb strings.Builder
	var titleBuf strings.Builder
	inSkip := 0
	inTitle := false

	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}

		switch t := tok.(type) {
		case xml.StartElement:
			name := strings.ToLower(t.Name.Local)
			if skipTags[name] {
				inSkip++
			} else if name == "title" && inSkip == 0 {
				inTitle = true
			}
		case xml.EndElement:
			name := strings.ToLower(t.Name.Local)
			if inSkip > 0 {
				if skipTags[name] {
					inSkip--
				}
				continue
			}
			if name == "title" {
				inTitle = false
			}
		case xml.CharData:
			if inSkip > 0 {
				continue
			}
			text := strings.TrimSpace(string(t))
			if text == "" {
				continue
			}
			if inTitle {
				titleBuf.WriteString(text)
			} else {
				sb.WriteString(text)
				sb.WriteByte('\n')
			}
		}
	}

	return titleBuf.String(), sb.String(), nil
}
