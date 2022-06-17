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
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "must be logged in with valid user to perform this function"})
	}

	var newSet models.Set

	if err := c.BindJSON(&newSet); err != nil {
		msg := models.BindingErrorToMessage(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": msg})
		return
	}

	// assigns ID to newSet
	id, err := s.db.AddSet(&newSet)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	newSet.ID = id
	newSet.UID = userID
	c.IndentedJSON(http.StatusCreated, newSet)
}
