package sets

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/hrand1005/training-notebook/data"
)

// New registers custom validators with the validator engine and returns the
// handler for the set resource.
func New() (*set, error) {
	// register set validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("movement", data.MovementValidator)
		return &set{}, nil
	}

	return nil, errors.New("failed to access validator engine")
}

func (s *set) RegisterHandlers(g *gin.RouterGroup) {
	// Register GET and DELETE requests on routerGroup
	g.GET("", s.ReadAll)
	g.GET("/:id", s.Read)
	g.DELETE("/:id", s.Delete)

	// create a subgroup with JSON validation
	validateGroup := g.Group("")
	validateGroup.Use(s.JSONValidator())
	validateGroup.POST("/", s.Create)
	validateGroup.PUT("/:id", s.Update)
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
