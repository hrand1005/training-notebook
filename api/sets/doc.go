// Package classification of Set API
//
// Documentation for Set API
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
package sets

import (
	"github.com/gin-gonic/gin"
	"github.com/hrand1005/training-notebook/data"
)

type set struct {
	db data.SetDB
}

// returns a set in the response
// swagger:response setResponse
type setResponse struct {
	// A single set
	// in: body
	Body data.Set
}

// returns sets in the response
// swagger:response setsResponse
type setsResponse struct {
	// A list of sets
	// in: body
	Body []data.Set
}

// returns generic error message as string
// swagger:response errorResponse
type errorResponse struct {
	// Description of the erro
	// in: body
	Body gin.H
}

// swagger:parameters readSet
// swagger:parameters updateSet
// swagger:parameters deleteSet
type setIDParameter struct {
	// The id of the set
	// in: required: true
	// required: true
	ID int `json:"id"`
}
