package users_service

import (
	"FreeLib/internal/core/domain"
	core_jwt "FreeLib/internal/core/jwt"
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	ID       int
	Email    string
	PassHash string
	IsAdmin  bool
}

func (s *UsersService) AuthorizationUser(ctx context.Context, email string, password string) (string, error) {
	if err := validate(email, password); err != nil {
		return "", fmt.Errorf("validate user domain: %w", err)
	}

	modalUser, err := s.usersRepository.GetUser(ctx, email)
	if err != nil {
		return "", fmt.Errorf("get user: %w", err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(modalUser.PassHash), []byte(password)); err != nil {
		return "", fmt.Errorf("compare hash and password: %w", err)
	}

	domainUser := domain.NewUser(modalUser.ID, modalUser.Email, modalUser.IsAdmin)

	jwtToken, err := core_jwt.GenerateToken(domainUser)
	if err != nil {
		return "", fmt.Errorf("generate jwt token: %w", err)
	}

	return jwtToken, nil
}
