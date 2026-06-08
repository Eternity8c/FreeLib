package book_postgres_repository

import (
	"FreeLib/internal/core/domain"
	"time"
)

type BookModel struct {
	ID        int
	Title     string
	Author    string
	Genre     string
	CreatedAt time.Time
}

func bookDomainsFromModels(books []BookModel) []domain.Book {
	bookDomains := make([]domain.Book, len(books))
	for i, book := range books {
		bookDomains[i] = domain.NewBook(
			book.ID,
			book.Title,
			book.Author,
			book.Genre,
			book.CreatedAt,
		)
	}

	return bookDomains
}
