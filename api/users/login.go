package users

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/hrand1005/training-notebook/models"
	"golang.org/x/crypto/bcrypt"
)

// token settings
const (
	LoginCookieName     = "token"
	LoginCookieMaxAge   = 3600
	LoginCookieSecure   = false
	LoginCookieHTTPOnly = true
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

	// TODO: get path and domain from configs or env
	c.SetCookie(LoginCookieName, token, LoginCookieMaxAge, "", "", LoginCookieSecure, LoginCookieHTTPOnly)
	c.IndentedJSON(http.StatusOK, gin.H{"message": "logged in successfully"})
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

type Claims struct {
	UserID models.UserID
	jwt.StandardClaims
}

func buildToken(user *models.User) (string, error) {
	// TODO: get this from configs
	signingKey := []byte(os.Getenv("SIGNING_KEY"))
	jwtClaims := &Claims{
		UserID: user.ID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	return token.SignedString(signingKey)
}
