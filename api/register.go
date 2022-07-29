package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/api/sets"
	"github.com/hrand1005/training-notebook/api/users"
	"github.com/hrand1005/training-notebook/data"
)

// RegisterAll uses registers all api endpoints on the given RouterGroup.
// Additionally, the provided db handle will be used for all resources.
func RegisterAll(db *sql.DB, g *gin.RouterGroup) error {
	setDB, err := data.NewSetDB(db)
	if err != nil {
		return err
	}
	setResource, err := sets.New(setDB)
	if err != nil {
		return err
	}
	setResource.RegisterHandlers(g)

	userDB, err := data.NewUserDB(db)
	if err != nil {
		return err
	}
	userResource, err := users.New(userDB)
	if err != nil {
		return err
	}
	userResource.RegisterHandlers(g)

	return nil
}
