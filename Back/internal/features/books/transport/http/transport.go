package books_transport_http

import (
	"FreeLib/internal/core/domain"
	core_http_middleware "FreeLib/internal/core/transport/http/middleware"
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
	}
}
