package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Set struct {
	ID string `json:"id"`
	Movement string `json:"movement"`
	Volume float64 `json:"volume"`
	Intensity float64 `json:"intensity"`
}

var sets = []Set{
	{
		ID: "1",
		Movement: "Squat",
		Volume: 5,
		Intensity: 80,
	},
}

func getSets(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, sets)
}

func getSetByID(c *gin.Context) {
	id := c.Param("id")

	for _, s := range sets {
		if s.ID == id {
			c.IndentedJSON(http.StatusOK, s)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "set not found"})
}

func postSets(c *gin.Context) {
	var newSet Set

	if err := c.BindJSON(&newSet); err != nil {
		log.Printf("could not bind json to set: %v", err)
		return
	}

	sets = append(sets, newSet)
	c.IndentedJSON(http.StatusCreated, newSet)
	return
}

func main() {
	router := gin.Default()
	router.GET("/sets", getSets)
	router.GET("/sets/:id", getSetByID)
	router.POST("/sets", postSets)

	router.Run("localhost:8080")
}