package book_service

import (
	core_errors "FreeLib/internal/core/errors"
	"fmt"
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
