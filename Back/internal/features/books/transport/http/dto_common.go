package books_transport_http

import "FreeLib/internal/core/domain"

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
