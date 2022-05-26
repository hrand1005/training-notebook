package data

import (
	"fmt"
	"os"
	"testing"
)

// TestAddSet calls the AddSet method on a db, and checks that the expected
// set gets added to the db, and does not produce an invalid id result.
func TestAddSet(t *testing.T) {
	testCases := []struct {
		name string
		set  *Set
	}{
		{
			name: "Nominal case adds set to db",
			set: &Set{
				Movement:  "Squat",
				Volume:    5,
				Intensity: 80,
			},
		},
		{
			name: "Invalid ID not used in add to db",
			set: &Set{
				ID:        -1,
				Movement:  "Press",
				Volume:    6,
				Intensity: 75,
			},
		},
	}

	for _, v := range testCases {
		sd := setupTestSetDB()
		defer teardownTestSetDB()

		id, err := sd.AddSet(v.set)
		if err != nil {
			t.Fatalf("Encountered unexpected error: %v", err)
		}
		if id <= 0 {
			t.Fatalf("Invalid set id: %v", id)
		}

		setExists, msg := checkSetInDB(sd, id, v.set)
		if !setExists {
			t.Fatalf("Failed db check: %v", msg)
		}
	}
}

const testSetDB = "testSetDB.sqlite"

func setupTestSetDB() *setDB {
	sd, err := newSetDB(testSetDB)
	if err != nil {
		panic("failed to create testSetDB.sqlite")
	}

	return sd
}

func teardownTestSetDB() {
	os.Remove(testSetDB)
}

// checkSetInDB is a testing utility that checks whehter the provided set exists in the
// provided database. If an error occurs, or the set is not found to exist, returns false
// and a description of the error. If a set is found and all checks pass, returns true and
// an empty string.
func checkSetInDB(sd *setDB, id int, s *Set) (bool, string) {
	var movement string
	var volume float64
	var intensity float64
	err := sd.handle.QueryRow(selectSetByID, id).Scan(&movement, &volume, &intensity)
	if err != nil {
		return false, fmt.Sprintf("error querying for set: %v", err)
	}

	if s.Movement != movement {
		return false, fmt.Sprintf("expected movement %q but got %q", s.Movement, movement)
	}
	if s.Volume != volume {
		return false, fmt.Sprintf("expected volume %v but got %v", s.Volume, volume)
	}
	if s.Intensity != intensity {
		return false, fmt.Sprintf("expected intensity %v but got %v", s.Intensity, intensity)
	}

	return true, ""
}
