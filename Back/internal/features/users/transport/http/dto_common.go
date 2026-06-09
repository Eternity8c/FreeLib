package users_transport_http

import "FreeLib/internal/core/domain"

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,contains=@"`
	Password string `json:"password" validate:"required"`
}

type CreateUserResponce struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type AuthorizationUserRequest struct {
	Email    string `json:"email" validate:"required,contains=@"`
	Password string `json:"password" validate:"required"`
}

type AuthorizationUserResponce struct {
	JWTToken string `json:"jwt_token"`
}

func dtoFromDomain(user domain.User) CreateUserResponce {
	return CreateUserResponce{
		ID:    user.ID,
		Email: user.Email,
	}
}
