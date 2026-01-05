package main

import (
	"FreeLib/internal/handlers"
	"FreeLib/internal/repository/postgres"
	"FreeLib/pkg/config"
	"FreeLib/pkg/database"
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
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

	log.Println("FreeLib server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))

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
