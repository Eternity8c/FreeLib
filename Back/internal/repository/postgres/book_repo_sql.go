package postgres

import (
	"FreeLib/internal/models"
	"FreeLib/internal/repository"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type bookRepository struct {
	pool *pgxpool.Pool
}

func NewBookRepository(pool *pgxpool.Pool) repository.BookRepository {
	return &bookRepository{pool: pool}
}

func (r *bookRepository) GetAll() ([]models.Book, error) {
	query := `SELECT * FROM books;`
	rows, err := r.pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.Description,
			&book.Genre,
			&book.Content,
			&book.CoverURL,
			&book.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return books, nil
}

func (r *bookRepository) GetByID(id uint) (*models.Book, error) {
	query := `SELECT id, title, author, description, genre, content, cover_url, created_at
	 FROM books WHERE id = $1;`
	var book models.Book
	err := r.pool.QueryRow(context.Background(),
		query, id).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.Description,
		&book.Genre,
		&book.Content,
		&book.CoverURL,
		&book.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &book, nil
}

func (r *bookRepository) Create(book *models.Book) error {
	ctx := context.Background()
	var createdAt time.Time
	err := r.pool.QueryRow(ctx, `
		INSERT INTO books (title, author, description, genre, content, cover_url)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`, book.Title, book.Author, book.Description, book.Genre, book.Content, book.CoverURL).Scan(&book.ID, &createdAt)
	if err != nil {
		return err
	}
	t := createdAt.UTC()
	book.CreatedAt = t
	return nil
}

func (r *bookRepository) Delete(id uint) error {
	query := `DELETE FROM books WHERE id = $1`
	result, err := r.pool.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return fmt.Errorf("books with id %d not found", id)
	}

	return nil
}

func (r *bookRepository) Update(book *models.Book) error {
	query := `UPDATE books
	SET title = $1, author = $2, description = $3, genre = $4, content = $5, cover_url = $6
	WHERE id = $7`

	result, err := r.pool.Exec(context.Background(), query,
		book.Title,
		book.Author,
		book.Description,
		book.Genre,
		book.Content,
		book.CoverURL,
		book.ID)

	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("book with id %d not found", book.ID)
	}

	return nil
}

func (r *bookRepository) AddFavorite(userID uint, bookID uint) error {
	query := `INSERT INTO favorite_books (user_id, book_id)
	VALUES ($1, $2)`

	_, err := r.pool.Exec(context.Background(), query,
		userID,
		bookID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *bookRepository) DeleteFavorite(userID uint, bookID uint) error {
	query := `DELETE FROM favorite_books WHERE user_id = $1 AND book_id = $2`
	_, err := r.pool.Exec(context.Background(), query, userID, bookID)

	if err != nil {
		return err
	}

	return nil
}

func (r *bookRepository) GetAllFavorite(userID uint) ([]models.Book, error) {
	query := `SELECT book_id FROM favorite_books WHERE user_id = $1`

	rows, err := r.pool.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var bookId int64
		if err := rows.Scan(&bookId); err != nil {
			return nil, err
		}

		book, err := r.GetByID(uint(bookId))
		if err != nil {
			return nil, err
		}

		books = append(books, *book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}
