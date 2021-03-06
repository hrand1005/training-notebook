package sets

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/hrand1005/training-notebook/api/users"
	"github.com/hrand1005/training-notebook/data"
	"github.com/hrand1005/training-notebook/models"
)

var ErrInvalidSetID = "Invalid set ID"

// Keys used to retrieve values from gin.Context
const (
	SetIDFromParamsKey = "paramSetID"
)

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
	// register RequireAuthorization middleware so that each request
	// on the sets requires token
	setGroup := g.Group("/sets")
	setGroup.Use(users.RequireAuthorization())
	setGroup.GET("/", s.ReadAll)
	setGroup.GET("/:"+SetIDFromParamsKey, s.Read)
	setGroup.DELETE("/:"+SetIDFromParamsKey, s.Delete)
	setGroup.POST("/", s.Create)
	setGroup.PUT("/:"+SetIDFromParamsKey, s.Update)
}

func SetIDFromParams(c *gin.Context) (models.SetID, error) {
	id, err := strconv.Atoi(c.Param(SetIDFromParamsKey))
	if err != nil {
		return data.InvalidSetID, err
	}
	if id < 0 {
		return data.InvalidSetID, fmt.Errorf("set id cannot be negative")
	}

	return models.SetID(id), nil
}
