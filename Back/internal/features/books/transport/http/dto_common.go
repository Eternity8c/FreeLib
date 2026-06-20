package books_transport_http

import (
	"time"

	"github.com/Eternity8c/FreeLib/internal/core/domain"
)

type BookDTOResponce struct {
	ID     int    `json:"id" example:"1"`
	Title  string `json:"title" example:"The Great Gatsby"`
	Author string `json:"author" example:"F. Scott Fitzgerald"`
	Genre  string `json:"genre" example:"Fiction"`
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
	Title  string `form:"title" example:"The Great Gatsby"`
	Author string `form:"author" example:"F. Scott Fitzgerald"`
	Genre  string `form:"genre" example:"Fiction"`
}

type CreateBookResponce BookDTOResponce

type FavoriteBookRequest struct {
	BookID int `json:"book_id" example:"1"`
}

type FavoriteBookResponce struct {
	UserID int             `json:"user_id" example:"1"`
	Book   BookDTOResponce `json:"book"`
}

type GetBookResponce BookDTOResponce

type GetBooksResponce []BookDTOResponce

type GetFavoriteBooksRecponce []BookDTOResponce

type GetNewBooksResponce []BookDTOResponce

type UpdateBookRequest struct {
	ID     int    `form:"id" example:"1"`
	Title  string `form:"title" example:"The Great Gatsby"`
	Author string `form:"author" example:"F. Scott Fitzgerald"`
	Genre  string `form:"genre" example:"Fiction"`
}

type UpdateBookResponce struct {
	ID     int    `json:"id" example:"1"`
	Title  string `json:"title" example:"The Great Gatsby"`
	Author string `json:"author" example:"F. Scott Fitzgerald"`
	Genre  string `json:"genre" example:"Fiction"`
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
