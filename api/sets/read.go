package sets

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/api/users"
	"github.com/hrand1005/training-notebook/data"
	"github.com/hrand1005/training-notebook/models"
)

// swagger:route GET /sets sets readAllSets
// Read all sets.
// responses:
//  200: setsResponse
// 	401: errorResponse
//  500: errorResponse

// ReadAll is the handler for read requests on the set resource where no id is
// specified. Returns all sets on this resource's data source.
func (s *set) ReadAll(c *gin.Context) {
	userID, err := users.UserIDFromContext(c)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "you must be logged in to perform this action"})
		return
	}

	sets, err := s.db.SetsByUserID(userID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if len(sets) == 0 {
		// if no sets are found, return an empty slice
		c.IndentedJSON(http.StatusOK, []*models.Set{})
		return
	}
	c.IndentedJSON(http.StatusOK, sets)
}

// swagger:route GET /sets/{id} sets readSet
// Read a set.
// responses:
//  200: setResponse
//  400: errorResponse
// 	401: errorResponse
//  404: errorResponse
//  500: errorResponse

// Read is the handler for read requests on the set resource where an id is
// specified.
func (s *set) Read(c *gin.Context) {
	userID, err := users.UserIDFromContext(c)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "you must be logged in to perform this action"})
		return
	}

	setID, err := SetIDFromParams(c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": ErrInvalidSetID})
		return
	}

	resultSet, err := s.db.SetByIDForUser(setID, userID)
	if err != nil {
		if err == data.ErrNotFound {
			msg := fmt.Sprintf("no such set with id %v for logged in user", setID)
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": msg})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, resultSet)
}
