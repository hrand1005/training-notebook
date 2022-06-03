package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/models"
)

// swagger:route POST /users users createUser
// Creates a user.
// responses:
//  201: userResponse
//  400: errorResponse
//  500: errorResponse

// Create is the handler for create requests on the users resource.
func (u *user) Create(c *gin.Context) {
	var newUser models.User

	if err := c.BindJSON(&newUser); err != nil {
		msg := models.BindingErrorToMessage(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": msg})
		return
	}

	// assigns ID to newUser
	id, err := u.db.AddUser(&newUser)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	newUser.ID = id
	c.IndentedJSON(http.StatusCreated, newUser)
}
