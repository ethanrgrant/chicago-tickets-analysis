package main

import (
	"database/sql"
	log "github.com/sirupsen/logrus"

	_ "github.com/mattn/go-sqlite3"
)

type dbAccessor struct {
	*sql.DB
}

const (
	addTicket    = "INSERT INTO ticket (ticketNum, zipcode, officer, issueDate, violationCode, fineAmt) VALUES (?, ?, ?, ?, ?, ?)"
	addViolation = "INSERT INTO violation (violationCode, violationDescription) VALUES (?, ?)"
)

func newDBAccessor(dbName string) (*dbAccessor, error) {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Error(err)
		panic("Could not open DB")
	}
	return &dbAccessor{db}, nil
}

func (d *dbAccessor) addTicket(ticket ticket) error {
	// todo shouldn't ignore this if it isn't a unique key error
	d.addViolation(ticket.violationCode, ticket.violationDescription)
	_, err := d.Exec(addTicket,
		ticket.ticketNumber,
		ticket.zipcode,
		ticket.officer,
		ticket.issueDate,
		ticket.violationCode,
		ticket.fineAmt)
	if err != nil {
		log.WithError(err).Error("Could not insert ticket")
		return err
	}
	return nil
}

func (d *dbAccessor) addViolation(code string, desc string) error {
	_, err := d.Exec(addViolation, code, desc)
	if err != nil {
		log.WithError(err).Error("Could not add violation")
		return err
	}
	return nil
}
