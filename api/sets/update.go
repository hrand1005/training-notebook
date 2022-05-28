package sets

import (
	"net/http"

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

	setID, err := setIDFromParams(c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": ErrInvalidSetID})
		return
	}

	if err := s.db.UpdateSet(setID, &newSet); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	newSet.ID = setID
	c.IndentedJSON(http.StatusOK, newSet)
}
