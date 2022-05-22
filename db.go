package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hrand1005/training-notebook/data"
)

func DBFromConfig(dbConf DBConfig) (data.SetDB, error) {
	// return TestSetData if using test config
	if dbConf.IsTest {
		log.Printf("Using test data\nDBConfig: %+v\n", dbConf)
		return data.TestSetData, nil
	}

	// Create new db file if one doesn't exist
	_, err := os.Stat(dbConf.Path)
	if os.IsNotExist(err) {
		file, err := os.Create(dbConf.Path)
		file.Close()
		if err != nil {
			return nil, fmt.Errorf("Failed to create db %s: %v", dbConf.Path, err)
		}
	}

	return data.NewSetDB(dbConf.Path)
}
