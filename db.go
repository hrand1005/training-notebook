package main

import (
	"log"

	"github.com/hrand1005/training-notebook/data"
)

func DBFromConfig(dbConf DBConfig) (data.SetDB, error) {
	log.Printf("DBConfig: %+v", dbConf)
	return data.TestSetData, nil
}
