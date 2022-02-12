package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/data"
)

func CreateBook(c *gin.Context) {
	var newBook data.Book

	if err := c.BindJSON(&newBook); err != nil {
		log.Printf("could not bind json to book: %v", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "could not bind json to book"})
		return
	}

	// assigns ID to newBook
	data.AddBook(&newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
	return
}

func ReadBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, data.Books())
	return
}

func ReadBook(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	s, err := data.BookByID(id)
	if err != nil {
		log.Printf("could not read book: %v", err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, s)
	return
}

func UpdateBook(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var newBook data.Book

	if err := c.BindJSON(&newBook); err != nil {
		log.Printf("could not bind json to book: %v", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "could not bind json to book"})
		return
	}

	// assigns newBook ID of id
	if err := data.UpdateBook(id, &newBook); err != nil {
		log.Printf("could not update book: %v", err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, newBook)
	return
}

func DeleteBook(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := data.DeleteBook(id); err != nil {
		log.Printf("could not delete book: %v", err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found"})
		return
	}

	c.IndentedJSON(http.StatusNoContent, gin.H{})
	return
}
