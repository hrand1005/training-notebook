// Package classification of Set API
//
// Documentation for Set API
//
//  Schemes: http
//  BasePath: /
//  Version: 1.0.0
//
//  Consumes:
//  - application/json
//
//  Produces:
//  - application/json
// swagger:meta
package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/hrand1005/training-notebook/data"
)

type set struct{}

// returns a set in the response
// swagger:response setResponse
type setResponse struct {
	// A single set
	// in: body
	Body data.Set
}

// returns sets in the response
// swagger:response setsResponse
type setsResponse struct {
	// A list of sets
	// in: body
	Body []data.Set
}

// swagger:parameters readSet
// swagger:parameters updateSet
// swagger:parameters deleteSet
type setIDParameter struct {
	// The id of the set
	// in: required: true
	// required: true
	ID int `json:"id"`
}

// NewSet registers custom validators with the validator engine and returns the
// handler for the set resource.
func NewSet() (*set, error) {
	// register set validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("movement", data.MovementValidator)
		return &set{}, nil
	}

	return nil, errors.New("failed to access validator engine")
}

// swagger:route POST /sets sets createSet
// Creates a set.
// responses:
//  201: setResponse

// Create is the handler for create requests on the set resource.
// Requires JSONValidator to be registered with the router group.
func (s *set) Create(c *gin.Context) {
	newSet := c.MustGet("newSet").(data.Set)

	// assigns ID to newSet
	data.AddSet(&newSet)
	c.IndentedJSON(http.StatusCreated, newSet)
	return
}

// swagger:route GET /sets sets readAllSets
// Read all sets.
// responses:
//  200: setsResponse

// ReadAll is the handler for read requests on the set resource where no id is
// specified.
func (s *set) ReadAll(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, data.Sets())
	return
}

// swagger:route GET /sets/{id} sets readSet
// Read a set.
// responses:
//  200: setResponse

// Read is the handler for read requests on the set resource where an id is
// specified.
func (s *set) Read(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	r, err := data.SetByID(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "set not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, r)
	return
}

// swagger:route PUT /sets/{id} sets updateSet
// Update a set.
// responses:
//  200: setResponse

// Update is the handler for update requests on the set resource. An id must be
// specified.
// Requires JSONValidator to be registered with the router group.
func (s *set) Update(c *gin.Context) {
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

// swagger:route DELETE /sets/{id} sets deleteSet
// Delete a set.
// responses:
//  204: noContent

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

// JSONValidator is middleware that validates set data in the request body.
// This must be registered with the router group in order for Creates and
// Updates on this resource.
func (s *set) JSONValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		var newSet data.Set

		if err := c.BindJSON(&newSet); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		// set newSet variable
		c.Set("newSet", newSet)
	}
}
