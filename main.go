package main

import (
	"flag"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	log.Info("Beginning ticket data analysis process")

	path := flag.String("path", "", "Path to parking ticket data, if empty parsing is skipped")
	db, err := newDBAccessor("chicago.db")
	flag.Parse()
	if err != nil {
		log.WithError(err).Error("Failed to open db, no point in trying to continue")
		panic(err)
	}
	defer db.Close()

	log.WithField("csv file", *path).Info("Attempting to parse file")
	if *path != "" {
		doParse(*path, db)
	}

	// get zipcode info
	zipMap, err := db.getZipcodeMap()
	for zip, count := range zipMap {
		fmt.Printf("Zip: %v, Count: %v\n", zip, count)
	}
}
