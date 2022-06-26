package users

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/models"
	"golang.org/x/crypto/bcrypt"
)

// swagger:route POST /users/signup users
// Signup a new user.
// responses:
//  201: userResponse
//  400: errorResponse
//  500: errorResponse

// Signup is the handler for post requests on the user resource.
// Creates a new user if fields are valid and returns the newly created user.
func (u *user) Signup(c *gin.Context) {
	var newUser models.User

	if err := c.BindJSON(&newUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := checkPasswordRequirements(newUser.Password); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// hash the password so that user information is protected
	hashedPassword, err := hashPassword(newUser.Password)
	if err != nil {
		// TODO: use configured logger for this server
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid password"})
	}
	newUser.Password = hashedPassword

	// assigns ID to newUser
	id, err := u.db.AddUser(&newUser)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// TODO: create a function that cleans the User model of sensitive data for
	// request building
	resp := &models.User{
		ID:   id,
		Name: newUser.Name,
	}

	c.IndentedJSON(http.StatusCreated, resp)
}

const minPasswordLength = 5

func checkPasswordRequirements(password string) error {
	if len(password) < minPasswordLength {
		return fmt.Errorf("password too short, must be at least %v characters", minPasswordLength)
	}
	return nil
}

func hashPassword(raw string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	return string(bytes), err
}
