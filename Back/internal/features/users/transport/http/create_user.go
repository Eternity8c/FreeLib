package users_transport_http

import (
	core_logger "FreeLib/internal/core/logger"
	"encoding/json"
	"fmt"
	"net/http"
)

type CreateUserRequest struct {
	UserName string `json:"user_name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type CreateUserResponce struct {
	ID       int    `json:"id,omitempty"`
	UserName string `json:"user_name,omitempty"`
	Email    string `json:"email,omitempty"`
}

func (h *UsersHTTPHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	log.Debug("invoke CreateUser handler")

	var request CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		fmt.Println("Error")
	}

	rw.WriteHeader(http.StatusOK)
}
