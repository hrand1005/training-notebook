package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/data"
)

/*
type Set struct {

}

func NewSet(l *log.Logger) *Set {
}*/

func CreateSet(c *gin.Context) {
	newSet := c.MustGet("newSet").(data.Set)

	// assigns ID to newSet
	data.AddSet(&newSet)
	c.IndentedJSON(http.StatusCreated, newSet)
	return
}

func ReadSets(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, data.Sets())
	return
}

func ReadSet(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	s, err := data.SetByID(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "set not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, s)
	return
}

func UpdateSet(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	newSet := c.MustGet("newSet").(data.Set)

	// assigns newSet ID of id
	if err := data.UpdateSet(id, &newSet); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "set not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, newSet)
	return
}

func DeleteSet(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := data.DeleteSet(id); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "set not found"})
		return
	}

	c.IndentedJSON(http.StatusNoContent, gin.H{})
	return
}

// Validator middleware validates provided set data
func JSONSetValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		var newSet data.Set

		if err := c.BindJSON(&newSet); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "could not bind json to set"})
			return
		}

		// set newSet variable
		c.Set("newSet", newSet)
	}
}
