package users

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/hrand1005/training-notebook/data"
	"github.com/hrand1005/training-notebook/models"
)

// TODO: authorization/permissions levels
func RequireAuthorization( /*l AuthorizationLevel*/ ) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(LoginCookieName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "must be logged in to perform this function"})
		}

		userID, err := parseUserIDFromToken(token)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("userID", userID)

		// user is authorized, continue chaining HandlerFuncs
		c.Next()
	}
}

func parseUserIDFromToken(token string) (models.UserID, error) {
	var claims Claims
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		// TODO: get this from configs
		return []byte(os.Getenv("SIGNING_KEY")), nil
	})
	if err != nil {
		return data.InvalidUserID, err
	}

	if !parsedToken.Valid {
		return data.InvalidUserID, fmt.Errorf("invalid token")
	}

	// TODO: can this return zero? how can I verify that it's initialized?
	return claims.UserID, nil
}

/*
func authorized(token string, l AuthorizationLevel) bool {
	level := getLevelFromToken(token)
	return meetsPermissions(l, level)
}
*/
