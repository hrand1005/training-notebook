package users

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
	// swagger:route POST /users createUser
	//
	// Creates a new user.
	//
	// Creates a new user in the system with the provided attributes.
	// Assigns the user an ID upon posting.
	//
	// Consumes:
	// - application/json
	//
	// Produces:
	// - application/json
	//
	// Schemes: http, https
	//
	// Security:
	// - api_key:
	//
	// responses:
	// 201: usersResponse
	// 400: errorsResponse
	// 500: errorsResponse
	r.Post("/users", u.Create)
	r.Get("/users/:userID", u.ReadByID)
}

// userFromRequest creates a user entity from an api request.
func userFromRequest(req *RequestBody) *app.User {
	return &app.User{
		FirstName: req.Data.Attributes.FirstName,
		LastName:  req.Data.Attributes.LastName,
		Email:     req.Data.Attributes.Email,
	}
}

const RequestTypeUser = "user"

// buildResponse builds an api response from a user entity.
func buildResponse(userID app.UserID, user *app.User) *ResponseBody {
	return &ResponseBody{
		Data: ResponseData{
			ID:   string(userID),
			Type: RequestTypeUser,
			Attributes: Attributes{
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Email:     user.Email,
			},
		},
	}
}
