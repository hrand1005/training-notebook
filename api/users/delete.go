package users

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/data"
)

// NOTE: The following method is used for debugging, and should not be reachable
// 	to normal users during production

// <SWAGGER-IGNORE>:route DELETE /users/{id} users deleteuser
// Delete a user.
// responses:
//  204: noContent
//  404: errorResponse
//  500: errorResponse

// Delete is the handler for delete requests on the user resource. An id must be
// specified.
func (u *user) Delete(c *gin.Context) {
	userID, err := UserIDFromParams(c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": ErrInvalidUserID})
		return
	}

	if err := u.db.DeleteUser(userID); err != nil {
		if err == data.ErrNotFound {
			msg := fmt.Sprintf("no such user with id %v", userID)
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": msg})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusNoContent, gin.H{})
}
