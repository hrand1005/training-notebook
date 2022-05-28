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

// TestSetByID calls the SetByID method on a db, and checks that the expected
// set is retrieved. A Set may be added to the db with AddSet, and is checked
// to have been assigned a valid id (>0)
func TestSetByID(t *testing.T) {
	sd := setupTestSetDB()
	defer teardownTestSetDB(sd)

	wantSet := &Set{
		Movement:  "Push ups",
		Volume:    25,
		Intensity: 45,
	}
	id, _ := sd.AddSet(wantSet)
	if id <= 0 {
		t.Fatalf("Invalid id assigned: %v", id)
	}
	// test nominal case
	gotSet, gotErr := sd.SetByID(id)
	if gotErr != nil {
		t.Fatalf("Encountered unexpected error in SetByID using id %v\nErr: %v", id, gotErr)
	}
	if !setsEqual(gotSet, wantSet) {
		t.Fatalf("Got Set: %+v\nWanted Set: %+v\n", gotSet, wantSet)
	}

	// test not found case
	gotSet, gotErr = sd.SetByID(-1)
	if gotErr != ErrNotFound {
		t.Fatalf("Expected error not found by SetByID(-1)")
	}
	if gotSet != nil {
		t.Fatalf("Expected nil set but got %+v", gotSet)
	}
}

// TestUpdateSet calls setDB's UpdateSet method and checks that the expected
// values are set in the DB, or that the expected error value is returned.
func TestUpdateSet(t *testing.T) {
	testCases := []struct {
		name string
		// validID determines whether to call UpdateSet with a valid ID
		validID       bool
		updateSet     *Set
		wantSetExists bool
		wantErr       error
	}{
		{
			name:    "Nominal case updates existing set",
			validID: true,
			updateSet: &Set{
				Movement:  "Hamstring Curl",
				Volume:    20,
				Intensity: 50,
			},
			wantSetExists: true,
		},
		{
			name:    "Nonexistent ID returns ErrNotFound",
			validID: false,
			updateSet: &Set{
				Movement:  "Dummy value",
				Volume:    10,
				Intensity: 10,
			},
			wantSetExists: false,
			wantErr:       ErrNotFound,
		},
	}
	for _, v := range testCases {
		sd := setupTestSetDB()

		// add an empty set to the db
		id, _ := sd.AddSet(&Set{})

		// set the id to invalid if we're testing the error case
		if !v.validID {
			id = -1
		}
		gotErr := sd.UpdateSet(id, v.updateSet)
		if gotErr != v.wantErr {
			t.Fatalf("Got error: %v\nWanted error: %v\n", gotErr, v.wantErr)
		}

		// check that the result of the update is equal to wantSetExists
		setExists, _ := checkSetInDB(sd, id, v.updateSet)
		if setExists != v.wantSetExists {
			setMsg := fmt.Sprintf("Set ID: %v\nSet: %+v", id, v.updateSet)
			t.Fatalf("Got that setExists is %v but wanted setExists to be %v\n%v", setExists, v.wantSetExists, setMsg)
		}

		teardownTestSetDB(sd)
	}
}

// TestSets calls the Sets method on a db, and checks that the expected
// sets have been retrieved. Sets are added to the db with AddSet, and are
// checked to have been assigned valid ids (>0)
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

func TestDeleteSet(t *testing.T) {
	testCases := []struct {
		name    string
		validID bool
		// set to be deleted
		deleteSet *Set
		wantErr   error
	}{
		{
			name:    "Nominal case removes set from db and returns nil error",
			validID: true,
			deleteSet: &Set{
				Movement:  "Tricep pulldowns",
				Volume:    12,
				Intensity: 60,
			},
		},
		{
			name:    "Invalid ID case found case returns ErrNotFound",
			validID: false,
			deleteSet: &Set{
				Movement:  "Tricep Pulldowns",
				Volume:    12,
				Intensity: 60,
			},
			wantErr: ErrNotFound,
		},
	}
	for _, v := range testCases {
		sd := setupTestSetDB()
		id, _ := sd.AddSet(v.deleteSet)

		if !v.validID {
			// use an invalidID to test the error case
			id = -1
		}
		gotErr := sd.DeleteSet(id)
		if gotErr != v.wantErr {
			t.Fatalf("Got error: %v\nWanted error: %v", gotErr, v.wantErr)
		}

		// check that the set is no longer in the database (same for error case)
		setExists, _ := checkSetInDB(sd, id, v.deleteSet)
		if setExists {
			t.Fatalf("Found unexpected set:\nSetID: %v\nSet: %+v", id, v.deleteSet)
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

// setsEqual returns true if all non-id fields of the set are equal, and false otherwise.
func setsEqual(set1, set2 *Set) bool {
	if set1.Movement != set2.Movement || set1.Volume != set2.Volume || set1.Intensity != set2.Intensity {
		return false
	}
	return true
}

// containsSet checks if the slice of sets contains the provided set.
func containsSet(sets []*Set, s *Set) bool {
	for _, v := range sets {
		if setsEqual(s, v) {
			return true
		}
	}
	return false
}

// checkSetInDB is a testing utility that checks whehter the provided set exists in the
// provided database. If an error occurs, or the set is not found to exist, returns false
// and a description of the error. If a set is found and all checks pass, returns true and
// an empty string.
func checkSetInDB(sd *setDB, id SetID, s *Set) (bool, string) {
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
