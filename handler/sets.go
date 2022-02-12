package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/data"
)

func CreateSet(c *gin.Context) {
	var newSet data.Set

	if err := c.BindJSON(&newSet); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "could not bind json to set"})
		return
	}

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
	var newSet data.Set

	if err := c.BindJSON(&newSet); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "could not bind json to set"})
		return
	}

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
