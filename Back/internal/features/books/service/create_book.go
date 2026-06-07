package book_service

import (
	"FreeLib/internal/core/domain"
	"context"
	"fmt"
)

func (s *BookService) CreateBook(ctx context.Context, book domain.Book) (domain.Book, error) {
	if err := book.Validate(); err != nil {
		return domain.Book{}, fmt.Errorf("validate book domain: %w", err)
	}

	domainBook, err := s.bookrepository.CreateBook(ctx, book)
	if err != nil {
		return domain.Book{}, fmt.Errorf("create book: %w", err)
	}

	return domainBook, nil
}
