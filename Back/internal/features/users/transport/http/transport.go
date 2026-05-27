package users_transport_http

import (
	core_http_server "FreeLib/internal/core/transport/http/server"
	"net/http"
)

type UsersHTTPHandler struct {
	userServices UserServices
}

type UserServices interface {
	// CreateUser(user *models.User) error
	// AuntificationUser(lr *models.LoginRequest) (*models.User, error)
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
			Path:    "/users",
			Handler: h.CreateUser,
		},
	}
}
