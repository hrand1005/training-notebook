package sets

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// swagger:route DELETE /sets/{id} sets deleteSet
// Delete a set.
// responses:
//  204: noContent
//  404: errorResponse

// Delete is the handler for delete requests on the set resource. An id must be
// specified.
func (s *set) Delete(c *gin.Context) {
	setID, err := setIDFromParams(c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": ErrInvalidSetID})
		return
	}

	if err := s.db.DeleteSet(setID); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusNoContent, gin.H{})
}
