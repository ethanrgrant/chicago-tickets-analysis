package main

import (
	"flag"
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

	path := flag.String("path", "data/parking_ticktes.csv", "Path to parking ticket data")
	doParse(*path)
}
