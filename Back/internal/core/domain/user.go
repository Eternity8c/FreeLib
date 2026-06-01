package domain

type User struct {
	ID      int    `json:"id" db:"id"`
	Email   string `json:"email" db:"email"`
	IsAdmin bool   `json:"isAdmin" db:"is_admin"`
}

func NewUserUninitialized(email string) User {
	return User{
		ID:      UninitializedID,
		Email:   email,
		IsAdmin: false,
	}
}

func NewUser(id int, email string, isAdmin bool) User {
	return User{
		ID:      id,
		Email:   email,
		IsAdmin: isAdmin,
	}
}
