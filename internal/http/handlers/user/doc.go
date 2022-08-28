// Package classification Users API.
//
//   Schemes: http, https
//   Host: localhost
//   BasePath: /api/v1/users
//   Version: 0.0.1
//   License: None
//   Contact: Herbie Rand<hjrand1005@gmail.com> http://hrand1005.github.io
//
//   Consumes:
//   - application/json
//
//   Produces:
//   -application/json
//
//   Security:
//   - api_key:
//
//   SecurityDefinitions:
//   api_key:
//      type: apiKey
//      name: KEY
//      in: header
//
// swagger:meta
package user

import (
	"github.com/hrand1005/training-notebook/internal/http/errors"
)

// swagger:parameters createUser
type userRequest struct {
	// A single user
	// in:body
	Body RequestBody
}

// returns a user in the response
// swagger:response userResponse
type userResponse struct {
	// A single user
	// in:body
	Body ResponseBody
}

// returns users in the response
// swagger:response usersResponse
type usersResponse struct {
	// A single user
	// in:body
	Body []ResponseBody
}

// returns errors in the response
// swagger:response errorsResponse
type errorsResponse struct {
	// api errors
	// in:body
	Body []ErrorsResponseBody
}

// --- EMBEDDED STRUCTS ---

type Attributes struct {
	FirstName string `json:"first-name" validate:"required,min=2,max=32"`
	LastName  string `json:"last-name" validate:"required,min=2,max=32"`
	Email     string `json:"email" validate:"required,email,min=6,max=32"`
}

type RequestData struct {
	Type       string     `json:"type" validate:"required"`
	Attributes Attributes `json:"attributes" validate:"required"`
}

type ResponseData struct {
	ID         string     `json:"id" validate:"required"`
	Type       string     `json:"type" validate:"required"`
	Attributes Attributes `json:"attributes" validate:"required"`
}

type RequestBody struct {
	Data RequestData `json:"data" validate:"required"`
}

type ResponseBody struct {
	Data ResponseData `json:"data" validate:"required"`
}

type ErrorsResponseBody struct {
	Errors []errors.FormattedError `json:"errors"`
}
