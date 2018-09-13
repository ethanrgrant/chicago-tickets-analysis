package main

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type ticket struct {
	ticketNumber         int
	zipcode              int
	officer              int
	issueDate            time.Time
	violationLocation    string
	violationCode        string
	violationDescription string
	fineAmt              float64
}

func (t *ticket) addValue(name string, value string) error {
	switch name {
	case "ticket_number":
		i, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		t.ticketNumber = i
	case "zipcode":
		i, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		t.zipcode = i
	case "officer":
		i, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		t.officer = i
	case "issue_date":
		// reference date: 2007-01-01 00:03:00
		layout := "2006-01-02 15:04:05"
		parsedTime, err := time.Parse(layout, value)
		if err != nil {
			return err
		}
		t.issueDate = parsedTime
	case "violation_description":
		t.violationDescription = value
	case "violation_location":
		t.violationLocation = value
	case "violation_code":
		t.violationCode = value
	case "fine_level1_amount":
		f, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return err
		}
		t.fineAmt = f
	default:
		return errors.New("Could not identify the string as a type")
	}
	return nil
}

func doParse(pathToData string) error {
	file, err := os.Open(pathToData)
	if err != nil {
		log.WithError(err).Error("Failed to find data. Path=%v", pathToData)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	db, err := newDBAccessor("Parking.db")
	if err != nil {
		log.WithError(err).Error("Failed to open db, no point in trying to parse")
		return err
	}

	goodColumns := map[string]bool{
		"ticket_number":         true,
		"issue_date":            true,
		"violation_location":    true,
		"zipcode":               true,
		"violation_code":        true,
		"violation_description": true,
		"fine_level1_amount":    true,
		"officer":               true,
	}
	_ = goodColumns
	columnTitles := strings.Split(scanner.Text(), ",")
	for scanner.Scan() {
		err = parseLine(scanner.Text(), columnTitles, db)
		if err != nil {
			log.WithError(err).Error("Failed to parse row!")
		}
	}
	return nil
}

type ticketAdder interface {
	addTicket(t ticket) error
}

func parseLine(input string, columnTitles []string, ticketAdder ticketAdder) error {
	columns := strings.Split(input, ",")
	var ticket ticket
	for _, val := range columns {
		_ = val
		continue
	}
	err := ticketAdder.addTicket(ticket)
	if err != nil {
		log.WithError(err).Error("Failed to add ticket!")
		return err
	}
	return nil
}
