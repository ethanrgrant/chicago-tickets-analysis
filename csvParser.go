package main

import (
	"bufio"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

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
	for scanner.Scan() {
		err = parseLine(scanner.Text(), db)
		if err != nil {
			log.WithError(err).Error("Failed to parse row!")
		}
	}
	return nil
}

func parseLine(input string, db *dbAccessor) error {
	columns := strings.Split(input, ",")
	_ = columns
	db.addRow("TODO")
	return nil

}
