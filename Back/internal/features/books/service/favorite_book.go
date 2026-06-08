package book_service

import (
	"FreeLib/internal/core/domain"
	core_errors "FreeLib/internal/core/errors"
	"context"
	"fmt"
)

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
