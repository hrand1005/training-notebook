package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/sets/data"
)

// swagger:route POST /sets sets createSet
// Creates a set.
// responses:
//  201: setResponse

// Create is the handler for create requests on the set resource.
// Requires JSONValidator to be registered with the router group.
func (s *set) Create(c *gin.Context) {
	newSet := c.MustGet("newSet").(data.Set)

	// assigns ID to newSet
	data.AddSet(&newSet)
	c.IndentedJSON(http.StatusCreated, newSet)
	return
}
