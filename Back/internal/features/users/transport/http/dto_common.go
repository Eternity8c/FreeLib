package users_transport_http

import "github.com/Eternity8c/FreeLib/internal/core/domain"

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,contains=@" example:"user@gmail.com"`
	Password string `json:"password" validate:"required" example:"password1234"`
}

type CreateUserResponce struct {
	ID    int    `json:"id" example:"1"`
	Email string `json:"email" example:"user@gmail.com"`
}

type AuthorizationUserRequest struct {
	Email    string `json:"email" validate:"required,contains=@" example:"user@gmail.com"`
	Password string `json:"password" validate:"required" example:"password1234"`
}

type AuthorizationUserResponce struct {
	JWTToken string `json:"jwt_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

func dtoFromDomain(user domain.User) CreateUserResponce {
	return CreateUserResponce{
		ID:    user.ID,
		Email: user.Email,
	}
}
