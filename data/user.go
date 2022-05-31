package data

import (
	"database/sql"
	"fmt"
	"os"
)

// UserID is the unique identifier for a user of the application
type UserID int

// User defines the model of the user resource
type User struct {
	ID   UserID `json:"id"`
	Name string `json:"name"`
}

const (
	InvalidUserID   UserID = -1
	createUserTable        = `
	CREATE TABLE IF NOT EXISTS users (
		id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		name TEXT
	);
	`
	insertUser     = `INSERT OR IGNORE INTO users(name) VALUES (?);`
	selectUserByID = `SELECT name FROM users WHERE id=?;`
	selectAllUsers = `SELECT * FROM users;`
)

type UserDB interface {
	AddUser(*User) (UserID, error)
	Users() ([]*User, error)
	/*
		UserByID(id UserID) (*User, error)
		UpdateUser(id UserID, u *User) error
		DeleteUser(id UserID) error
		Close() error
	*/
}

type userDB struct {
	handle *sql.DB
}

// NewUserDB loads the data from the given file and returns a UserDB interface or error
func NewUserDB(filename string) (UserDB, error) {
	return newUserDB(filename)
}

// newUserDB returns the underlying userDB and error created from the given filename.
func newUserDB(filename string) (*userDB, error) {
	// Create new db file if one doesn't exist
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		file, err := os.Create(filename)
		file.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to create db %s: %v", filename, err)
		}
	}

	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, fmt.Errorf("failed to load %s: %v", filename, err)
	}

	_, err = db.Exec(createUserTable)
	if err != nil {
		return nil, fmt.Errorf("couldn't prepare SQL statement:\n%s\nerr: %v", createUserTable, err)
	}

	return &userDB{
		handle: db,
	}, nil
}

// AddUser implements the UserDB interface method for adding a user to the database.
// User ID is automatically assigned at the time that the user is inserted into the DB.
// Returns the assigned id upon successfully inserting the provided user, and nil error.
// If an error occurs, returns -1 for the id and the error value.
func (ud *userDB) AddUser(u *User) (UserID, error) {
	result, err := ud.handle.Exec(insertUser, u.Name)
	if err != nil {
		return InvalidUserID, fmt.Errorf("encountered error executing SQL statement: %v", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return InvalidUserID, fmt.Errorf("encountered error retrieving last inserted id: %v", err)
	}

	return UserID(userID), nil
}

// Users implements the UserDB interface method for retrieving all users from the database.
// An empty slice of users is considered a valid result of the database query.
func (ud *userDB) Users() ([]*User, error) {
	rows, err := ud.handle.Query(selectAllUsers)
	if err != nil {
		return nil, fmt.Errorf("error executing SQL query: %v", err)
	}
	defer rows.Close()

	users := make([]*User, 0)
	for rows.Next() {
		var id UserID
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, fmt.Errorf("encountered error scanning row: %v", err)
		}

		users = append(users, &User{
			ID:   id,
			Name: name,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("encountered error after scanning rows: %v", err)
	}

	return users, nil
}
