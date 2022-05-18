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

type SetDB interface {
	AddSet(s *Set)
	Sets() []*Set
	SetByID(id int) (*Set, error)
	UpdateSet(id int, s *Set) error
	DeleteSet(id int) error
}

// setData contains a slice of Sets and implements the SetData interface
type setData struct {
	sets []*Set
}

// Replace with DB logic
func (sd *setData) AddSet(s *Set) {
	if len(sd.sets) == 0 {
		s.ID = 1
		sd.sets = append(sd.sets, s)
		return
	}

	maxID := sd.sets[len(sd.sets)-1].ID
	s.ID = maxID + 1
	sd.sets = append(sd.sets, s)

	return
}

func (sd *setData) Sets() []*Set {
	return sd.sets
}

func (sd *setData) SetByID(id int) (*Set, error) {
	for _, s := range sd.sets {
		if s.ID == id {
			return s, nil
		}
	}

	return nil, ErrNotFound
}

func (sd *setData) UpdateSet(id int, s *Set) error {
	for i, v := range sd.sets {
		if id == v.ID {
			s.ID = id
			sd.sets[i] = s
			return nil
		}
	}

	return ErrNotFound
}

func (sd *setData) DeleteSet(id int) error {
	for i, s := range sd.sets {
		if s.ID == id {
			sd.sets = append(sd.sets[:i], sd.sets[i+1:]...)
			return nil
		}
	}

	return ErrNotFound
}
