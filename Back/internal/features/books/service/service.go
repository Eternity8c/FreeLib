package book_service

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/Eternity8c/FreeLib/internal/core/domain"
	core_errors "github.com/Eternity8c/FreeLib/internal/core/errors"
)

type BookService struct {
	bookRepository   BookRepository
	bookS3Repository BookS3Repository
}

type BookRepository interface {
	CreateBook(ctx context.Context, book domain.Book, fileURL string) (domain.Book, error)
	GetBooks(ctx context.Context, limit *int, offset *int) ([]domain.Book, error)
	GetNewBooks(ctx context.Context, limit *int, offset *int) ([]domain.Book, error)
	GetBook(ctx context.Context, id int) (domain.Book, error)
	FavoriteBook(ctx context.Context, userID int, bookID int) (int, domain.Book, error)
	GetFavoriteBooks(ctx context.Context, userID int) ([]domain.Book, error)
	GetBooksByGenre(ctx context.Context, genre string) ([]domain.Book, error)
	UpdateBook(ctx context.Context, book domain.Book) (domain.Book, error)
	DeleteBook(ctx context.Context, bookID int) error
}

type BookS3Repository interface {
	SaveBookFile(ctx context.Context, file multipart.File, fileName string) (string, error)
}

func NewBookService(bookRepository BookRepository, bookS3Repository BookS3Repository) *BookService {
	return &BookService{
		bookRepository:   bookRepository,
		bookS3Repository: bookS3Repository,
	}
}

func (s *BookService) CreateBook(
	ctx context.Context,
	book domain.Book,
	file multipart.File,
	fileHeader *multipart.FileHeader,
) (domain.Book, error) {
	if err := book.Validate(); err != nil {
		return domain.Book{}, fmt.Errorf("validate book domain: %w", err)
	}

	extencion := filepath.Ext(fileHeader.Filename)
	if extencion != ".epub" {
		return domain.Book{}, fmt.Errorf("failed extencion: %s: %w", extencion, core_errors.ErrInvalidArgumment)
	}

	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), fileHeader.Filename)

	fileURL, err := s.bookS3Repository.SaveBookFile(ctx, file, fileName)
	if err != nil {
		return domain.Book{}, fmt.Errorf("save book repository: %w", err)
	}

	domainBook, err := s.bookRepository.CreateBook(ctx, book, fileURL)
	if err != nil {
		return domain.Book{}, fmt.Errorf("create book: %w", err)
	}

	return domainBook, nil
}

func (s *BookService) FavoriteBook(ctx context.Context, userID int, bookID int) (int, domain.Book, error) {
	if userID == domain.UninitializedID {
		return domain.UninitializedID, domain.Book{}, fmt.Errorf("userID must be non-negative: %w", core_errors.ErrInvalidArgumment)
	}
	if bookID <= 0 {
		return domain.UninitializedID, domain.Book{}, fmt.Errorf("bookID must be non-negative: %w", core_errors.ErrInvalidArgumment)
	}

	uID, domainBook, err := s.bookRepository.FavoriteBook(ctx, userID, bookID)
	if err != nil {
		return domain.UninitializedID, domain.Book{}, fmt.Errorf("favorite book from repository: %w", err)
	}

	return uID, domainBook, err
}

func (s *BookService) GetBook(ctx context.Context, id int) (domain.Book, error) {
	if err := validateID(id); err != nil {
		return domain.Book{}, fmt.Errorf("validate ID: %w", err)
	}

	book, err := s.bookRepository.GetBook(ctx, id)
	if err != nil {
		return domain.Book{}, fmt.Errorf("get book from repository: %w", err)
	}

	return book, nil
}

func (s *BookService) GetBooks(ctx context.Context, limit *int, offset *int) ([]domain.Book, error) {
	if err := validateLimitOffset(limit, offset); err != nil {
		return nil, fmt.Errorf("validate limit offset: %w", err)
	}

	books, err := s.bookRepository.GetBooks(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get books from repository: %w", err)
	}

	return books, nil
}

func (s *BookService) GetFavoriteBooks(ctx context.Context, userID int) ([]domain.Book, error) {
	bookDomains, err := s.bookRepository.GetFavoriteBooks(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get book domains from repository: %w", err)
	}

	return bookDomains, nil
}

func (s *BookService) GetNewBooks(ctx context.Context, limit *int, offset *int) ([]domain.Book, error) {
	if err := validateLimitOffset(limit, offset); err != nil {
		return nil, fmt.Errorf("validate limit offset: %w", err)
	}

	books, err := s.bookRepository.GetNewBooks(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get new books from repository: %w", err)
	}

	return books, nil
}

func (s *BookService) GetBooksByGenre(ctx context.Context, genre string) ([]domain.Book, error) {
	genreLenght := len([]rune(genre))
	if genreLenght < 3 {
		return nil, fmt.Errorf("genre len: %d: %w", genreLenght, core_errors.ErrInvalidArgumment)
	}

	bookDomains, err := s.bookRepository.GetBooksByGenre(ctx, genre)
	if err != nil {
		return nil, fmt.Errorf("get book by genre from repository: %w", err)
	}

	return bookDomains, nil
}

func (s *BookService) UpdateBook(ctx context.Context, book domain.Book) (domain.Book, error) {
	if err := book.Validate(); err != nil {
		return domain.Book{}, fmt.Errorf("validate book domain: %w", err)
	}

	bookDomain, err := s.bookRepository.UpdateBook(ctx, book)
	if err != nil {
		return domain.Book{}, fmt.Errorf("update book: %w", err)
	}

	return bookDomain, nil
}

func (s *BookService) DeleteBook(ctx context.Context, bookID int) error {
	err := s.bookRepository.DeleteBook(ctx, bookID)
	if err != nil {
		return fmt.Errorf("delete book from repository: %w", err)
	}

	return nil
}
