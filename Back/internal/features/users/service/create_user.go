package users_service

import (
	"FreeLib/internal/core/domain"
	core_errors "FreeLib/internal/core/errors"
	"context"
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

func (s *UsersService) CreateUser(ctx context.Context, email string, password string) (domain.User, error) {

	if err := validate(email, password); err != nil {
		return domain.User{}, fmt.Errorf("validate user domain: %w", err)
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, fmt.Errorf("generate password hash: %w", err)
	}

	user := domain.NewUserUninitialized(email)

	user, err = s.usersRepository.CreateUser(ctx, user, string(passHash))
	if err != nil {
		return domain.User{}, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

func validate(email string, password string) error {
	emailLenght := len([]rune(email))
	if emailLenght < 4 || emailLenght > 100 {
		return fmt.Errorf(
			"invalid `Email` len: %d: %w",
			emailLenght,
			core_errors.ErrInvalidArgumment,
		)
	}

	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	if !re.MatchString(email) {
		return fmt.Errorf(
			"invalid `Email` format: %w",
			core_errors.ErrInvalidArgumment,
		)
	}

	passwordLenght := len([]rune(password))
	if passwordLenght < 7 {
		return fmt.Errorf(
			"invalid `Password` len: %d: %w",
			passwordLenght,
			core_errors.ErrInvalidArgumment,
		)
	}

	return nil
}
