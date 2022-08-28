package errors

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// Response defines an error response.
type Response struct {
	InvalidField string `json:"invalid-field"`
	Tag          string `json:"tag"`
	Value        string `json:"value"`
}

// ValidateRequestBody checks the validate struct tages for the provided
// request body. If invalid fields are encountered, they are returned in the
// slice of error responses.
func ValidateRequestBody(body interface{}) []*Response {
	var errorResponses []*Response

	err := validate.Struct(body)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, v := range validationErrors {
			errorResponses = append(errorResponses, &Response{
				InvalidField: v.StructNamespace(),
				Tag:          v.Tag(),
				Value:        v.Param(),
			})
		}
	}

	return errorResponses
}
