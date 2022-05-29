package data

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const (
	// TODO: improve typing, add CreatedOn and LastUpdatedOn datetimes
	createSetTable = `
	CREATE TABLE IF NOT EXISTS sets (
		id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		movement TEXT,
		volume FLOAT,
		intensity FLOAT
	);
	`
	insertSet     = `INSERT OR IGNORE INTO sets(movement, volume, intensity) VALUES (?, ?, ?);`
	selectSetByID = `SELECT movement, volume, intensity FROM sets WHERE id=?;`
	selectAllSets = `SELECT * FROM sets;`
	updateSetByID = `UPDATE sets SET movement=?, volume=?, intensity=? WHERE id=?;`
	deleteSetByID = `DELETE FROM sets WHERE id=?;`
)

// SetDB defines the interface for accessing/manipulating set data
type SetDB interface {
	AddSet(s *Set) (SetID, error)
	Sets() ([]*Set, error)
	SetByID(id SetID) (*Set, error)
	UpdateSet(id SetID, s *Set) error
	DeleteSet(id SetID) error
	Close() error
}

// setDB contains a handle to the underlying sql database, and implements SetDB
type setDB struct {
	handle *sql.DB
}

// NewSetDB loads the data from the given file and returns a SetDB interface or error
func NewSetDB(filename string) (SetDB, error) {
	return newSetDB(filename)
}

// newSetDB returns the underlying setDB and error created from the given filename.
func newSetDB(filename string) (*setDB, error) {
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

	_, err = db.Exec(createSetTable)
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
func (sd *setDB) AddSet(s *Set) (SetID, error) {
	result, err := sd.handle.Exec(insertSet, s.Movement, s.Volume, s.Intensity)
	if err != nil {
		return -1, fmt.Errorf("encountered error executing SQL statement: %v", err)
	}

	setID, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("encountered error retrieving last inserted id: %v", err)
	}

	return SetID(setID), nil
}

// Sets implements the SetDB interface method for retrieving all sets from the database.
// An empty slice of sets is considered a valid result of the database query.
func (sd *setDB) Sets() ([]*Set, error) {
	rows, err := sd.handle.Query(selectAllSets)
	if err != nil {
		return nil, fmt.Errorf("error executing SQL query: %v", err)
	}
	defer rows.Close()

	sets := make([]*Set, 0)
	for rows.Next() {
		var id SetID
		var movement string
		var volume float64
		var intensity float64
		if err := rows.Scan(&id, &movement, &volume, &intensity); err != nil {
			return nil, fmt.Errorf("encountered error scanning row: %v", err)
		}

		sets = append(sets, &Set{
			ID:        id,
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

// SetByID implements the SetDB interface method for finding a particular set in the database.
// Returns a set matching the given ID in the database.
// If no set with the given id is found, returns ErrNotFound.
func (sd *setDB) SetByID(id SetID) (*Set, error) {
	var movement string
	var volume float64
	var intensity float64
	err := sd.handle.QueryRow(selectSetByID, id).Scan(&movement, &volume, &intensity)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error executing SQL query: %v", err)
	}

	return &Set{
		ID:        id,
		Movement:  movement,
		Volume:    volume,
		Intensity: intensity,
	}, nil
}

// UpdateSet implements the SetDB interface method for updating a particular set in the database.
// Updates the columns of the set matching the given id with the fields of the given set.
// If no set with the given id is found, returns ErrNotFound.
func (sd *setDB) UpdateSet(id SetID, s *Set) error {
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
func (sd *setDB) DeleteSet(id SetID) error {
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