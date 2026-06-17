package book_service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"

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

func CalculateFileHash(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("file open: %w", err)
	}
	defer file.Close()

	hasher := sha256.New()

	if _, err := io.Copy(hasher, file); err != nil {
		return "", fmt.Errorf("cope file: %w", err)
	}

	hashInBytes := hasher.Sum(nil)

	return hex.EncodeToString(hashInBytes), nil
}
