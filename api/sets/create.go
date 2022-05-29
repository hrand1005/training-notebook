package sets

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/data"
)

// swagger:route POST /sets sets createSet
// Creates a set.
// responses:
//  201: setResponse
//  400: errorResponse

// Create is the handler for create requests on the set resource.
func (s *set) Create(c *gin.Context) {
	var newSet data.Set

	if err := c.BindJSON(&newSet); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// assigns ID to newSet
	id, err := s.db.AddSet(&newSet)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	newSet.ID = id
	c.IndentedJSON(http.StatusCreated, newSet)
}
