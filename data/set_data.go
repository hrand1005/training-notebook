// This file defines setData, a stand-in (and testing utility) for a database.
// setData implements the SetDB interface.
package data

// NewSetData returns a new setData object initialized with the provided slice of
// sets. To initialize an empty setData, provide nil.
func NewSetData(sets []*Set) *setData {
	return &setData{sets: sets}
}

// setData contains a slice of Sets and implements the SetData interface
type setData struct {
	sets []*Set
}

// Replace with DB logic
func (sd *setData) AddSet(s *Set) error {
	if len(sd.sets) == 0 {
		s.ID = 1
		sd.sets = append(sd.sets, s)
		return nil
	}

	maxID := sd.sets[len(sd.sets)-1].ID
	s.ID = maxID + 1
	sd.sets = append(sd.sets, s)
	return nil
}

func (sd *setData) Sets() ([]*Set, error) {
	return sd.sets, nil
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

// Example data, to be replaced with db
var TestSetData = &setData{
	sets: []*Set{
		{
			ID:        1,
			Movement:  "Squat",
			Volume:    5,
			Intensity: 80,
		},
	},
}
