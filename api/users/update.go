package users

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/data"
	"github.com/hrand1005/training-notebook/models"
)

// <SWAGGER-IGNORE>:route PUT /users/{id} users updateUser
// Update a user.
// responses:
//  200: userResponse
//  400: errorResponse
//  404: errorResponse
//  500: errorResponse

// Update is the handler for update requests on the user resource. An id must be
// specified.
func (u *user) Update(c *gin.Context) {
	var newUser models.User

	if err := c.BindJSON(&newUser); err != nil {
		msg := models.BindingErrorToMessage(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": msg})
		return
	}

	userID, err := UserIDFromParams(c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": ErrInvalidUserID})
		return
	}

	if err := u.db.UpdateUser(userID, &newUser); err != nil {
		if err == data.ErrNotFound {
			msg := fmt.Sprintf("no such user with id %v", userID)
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": msg})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	newUser.ID = userID
	c.IndentedJSON(http.StatusOK, newUser)
}
