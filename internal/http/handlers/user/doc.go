// Package classification of Users API.
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

type responseFields struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// swagger:parameters createUser
type userObject struct {
	ID         string         `json:"id,omitempty"`
	Type       string         `json:"type"`
	Attributes userAttributes `json:"attributes"`
}

type userAttributes struct {
	FirstName string `json:"first-name"`
	LastName  string `json:"last-name"`
	Email     string `json:"email"`
}

type userRequest struct {
	Data userObject `json:"data"`
}

// returns a user in the response
// swagger:response userResponse
type userResponse struct {
	// A single user
	// in:body
	Data userObject `json:"data"`
}

// returns users in the response
// swagger:response userResponse
type usersResponse struct {
	// Multiple users
	// in: body
	Data []userObject `json:"data"`
}
