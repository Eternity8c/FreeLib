package users_service

import (
	"context"
	"fmt"
	"regexp"

	"github.com/Eternity8c/FreeLib/internal/core/domain"
	core_errors "github.com/Eternity8c/FreeLib/internal/core/errors"
	core_jwt "github.com/Eternity8c/FreeLib/internal/core/jwt"
	users_postgres_repository "github.com/Eternity8c/FreeLib/internal/features/users/repository/postrgres"

	"golang.org/x/crypto/bcrypt"
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
