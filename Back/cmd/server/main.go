package main

import (
	core_logger "FreeLib/internal/core/logger"
	core_http_middleware "FreeLib/internal/core/transport/http/middleware"
	core_http_server "FreeLib/internal/core/transport/http/server"
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

	logger.Debug("Starting FreeLib application")

	usersTransportHTTP := users_transport_http.NewUserHTTPHandler(nil)
	usersRoutes := usersTransportHTTP.Routes()

	router := core_http_server.NewRouter()
	router.RegisterRoutes(usersRoutes...)

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Panic(),
		core_http_middleware.Trace(),
	)
	httpServer.RegisterAPIRoutes(router)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
	// cfg := config.NewConfig()
	// ctx := context.Background()
	// pool, err := database.ConnectDB(ctx, cfg.DBAddr)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer pool.Close()
	// log.Println("Connect db")

	// bookRepo := postgres.NewBookRepository(pool)

	// bookHandler := handlers.NewBookHandler(bookRepo)

	// userRepo := postgres.NewUserRepository(pool)

	// userHandler := handlers.NewUserHandler(userRepo)

	// r := mux.NewRouter()

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

	// handler := withCORS(r)

	// log.Println("FreeLib server starting on :8080")
	// log.Fatal(http.ListenAndServe(":8080", handler))

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
