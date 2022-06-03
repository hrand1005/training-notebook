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
