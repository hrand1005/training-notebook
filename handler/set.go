package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/hrand1005/training-notebook/data"
)

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
