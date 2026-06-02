package users_postgres_repository

import (
	"context"
	"fmt"
)

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
