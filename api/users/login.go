package users

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/hrand1005/training-notebook/models"
	"golang.org/x/crypto/bcrypt"
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
	if checkPasswordHash(user.Password, credentials.Password) {
		return buildToken(user)
	}

	return "", fmt.Errorf("incorrect password")
}

func checkPasswordHash(hash, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}
	return true
}

type claims struct {
	UserID models.UserID
	jwt.StandardClaims
}

func buildToken(user *models.User) (string, error) {
	// TODO: get the signing key from super secret env var
	signingKey := []byte("oi bruv here's ur key")
	jwtClaims := &claims{
		UserID: user.ID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	tokenString, err := token.SignedString(signingKey)
	log.Printf("Error: %v", err)
	return tokenString, err
}
