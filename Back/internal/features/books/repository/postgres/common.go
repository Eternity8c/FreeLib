package book_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Eternity8c/FreeLib/internal/core/domain"
	core_errors "github.com/Eternity8c/FreeLib/internal/core/errors"
	"github.com/jackc/pgx/v5"
)

func (r *BookRepositry) queryBooks(ctx context.Context, query string, arg ...any) ([]domain.Book, error) {
	rows, err := r.pool.Query(ctx, query, arg...)
	if err != nil {
		return nil, fmt.Errorf("select users: %w", err)
	}
	defer rows.Close()

	var bookModels []BookModel
	for rows.Next() {
		var bookModel BookModel

		err := rows.Scan(
			&bookModel.ID,
			&bookModel.Title,
			&bookModel.Author,
			&bookModel.Genre,
			&bookModel.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan books: %w", err)
		}

		bookModels = append(bookModels, bookModel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	bookDomains := bookDomainsFromModels(bookModels)

	return bookDomains, nil
}

func (r *BookRepositry) GetFileHashFromBook(ctx context.Context, id int) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT b.file_hash
	FROM freelib.books b
	WHERE book_id = $1;
	`

	row := r.pool.QueryRow(ctx, query, id)

	var fileHash string
	err := row.Scan(&fileHash)
	if err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return "", fmt.Errorf("row scan: %w", core_errors.ErrInvalidArgumment)
		}
		return "", fmt.Errorf("row scan: %w", err)
	}

	return fileHash, nil
}

func (r *BookRepositry) GetS3URLFromBook(ctx context.Context, id int) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT b.s3_url
	FROM freelib.books b
	WHERE book_id = $1;
	`

	row := r.pool.QueryRow(ctx, query, id)

	var s3URL string
	err := row.Scan(&s3URL)
	if err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return "", fmt.Errorf("row scan: %w", core_errors.ErrInvalidArgumment)
		}
		return "", fmt.Errorf("row scan: %w", err)
	}

	return s3URL, nil
}
