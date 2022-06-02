package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/data"
)

// swagger:route GET /users users readAllUsers
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
		c.IndentedJSON(http.StatusOK, []*data.User{})
		return
	}
	c.IndentedJSON(http.StatusOK, users)

}
