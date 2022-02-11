package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/data"
)

func Home(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, data.Home())
	return
}
