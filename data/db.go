package data

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// LoadSqliteFile creates a sql.DB handle using the sqlite3 driver and the filename.
// If the file doesn't exist, it is created
func SqliteDB(f string) (*sql.DB, error) {
	// Create new db file if one doesn't exist
	_, err := os.Stat(f)
	if os.IsNotExist(err) {
		file, err := os.Create(f)
		file.Close()
		if err != nil {
			return nil, err
		}
	}

	db, err := sql.Open("sqlite3", f)
	if err != nil {
		return nil, err
	}

	return db, nil
}
