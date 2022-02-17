package data

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"regexp"
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

	if len(matches) != 1 {
		return false
	}

	return true
}

// Replace with DB logic
func AddSet(s *Set) {
	if len(sets) == 0 {
		s.ID = 1
		sets = append(sets, s)
		return
	}

	maxID := sets[len(sets)-1].ID
	s.ID = maxID + 1
	sets = append(sets, s)

	return
}

func Sets() []*Set {
	return sets
}

func SetByID(id int) (*Set, error) {
	for _, s := range sets {
		if s.ID == id {
			return s, nil
		}
	}

	return nil, ErrNotFound
}

func UpdateSet(id int, s *Set) error {
	for i, v := range sets {
		if id == v.ID {
			s.ID = id
			sets[i] = s
			return nil
		}
	}

	return ErrNotFound
}

func DeleteSet(id int) error {
	for i, s := range sets {
		if s.ID == id {
			sets = append(sets[:i], sets[i+1:]...)
			return nil
		}
	}

	return ErrNotFound
}
