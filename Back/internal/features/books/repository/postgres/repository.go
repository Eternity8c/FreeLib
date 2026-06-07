package book_postgres_repository

import core_postgres_pool "FreeLib/internal/core/repository/postgres/pool"

type BookRepositry struct {
	pool core_postgres_pool.Pool
}

func NewBookRepository(pool core_postgres_pool.Pool) *BookRepositry {
	return &BookRepositry{
		pool: pool,
	}
}
