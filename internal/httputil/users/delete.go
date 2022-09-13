package users

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/hrand1005/training-notebook/internal/app"
)

// DeleteByID handles user deletes by parsing the userID from the url and
// calling the user service DeleteByID method
func (u *UserHandler) DeleteByID(c *fiber.Ctx) error {
	userID := app.UserID(c.Params("userID"))

	if err := u.service.DeleteByID(userID); err != nil {
		if errors.Is(err, app.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(&ErrorsResponseBody{
				Errors: []FormattedError{
					{
						Message: "user not found",
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

	return c.Status(fiber.StatusNoContent).Send(nil)
}
