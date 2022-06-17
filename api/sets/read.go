package sets

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/data"
	"github.com/hrand1005/training-notebook/models"
)

// swagger:route GET /sets sets readAllSets
// Read all sets.
// responses:
//  200: setsResponse
//  500: errorResponse

// ReadAll is the handler for read requests on the set resource where no id is
// specified. Returns all sets on this resource's data source.
func (s *set) ReadAll(c *gin.Context) {
	sets, err := s.db.Sets()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if len(sets) == 0 {
		// if no sets are found, return an empty slice
		c.IndentedJSON(http.StatusOK, []*models.Set{})
		return
	}
	c.IndentedJSON(http.StatusOK, sets)
}

// swagger:route GET /sets/{id} sets readSet
// Read a set.
// responses:
//  200: setResponse
//  400: errorResponse
//  404: errorResponse
//  500: errorResponse

// Read is the handler for read requests on the set resource where an id is
// specified.
func (s *set) Read(c *gin.Context) {
	setID, err := SetIDFromParams(c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": ErrInvalidSetID})
		return
	}

	r, err := s.db.SetByID(setID)
	if err != nil {
		if err == data.ErrNotFound {
			msg := fmt.Sprintf("no such set with id %v", setID)
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": msg})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, r)
}
