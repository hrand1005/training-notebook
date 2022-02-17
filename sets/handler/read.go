package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/sets/data"
)

// swagger:route GET /sets sets readAllSets
// Read all sets.
// responses:
//  200: setsResponse

// ReadAll is the handler for read requests on the set resource where no id is
// specified.
func (s *set) ReadAll(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, data.Sets())
	return
}

// swagger:route GET /sets/{id} sets readSet
// Read a set.
// responses:
//  200: setResponse
//  404: errorResponse

// Read is the handler for read requests on the set resource where an id is
// specified.
func (s *set) Read(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	r, err := data.SetByID(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "set not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, r)
	return
}
