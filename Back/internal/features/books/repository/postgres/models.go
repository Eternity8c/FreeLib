package book_postgres_repository

import "time"

type BookModel struct {
	ID        int
	Title     string
	Author    string
	Genre     string
	CreatedAt time.Time
}
