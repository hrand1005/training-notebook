package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/data"
)

func ReadUser(c *gin.Context) {
	userID := c.Param("userID")

	// more may need to be contained within the body of the request to produce a
	// valid response...

	u, err := data.UserByUserID(userID)
	if err != nil {
		log.Printf("could not read set: %v", err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "set not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, u)
	return
}
