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
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"movement" TEXT,
		"volume" TEXT,
		"intensity" TEXT
	)
	`
	insertSet = `INSERT OR IGNORE INTO sets(movement, volume, intensity) VALUES (?, ?, ?)`
)

// SetDB defines the interface for accessing/manipulating set data
type SetDB interface {
	AddSet(s *Set) error
	Sets() []*Set
	SetByID(id int) (*Set, error)
	UpdateSet(id int, s *Set) error
	DeleteSet(id int) error
}

// setDB contains a handle to the underlying sql database, and implements SetDB
type setDB struct {
	handle *sql.DB
}

// NewSetDB loads the data from the given file and returns a SetDB interface or error
func NewSetDB(filename string) (SetDB, error) {
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

	statement, err := db.Prepare(createSetTable)
	if err != nil {
		return nil, fmt.Errorf("couldn't prepare SQL statement:\n%s\nerr: %v", createSetTable, err)
	}
	statement.Exec()

	return &setDB{
		handle: db,
	}, nil
}

// AddSet implements the SetDB interface method for adding a set to the database.
// Set ID is automatically assigned at the time that the set is inserted into the DB.
func (sd *setDB) AddSet(s *Set) error {
	statement, err := sd.handle.Prepare(insertSet)
	if err != nil {
		return fmt.Errorf("couldn't prepare SQL statement:\n%s\nerr: %v", insertSet, err)
	}
	statement.Exec(s.Movement, s.Volume, s.Intensity)

	return nil
}

func (sd *setDB) Sets() []*Set {
	return []*Set{}
}

func (sd *setDB) SetByID(id int) (*Set, error) {
	return nil, nil
}

func (sd *setDB) UpdateSet(id int, s *Set) error {
	return nil
}

func (sd *setDB) DeleteSet(id int) error {
	return nil
}
