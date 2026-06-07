package book_postgres_repository

import (
	"FreeLib/internal/core/domain"
	"context"
	"fmt"
)

func (r *BookRepositry) CreateBook(
	ctx context.Context,
	book domain.Book,
) (domain.Book, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.Book{}, fmt.Errorf("create transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var authorID int
	var genreID int

	err = tx.QueryRow(ctx, `
		INSERT INTO freelib.author (name_author)
		VALUES ($1)
		ON CONFLICT (name_author) DO UPDATE SET name_author = EXCLUDED.name_author
		RETURNING author_id;`, book.Author,
	).Scan(&authorID)

	if err != nil {
		return domain.Book{}, fmt.Errorf("scan author ID: %w", err)
	}

	err = tx.QueryRow(ctx, `
		INSERT INTO freelib.genre (name_genre)
		VALUES ($1)
		ON CONFLICT (name_genre) DO UPDATE SET name_genre = EXCLUDED.name_genre
		RETURNING genre_id;
	`, book.Genre,
	).Scan(&genreID)

	if err != nil {
		return domain.Book{}, fmt.Errorf("scan genre ID: %w", err)
	}

	var bookModel BookModel

	err = tx.QueryRow(ctx, `
    WITH inserted_book AS (
        INSERT INTO freelib.books (title, author_id, genre_id, created_at)
        VALUES ($1, $2, $3, $4)
        RETURNING book_id, title, author_id, genre_id, created_at
    )
    SELECT 
        ib.book_id, 
        ib.title, 
        a.name_author, 
        g.name_genre, 
        ib.created_at
    FROM inserted_book ib
    JOIN freelib.author a ON ib.author_id = a.author_id
    JOIN freelib.genre g ON ib.genre_id = g.genre_id;
`, book.Title, authorID, genreID, book.CreatedAt,
	).Scan(
		&bookModel.ID,
		&bookModel.Title,
		&bookModel.Author,
		&bookModel.Genre,
		&bookModel.CreatedAt,
	)

	if err != nil {
		return domain.Book{}, fmt.Errorf("scan book model: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return domain.Book{}, fmt.Errorf("transaction commit: %w", err)
	}

	return domain.NewBook(
		bookModel.ID,
		bookModel.Title,
		bookModel.Author,
		bookModel.Genre,
		bookModel.CreatedAt,
	), nil
}
