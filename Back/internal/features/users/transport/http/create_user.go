package users_transport_http

import (
	"FreeLib/internal/core/domain"
	core_logger "FreeLib/internal/core/logger"
	core_http_request "FreeLib/internal/core/transport/http/request"
	core_http_responce "FreeLib/internal/core/transport/http/responce"
	"net/http"
)

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,contains=@"`
	Password string `json:"password" validate:"required"`
}

type CreateUserResponce struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func (h *UsersHTTPHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responceHandler := core_http_responce.NewHTTPResponceHandler(log, rw)

	log.Debug("invoke CreateUser handler")

	var request CreateUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responceHandler.ErrorResponce(err, "failed to decode and validate HTTP request")
		return
	}

	userDomain, err := h.userServices.CreateUser(ctx, request.Email, request.Password)
	if err != nil {
		responceHandler.ErrorResponce(err, "failed to create user")
		return
	}

	responce := dtoFromDomain(userDomain)
	responceHandler.JSONResponce(responce, http.StatusCreated)
}

func dtoFromDomain(user domain.User) CreateUserResponce {
	return CreateUserResponce{
		ID:    user.ID,
		Email: user.Email,
	}
}
