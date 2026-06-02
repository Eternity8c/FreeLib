package domain

import "time"

type Book struct {
	ID          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Author      string    `json:"author" db:"author"`
	Description string    `json:"description" db:"description"`
	Genre       string    `json:"genre" db:"genre"`
	Content     string    `json:"content" db:"content"`
	CoverURL    string    `json:"cover_url" db:"cover_url"`
	CreatedAt   time.Time `json:"created_at,omitempty" db:"created_at"`
}
