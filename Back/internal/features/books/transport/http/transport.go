package books_transport_http

import (
	"FreeLib/internal/core/domain"
	core_logger "FreeLib/internal/core/logger"
	core_http_middleware "FreeLib/internal/core/transport/http/middleware"
	core_http_request "FreeLib/internal/core/transport/http/request"
	core_http_responce "FreeLib/internal/core/transport/http/responce"
	core_http_server "FreeLib/internal/core/transport/http/server"
	"context"
	"net/http"
)

type BooksHTTPHandler struct {
	bookServices BookServices
}

type BookServices interface {
	CreateBook(ctx context.Context, book domain.Book) (domain.Book, error)
	GetBooks(ctx context.Context, limit *int, offset *int) ([]domain.Book, error)
	GetNewBooks(ctx context.Context, limit *int, offset *int) ([]domain.Book, error)
	GetBook(ctx context.Context, id int) (domain.Book, error)
	FavoriteBook(ctx context.Context, userID int, bookID int) (int, domain.Book, error)
	GetFavoriteBooks(ctx context.Context, userID int) ([]domain.Book, error)
	GetBooksByGenre(ctx context.Context, genre string) ([]domain.Book, error)
	UpdateBook(ctx context.Context, book domain.Book) (domain.Book, error)
	DeleteBook(ctx context.Context, bookID int) error
}

func NewBookHTTPHandler(bookServices BookServices) *BooksHTTPHandler {
	return &BooksHTTPHandler{
		bookServices: bookServices,
	}
}

func (h *BooksHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/books",
			Handler: core_http_middleware.AdminOnly(h.CreateBook),
		},
		{
			Method:  http.MethodGet,
			Path:    "/books",
			Handler: h.GetBooks,
		},
		{
			Method:  http.MethodGet,
			Path:    "/books/new",
			Handler: h.GetNewBooks,
		},
		{
			Method:  http.MethodGet,
			Path:    "/book",
			Handler: h.GetBook,
		},
		{
			Method: http.MethodPost,
			Path:   "/book",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				core_http_middleware.Authorization()(http.HandlerFunc(h.FavoriteBook)).ServeHTTP(w, r)
			},
		},
		{
			Method: http.MethodGet,
			Path:   "/books/favorite",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				core_http_middleware.Authorization()(http.HandlerFunc(h.GetFavoriteBooks)).ServeHTTP(w, r)
			},
		},
		{
			Method:  http.MethodPut,
			Path:    "/book",
			Handler: core_http_middleware.AdminOnly(h.UpdateBook),
		},
		{
			Method:  http.MethodDelete,
			Path:    "/book",
			Handler: core_http_middleware.AdminOnly(h.DeleteBook),
		},
	}
}

func (h *BooksHTTPHandler) CreateBook(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	responceHandler := core_http_responce.NewHTTPResponceHandler(log, rw)

	log.Debug("invoke CreateBook handler")

	var request CreateBookRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responceHandler.ErrorResponce(err, "failed to validate and decode HTTP request")
		return
	}

	bookDomain, err := h.bookServices.CreateBook(ctx, createBookDomainFromDTO(request))
	if err != nil {
		responceHandler.ErrorResponce(err, "failed to create book")
		return
	}

	responce := CreateBookResponce(bookDTOFromDomain(bookDomain))
	responceHandler.JSONResponce(responce, http.StatusCreated)
}

func (h *BooksHTTPHandler) FavoriteBook(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responceHandler := core_http_responce.NewHTTPResponceHandler(log, rw)

	log.Debug("invoke favorite book handler")

	userID, err := idFromJWTToken(ctx)
	if err != nil {
		responceHandler.ErrorResponce(err, "failed get id from JWT token")
		return
	}

	var request FavoriteBookRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responceHandler.ErrorResponce(err, "failed decode and validate request")
		return
	}

	uID, bookDomain, err := h.bookServices.FavoriteBook(ctx, userID, request.BookID)
	if err != nil {
		responceHandler.ErrorResponce(err, "failed favorite book")
		return
	}

	responce := FavoriteBookResponce{
		UserID: uID,
		Book:   bookDTOFromDomain(bookDomain),
	}
	responceHandler.JSONResponce(responce, http.StatusOK)
}

