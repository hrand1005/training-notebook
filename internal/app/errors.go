package app

import (
	"errors"
	"fmt"
	// "log"

	"github.com/go-playground/validator/v10"
)

// ErrNotFound should be used for failures to obtain a resource by request parameters
var ErrNotFound = errors.New("failed to find the resource")

// ErrInvalidField should be used when an entity contains invalid fields
var ErrInvalidField = errors.New("invalid field")

// ErrServiceFailure should be used for failures that are specific to a service's implementation
var ErrServiceFailure = errors.New("the service failed due to internal reasons")

var validate = validator.New()

// ValiateEntity attempts to validate the provided interface struct.
// Returns error immediately if one is encountered.
func ValidateEntity(entity interface{}) error {
	// log.Printf("Encountered entity: %+v\n", entity)
	if err := validate.Struct(entity); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return formatFieldError(validationErrors[0])
	}

	return nil
}

// ValidateAllEntityFields attempts to validate all fields in the provided interface struct.
// Returns all errors encountered during validation.
func ValidateAllEntityFields(entity interface{}) []error {
	// log.Printf("Encountered entity: %+v\n", entity)
	if err := validate.Struct(entity); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return formatValidationErrors(validationErrors)
	}

	return nil
}

func formatValidationErrors(ve validator.ValidationErrors) []error {
	errs := make([]error, 0, len(ve))
	for _, v := range ve {
		errs = append(errs, formatFieldError(v))
	}
	return errs
}

// formatFieldError turns validation errors into client-readable messages.
func formatFieldError(f validator.FieldError) error {
	switch f.Value().(type) {
	case string:
		return formatStringFieldError(f)
	}
	return fmt.Errorf("%w: '%s'", ErrInvalidField, f)
}

func formatStringFieldError(f validator.FieldError) error {
	var msg string
	switch f.Tag() {
	case "min":
		msg = fmt.Sprintf("'%s' must be at least %v characters", f.Field(), f.Param())
	case "max":
		msg = fmt.Sprintf("'%s' cannot exceed %v characters", f.Field(), f.Param())
	case "email":
		msg = fmt.Sprintf("invalid email address")
	default:
		msg = fmt.Sprintf("'%s'", f.Field())
	}
	return fmt.Errorf("%w: %s", ErrInvalidField, msg)
}
