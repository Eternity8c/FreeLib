package book_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Eternity8c/FreeLib/internal/core/domain"
	core_errors "github.com/Eternity8c/FreeLib/internal/core/errors"
	core_postgres_pool "github.com/Eternity8c/FreeLib/internal/core/repository/postgres/pool"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type BookRepositry struct {
	pool core_postgres_pool.Pool
}

func NewBookRepository(pool core_postgres_pool.Pool) *BookRepositry {
	return &BookRepositry{
		pool: pool,
	}
}

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

	bookDomains, err := r.queryBooks(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get query books: %w", err)
	}

	return bookDomains, nil
}

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

func (r *BookRepositry) GetFavoriteBooks(ctx context.Context, userID int) ([]domain.Book, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT 
		b.book_id, 
		b.title, 
		a.name_author, 
		g.name_genre, 
		b.created_at
	FROM freelib.favorite_book fb
	JOIN freelib.books b ON fb.book_id = b.book_id
	JOIN freelib.author a ON b.author_id = a.author_id
	JOIN freelib.genre g ON b.genre_id = g.genre_id
	WHERE fb.user_id = $1;
	`

	bookDomains, err := r.queryBooks(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("get query books: %w", err)
	}

	return bookDomains, nil
}

func (r *BookRepositry) FavoriteBook(ctx context.Context, userID int, bookID int) (int, domain.Book, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO freelib.favorite_book (user_id, book_id)
	VALUES ($1, $2)
	ON CONFLICT (user_id, book_id) DO NOTHING
	RETURNING user_id, book_id;
	`

	row := r.pool.QueryRow(ctx, query, userID, bookID)

	var returnedUserID int
	var returnedBookID int
	err := row.Scan(
		&returnedUserID,
		&returnedBookID,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" {
				return domain.UninitializedID, domain.Book{}, fmt.Errorf("invalid input: %w", core_errors.ErrInvalidArgumment)
			}
		}
		if errors.Is(pgx.ErrNoRows, err) {
			return domain.UninitializedID, domain.Book{}, fmt.Errorf("the book is already in favorite: %w", core_errors.ErrInvalidArgumment)
		}
		return domain.UninitializedID, domain.Book{}, fmt.Errorf("row scan: %w", err)
	}

	bookDomain, err := r.GetBook(ctx, returnedBookID)
	if err != nil {
		return domain.UninitializedID, domain.Book{}, fmt.Errorf("get book: %w", err)
	}
	return returnedUserID, bookDomain, nil
}

func (r *BookRepositry) GetBooksByGenre(ctx context.Context, genre string) ([]domain.Book, error) {
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
	WHERE name_genre = $1
	`

	bookDomains, err := r.queryBooks(ctx, query, genre)
	if err != nil {
		return nil, fmt.Errorf("get query books: %w", err)
	}

	return bookDomains, nil
}

func (r *BookRepositry) UpdateBook(ctx context.Context, book domain.Book) (domain.Book, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.Book{}, fmt.Errorf("begin  transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var authorID, genreID int
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
		UPDATE freelib.books
		SET title = $1, author_id = $2, genre_id = $3, created_at = $4
		WHERE book_id = $5
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
`, book.Title, authorID, genreID, book.CreatedAt, book.ID,
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

func (r *BookRepositry) DeleteBook(ctx context.Context, bookID int) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	DELETE FROM freelib.books WHERE book_id = $1;
	`

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	tag, err := tx.Exec(ctx, query, bookID)
	if err != nil {
		return fmt.Errorf("transaction exec: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("rows affected %d: %w", tag.RowsAffected(), core_errors.ErrInvalidArgumment)
	}

	tx.Commit(ctx)

	return nil
}
