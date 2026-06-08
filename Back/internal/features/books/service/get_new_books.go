package book_service

import (
	"FreeLib/internal/core/domain"
	"context"
	"fmt"
)

func (s *BookService) GetNewBooks(ctx context.Context, limit *int, offset *int) ([]domain.Book, error) {
	if err := validateLimitOffset(limit, offset); err != nil {
		return nil, fmt.Errorf("validate limit offset: %w", err)
	}

	books, err := s.bookRepository.GetNewBooks(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get new books from repository")
	}

	return books, nil
}
