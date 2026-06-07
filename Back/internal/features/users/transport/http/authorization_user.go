package users_transport_http

import (
	core_logger "FreeLib/internal/core/logger"
	core_http_request "FreeLib/internal/core/transport/http/request"
	core_http_responce "FreeLib/internal/core/transport/http/responce"
	"net/http"
)

type AuthorizationUserRequest struct {
	Email    string `json:"email" validate:"required,contains=@"`
	Password string `json:"password" validate:"required"`
}

type AuthorizationUserResponce struct {
	JWTToken string `json:"jwt_token"`
}

func (h *UsersHTTPHandler) AuthorizationUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responceHandler := core_http_responce.NewHTTPResponceHandler(log, rw)

	log.Debug("invoke AuthorizationUser handler")

	var request AuthorizationUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responceHandler.ErrorResponce(err, "failde decode and validate HTTP request")
		return
	}

	jwtToken, err := h.userServices.AuthorizationUser(ctx, request.Email, request.Password)
	if err != nil {
		responceHandler.ErrorResponce(err, "failed authorization user")
		return
	}

	responceHandler.JSONResponce(jwtToken, http.StatusOK)
}
