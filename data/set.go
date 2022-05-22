package data

import (
	"errors"
	"regexp"

	"github.com/go-playground/validator/v10"
)

// ErrNotFound should be returned when the resource does not exist.
var ErrNotFound = errors.New("resource not found")

// Set defines the structure for an API set
// swagger:model
type Set struct {
	// the id for this set
	//
	// required: true
	// min: 1
	ID            int     `json:"id"`
	Movement      string  `json:"movement" binding:"movement"`
	Volume        float64 `json:"volume" binding:"gt=0"`
	Intensity     float64 `json:"intensity" binding:"gt=0,lte=100"`
	CreatedOn     string  `json:"-"`
	LastUpdatedOn string  `json:"-"`
}

// MovementValidator validates the movement field in a Set.
var MovementValidator validator.Func = func(fl validator.FieldLevel) bool {
	// rule for movement string checks for unicode characters
	rule := regexp.MustCompile(`^\w+(\s+\w+)*$`)
	matches := rule.FindAllString(fl.Field().String(), -1)

	// returns true if valid, else false
	return len(matches) == 1
}
