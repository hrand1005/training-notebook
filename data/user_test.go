package data

import (
	"fmt"
	"os"
	"testing"

	"github.com/hrand1005/training-notebook/models"
)

// TestAddUser calls the AddUser method on a db, and checks that the expected
// user gets added to the db, and does not produce an invalid id result.
func TestAddUser(t *testing.T) {
	tests := []struct {
		name string
		user *models.User
	}{
		{
			name: "Nominal case adds user to db",
			user: &models.User{
				Name: "Horbus",
			},
		},
		{
			name: "Invalid ID not used in add to db",
			user: &models.User{
				ID:   InvalidUserID,
				Name: "Yorbus",
			},
		},
	}

	for _, v := range tests {
		ud := setupTestUserDB()

		id, err := ud.AddUser(v.user)
		if err != nil {
			t.Fatalf("Encountered unexpected error: %v", err)
		}
		if id <= 0 {
			t.Fatalf("Invalid user id: %v", id)
		}

		userExists, msg := checkUserInDB(ud, id, v.user)
		if !userExists {
			t.Fatalf("Failed db check: %v", msg)
		}

		teardownTestUserDB(ud)
	}
}

// TestUsers calls the Users method on a db, and checks that the expected
// sets have been retrieved. Users are added to the db with AddUser, and are
// checked to have been assigned valid ids (>0)
func TestUsers(t *testing.T) {
	testCases := []struct {
		name      string
		addUsers  []*models.User
		wantUsers []*models.User
	}{
		{
			name: "No added sets returns empty slice",
		},
		{
			name: "ID value before insertion does not affect Users result",
			addUsers: []*models.User{
				{
					ID:   InvalidUserID,
					Name: "Bort",
				},
			},
			wantUsers: []*models.User{
				{
					Name: "Bort",
				},
			},
		},
		{
			name: "Multiple added sets each appear in returned Users",
			addUsers: []*models.User{
				{
					ID:   InvalidUserID,
					Name: "Paul",
				},
				{
					ID:   InvalidUserID,
					Name: "Liam",
				},
				{
					ID:   InvalidUserID,
					Name: "Herb",
				},
			},
			wantUsers: []*models.User{
				{
					Name: "Paul",
				},
				{
					Name: "Liam",
				},
				{
					Name: "Herb",
				},
			},
		},
	}
	for _, v := range testCases {
		sd := setupTestUserDB()

		// add 'addUsers' to the database
		for _, s := range v.addUsers {
			_, err := sd.AddUser(s)
			if err != nil {
				t.Fatalf("Encountered unexpected error: %v", err)
			}
		}

		// check that each of 'wantUsers' is returned
		users, err := sd.Users()
		if err != nil {
			t.Fatalf("Encountered unexpected error: %v", err)
		}
		if len(users) != len(v.wantUsers) {
			t.Fatalf("Wanted length %v but got %v\nWant: %+v\nGot: %+v", len(v.wantUsers), len(users), v.wantUsers, users)
		}
		for _, wantUser := range v.wantUsers {
			if !models.ContainsUser(users, wantUser) {
				t.Fatalf("Did not find expected user in users response\nWanted: %+v\nFull Response: %+v", wantUser, users)
			}
		}

		// check that ids were not assigned invalid values
		for _, s := range users {
			if s.ID <= 0 {
				t.Fatalf("Invalid user id: %v", s.ID)
			}
		}

		teardownTestUserDB(sd)
	}

}

// TestUserByID calls the UserByID method on a db, and checks that the expected
// user is retrieved. A User may be added to the db with AddUser, and is checked
// to have been assigned a valid id (>0)
func TestUserByID(t *testing.T) {
	sd := setupTestUserDB()
	defer teardownTestUserDB(sd)

	wantUser := &models.User{
		Name: "Hubie",
	}
	id, _ := sd.AddUser(wantUser)
	if id <= 0 {
		t.Fatalf("Invalid id assigned: %v", id)
	}
	// test nominal case
	gotUser, gotErr := sd.UserByID(id)
	if gotErr != nil {
		t.Fatalf("Encountered unexpected error in UserByID using id %v\nErr: %v", id, gotErr)
	}
	if !models.UsersEqual(gotUser, wantUser) {
		t.Fatalf("Got User: %+v\nWanted User: %+v\n", gotUser, wantUser)
	}

	// test not found case
	gotUser, gotErr = sd.UserByID(InvalidUserID)
	if gotErr != ErrNotFound {
		t.Fatalf("Expected error not found by UserByID(InvalidUserID)")
	}
	if gotUser != nil {
		t.Fatalf("Expected nil user but got %+v", gotUser)
	}
}

// TestUpdateUser calls userDB's UpdateUser method and checks that the expected
// values are user in the DB, or that the expected error value is returned.
func TestUpdateUser(t *testing.T) {
	testCases := []struct {
		name string
		// validID determines whether to call UpdateUser with a valid ID
		validID        bool
		updateUser     *models.User
		wantUserExists bool
		wantErr        error
	}{
		{
			name:    "Nominal case updates existing user",
			validID: true,
			updateUser: &models.User{
				Name: "Yert",
			},
			wantUserExists: true,
		},
		{
			name:    "Nonexistent ID returns ErrNotFound",
			validID: false,
			updateUser: &models.User{
				Name: "Yorb",
			},
			wantUserExists: false,
			wantErr:        ErrNotFound,
		},
	}
	for _, v := range testCases {
		ud := setupTestUserDB()

		// add an empty user to the db
		id, _ := ud.AddUser(&models.User{})

		// user the id to invalid if we're testing the error case
		if !v.validID {
			id = InvalidUserID
		}
		gotErr := ud.UpdateUser(id, v.updateUser)
		if gotErr != v.wantErr {
			t.Fatalf("Got error: %v\nWanted error: %v\n", gotErr, v.wantErr)
		}

		// check that the result of the update is equal to wantUserExists
		userExists, _ := checkUserInDB(ud, id, v.updateUser)
		if userExists != v.wantUserExists {
			userMsg := fmt.Sprintf("User ID: %v\nUser: %+v", id, v.updateUser)
			t.Fatalf("Got that userExists is %v but wanted userExists to be %v\n%v", userExists, v.wantUserExists, userMsg)
		}

		teardownTestUserDB(ud)
	}
}

func teardownTestUserDB(ud *userDB) {
	ud.handle.Close()
	os.Remove(testUserDB)
}

// checkUserInDB is a testing utility that checks whehter the provided user exists in the
// provided database. If an error occurs, or the user is not found to exist, returns false
// and a description of the error. If a user is found and all checks pass, returns true and
// an empty string.
func checkUserInDB(ud *userDB, id models.UserID, u *models.User) (bool, string) {
	var name string
	err := ud.handle.QueryRow(selectUserByID, id).Scan(&name)
	if err != nil {
		return false, fmt.Sprintf("error querying for user: %v", err)
	}

	if u.Name != name {
		return false, fmt.Sprintf("expected name %q but got %q", u.Name, name)
	}

	return true, ""
}

const testUserDB = "testUserDB.sqlite"

func setupTestUserDB() *userDB {
	ud, err := newUserDB(testUserDB)
	if err != nil {
		msg := fmt.Sprintf("failed to setup test db: %v, err: %v", testUserDB, err)
		panic(msg)
	}

	return ud
}
