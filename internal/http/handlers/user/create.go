package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hrand1005/training-notebook/internal/http/errors"
)

// Create handles user creation by parsing the request into a User to
// call the UserService Create method.
func (u *UserHandler) Create(c *fiber.Ctx) error {
	var req RequestBody
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid json",
			"error":   err,
		})
	}

	validationErrors := errors.ValidateRequestBody(&req)
	if validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(validationErrors)
	}

	user := userFromRequest(&req)
	userID, err := u.service.Create(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to create user",
			"error":   err, // TODO: make this client appropriate
		})
	}

	resp := buildResponse(userID, user)
	return c.Status(fiber.StatusCreated).JSON(resp)
}
