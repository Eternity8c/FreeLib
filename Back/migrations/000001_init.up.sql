CREATE SCHEMA freelib;

CREATE TABLE freelib.author (
    author_id SERIAL PRIMARY KEY,
    name_author VARCHAR(50)
);

CREATE TABLE freelib.genre (
    genre_id SERIAL PRIMARY KEY,
    name_genre VARCHAR(50)
);

CREATE TABLE freelib.books (
    book_id SERIAL PRIMARY KEY,
    title VARCHAR(50),
    author_id INT NOT NULL,
    genre_id INT NOT NULL,
    created_at TIMESTAMPTZ,
    FOREIGN KEY (author_id) REFERENCES freelib.author (author_id),
    FOREIGN KEY (genre_id) REFERENCES freelib.genre (genre_id)
);

CREATE TABLE freelib.users (
    user_id SERIAL PRIMARY KEY,
    email VARCHAR(50),
    pass_hash VARCHAR(50),
    is_admin BOOLEAN DEFAULT FALSE
);

CREATE TABLE freelib.favorite_book (
    user_id INT NOT NULL,
    book_id INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES freelib.users (user_id),
    FOREIGN KEY (book_id) REFERENCES freelib.books (book_id)
);