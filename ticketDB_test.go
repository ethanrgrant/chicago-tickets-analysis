package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type DBTestSuite struct {
	db *dbAccessor
	suite.Suite
}

func (d *DBTestSuite) SetupSuite() {
	var err error
	// d.db, err = newDBAccessor("testTickets.db")
	d.NoError(err)
}

func (d *DBTestSuite) TearDownTest() {
	// _, err := d.db.Exec("DELETE FROM album")
	var err error
	d.NoError(err, "Failed to clear album table")
}

func (d *DBTestSuite) Test_AddRow_RowExists() {
	return
}

func (d *DBTestSuite) Test_AddMultipleRows_RowsExist() {
	return
}

func TestDBTestSuite(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
}
