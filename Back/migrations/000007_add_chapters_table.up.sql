CREATE TABLE freelib.chapters (
    chapters_id SERIAL PRIMARY KEY,
    book_id INT NOT NULL,
    chapters_number INT NOT NULL,
    title VARCHAR(255),
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ,

    FOREIGN KEY (book_id) REFERENCES freelib.books (book_id) ON DELETE CASCADE,

    CONSTRAINT unique_book_chapter UNIQUE (book_id, chapters_number)
);