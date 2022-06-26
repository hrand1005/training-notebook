package users

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/data"
	"github.com/hrand1005/training-notebook/models"
)

// NOTE: The following method is used for debugging, and should not be reachable
// 	to normal users during production

// <SWAGGER-IGNORE>:route GET /users users readAllUsers
// Read all users.
// responses:
//  200: usersResponse
//  500: errorResponse

// ReadAll is the handler for read requests on the users resource where no id is
// specified. Returns all users on this resource's data source.
func (u *user) ReadAll(c *gin.Context) {
	users, err := u.db.Users()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if len(users) == 0 {
		// if no sets are found, return an empty slice
		c.IndentedJSON(http.StatusOK, []*models.User{})
		return
	}

	// TODO: clean r of personal data if this is going to be a public api endpoint
	c.IndentedJSON(http.StatusOK, users)
}

// swagger:route GET /users/{id} users readUser
// Read a user.
// responses:
//  200: userResponse
//  400: errorResponse
//  404: errorResponse
//  500: errorResponse

// Read is the handler for read requests on the user resource where an id is
// specified.
func (s *user) Read(c *gin.Context) {
	userID, err := UserIDFromParams(c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": ErrInvalidUserID})
		return
	}

	r, err := s.db.UserByID(userID)
	if err != nil {
		if err == data.ErrNotFound {
			msg := fmt.Sprintf("no such user with id %v", userID)
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": msg})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// TODO: clean r of personal data if this is going to be a public api endpoint
	c.IndentedJSON(http.StatusOK, r)
}
