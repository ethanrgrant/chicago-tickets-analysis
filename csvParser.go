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

type addTicketer interface {
	addTicket(ticket) error
}

func doParse(pathToData string, addTicketer addTicketer) error {
	file, err := os.Open(pathToData)
	if err != nil {
		log.WithError(err).Error("Failed to find data. Path=%v", pathToData)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// discard name of columns
	_ = scanner.Text()

	lineChan := make(chan string, 0)

	// generate pool of line parsers
	workerCount := 25
	for i := 0; i < workerCount; i++ {
		go parseLine(lineChan, addTicketer)
	}

	// go through remaining rows
	totalLines := 0
	for scanner.Scan() {
		lineChan <- scanner.Text()
		totalLines += 1
	}
	log.WithField("total tickets", totalLines).Info("Read csv")
	close(lineChan)

	// send ticket to db to be processed
	for i := 0; i < totalLines; i++ {
	}
	return nil
}

var (
	goodColumns = map[int]string{
		0:  "ticket_number",
		1:  "issue_date",
		2:  "violation_location",
		6:  "zipcode",
		7:  "violation_code",
		8:  "violation_description",
		12: "fine_level1_amount",
		21: "officer",
	}
)

func parseLine(lines chan string, addTicketer addTicketer) error {
	for line := range lines {
		columns := strings.Split(line, ",")
		tic := &ticket{}
		for i, val := range columns {
			if columnName, ok := goodColumns[i]; ok {
				err := tic.addValue(columnName, val)
				if err != nil {
					tic = &ticket{}
					break
				}
			}
		}
		err := addTicketer.addTicket(*tic)
		if err != nil {
			log.WithError(err).Error("Could not add ticket to db, ignoring")
		}
	}
	return nil
}
