package users_transport_http

import (
	"FreeLib/internal/core/domain"
	core_http_server "FreeLib/internal/core/transport/http/server"
	"context"
	"net/http"
)

type UsersHTTPHandler struct {
	userServices UserServices
}

type UserServices interface {
	CreateUser(ctx context.Context, email string, password string) (domain.User, error)
	AuthorizationUser(ctx context.Context, email string, password string) (string, error)
}

func NewUserHTTPHandler(userServices UserServices) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		userServices: userServices,
	}
}

func (h *UsersHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/register",
			Handler: h.CreateUser,
		},
		{
			Method:  http.MethodPost,
			Path:    "/login",
			Handler: h.AuthorizationUser,
		},
	}
}
