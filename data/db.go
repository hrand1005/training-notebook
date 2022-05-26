package data

import (
	"database/sql"
	"fmt"
	"log"
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
	AddSet(s *Set) (int, error)
	Sets() ([]*Set, error)
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
func (sd *setDB) AddSet(s *Set) (int, error) {
	result, err := sd.handle.Exec(insertSet, s.Movement, s.Volume, s.Intensity)
	if err != nil {
		return -1, fmt.Errorf("encountered error executing SQL statement: %v", err)
	}

	setID, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("encountered error retrieving last inserted id: %v", err)
	}
	log.Printf("Newly assigned set assigned id: %v\n", setID)

	return int(setID), nil
}

// Sets implements the SetDB interface method for retrieving all sets from the database.
// An empty slice of sets is considered a valid result of the database query.
func (sd *setDB) Sets() ([]*Set, error) {
	// log.Println("In Sets")
	statement, err := sd.handle.Prepare(selectAllSets)
	if err != nil {
		// log.Printf("encountered error on line 86 in db.go: %v", err)
		return nil, fmt.Errorf("couldn't prepare SQL statement:\n%s\nerr: %v", selectAllSets, err)
	}
	defer statement.Close()

	rows, err := statement.Query()
	if err != nil {
		// log.Printf("encountered error on line 91 in db.go: %v", err)
		return nil, fmt.Errorf("error executing SQL query: %v", err)
	}
	defer rows.Close()

	sets := make([]*Set, 0)
	for rows.Next() {
		var id int
		var movement string
		var volume float64
		var intensity float64
		if err := rows.Scan(&id, &movement, &volume, &intensity); err != nil {
			// log.Printf("encountered error on line 103 in db.go: %v", err)
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
		// log.Printf("encountered error on line 116 in db.go: %v", err)
		return nil, fmt.Errorf("encountered error after scanning rows: %v", err)
	}

	// log.Printf("All sets: %+v\n", sets)

	return sets, nil
}

// SetByID implements the SetDB interface method for finding a particular set in the database.
// If no set with the given id is found, returns ErrNotFound.
func (sd *setDB) SetByID(id int) (*Set, error) {
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

func (sd *setDB) UpdateSet(id int, s *Set) error {
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

func (sd *setDB) DeleteSet(id int) error {
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
