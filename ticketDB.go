package main

import (
	"database/sql"
	"errors"
	log "github.com/sirupsen/logrus"

	sq "github.com/mattn/go-sqlite3"
)

type dbAccessor struct {
	*sql.DB
}

const (
	addTicket = `
	INSERT INTO ticket (ticketNum, zipcode, officer, issueDate, violationCode, fineAmt)
	VALUES (?, ?, ?, ?, ?, ?)`
	addViolation = `
	INSERT INTO violation (violationCode, violationDescription) 
	VALUES (?, ?)`
	getZips = `
	SELECT zipcode, COUNT(1)
	FROM ticket 
	GROUP BY zipcode`
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
	// don't bother trying to add empty tickets
	if ticket.ticketNumber == 0 {
		return errors.New("Empty tickets are not accepted")
	}
	d.addViolation(ticket.violationCode, ticket.violationDescription)
	_, err := d.Exec(addTicket,
		ticket.ticketNumber,
		ticket.zipcode,
		ticket.officer,
		ticket.issueDate,
		ticket.violationCode,
		ticket.fineAmt)
	if err != nil {
		return err
	}
	return nil
}

func (d *dbAccessor) addViolation(code string, desc string) error {
	_, err := d.Exec(addViolation, code, desc)
	if err == sq.ErrConstraintUnique { // https://godoc.org/github.com/mattn/go-sqlite3#pkg-files
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

func (d *dbAccessor) getZipcodeMap() (map[int]int, error) {
	rows, err := d.Query(getZips)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	zipTicketNum := make(map[int]int)
	for rows.Next() {
		var zip, count int
		err = rows.Scan(&zip, &count)
		if err != nil {
			log.WithError(err).Error("diregarding a zip")
			continue
		}
		zipTicketNum[zip] = count
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return zipTicketNum, nil
}
