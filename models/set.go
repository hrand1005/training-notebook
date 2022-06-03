package models

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

// SetID is the unique int identifier assigned to sets when added to the SetDB
type SetID int

// Set is the model representation of a sets resource
// swagger:model
type Set struct {
	// the id for this set
	//
	// required: true
	// min: 1
	ID            SetID   `json:"id"`
	Movement      string  `json:"movement" binding:"movement"`
	Volume        float64 `json:"volume" binding:"gt=0"`
	Intensity     float64 `json:"intensity" binding:"gt=0,lte=100"`
	CreatedOn     string  `json:"-"`
	LastUpdatedOn string  `json:"-"`
}

// MovementValidator validates the movement field in a Set.
// Returns true if valid, else false.
var MovementValidator validator.Func = func(fl validator.FieldLevel) bool {
	// rule for movement string checks for unicode characters
	rule := regexp.MustCompile(`^\w+(\s+\w+)*$`)
	matches := rule.FindAllString(fl.Field().String(), -1)

	return len(matches) == 1
}

// SetsEqual returns true if all non-id fields of the set are equal, and false otherwise.
func SetsEqual(set1, set2 *Set) bool {
	if set1.Movement != set2.Movement || set1.Volume != set2.Volume || set1.Intensity != set2.Intensity {
		return false
	}
	return true
}

// ContainsSet checks if the slice of sets contains the provided set.
func ContainsSet(sets []*Set, s *Set) bool {
	for _, v := range sets {
		if SetsEqual(s, v) {
			return true
		}
	}
	return false
}

/*
* The following code exists to translate binding errors into readable messages
* for the API client. Use BindingErrorToMessage to produce readable messages
* from errors binding json to Sets.
* TODO: put this in util/ or pkg/ ?
 */

// BindingErrorToMessage turns a binding error from a set field into a
// a readable message on the field's constraints.
func BindingErrorToMessage(err error) string {
	var msg string
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, v := range ve {
			msg += fmt.Sprintf("'%s' field %v.", v.Field(), fieldErrorToMessage(v))
		}
	}
	return msg
}

// customTagToMessage gets condition for custom validators
var customTagToMessage = map[string]string{
	"movement": "must use unicode characters",
}

// standardTagToMessage gets conditions for standard validators
var standardTagToMessage = map[string]string{
	"gt":  "must be greater than",
	"lte": "must be no more than",
}

// fieldErrorToMessage gets either standard or custom error messages
func fieldErrorToMessage(fieldError validator.FieldError) string {
	if msg, ok := standardTagToMessage[fieldError.Tag()]; ok {
		return msg + " " + fieldError.Param()
	}

	if msg, ok := customTagToMessage[fieldError.Tag()]; ok {
		return msg
	}

	return fieldError.Error()
}
