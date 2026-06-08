package book_service

import (
	"FreeLib/internal/core/domain"
	"context"
)

type BookService struct {
	bookRepository BookRepository
}

type BookRepository interface {
	CreateBook(ctx context.Context, book domain.Book) (domain.Book, error)
	GetBooks(ctx context.Context, limit *int, offset *int) ([]domain.Book, error)
	GetNewBooks(ctx context.Context, limit *int, offset *int) ([]domain.Book, error)
}

func NewBookService(bookRepository BookRepository) *BookService {
	return &BookService{
		bookRepository: bookRepository,
	}
}
