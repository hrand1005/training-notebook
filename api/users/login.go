package users

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/models"
)

func (u *user) Login(c *gin.Context) {
	var credentials models.Credentails

	if err := c.BindJSON(&credentials); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, err := u.db.UserByID(credentials.UID)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	token, err := AuthenticateUser(user, credentials)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, token)
}

func AuthenticateUser(user *models.User, credentials models.Credentails) (string, error) {
	if hashPassword(credentials.Password) == user.Password {
		return buildToken(user), nil
	}

	return "", fmt.Errorf("incorrect password")
}

func buildToken(user *models.User) string {
	return user.Name + user.Password
}
