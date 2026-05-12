package handlers

import (
	"FreeLib/internal/models"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type BookRepository interface {
	GetAll() ([]models.Book, error)
	GetByID(id uint) (*models.Book, error)
	Create(book *models.Book) error
	Delete(id uint) error
	Update(book *models.Book) error
	AddFavorite(userID uint, bookID uint) error
	DeleteFavorite(userID uint, bookID uint) error
	GetAllFavorite(userID uint) ([]models.Book, error)
}

type BookHandler struct {
	repo BookRepository
}

func NewBookHandler(repo BookRepository) *BookHandler {
	return &BookHandler{
		repo: repo,
	}
}

func (h *BookHandler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"service": "FreeLib API",
		"version": "1.0.0",
	})

}

func (h *BookHandler) GetBooksHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetBooksHandler called from=%s", r.RemoteAddr)
	books, err := h.repo.GetAll()
	if err != nil {
		log.Println("GetBooksHandler: repo error:", err)
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(books)
	log.Printf("GetBooksHandler done count=%d", len(books))
}

func (h *BookHandler) GetByIDHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetByIDHandler called from=%s", r.RemoteAddr)
	idStr := r.URL.Query().Get("id")

	idUint, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("GetByIDHandler: invalid id:", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	book, err := h.repo.GetByID(uint(idUint))
	if err != nil {
		log.Println("GetByIDHandler: repo error:", err)
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(book)
	log.Printf("GetByIDHandler done id=%d title=%q author=%q", idUint, book.Title, book.Author)
}

func (h *BookHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("CreateHandler called from=%s", r.RemoteAddr)
	var book models.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		log.Println("CreateHandler: invalid JSON:", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.repo.Create(&book); err != nil {
		log.Println("CreateHandler: repo error:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)

	log.Printf("Book created id=%d", book.ID)
}

func (h *BookHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("DeleteHandler called from=%s", r.RemoteAddr)
	strID := r.URL.Query().Get("id")
	intId, err := strconv.Atoi(strID)
	if err != nil {
		log.Println("DeleteHandler: invalid id:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(uint(intId)); err != nil {
		log.Println("DeleteHandler: repo error:", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Printf("Book deleted id=%d", intId)
}

func (h *BookHandler) UpdateBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	book, err := h.repo.GetByID(uint(id))
	if err != nil {
		log.Println("get by id failed:", err)
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	type partialReq struct {
		Title       *string `json:"title"`
		Author      *string `json:"author"`
		Description *string `json:"description"`
		Genre       *string `json:"genre"`
		Content     *string `json:"content"`
		CoverURL    *string `json:"coverUrl"`
		PublishYear *int    `json:"publishYear"`
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("read body:", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	r.Body = io.NopCloser(bytes.NewReader(body))

	var rawMap map[string]interface{}
	if err := json.Unmarshal(body, &rawMap); err != nil {
		rawMap = map[string]interface{}{}
	}

	var req partialReq
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&req); err != nil {
		log.Println("decode body to struct:", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.CoverURL == nil {
		if v, ok := rawMap["cover_url"]; ok {
			if s, ok2 := v.(string); ok2 {
				req.CoverURL = &s
			}
		} else if v, ok := rawMap["coverURL"]; ok {
			if s, ok2 := v.(string); ok2 {
				req.CoverURL = &s
			}
		}
	}

	if req.Title != nil {
		book.Title = *req.Title
	}
	if req.Author != nil {
		book.Author = *req.Author
	}
	if req.Description != nil {
		book.Description = *req.Description
	}
	if req.Genre != nil {
		book.Genre = *req.Genre
	}
	if req.Content != nil {
		book.Content = *req.Content
	}
	if req.CoverURL != nil {
		book.CoverURL = *req.CoverURL
	}

	if err := h.repo.Update(book); err != nil {
		log.Println("update failed:", err)
		http.Error(w, "Failed to update book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(book); err != nil {
		log.Println("encode response err:", err)
	}

	log.Printf("Book updated id=%d", book.ID)
}

func (h *BookHandler) AddFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idUserStr := vars["id"]
	idUserInt, err := strconv.Atoi(idUserStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	type addFavReq struct {
		BookID uint `json:"book_id"`
	}

	var req addFavReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	if err := h.repo.AddFavorite(uint(idUserInt), uint(req.BookID)); err != nil {
		http.Error(w, "Failed to add favorite", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"ok":true}`))
	log.Printf("Favorite added user_id=%v book_id=%v from=%s", idUserInt, req.BookID, r.RemoteAddr)
}

func (h *BookHandler) DeleteFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	idUserStr := vars["user_id"]
	idUserInt, err := strconv.Atoi(idUserStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	idBookStr := vars["book_id"]
	idBookInt, err := strconv.Atoi(idBookStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = h.repo.DeleteFavorite(uint(idUserInt), uint(idBookInt))
	if err != nil {
		log.Printf("DeleteFavorite error: %v\n", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Printf("Favorite removed user_id=%v book_id=%v from=%s", idUserInt, idBookInt, r.RemoteAddr)
}

func (h *BookHandler) GetAllFavotiteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	idUserStr := vars["user_id"]
	idUserInt, err := strconv.Atoi(idUserStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err != nil {
		log.Printf("GetAllFavorite failed for user %d: %v\n", idUserInt, err)
		http.Error(w, "Failed to load favorites", http.StatusInternalServerError)
		return
	}

	books, err := h.repo.GetAllFavorite(uint(idUserInt))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(books); err != nil {
		log.Printf("encode favorites response failed: %v", err)
	}
	log.Printf("Favorites returned user_id=%v count=%d from=%s", idUserInt, len(books), r.RemoteAddr)
}
