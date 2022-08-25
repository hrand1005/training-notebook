package handlers

import (
	"log"

	// "github.com/gofiber/fiber/v2"
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

/*

// Create handles user creation by parsing the request into a User to
// call the UserService Create method.
func (u *UserHandler) Create(c fiber.Ctx) error {
  var req UserRequest
  if err := c.BodyParser(&req); err != nil {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
      "message": "invalid json",
      "error": err,
    })
  }

  // create model from user request
  user := UserFromRequest(req)
  userID, err := u.service.Create(user)
  if err != nil {
    return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
      "message": "failed to create user",
      "error": err, // TODO: make this client appropriate
    })
  }

  // create response from model
  resp := BuildUserResponse(userID, user)
  return c.Status(fiber.StatusCreated).JSON(resp)
}

*/
