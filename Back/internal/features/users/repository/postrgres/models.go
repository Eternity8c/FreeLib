package users_postgres_repository

type UserModel struct {
	ID      int
	Email   string
	IsAdmin bool
}
