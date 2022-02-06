package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/handler"
)

func main() {
	router := gin.Default()
	router.GET("/sets", handler.ReadSets)
	router.GET("/sets/:id", handler.ReadSet)
	router.POST("/sets", handler.CreateSet)

	router.Run("localhost:8080")
}
