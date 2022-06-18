package data

import (
	"database/sql"
	"fmt"

	"github.com/hrand1005/training-notebook/models"
)

const (
	// InvalidSetID represents the special int value returned as ID under error conditions
	InvalidSetID models.SetID = -1
	// TODO: improve typing, add CreatedOn and LastUpdatedOn datetimes
	createSetTable = `
	CREATE TABLE IF NOT EXISTS sets (
		id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		userid INT,
		movement TEXT,
		volume FLOAT,
		intensity FLOAT
	);
	`
	insertSet              = `INSERT OR IGNORE INTO sets(userid, movement, volume, intensity) VALUES (?, ?, ?, ?);`
	selectSetByID          = `SELECT userid, movement, volume, intensity FROM sets WHERE id=?;`
	selectSetByIDAndUserID = `SELECT movement, volume, intensity FROM sets WHERE id=? and userid=?;`
	selectAllSets          = `SELECT id, userid, movement, volume, intensity FROM sets;`
	selectSetsByUserID     = `SELECT id, movement, volume, intensity FROM sets WHERE userid=?;`
	updateSetByID          = `UPDATE sets SET movement=?, volume=?, intensity=? WHERE id=?;`
	deleteSetByID          = `DELETE FROM sets WHERE id=?;`
)

// SetDB defines the interface for accessing/manipulating set data
type SetDB interface {
	AddSet(s *models.Set) (models.SetID, error)
	Sets() ([]*models.Set, error)
	SetsByUserID(models.UserID) ([]*models.Set, error)
	SetByID(id models.SetID) (*models.Set, error)
	SetByIDForUser(models.SetID, models.UserID) (*models.Set, error)
	UpdateSet(id models.SetID, s *models.Set) error
	DeleteSet(id models.SetID) error
	Close() error
}

// setDB contains a handle to the underlying sql database, and implements SetDB
type setDB struct {
	handle *sql.DB
}

// NewSetDB loads the data from the given sql db handle and returns a SetDB interface or error
func NewSetDB(db *sql.DB) (SetDB, error) {
	return newSetDB(db)
}

// newSetDB returns the underlying setDB and error created from the given sql db handle.
func newSetDB(db *sql.DB) (*setDB, error) {
	_, err := db.Exec(createSetTable)
	if err != nil {
		return nil, fmt.Errorf("couldn't prepare SQL statement:\n%s\nerr: %v", createSetTable, err)
	}

	return &setDB{
		handle: db,
	}, nil
}

// AddSet implements the SetDB interface method for adding a set to the database.
// Set ID is automatically assigned at the time that the set is inserted into the DB.
// Returns the assigned id upon successfully inserting the provided set, and nil error.
// If an error occurs, returns -1 for the id and the error value.
func (sd *setDB) AddSet(s *models.Set) (models.SetID, error) {
	result, err := sd.handle.Exec(insertSet, s.UID, s.Movement, s.Volume, s.Intensity)
	if err != nil {
		return InvalidSetID, fmt.Errorf("encountered error executing SQL statement: %v", err)
	}

	setID, err := result.LastInsertId()
	if err != nil {
		return InvalidSetID, fmt.Errorf("encountered error retrieving last inserted id: %v", err)
	}

	return models.SetID(setID), nil
}

// Sets implements the SetDB interface method for retrieving all sets from the database.
// An empty slice of sets is considered a valid result of the database query.
func (sd *setDB) Sets() ([]*models.Set, error) {
	rows, err := sd.handle.Query(selectAllSets)
	if err != nil {
		return nil, fmt.Errorf("error executing SQL query: %v", err)
	}
	defer rows.Close()

	sets := make([]*models.Set, 0, 10)
	for rows.Next() {
		var id models.SetID
		var userID models.UserID
		var movement string
		var volume float64
		var intensity float64
		if err := rows.Scan(&id, &userID, &movement, &volume, &intensity); err != nil {
			return nil, fmt.Errorf("encountered error scanning row: %v", err)
		}

		sets = append(sets, &models.Set{
			ID:        id,
			UID:       userID,
			Movement:  movement,
			Volume:    volume,
			Intensity: intensity,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("encountered error after scanning rows: %v", err)
	}

	return sets, nil
}

// SetByUserID implements the SetDB interface method for finding a particular set in the database.
// Returns sets with the matching UserID from the database.
// An empty slice of sets is considered a valid result of the database query.
func (sd *setDB) SetsByUserID(userID models.UserID) ([]*models.Set, error) {
	rows, err := sd.handle.Query(selectSetsByUserID, userID)
	if err != nil {
		return nil, fmt.Errorf("error executing SQL query: %v", err)
	}
	defer rows.Close()

	sets := make([]*models.Set, 0, 10)
	for rows.Next() {
		var id models.SetID
		var userID models.UserID
		var movement string
		var volume float64
		var intensity float64
		if err := rows.Scan(&id, &userID, &movement, &volume, &intensity); err != nil {
			return nil, fmt.Errorf("encountered error scanning row: %v", err)
		}

		sets = append(sets, &models.Set{
			ID:        id,
			UID:       userID,
			Movement:  movement,
			Volume:    volume,
			Intensity: intensity,
		})
	}

	return sets, nil
}

// SetByID implements the SetDB interface method for finding a particular set in the database.
// Returns a set matching the given ID in the database.
// If no set with the given id is found, returns ErrNotFound.
func (sd *setDB) SetByID(id models.SetID) (*models.Set, error) {
	var userID int
	var movement string
	var volume float64
	var intensity float64
	err := sd.handle.QueryRow(selectSetByID, id).Scan(&userID, &movement, &volume, &intensity)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error executing SQL query: %v", err)
	}

	return &models.Set{
		ID:        id,
		UID:       models.UserID(userID),
		Movement:  movement,
		Volume:    volume,
		Intensity: intensity,
	}, nil
}

// SetByIDForUser implements the SetDB interface method for finding a particular set in the database.
// Returns a set matching the given ID in the database.
// If no set with the given id is found, returns ErrNotFound.
func (sd *setDB) SetByIDForUser(setID models.SetID, userID models.UserID) (*models.Set, error) {
	var movement string
	var volume float64
	var intensity float64
	err := sd.handle.QueryRow(selectSetByIDAndUserID, setID, userID).Scan(&movement, &volume, &intensity)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error executing SQL query: %v", err)
	}

	return &models.Set{
		ID:        setID,
		UID:       userID,
		Movement:  movement,
		Volume:    volume,
		Intensity: intensity,
	}, nil
}

// UpdateSet implements the SetDB interface method for updating a particular set in the database.
// Updates the columns of the set matching the given id with the fields of the given set.
// If no set with the given id is found, returns ErrNotFound.
func (sd *setDB) UpdateSet(id models.SetID, s *models.Set) error {
	result, err := sd.handle.Exec(updateSetByID, s.Movement, s.Volume, s.Intensity, id)
	if err != nil {
		return fmt.Errorf("failed to update set: %v", err)
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

// DeleteSet implements the SetDB interface method for removing a particular set from the database.
// Deletes the record of the set matching the given id.
// If no set with the given id is found, returns ErrNotFound.
func (sd *setDB) DeleteSet(id models.SetID) error {
	result, err := sd.handle.Exec(deleteSetByID, id)
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
func (sd *setDB) Close() error {
	return sd.handle.Close()
}
