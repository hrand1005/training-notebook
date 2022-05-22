package sets

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/data"
)

// swagger:route PUT /sets/{id} sets updateSet
// Update a set.
// responses:
//  200: setResponse
//  400: errorResponse
//  404: errorResponse

// Update is the handler for update requests on the set resource. An id must be
// specified.
func (s *set) Update(c *gin.Context) {
	var newSet data.Set

	if err := c.BindJSON(&newSet); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))

	// assigns newSet ID of id
	if err := s.db.UpdateSet(id, &newSet); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "set not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, newSet)
}
