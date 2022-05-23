package main

import (
	"log"

	"github.com/hrand1005/training-notebook/data"
)

func DBFromConfig(dbConf DBConfig) (data.SetDB, error) {
	// return TestSetData if using test config
	if dbConf.IsTest {
		log.Printf("Using test data\nDBConfig: %+v\n", dbConf)
		return data.TestSetData, nil
	}

	return data.NewSetDB(dbConf.Path)
}
