package book_service

import (
	"fmt"

	core_errors "github.com/Eternity8c/FreeLib/internal/core/errors"
)

func validateLimitOffset(limit *int, offset *int) error {
	if limit != nil && *limit < 0 {
		return fmt.Errorf(
			"limit must be non-negative: %w",
			core_errors.ErrInvalidArgumment,
		)
	}

	if offset != nil && *offset < 0 {
		return fmt.Errorf(
			"offset must be non-negative: %w",
			core_errors.ErrInvalidArgumment,
		)
	}

	return nil
}

func validateID(id int) error {
	if id <= 0 {
		return fmt.Errorf("ID must be non-negstive: %w", core_errors.ErrInvalidArgumment)
	}

	return nil
}
