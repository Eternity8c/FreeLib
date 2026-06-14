package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/Eternity8c/FreeLib/internal/core/logger"
	core_postgres_pool "github.com/Eternity8c/FreeLib/internal/core/repository/postgres/pool"
	core_http_middleware "github.com/Eternity8c/FreeLib/internal/core/transport/http/middleware"
	core_http_server "github.com/Eternity8c/FreeLib/internal/core/transport/http/server"
	book_postgres_repository "github.com/Eternity8c/FreeLib/internal/features/books/repository/postgres"
	book_service "github.com/Eternity8c/FreeLib/internal/features/books/service"
	books_transport_http "github.com/Eternity8c/FreeLib/internal/features/books/transport/http"
	users_postgres_repository "github.com/Eternity8c/FreeLib/internal/features/users/repository/postrgres"
	users_service "github.com/Eternity8c/FreeLib/internal/features/users/service"
	users_transport_http "github.com/Eternity8c/FreeLib/internal/features/users/transport/http"

	"go.uber.org/zap"

	_ "github.com/Eternity8c/FreeLib/docs"
)

// @title       Goland FreeLib API
// @version     1.0
// @description FreeLib Aplication REST-API schema
// @host        localhost:8080
// @BasePath    /
func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("initializing postgres connection pool")
	pool, err := core_postgres_pool.NewConnectionPool(
		ctx,
		core_postgres_pool.NewConfigMust(),
	)

	if err != nil {
		logger.Fatal(
			"failed to init postges connection pool",
			zap.Error(err),
		)
	}
	defer pool.Close()

	logger.Debug("initializing feature", zap.String("feature", "users"))

	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	userService := users_service.NewUsersService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUserHTTPHandler(userService)

	logger.Debug("initializing feature", zap.String("feature", "books"))

	booksReposirory := book_postgres_repository.NewBookRepository(pool)
	booksService := book_service.NewBookService(booksReposirory)
	booksTransportHTTP := books_transport_http.NewBookHTTPHandler(booksService)

	logger.Debug("initializing HTTP server")

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	router := core_http_server.NewRouter()
	router.RegisterRoutes(usersTransportHTTP.Routes()...)
	router.RegisterRoutes(booksTransportHTTP.Routes()...)
	httpServer.RegisterAPIRoutes(router)

	httpServer.RegisterSwagger()

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
