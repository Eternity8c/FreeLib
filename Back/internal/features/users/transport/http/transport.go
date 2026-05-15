package users_transport_http

type UsersHTTPHandler struct {
	userServices UserServices
}

type UserServices interface {
	CreateUser(user *models.User) error
	AuntificationUser(lr *models.LoginRequest) (*models.User, error)
}

func NewUserHTTPHandler(userServices UserServices) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		userServices: userServices,
	}
}
