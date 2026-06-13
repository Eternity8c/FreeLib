package books_transport_http

import (
	"FreeLib/internal/core/domain"
	"time"
)

type BookDTOResponce struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Genre  string `json:"genre"`
}

func bookDTOFromDomain(book domain.Book) BookDTOResponce {
	return BookDTOResponce{
		ID:     book.ID,
		Title:  book.Title,
		Genre:  book.Genre,
		Author: book.Author,
	}
}

func booksDTOFromDomains(books []domain.Book) []BookDTOResponce {
	bookDTO := make([]BookDTOResponce, len(books))
	for i, book := range books {
		bookDTO[i] = bookDTOFromDomain(book)
	}

	return bookDTO
}

type CreateBookRequest struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Genre  string `json:"genre"`
}

type CreateBookResponce BookDTOResponce

type FavoriteBookRequest struct {
	BookID int `json:"book_id"`
}

type FavoriteBookResponce struct {
	UserID int             `json:"user_id"`
	Book   BookDTOResponce `json:"book"`
}

type GetBookResponce BookDTOResponce

type GetBooksResponce []BookDTOResponce

type GetFavoriteBooksRecponce []BookDTOResponce

type GetNewBooksResponce []BookDTOResponce

type UpdateBookRequest struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Genre  string `json:"genre"`
}

type UpdateBookResponce struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Genre  string `json:"genre"`
}

func createBookDomainFromDTO(request CreateBookRequest) domain.Book {
	return domain.NewBookUninitialized(request.Title, request.Author, request.Genre)
}

func updateBookDomainFromDTO(request UpdateBookRequest) domain.Book {
	return domain.NewBook(
		request.ID,
		request.Title,
		request.Author,
		request.Genre,
		time.Now().UTC(),
	)
}
