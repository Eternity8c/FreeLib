package users_service

import (
	"FreeLib/internal/core/domain"
	users_postgres_repository "FreeLib/internal/features/users/repository/postrgres"
	"context"
)

type UsersService struct {
	usersRepository UsersRepository
}

type UsersRepository interface {
	CreateUser(ctx context.Context, user domain.User, passHash string) (domain.User, error)
	GetUser(ctx context.Context, email string) (users_postgres_repository.UserModel, error)
}

func NewUsersService(usersRepository UsersRepository) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
	}
}
