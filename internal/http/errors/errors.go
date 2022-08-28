package errors

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func init() {
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "" {
			return ""
		}
		return name
	})
}

type FormattedError struct {
	Message string `json:"message"`
}

// ValidateRequestBody checks the validate struct tages for the provided
// request body. If invalid fields are encountered, they are returned in the
// slice of error responses.
func ValidateRequestBody(body interface{}, errorFactory func(validator.FieldError) FormattedError) []FormattedError {
	var errorResponses []FormattedError

	err := validate.Struct(body)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, v := range validationErrors {
			errorResponses = append(errorResponses, errorFactory(v))
		}
	}

	return errorResponses
}
