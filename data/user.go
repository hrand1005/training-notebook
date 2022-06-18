package data

import (
	"database/sql"
	"fmt"

	"github.com/hrand1005/training-notebook/models"
)

const (
	InvalidUserID   models.UserID = -1
	createUserTable               = `
	CREATE TABLE IF NOT EXISTS users (
		id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		password TEXT
	);
	`
	insertUser     = `INSERT OR IGNORE INTO users(name, password) VALUES (?, ?);`
	selectUserByID = `SELECT name, password FROM users WHERE id=?;`
	selectAllUsers = `SELECT * FROM users;`
	updateUserByID = `UPDATE users SET name=?, password=? WHERE id=?;`
	deleteUserByID = `DELETE FROM users WHERE id=?;`
)

type UserDB interface {
	AddUser(*models.User) (models.UserID, error)
	Users() ([]*models.User, error)
	UserByID(id models.UserID) (*models.User, error)
	UpdateUser(models.UserID, *models.User) error
	DeleteUser(id models.UserID) error
	Close() error
}

type userDB struct {
	handle *sql.DB
}

// NewUserDB loads the data from the given file and returns a UserDB interface or error
func NewUserDB(handle *sql.DB) (UserDB, error) {
	return newUserDB(handle)
}

// newUserDB returns the underlying userDB and error created from the given filename.
func newUserDB(db *sql.DB) (*userDB, error) {
	_, err := db.Exec(createUserTable)
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
func (ud *userDB) AddUser(u *models.User) (models.UserID, error) {
	result, err := ud.handle.Exec(insertUser, u.Name, u.Password)
	if err != nil {
		return InvalidUserID, fmt.Errorf("encountered error executing SQL statement: %v", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return InvalidUserID, fmt.Errorf("encountered error retrieving last inserted id: %v", err)
	}

	return models.UserID(userID), nil
}

// Users implements the UserDB interface method for retrieving all users from the database.
// An empty slice of users is considered a valid result of the database query.
func (ud *userDB) Users() ([]*models.User, error) {
	rows, err := ud.handle.Query(selectAllUsers)
	if err != nil {
		return nil, fmt.Errorf("error executing SQL query: %v", err)
	}
	defer rows.Close()

	users := make([]*models.User, 0)
	for rows.Next() {
		var id models.UserID
		var name string
		var password string
		if err := rows.Scan(&id, &name, &password); err != nil {
			return nil, fmt.Errorf("encountered error scanning row: %v", err)
		}

		users = append(users, &models.User{
			ID:       id,
			Name:     name,
			Password: password,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("encountered error after scanning rows: %v", err)
	}

	return users, nil
}

// UserByID implements the UserDB interface method for finding a particular user in the database.
// Returns a user matching the given ID in the database.
// If no user with the given id is found, returns ErrNotFound.
func (ud *userDB) UserByID(id models.UserID) (*models.User, error) {
	var name string
	var password string
	err := ud.handle.QueryRow(selectUserByID, id).Scan(&name, &password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error executing SQL query: %v", err)
	}

	return &models.User{
		ID:       id,
		Name:     name,
		Password: password,
	}, nil
}

// UpdateUser implements the UserDB interface method for updating a particular user in the database.
// Updates the columns of the user matching the given id with the fields of the given user.
// If no user with the given id is found, returns ErrNotFound.
func (ud *userDB) UpdateUser(id models.UserID, u *models.User) error {
	result, err := ud.handle.Exec(updateUserByID, u.Name, u.Password, id)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get RowsAffected by update: %v", err)
	}

	if rowsAffected != 1 {
		if rowsAffected == 0 {
			return ErrNotFound
		}
		return fmt.Errorf("unexpected number of affected rows: %v", rowsAffected)
	}

	return nil
}

// DeleteUser implements the UserDB interface method for removing a particular user from the database.
// Deletes the record of the user matching the given id.
// If no user with the given id is found, returns ErrNotFound.
func (sd *userDB) DeleteUser(id models.UserID) error {
	result, err := sd.handle.Exec(deleteUserByID, id)
	if err != nil {
		return fmt.Errorf("error executing SQL statement: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("encountered error checking rows affected: %v", err)
	}
	if rowsAffected != 1 {
		if rowsAffected == 0 {
			return ErrNotFound
		}
		return fmt.Errorf("unexpected number of affected rows: %v", rowsAffected)
	}

	return nil
}

// Close calls close on the underlying sql.DB
func (ud *userDB) Close() error {
	return ud.handle.Close()
}
