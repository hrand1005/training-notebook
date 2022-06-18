package sets

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/api/users"
	"github.com/hrand1005/training-notebook/data"
	"github.com/hrand1005/training-notebook/models"
)

// swagger:route PUT /sets/{id} sets updateSet
// Update a set.
// responses:
//  200: setResponse
//  400: errorResponse
//  404: errorResponse
//  500: errorResponse

// Update is the handler for update requests on the set resource. An id must be
// specified.
func (s *set) Update(c *gin.Context) {
	userID, err := users.UserIDFromContext(c)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "you must be logged in to perform this action"})
		return
	}

	var newSet models.Set

	if err := c.BindJSON(&newSet); err != nil {
		msg := models.BindingErrorToMessage(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": msg})
		return
	}

	setID, err := SetIDFromParams(c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": ErrInvalidSetID})
		return
	}

	if err := s.db.UpdateSetForUser(setID, userID, &newSet); err != nil {
		if err == data.ErrNotFound {
			msg := fmt.Sprintf("no such set with id %v", setID)
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": msg})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	newSet.ID = setID
	newSet.UID = userID
	c.IndentedJSON(http.StatusOK, newSet)
}
