package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/Eternity8c/FreeLib/internal/core/domain"
	core_postgres_pool "github.com/Eternity8c/FreeLib/internal/core/repository/postgres/pool"
)

type UsersRepository struct {
	pool core_postgres_pool.Pool
}

func NewUsersRepository(pool core_postgres_pool.Pool) *UsersRepository {
	return &UsersRepository{
		pool: pool,
	}
}

func (r *UsersRepository) CreateUser(
	ctx context.Context,
	user domain.User,
	passHash string,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO freelib.users (email, pass_hash)
	VALUES ($1, $2)
	RETURNING user_id, email, is_admin;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		user.Email,
		passHash,
	)

	var userModel UserModel
	err := row.Scan(
		&userModel.ID,
		&userModel.Email,
		&userModel.IsAdmin,
	)
	if err != nil {
		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	userDomain := domain.NewUser(
		userModel.ID,
		userModel.Email,
		userModel.IsAdmin,
	)

	return userDomain, nil
}

func (r *UsersRepository) GetUser(ctx context.Context, email string) (UserModel, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT user_id, email, pass_hash, is_admin
	FROM freelib.users
	WHERE email = $1;
	`

	row := r.pool.QueryRow(ctx, query, email)

	var userModal UserModel
	if err := row.Scan(&userModal.ID, &userModal.Email, &userModal.PassHash, &userModal.IsAdmin); err != nil {
		return UserModel{}, fmt.Errorf("scan error: %w", err)
	}

	return userModal, nil
}
