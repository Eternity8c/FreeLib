package domain

type User struct {
	ID           uint   `json:"id" db:"id"`
	Username     string `json:"username" db:"username"`
	Email        string `json:"email" db:"email"`
	IsAdmin      bool   `json:"isAdmin" db:"is_admin"`
	PasswordHash string `json:"-" db:"password_hash"`
}
