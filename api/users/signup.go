package users

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/models"
)

func (u *user) Signup(c *gin.Context) {
	var newUser models.User

	if err := c.BindJSON(&newUser); err != nil {
		log.Printf("debug: failed to bind in signup: %s", err.Error())
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := checkPasswordRequirements(newUser.Password); err != nil {
		log.Printf("debug: password requirements failed in signup: %s", err.Error())
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// hash the password so that user information is protected
	hashedPassword := hashPassword(newUser.Password)
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

func hashPassword(raw string) string {
	return raw + "hi"
}
