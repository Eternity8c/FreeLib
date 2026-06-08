package book_postgres_repository

import (
	"FreeLib/internal/core/domain"
	"context"
	"fmt"
)

func (r *BookRepositry) FavoriteBook(ctx context.Context, userID int, bookID int) (int, domain.Book, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO freelib.favorite_book (user_id, book_id)
	VALUES
	(
		(SELECT user_id FROM freelib.users WHERE user_id = $1),
		(SELECT book_id FROM freelib.books WHERE book_id = $2)
	)
	RETURNING user_id;
	`

	row := r.pool.QueryRow(ctx, query, userID, bookID)

	var uID int
	err := row.Scan(
		&uID,
	)
	if err != nil {
		return domain.UninitializedID, domain.Book{}, fmt.Errorf("row scan: %w", err)
	}

	bookDomain, err := r.GetBook(ctx, bookID)
	if err != nil {
		return domain.UninitializedID, domain.Book{}, fmt.Errorf("get book: %w", err)
	}
	return uID, bookDomain, nil
}
