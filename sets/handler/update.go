package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/sets/data"
)

// swagger:route PUT /sets/{id} sets updateSet
// Update a set.
// responses:
//  200: setResponse
//  400: errorResponse
//  404: errorResponse

// Update is the handler for update requests on the set resource. An id must be
// specified.
// Requires JSONValidator to be registered with the router group.
func (s *set) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	newSet := c.MustGet("newSet").(data.Set)

	// assigns newSet ID of id
	if err := data.UpdateSet(id, &newSet); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "set not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, newSet)
	return
}
