package book_service

import (
	"FreeLib/internal/core/domain"
	"context"
	"fmt"
)

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
