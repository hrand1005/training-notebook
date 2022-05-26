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

		teardownTestSetDB(sd)
	}
}

func TestSets(t *testing.T) {
	testCases := []struct {
		name     string
		addSets  []*Set
		wantSets []*Set
	}{
		{
			name: "No added sets returns empty slice",
		},
		{
			name: "ID value before insertion does not affect Sets result",
			addSets: []*Set{
				{
					ID:        -1,
					Movement:  "High jump",
					Volume:    1,
					Intensity: 90,
				},
			},
			wantSets: []*Set{
				{
					Movement:  "High jump",
					Volume:    1,
					Intensity: 90,
				},
			},
		},
		{
			name: "Multiple added sets each appear in returned Sets",
			addSets: []*Set{
				{
					ID:        -1,
					Movement:  "Yeet ball",
					Volume:    10,
					Intensity: 50,
				},
				{
					ID:        -1,
					Movement:  "Farmer's carry",
					Volume:    30,
					Intensity: 20,
				},
				{
					ID:        -1,
					Movement:  "Lateral raise",
					Volume:    12,
					Intensity: 60,
				},
			},
			wantSets: []*Set{
				{
					Movement:  "Yeet ball",
					Volume:    10,
					Intensity: 50,
				},
				{
					Movement:  "Farmer's carry",
					Volume:    30,
					Intensity: 20,
				},
				{
					Movement:  "Lateral raise",
					Volume:    12,
					Intensity: 60,
				},
			},
		},
	}
	for _, v := range testCases {
		sd := setupTestSetDB()

		// add 'addSets' to the database
		for _, s := range v.addSets {
			_, err := sd.AddSet(s)
			if err != nil {
				t.Fatalf("Encountered unexpected error: %v", err)
			}
		}

		// check that each of 'wantSets' is returned
		sets, err := sd.Sets()
		if err != nil {
			t.Fatalf("Encountered unexpected error: %v", err)
		}
		if len(sets) != len(v.wantSets) {
			t.Fatalf("Wanted length %v but got %v\nWant: %+v\nGot: %+v", len(v.wantSets), len(sets), v.wantSets, sets)
		}
		for _, wantSet := range v.wantSets {
			if !containsSet(sets, wantSet) {
				t.Fatalf("Did not find expected set in sets response\nWanted: %+v\nFull Response: %+v", wantSet, sets)
			}
		}

		// check that ids were not assigned invalid values
		for _, s := range sets {
			if s.ID <= 0 {
				t.Fatalf("Invalid set id: %v", s.ID)
			}
		}

		teardownTestSetDB(sd)
	}

}

const testSetDB = "testSetDB.sqlite"

func setupTestSetDB() *setDB {
	sd, err := newSetDB(testSetDB)
	if err != nil {
		msg := fmt.Sprintf("failed to setup test db: %v, err: %v", testSetDB, err)
		panic(msg)
	}

	return sd
}

func teardownTestSetDB(sd *setDB) {
	sd.handle.Close()
	os.Remove(testSetDB)
}

// containsSet checks if the slice of sets contains the provided set.
func containsSet(sets []*Set, s *Set) bool {
	for _, v := range sets {
		if s.Movement == v.Movement && s.Volume == v.Volume && s.Intensity == v.Intensity {
			return true
		}
	}
	return false
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
