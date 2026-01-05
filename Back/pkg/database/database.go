package database

import (
	"FreeLib/pkg/config"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB(ctx context.Context, cfg config.Config) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBname, cfg.SSLMode,
	)

	var pool *pgxpool.Pool
	var err error

	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		log.Printf("Attempt %d to connect to database...", attempt)

		// Создаем контекст с таймаутом для каждой попытки
		attemptCtx, cancel := context.WithTimeout(ctx, 5*time.Second)

		pool, err = pgxpool.New(attemptCtx, dsn)
		cancel() // Освобождаем ресурсы контекста сразу после использования

		if err == nil {
			// Проверяем, что соединение действительно работает
			if err = pool.Ping(ctx); err == nil {
				log.Println("Successfully connected to PostgreSQL")
				return pool, nil
			}
		}

		log.Printf("Failed to connect to PostgreSQL (attempt %d/%d): %v",
			attempt, cfg.MaxAttempts, err)

		// Если это не последняя попытка, ждем перед следующей
		if attempt < cfg.MaxAttempts {
			waitTime := time.Duration(attempt) * time.Second // Экспоненциальная задержка
			log.Printf("Waiting %v before next attempt...", waitTime)
			time.Sleep(waitTime)
		}
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts: %w",
		cfg.MaxAttempts, err)
}
