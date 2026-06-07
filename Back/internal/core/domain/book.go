package domain

import (
	"fmt"
	"time"
	"unicode"
)

type Book struct {
	ID        int       `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Author    string    `json:"author" db:"author"`
	Genre     string    `json:"genre" db:"genre"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
}

func NewBookUninitialized(title string, author string, genre string) Book {
	return Book{
		ID:        UninitializedID,
		Title:     title,
		Author:    author,
		Genre:     genre,
		CreatedAt: time.Now().UTC(),
	}
}

func NewBook(id int, title string, author string, genre string, createdAt time.Time) Book {
	return Book{
		ID:        id,
		Title:     title,
		Author:    author,
		Genre:     genre,
		CreatedAt: createdAt,
	}
}

func (b *Book) Validate() error {
	titleLenght := len([]rune(b.Title))
	if titleLenght < 1 {
		return fmt.Errorf("len `title`: %d", titleLenght)
	}

	authorLenght := len([]rune(b.Author))
	if authorLenght < 3 {
		return fmt.Errorf("len `author`: %d", titleLenght)
	}

	runes := []rune(b.Author)
	if !unicode.IsUpper(runes[0]) {
		return fmt.Errorf("the first letter author not upper: %c", runes[0])
	}

	runes = []rune(b.Genre)
	if !unicode.IsUpper(runes[0]) {
		return fmt.Errorf("the first letter genre not upper: %c", runes[0])
	}

	return nil
}
