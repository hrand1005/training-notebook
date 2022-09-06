package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hrand1005/training-notebook/internal/httpserver/apierror"
)

// Create handles user creation by parsing the request into a User to
// call the UserService Create method.
func (u *UserHandler) Create(c *fiber.Ctx) error {
	var req RequestBody
	if err := c.BodyParser(&req); err != nil {
		// u.logger.Printf("Recieved Body:\n%v", string(c.Body()))
		return c.Status(fiber.StatusBadRequest).JSON(&ErrorsResponseBody{
			Errors: []apierror.FormattedError{
				{
					Message: "invalid json: " + err.Error(),
				},
			},
		})
	}

	validationErrors := apierror.ValidateRequestBody(&req, errorFormatter)
	if validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&ErrorsResponseBody{
			Errors: validationErrors,
		})
	}

	user := userFromRequest(&req)
	userID, err := u.service.Create(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&ErrorsResponseBody{
			Errors: []apierror.FormattedError{
				{
					Message: "failed to create user",
				},
			},
		})
	}

	resp := buildResponse(userID, user)
	return c.Status(fiber.StatusCreated).JSON(resp)
}
