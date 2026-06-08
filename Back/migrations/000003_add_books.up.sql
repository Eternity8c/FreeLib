INSERT INTO freelib.books (title, author_id, genre_id, created_at)
VALUES 
    (
        'Анна Каренина', 
        (SELECT author_id FROM freelib.author WHERE name_author = 'Лев Толстой'),
        (SELECT genre_id FROM freelib.genre WHERE name_genre = 'Классика'),
        NOW()
    ),
    (
        'Идиот', 
        (SELECT author_id FROM freelib.author WHERE name_author = 'Федор Достоевский'),
        (SELECT genre_id FROM freelib.genre WHERE name_genre = 'Классика'),
        NOW()
    ),
    (
        'Палата №6', 
        (SELECT author_id FROM freelib.author WHERE name_author = 'Антон Чехов'),
        (SELECT genre_id FROM freelib.genre WHERE name_genre = 'Классика'),
        NOW()
    ),
    (
        'Скотный двор', 
        (SELECT author_id FROM freelib.author WHERE name_author = 'Джордж Оруэлл'),
        (SELECT genre_id FROM freelib.genre WHERE name_genre = 'Фантастика'),
        NOW()
    ),
    (
        'Медный всадник', 
        (SELECT author_id FROM freelib.author WHERE name_author = 'Александр Пушкин'),
        (SELECT genre_id FROM freelib.genre WHERE name_genre = 'Поэзия'),
        NOW()
    );