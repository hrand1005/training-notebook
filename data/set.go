package data

import (
    "errors"
)

var ErrNotFound = errors.New("resource not found")

type Set struct {
	ID int `json:"id"`
	Movement string `json:"movement"`
	Volume float64 `json:"volume"`
	Intensity float64 `json:"intensity"`
}

// Example data, to be replaced with db
var sets = []*Set{
	{
		ID: 1,
		Movement: "Squat",
		Volume: 5,
		Intensity: 80,
	},
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

func UpdateSet(s *Set) error {
    for i, v := range sets {
        if s.ID == v.ID {
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
