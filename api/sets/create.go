package sets

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/api/users"
	"github.com/hrand1005/training-notebook/models"
)

// swagger:route POST /sets sets createSet
// Creates a set.
// responses:
//  201: setResponse
//  400: errorResponse
//  500: errorResponse

// Create is the handler for create requests on the set resource.
func (s *set) Create(c *gin.Context) {
	userID, err := users.UserIDFromContext(c)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "you must be logged in with valid user to perform this action"})
	}

	var newSet models.Set

	if err := c.BindJSON(&newSet); err != nil {
		msg := models.BindingErrorToMessage(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": msg})
		return
	}

	// created set should always have the id of the logged in user
	newSet.UID = userID

	// assigns ID to newSet upon entry
	id, err := s.db.AddSet(&newSet)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	newSet.ID = id
	c.IndentedJSON(http.StatusCreated, newSet)
}
