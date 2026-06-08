package book_postgres_repository

import (
	"FreeLib/internal/core/domain"
	"context"
	"fmt"
)

func (r *BookRepositry) GetNewBooks(ctx context.Context, limit *int, offset *int) ([]domain.Book, error) {
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
	WHERE b.created_at >= NOW() - INTERVAL '7 days'
	LIMIT $1
	OFFSET $2;
	`

	bookDomains, err := r.queryBooks(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("getBooksLimitOffset: %w", err)
	}

	return bookDomains, nil
}
