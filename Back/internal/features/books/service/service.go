package book_service

import (
	"FreeLib/internal/core/domain"
	"context"
)

type BookService struct {
	bookrepository BookRepository
}

type BookRepository interface {
	CreateBook(ctx context.Context, book domain.Book) (domain.Book, error)
}

func NewBookService(bookRepository BookRepository) *BookService {
	return &BookService{
		bookrepository: bookRepository,
	}
}
