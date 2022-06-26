// Package classification of User API
//
// Documentation for User API
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
package users

import (
	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/data"
	"github.com/hrand1005/training-notebook/models"
)

type user struct {
	db data.UserDB
}

// returns a user in the response
// swagger:response userResponse
type userResponse struct {
	// A single user
	// in: body
	Body models.User
}

// returns users in the response
// swagger:response usersResponse
type usersResponse struct {
	// A list of users
	// in: body
	Body []models.User
}

// returns generic error message as string
// swagger:response errorResponse
type errorResponse struct {
	// Description of the error
	// in: body
	Body gin.H
}

// swagger:parameters readuser
// swagger:parameters updateuser
// swagger:parameters deleteuser
type userIDParameter struct {
	// The id of the user
	// in: required: true
	// required: true
	ID int `json:"user-id"`
}
