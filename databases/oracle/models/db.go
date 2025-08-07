package models

import (
	"database/sql"
	"log"
	"strings"
)

type OracleConfig struct {
	DSN string
}

var OracleConfigs = make(map[string]OracleConfig)

func GetOracleConnection(segment_no string) (*sql.DB, error) {
	// the config will always be available because it is set in main
	config := OracleConfigs[strings.ToUpper(segment_no)]
	db, err := sql.Open("oracle", config.DSN)
	if err != nil {
		log.Printf("Error opening Oracle connection: %v", err)

		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Printf("Error pinging Oracle database: %v", err)
		return nil, err
	}
	return db, nil
}
