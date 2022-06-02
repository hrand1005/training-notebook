package users

import (
	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/data"
)

type user struct {
	db data.UserDB
}

func New(db data.UserDB) (*user, error) {
	return &user{db: db}, nil
}

func (u *user) RegisterHandlers(g *gin.RouterGroup) {
	g.GET("/", u.ReadAll)
	//	g.POST("/", u.Create)
}
