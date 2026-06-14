package book_postgres_repository

import (
	"time"

	"github.com/Eternity8c/FreeLib/internal/core/domain"
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

func bookDomainFromModel(book BookModel) domain.Book {
	return domain.NewBook(
		book.ID,
		book.Title,
		book.Author,
		book.Genre,
		book.CreatedAt,
	)
}
