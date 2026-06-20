package users_transport_http

import (
	"context"
	"net/http"

	"github.com/Eternity8c/FreeLib/internal/core/domain"
	core_logger "github.com/Eternity8c/FreeLib/internal/core/logger"
	core_http_request "github.com/Eternity8c/FreeLib/internal/core/transport/http/request"
	core_http_responce "github.com/Eternity8c/FreeLib/internal/core/transport/http/responce"
	core_http_server "github.com/Eternity8c/FreeLib/internal/core/transport/http/server"
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

// CreateUser    godoc
// @Summary      Создать пользователя
// @Description  Создать нового пользователя в системе
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request    body CreateUserRequest true "CreateUser тело запроса"
// @Success      201        {object} CreateUserResponce "Успешно созданный пользователь"
// @Failure      400       {object} core_http_responce.ErrorResponce "BadRequest"
// @Failure      500       {object} core_http_responce.ErrorResponce "Internal server error"
// @Router      /register  [post]
func (h *UsersHTTPHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responceHandler := core_http_responce.NewHTTPResponceHandler(log, rw)

	log.Debug("invoke CreateUser handler")

	var request CreateUserRequest
	if err := core_http_request.DecodeAndValidateJSONRequest(r, &request); err != nil {
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

// AuthorizationUser	godoc
// @Summary		Авторизовать пользователя
// @Description	Авторизовать пользователя в системе и получить JWT токен
// @Tags		users
// @Accept		json
// @Produce		json
// @Param		request	body AuthorizationUserRequest true "AuthorizationUser тело запроса"
// @Success		200	{object} AuthorizationUserResponce "Успешная авторизация"
// @Failure		400	{object} core_http_responce.ErrorResponce "BadRequest"
// @Failure		401	{object} core_http_responce.ErrorResponce "Unauthorized"
// @Failure		500	{object} core_http_responce.ErrorResponce "Internal server error"
// @Router		/login	[post]
func (h *UsersHTTPHandler) AuthorizationUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responceHandler := core_http_responce.NewHTTPResponceHandler(log, rw)

	log.Debug("invoke AuthorizationUser handler")

	var request AuthorizationUserRequest
	if err := core_http_request.DecodeAndValidateJSONRequest(r, &request); err != nil {
		responceHandler.ErrorResponce(err, "failde decode and validate HTTP request")
		return
	}

	jwtToken, err := h.userServices.AuthorizationUser(ctx, request.Email, request.Password)
	if err != nil {
		responceHandler.ErrorResponce(err, "failed authorization user")
		return
	}

	responce := AuthorizationUserResponce{JWTToken: jwtToken}
	responceHandler.JSONResponce(responce, http.StatusOK)
}
