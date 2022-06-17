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
	// TODO: REQUIRE ADMIN PERMS
	//g.GET("/", u.ReadAll)
	//g.POST("/", u.Create)

	// Authorization NOT required
	g.POST("/signup", u.Signup)
	g.POST("/login", u.Login)

	// Require User Authentication for Read/Updates on a user
	authGroup := g.Group("")
	authGroup.Use(RequireAuthorization())
	authGroup.GET("/:userID", u.Read)
	authGroup.PUT("/:userID", u.Update)

}

func userIDFromParams(c *gin.Context) (models.UserID, error) {
	id, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		return data.InvalidUserID, err
	}
	if id < 0 {
		return data.InvalidUserID, fmt.Errorf("user id cannot be negative")
	}

	return models.UserID(id), nil
}
