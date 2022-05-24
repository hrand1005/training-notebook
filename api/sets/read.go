package sets

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/data"
)

// swagger:route GET /sets sets readAllSets
// Read all sets.
// responses:
//  200: setsResponse

// ReadAll is the handler for read requests on the set resource where no id is
// specified.
func (s *set) ReadAll(c *gin.Context) {
	sets, err := s.db.Sets()
	if err != nil {
		msg := fmt.Sprintf("failed to fetch data: %v", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": msg})
	}
	if len(sets) == 0 {
		// if no sets are found, return an empty slice
		c.IndentedJSON(http.StatusOK, []*data.Set{})
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

// Read is the handler for read requests on the set resource where an id is
// specified.
func (s *set) Read(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid request parameters"})
		return
	}

	r, err := s.db.SetByID(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "set not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, r)
}
