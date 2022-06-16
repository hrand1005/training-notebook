package sets

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/hrand1005/training-notebook/data"
	"github.com/hrand1005/training-notebook/models"
)

var ErrInvalidSetID = "Invalid set ID"

// New registers custom validators with the validator engine and returns the
// handler for the set resource.
func New(db data.SetDB) (*set, error) {
	// register set validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// TODO: remove this fancy business
		v.RegisterValidation("movement", models.MovementValidator)
		return &set{db: db}, nil
	}

	return nil, errors.New("failed to access validator engine")
}

func (s *set) RegisterHandlers(g *gin.RouterGroup) {
	// g.Use(RequireAuthorization)
	g.GET("/", s.ReadAll)
	g.GET("/:setID", s.Read)
	g.DELETE("/:setID", s.Delete)
	g.POST("/", s.Create)
	g.PUT("/:setID", s.Update)
}

/*
func (s *set) RegisterHandlersWithUserAuthentication(g *gin.RouterGroup) {
	// Register set CRUD operations on a group with the UserAuthentication middleware
	// The client must provide a jwt in the request header
	//g.Use(UserAuthentication)

	// GET, DELETE, and PUT requests must verify that the retrieved set belongs to the user
	// GET for ReadAll must return all sets matching the requesting user's id
	g.GET("/", s.ReadAll)
}
*/

func setIDFromParams(c *gin.Context) (models.SetID, error) {
	id, err := strconv.Atoi(c.Param("setID"))
	if err != nil {
		return data.InvalidSetID, err
	}
	if id < 0 {
		return data.InvalidSetID, fmt.Errorf("set id cannot be negative")
	}

	return models.SetID(id), nil
}
