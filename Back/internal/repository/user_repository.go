package repository

import "FreeLib/internal/models"

type UserRepository interface {
	CreateUser(user *models.User) error
	AuntificationUser(lr *models.LoginRequest) (*models.User, error)
}
