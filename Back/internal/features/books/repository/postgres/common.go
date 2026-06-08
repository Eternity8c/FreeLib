package book_postgres_repository

import (
	"FreeLib/internal/core/domain"
	"context"
	"fmt"
)

func (r *BookRepositry) queryBooks(ctx context.Context, query string, limit *int, offset *int) ([]domain.Book, error) {
	rows, err := r.pool.Query(ctx, query, limit, offset)
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
