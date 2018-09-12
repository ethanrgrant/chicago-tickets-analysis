package main

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
)

type dbAccessor struct {
	*sql.DB
}

func newDBAccessor(dbName string) (*dbAccessor, error) {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Error(err)
		panic("Could not open DB")
	}
	return &dbAccessor{db}, nil
}

func (d *dbAccessor) addRow(test string) {
	// todo
	return
}
