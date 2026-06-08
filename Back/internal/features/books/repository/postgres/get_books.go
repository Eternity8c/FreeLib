package book_postgres_repository

import (
	"FreeLib/internal/core/domain"
	"context"
	"fmt"
)

func (r *BookRepositry) GetBooks(ctx context.Context, limit *int, offset *int) ([]domain.Book, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT 
        b.book_id, 
        b.title, 
        a.name_author, 
        g.name_genre, 
        b.created_at
    FROM freelib.books b
    JOIN freelib.author a ON b.author_id = a.author_id
    JOIN freelib.genre g ON b.genre_id = g.genre_id
	LIMIT $1
	OFFSET $2;
	`

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
