package repository_test

import (
	"FreeLib/pkg/config"
	"FreeLib/pkg/database"
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
)

func StartPostgreSQL() (cfng config.Config, cleaner func()) {
pool, err := dockertest.NewPool("")
    if err != nil {
        log.Fatalf("Could not connect to docker: %s", err)
    }

    pool.MaxWait = 60 * time.Second

    resource, err := pool.Run("postgres", "13", []string{
        "POSTGRES_DB=testdb",
        "POSTGRES_USER=postgres", 
        "POSTGRES_PASSWORD=test",
    })
    if err != nil {
        log.Fatalf("Could not start resource: %s", err)
    }

    port := resource.GetPort("5432/tcp")
    connString := fmt.Sprintf("postgres://postgres:test@localhost:%s/testdb?sslmode=disable", port)

	cnfg := config.Config{
        Host:        "localhost",
        Port:        port,
        Username:    "postgres",
        Password:    "test",
        DBname:      "testdb", 
        SSLMode:     "disable",
        MaxAttempts: 5,
    }

    // Используем встроенный Retry вместо ручного цикла
    if err := pool.Retry(func() error {
        conn, err := pgx.Connect(context.Background(), connString)
        if err != nil {
            return err
        }
        defer conn.Close(context.Background())
        
        // Дополнительная проверка - выполняем простой запрос
        _, err = conn.Exec(context.Background(), "SELECT 1")
        return err
    }); err != nil {
        resource.Close()
        log.Fatalf("Could not connect to database: %s", err)
    }

	cleanerFunc := func() {
		err := pool.Purge(resource)
		if err != nil {
			log.Fatalf("pool.Purge failed: %v", err)
		}
	}

	return cnfg, cleanerFunc
}

func TestConnectDB(t *testing.T) {
	cnfg, cleanup := StartPostgreSQL()
	defer cleanup()
	ctx := context.Background()
	_, err := database.ConnectDB(ctx, cnfg)

	assert.Nil(t, err)
}
