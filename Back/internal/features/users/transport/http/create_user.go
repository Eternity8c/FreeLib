package users_transport_http

import "net/http"

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

}
