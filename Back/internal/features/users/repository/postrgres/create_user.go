package users_postgres_repository

import (
	"FreeLib/internal/core/domain"
	"context"
	"fmt"
)

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
