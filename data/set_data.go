// This file defines setData, a stand-in (and testing utility) for a database.
// setData implements the SetDB interface.
package data

import "github.com/hrand1005/training-notebook/models"

// NewSetData returns a new setData object initialized with the provided slice of
// sets. To initialize an empty setData, provide nil.
func NewSetData(sets []*models.Set) *setData {
	return &setData{sets: sets}
}

// setData contains a slice of Sets and implements the SetData interface
type setData struct {
	sets []*models.Set
}

// Replace with DB logic
func (sd *setData) AddSet(s *models.Set) (models.SetID, error) {
	if len(sd.sets) == 0 {
		s.ID = 1
		sd.sets = append(sd.sets, s)
		return 1, nil
	}

	maxID := sd.sets[len(sd.sets)-1].ID
	s.ID = maxID + 1
	sd.sets = append(sd.sets, s)
	return s.ID, nil
}

func (sd *setData) Sets() ([]*models.Set, error) {
	return sd.sets, nil
}

func (sd *setData) SetByID(id models.SetID) (*models.Set, error) {
	for _, s := range sd.sets {
		if s.ID == id {
			return s, nil
		}
	}

	return nil, ErrNotFound
}

func (sd *setData) UpdateSet(id models.SetID, s *models.Set) error {
	for i, v := range sd.sets {
		if id == v.ID {
			s.ID = id
			sd.sets[i] = s
			return nil
		}
	}

	return ErrNotFound
}

func (sd *setData) DeleteSet(id models.SetID) error {
	for i, s := range sd.sets {
		if s.ID == id {
			sd.sets = append(sd.sets[:i], sd.sets[i+1:]...)
			return nil
		}
	}

	return ErrNotFound
}

// Close cleans up any outstanding resources. In this case, returns nil
func (sd *setData) Close() error {
	// always returns nil
	return nil
}

// Example data, to be replaced with db
var TestSetData = &setData{
	sets: []*models.Set{
		{
			ID:        1,
			Movement:  "Squat",
			Volume:    5,
			Intensity: 80,
		},
	},
}
