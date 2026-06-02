package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB(ctx context.Context, dbAddr string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dbAddr)
	if err != nil {
		return nil, err
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, err
	}
	return pool, nil
}
