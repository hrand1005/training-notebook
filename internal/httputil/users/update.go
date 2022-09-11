package users

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/hrand1005/training-notebook/internal/app"
)

// UpdateByID handles user updates by parsing url param into a user id,
// and the request into a User, and updating the existing user with the
// provided fields.
func (u *UserHandler) UpdateByID(c *fiber.Ctx) error {
	userID := app.UserID(c.Params("userID"))

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

	if err := u.service.UpdateByID(userID, user); err != nil {
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

	return c.Status(fiber.StatusOK).JSON(buildResponse(userID, user))
}
