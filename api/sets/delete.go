package sets

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/data"
)

// swagger:route DELETE /sets/{id} sets deleteSet
// Delete a set.
// responses:
//  204: noContent
//  404: errorResponse
//  500: errorResponse

// Delete is the handler for delete requests on the set resource. An id must be
// specified.
func (s *set) Delete(c *gin.Context) {
	setID, err := SetIDFromParams(c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": ErrInvalidSetID})
		return
	}

	if err := s.db.DeleteSet(setID); err != nil {
		if err == data.ErrNotFound {
			msg := fmt.Sprintf("no such set with id %v", setID)
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": msg})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusNoContent, gin.H{})
}
