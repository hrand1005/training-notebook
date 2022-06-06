package users

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/data"
	"github.com/hrand1005/training-notebook/models"
)

var ErrInvalidUserID = "Invalid user ID"

type user struct {
	db data.UserDB
}

func New(db data.UserDB) (*user, error) {
	return &user{db: db}, nil
}

func (u *user) RegisterHandlers(g *gin.RouterGroup) {
	g.GET("/", u.ReadAll)
	g.GET("/:id", u.Read)
	g.POST("/", u.Create)
	g.PUT("/:id", u.Update)
}

func userIDFromParams(c *gin.Context) (models.UserID, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return data.InvalidUserID, err
	}
	if id < 0 {
		return data.InvalidUserID, fmt.Errorf("user id cannot be negative")
	}

	return models.UserID(id), nil
}
