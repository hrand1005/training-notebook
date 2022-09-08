package app

import (
	"errors"
	"log"

	"github.com/go-playground/validator/v10"
)

// ErrNotFound should be used for failures to obtain a resource by request parameters
var ErrNotFound = errors.New("failed to find the resource")

// ErrServiceFailure should be used for failures that are specific to a service's implementation
var ErrServiceFailure = errors.New("the service failed due to internal reasons")

var validate = validator.New()

// ValidateRequestBody checks the validate struct tags for the provided
// request body. If invalid fields are encountered, they are returned in the
// slice of error responses.
func ValidateEntity(entity interface{}) validator.ValidationErrors {
	log.Printf("Encountered entity: %+v\n", entity)
	if err := validate.Struct(entity); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return validationErrors
	}

	return nil
}
