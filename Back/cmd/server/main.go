package main

import (
	"FreeLib/internal/handlers"
	"FreeLib/internal/repository/postgres"
	"FreeLib/pkg/config"
	"FreeLib/pkg/database"
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	cfg := config.LoadConfig()
	ctx := context.Background()
	pool, err := database.ConnectDB(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	log.Println("Connect db")

	bookRepo := postgres.NewBookRepository(pool)

	bookHandler := handlers.NewBookHandler(bookRepo)

	userRepo := postgres.NewUserRepository(pool)

	userHandler := handlers.NewUserHandler(userRepo)

	r := mux.NewRouter()

	//Books Handlers
	r.HandleFunc("/api/health", bookHandler.HealthHandler).Methods("GET")
	r.HandleFunc("/api/books", bookHandler.GetBooksHandler).Methods("GET")
	r.HandleFunc("/api/book", bookHandler.GetByIDHandler).Methods("GET")
	r.HandleFunc("/api/create", bookHandler.CreateHandler).Methods("POST")
	r.HandleFunc("/api/book", bookHandler.DeleteHandler).Methods("DELETE")
	r.HandleFunc("/api/book/{id}", bookHandler.UpdateBookHandler).Methods("PATCH")
	r.HandleFunc("/api/users/{id}/favorites", bookHandler.AddFavoriteHandler).Methods("POST")
	r.HandleFunc("/api/users/{user_id}/favorites/book/{book_id}", bookHandler.DeleteFavoriteHandler).Methods("DELETE")
	r.HandleFunc("/api/users/{user_id}/favorites", bookHandler.GetAllFavotiteHandler).Methods("GET")

	//User Handlers
	r.HandleFunc("/api/register", userHandler.RegisterHandler).Methods("POST")
	r.HandleFunc("/api/login", userHandler.AuntificationHandler).Methods("POST")

	handler := withCORS(r)

	// Получаем порт из переменной окружения
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("FreeLib server starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))

}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получаем CORS настройки из переменных окружения
		allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
		if allowedOrigins == "" {
			allowedOrigins = "*"
		}
		allowedMethods := os.Getenv("CORS_ALLOWED_METHODS")
		if allowedMethods == "" {
			allowedMethods = "GET, POST, PUT, PATCH, DELETE, OPTIONS"
		}
		allowedHeaders := os.Getenv("CORS_ALLOWED_HEADERS")
		if allowedHeaders == "" {
			allowedHeaders = "Content-Type, Authorization"
		}

		w.Header().Set("Access-Control-Allow-Origin", allowedOrigins)
		w.Header().Set("Access-Control-Allow-Methods", allowedMethods)
		w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