func (h *BooksHTTPHandler) GetBook(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responceHandler := core_http_responce.NewHTTPResponceHandler(log, rw)

	log.Debug("invoke get book")
	id, err := getIDQueryParam(r)
	if err != nil {
		responceHandler.ErrorResponce(err, "failed to get ID query param")
		return
	}

	bookDomain, err := h.bookServices.GetBook(ctx, id)
	if err != nil {
		responceHandler.ErrorResponce(err, "failed to get book")
		return
	}

	responce := GetBookResponce(bookDTOFromDomain(bookDomain))
	responceHandler.JSONResponce(responce, http.StatusOK)
}

func (h *BooksHTTPHandler) GetBooks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responceHandler := core_http_responce.NewHTTPResponceHandler(log, rw)

	log.Debug("invoke GetBooks handler")

	limit, offset, err := getLimitOffsetQueryParams(r)
	if err != nil {
		responceHandler.ErrorResponce(err, "failed to get limit offset query param")
		return
	}

	var bookDomains []domain.Book
	genre := getGenreQueryParam(r)
	if genre == "" {
		bookDomains, err = h.bookServices.GetBooks(ctx, limit, offset)
		if err != nil {
			responceHandler.ErrorResponce(err, "failed to get books")
			return
		}
	} else {
		bookDomains, err = h.bookServices.GetBooksByGenre(ctx, genre)
		if err != nil {
			responceHandler.ErrorResponce(err, "failed get books by genre")
		}
	}

	responce := GetBooksResponce(booksDTOFromDomains(bookDomains))
	responceHandler.JSONResponce(responce, http.StatusOK)
}

func (h *BooksHTTPHandler) GetFavoriteBooks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responceHandler := core_http_responce.NewHTTPResponceHandler(log, rw)

	log.Debug("invoke GetFavoriteBooks")

	userID, err := idFromJWTToken(ctx)
	if err != nil {
		responceHandler.ErrorResponce(err, "failed get ID from JWT token")
		return
	}

	bookDomains, err := h.bookServices.GetFavoriteBooks(ctx, userID)
	if err != nil {
		responceHandler.ErrorResponce(err, "failed to get favorite books")
		return
	}

	responce := GetFavoriteBooksRecponce(booksDTOFromDomains(bookDomains))
	responceHandler.JSONResponce(responce, http.StatusOK)
}

func (h *BooksHTTPHandler) GetNewBooks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responceHandler := core_http_responce.NewHTTPResponceHandler(log, rw)

	log.Debug("invoke GetNewBooks handler")

	limit, offset, err := getLimitOffsetQueryParams(r)
	if err != nil {
		responceHandler.ErrorResponce(err, "failed to get limit offset query param")
		return
	}

	booksDomain, err := h.bookServices.GetNewBooks(ctx, limit, offset)
	if err != nil {
		responceHandler.ErrorResponce(err, "failed to get new books")
		return
	}

	responce := GetNewBooksResponce(booksDTOFromDomains(booksDomain))
	responceHandler.JSONResponce(responce, http.StatusOK)
}

func (h *BooksHTTPHandler) UpdateBook(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responceHandler := core_http_responce.NewHTTPResponceHandler(log, rw)

	log.Debug("invoke update book handler")

	var request UpdateBookRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responceHandler.ErrorResponce(err, "failed decode and validate request")
		return
	}

	bookDomain, err := h.bookServices.UpdateBook(ctx, updateBookDomainFromDTO(request))
	if err != nil {
		responceHandler.ErrorResponce(err, "failed update book from repository")
		return
	}

	responce := UpdateBookResponce(bookDTOFromDomain(bookDomain))
	responceHandler.JSONResponce(responce, http.StatusOK)
}

func (h *BooksHTTPHandler) DeleteBook(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responceHandler := core_http_responce.NewHTTPResponceHandler(log, rw)

	log.Debug("invoke delete book handler")

	bookID, err := getIDQueryParam(r)
	if err != nil {
		responceHandler.ErrorResponce(err, "failed get ID query param")
		return
	}

	err = h.bookServices.DeleteBook(ctx, bookID)
	if err != nil {
		responceHandler.ErrorResponce(err, "failed delete book from repository")
		return
	}

	responceHandler.NoContentResponce()
}
