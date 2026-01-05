package handlers

import (
	"FreeLib/internal/models"
	"FreeLib/internal/repository"
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	repo repository.UserRepository
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{
		repo: repo,
	}
}

func (h *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("RegisterHandler called from=%s", r.RemoteAddr)

	var req models.RegisteRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		log.Printf("RegisterHandler: invalid JSON from=%s err=%v", r.RemoteAddr, err)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Password hashing error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	user := models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(passwordHash),
	}

	err = h.repo.CreateUser(&user)
	if err != nil {
		log.Println("Create user error:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("User registered id=%d username=%s email=%s", user.ID, user.Username, user.Email)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"message":  "User registred successfully",
	})
}

func (h *UserHandler) AuntificationHandler(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("Invalid JSON: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.repo.AuntificationUser(&req)
	if err != nil {
		log.Println("Create user error:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"isAdmin":  user.IsAdmin,
		},
	})

	log.Printf("User authenticated id=%d email=%s", user.ID, user.Email)
}
