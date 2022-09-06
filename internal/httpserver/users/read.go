package users

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/hrand1005/training-notebook/internal/app"
	"github.com/hrand1005/training-notebook/internal/httpserver/apierror"
)

// ReadByID handles user reads by parsing the UserID as a URL param and
// calling the UserService ReadByID method.
func (u *UserHandler) ReadByID(c *fiber.Ctx) error {
	userID := app.UserID(c.Params("userID"))

	user, err := u.service.ReadByID(userID)
	if err != nil {
		u.logger.Printf("ReadByID: %v", err)
		if errors.Is(err, app.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(&ErrorsResponseBody{
				Errors: []apierror.FormattedError{
					{
						Message: "user not found",
					},
				},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(&ErrorsResponseBody{
			Errors: []apierror.FormattedError{
				{
					Message: app.ErrServiceFailure.Error(),
				},
			},
		})
	}

	resp := buildResponse(userID, user)
	return c.Status(fiber.StatusOK).JSON(resp)
}
