package user

import (
	"github.com/gofiber/fiber/v2"
)

// Create handles user creation by parsing the request into a User to
// call the UserService Create method.
// swagger:route POST /users createUser
// Creates a user.
// responses:
// 201: usersResponse
func (u *UserHandler) Create(c *fiber.Ctx) error {
	var req userRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid json",
			"error":   err,
		})
	}
	u.logger.Printf("Parsed request:\n%#v", req)

	// create model from user request
	user := UserFromRequest(&req)
	userID, err := u.service.Create(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to create user",
			"error":   err, // TODO: make this client appropriate
		})
	}

	// create response from model
	resp := BuildUserResponse(userID, user)
	return c.Status(fiber.StatusCreated).JSON(resp)
}
