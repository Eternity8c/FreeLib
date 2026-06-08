package book_service

import (
	"FreeLib/internal/core/domain"
	core_errors "FreeLib/internal/core/errors"
	"context"
	"fmt"
)

func (s *BookService) GetBooks(ctx context.Context, limit *int, offset *int) ([]domain.Book, error) {
	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf(
			"limit must be non-negative: %w",
			core_errors.ErrInvalidArgumment,
		)
	}

	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf(
			"offset must be non-negative: %w",
			core_errors.ErrInvalidArgumment,
		)
	}

	books, err := s.bookRepository.GetBooks(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get books from repository: %w", err)
	}

	return books, nil
}
