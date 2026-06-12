package main

import (
	core_logger "FreeLib/internal/core/logger"
	core_postgres_pool "FreeLib/internal/core/repository/postgres/pool"
	core_http_middleware "FreeLib/internal/core/transport/http/middleware"
	core_http_server "FreeLib/internal/core/transport/http/server"
	book_postgres_repository "FreeLib/internal/features/books/repository/postgres"
	book_service "FreeLib/internal/features/books/service"
	books_transport_http "FreeLib/internal/features/books/transport/http"
	users_postgres_repository "FreeLib/internal/features/users/repository/postrgres"
	users_service "FreeLib/internal/features/users/service"
	users_transport_http "FreeLib/internal/features/users/transport/http"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

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

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
	// //Books Handlers
	// r.HandleFunc("/api/health", bookHandler.HealthHandler).Methods("GET")
	// r.HandleFunc("/api/books", bookHandler.GetBooksHandler).Methods("GET")
	// r.HandleFunc("/api/book", bookHandler.GetByIDHandler).Methods("GET")
	// r.HandleFunc("/api/create", bookHandler.CreateHandler).Methods("POST")
	// r.HandleFunc("/api/book", bookHandler.DeleteHandler).Methods("DELETE")
	// r.HandleFunc("/api/book/{id}", bookHandler.UpdateBookHandler).Methods("PATCH")
	// r.HandleFunc("/api/users/{id}/favorites", bookHandler.AddFavoriteHandler).Methods("POST")
	// r.HandleFunc("/api/users/{user_id}/favorites/book/{book_id}", bookHandler.DeleteFavoriteHandler).Methods("DELETE")
	// r.HandleFunc("/api/users/{user_id}/favorites", bookHandler.GetAllFavotiteHandler).Methods("GET")

	// //User Handlers
	// r.HandleFunc("/api/register", userHandler.RegisterHandler).Methods("POST")
	// r.HandleFunc("/api/login", userHandler.AuntificationHandler).Methods("POST")
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
