package users

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/hrand1005/training-notebook/internal/app"
)

// Create handles user creation by parsing the request into a User to
// call the UserService Create method.
func (u *UserHandler) Create(c *fiber.Ctx) error {
	var req RequestBody
	if err := c.BodyParser(&req); err != nil {
		// u.logger.Printf("Recieved Body:\n%v", string(c.Body()))
		return c.Status(fiber.StatusBadRequest).JSON(&ErrorsResponseBody{
			Errors: []FormattedError{
				{
					Message: "invalid json: " + err.Error(),
				},
			},
		})
	}

	user := userFromRequest(&req)
	userID, err := u.service.Create(user)
	if err != nil {
		if errors.Is(err, app.ErrInvalidField) {
			return c.Status(fiber.StatusBadRequest).JSON(&ErrorsResponseBody{
				Errors: []FormattedError{
					{
						Message: err.Error(),
					},
				},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(&ErrorsResponseBody{
			Errors: []FormattedError{
				{
					Message: app.ErrServiceFailure.Error(),
				},
			},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(buildResponse(userID, user))
}
