package sets

import (
	"log"
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
	log.Println("In Create")
	var newSet data.Set

	if err := c.BindJSON(&newSet); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// assigns ID to newSet
	s.db.AddSet(&newSet)
	// TODO: return newly created set ID /object representation
	c.IndentedJSON(http.StatusCreated, gin.H{"message": "set created!"})
}
