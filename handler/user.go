package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/data"
)

func CreateUser(c *gin.Context) {
	var newUser data.User

	if err := c.BindJSON(&newUser); err != nil {
		log.Printf("could not bind json to set: %v", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "could not bind json to set"})
		return
	}

	// assigns ID to newUser
	if err := data.AddUser(&newUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err})
	}
	c.IndentedJSON(http.StatusCreated, newUser)
	return
}

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

func UpdateUser(c *gin.Context) {
	userID := c.Param("userID")
	var newUser data.User

	if err := c.BindJSON(&newUser); err != nil {
		log.Printf("could not bind json to set: %v", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "could not bind json to set"})
		return
	}

	// assigns newUser ID of id
	if err := data.UpdateUser(userID, &newUser); err != nil {
		log.Printf("could not update set: %v", err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "set not found"})
		return
	}

	newUser.UserID = userID
	c.IndentedJSON(http.StatusOK, newUser)
	return
}

func DeleteUser(c *gin.Context) {
	userID := c.Param("userID")

	if err := data.DeleteUser(userID); err != nil {
		log.Printf("could not delete set: %v", err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "set not found"})
		return
	}

	c.IndentedJSON(http.StatusNoContent, gin.H{})
	return
}
