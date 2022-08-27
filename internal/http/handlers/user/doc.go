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

type Attributes struct {
	FirstName string `json:"first-name"`
	LastName  string `json:"last-name"`
	Email     string `json:"email"`
}

type RequestData struct {
	Type       string     `json:"type"`
	Attributes Attributes `json:"attributes"`
}

type ResponseData struct {
	ID         string     `json:"id"`
	Type       string     `json:"type"`
	Attributes Attributes `json:"attributes"`
}

type RequestBody struct {
	Data RequestData `json:"data"`
}

type ResponseBody struct {
	Data ResponseData `json:"data"`
}
