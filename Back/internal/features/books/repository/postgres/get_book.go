package book_postgres_repository

import (
	"FreeLib/internal/core/domain"
	core_errors "FreeLib/internal/core/errors"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (r *BookRepositry) GetBook(ctx context.Context, id int) (domain.Book, error) {
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
	WHERE book_id = $1;
	`

	row := r.pool.QueryRow(ctx, query, id)

	var bookModel BookModel
	err := row.Scan(
		&bookModel.ID,
		&bookModel.Title,
		&bookModel.Author,
		&bookModel.Genre,
		&bookModel.CreatedAt,
	)
	if err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return domain.Book{}, fmt.Errorf("row scan: %w", core_errors.ErrInvalidArgumment)
		}
		return domain.Book{}, fmt.Errorf("row scan: %w", err)
	}

	bookDomain := bookDomainFromModel(bookModel)
	return bookDomain, nil
}
