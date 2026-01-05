package models

type User struct {
	ID           uint   `json:"id" db:"id"`
	Username     string `json:"username" db:"username"`
	Email        string `json:"email" db:"email"`
	IsAdmin      bool   `json:"isAdmin" db:"is_admin"`
	PasswordHash string `json:"-" db:"password_hash"`
}

type RegisteRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
