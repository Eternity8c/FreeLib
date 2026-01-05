package postgres

import (
	"FreeLib/internal/models"
	"FreeLib/internal/repository"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) repository.UserRepository {
	return &userRepository{pool: pool}
}

func (r *userRepository) CreateUser(user *models.User) error {
	query := `INSERT INTO users(
	username,
	email,
	password_hash)
	VALUES($1, $2, $3)
	RETURNING id`

	err := r.pool.QueryRow(context.Background(),query, 
	user.Username,
	user.Email,
	user.PasswordHash).Scan(&user.ID)

	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) AuntificationUser(lr *models.LoginRequest) (*models.User, error) {
	query := `SELECT id, username, email, is_admin, password_hash FROM users
	WHERE email = $1`
	var user models.User
	err := r.pool.QueryRow(context.Background(), query, lr.Email).Scan(
		&user.ID, &user.Username, &user.Email,
		&user.IsAdmin, &user.PasswordHash)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(lr.Password)); err != nil {
		return nil, err
	}

	return &user, nil
}