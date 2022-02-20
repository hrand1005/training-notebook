package sets

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/data"
)

// swagger:route DELETE /sets/{id} sets deleteSet
// Delete a set.
// responses:
//  204: noContent
//  404: errorResponse

// Delete is the handler for delete requests on the set resource. An id must be
// specified.
func (s *set) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := data.DeleteSet(id); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "set not found"})
		return
	}

	c.IndentedJSON(http.StatusNoContent, gin.H{})
	return
}
