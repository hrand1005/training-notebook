package user

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/hrand1005/training-notebook/internal/app"
)

// UserHandler contains the service and logger required to
// handle http requests.
type UserHandler struct {
	service app.UserService
	logger  *log.Logger
}

// NewUserHandler is a constructor for UserHandlers.
func NewUserHandler(s app.UserService, l *log.Logger) *UserHandler {
	return &UserHandler{
		service: s,
		logger:  l,
	}
}

// Register adds the users endpoints to the provided router.
func (u *UserHandler) Register(r fiber.Router) {
	r.Post("/users", u.Create)
}

// UserFromRequest creates a user entity from an api request.
func UserFromRequest(req *userRequest) *app.User {
	return &app.User{
		FirstName: req.Data.Attributes.FirstName,
		LastName:  req.Data.Attributes.LastName,
		Email:     req.Data.Attributes.Email,
	}
}

// BuildUserResponse builds an api response from a user entity.
func BuildUserResponse(userID app.UserID, user *app.User) *userResponse {
	return &userResponse{
		Data: userObject{
			ID: string(userID),
			Attributes: userAttributes{
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Email:     user.Email,
			},
		},
	}
}
