package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type DBTestSuite struct {
	db *dbAccessor
	suite.Suite
}

func (d *DBTestSuite) SetupSuite() {
	var err error
	d.db, err = newDBAccessor("chicago_test.db")
	d.NoError(err)
}

func (d *DBTestSuite) TearDownTest() {
	_, err := d.db.Exec("DELETE FROM violation")
	d.NoError(err, "failed to clear from violation table")
	_, err = d.db.Exec("DELETE FROM ticket")
	d.NoError(err, "Failed to clear ticket table")
}

func (d *DBTestSuite) TearDownSuite() {
	d.db.Close()
}

func (d *DBTestSuite) Test_AddViolation_RowExists() {
	err := d.db.addViolation("code1", "desc1")
	d.NoError(err)
	// would return no rows error if didn't exist
	_, err = d.db.Query("SELECT violationCode, violationDescription FROM violation WHERE violationCode=?", "code1")
	d.NoError(err)
}

func (d *DBTestSuite) Test_AddSameViolation_GetError() {
	err := d.db.addViolation("code1", "desc1")
	d.NoError(err)
	err = d.db.addViolation("code1", "desc1")
	d.Error(err)
}

var (
	goodTicket = ticket{
		ticketNumber:         51551278,
		zipcode:              60638,
		officer:              15227,
		issueDate:            time.Date(2007, time.January, 1, 0, 0, 0, 0, time.UTC),
		violationDescription: "EXPIRED PLATES OR TEMPORARY REGISTRATION",
		violationCode:        "0976160F",
		violationLocation:    "6014 W 64TH ST",
		fineAmt:              float64(50),
	}
)

func (d *DBTestSuite) Test_AddTicket_TicketExists() {
	err := d.db.addTicket(goodTicket)
	d.NoError(err)
	_, err = d.db.Query("SELECT ticketNum from ticket WHERE ticketNum=?", goodTicket.ticketNumber)
	d.NoError(err)
}

func (d *DBTestSuite) Test_AddSameTicket_Error() {
	err := d.db.addTicket(goodTicket)
	d.NoError(err)
	err = d.db.addTicket(goodTicket)
	d.Error(err)
}
func (d *DBTestSuite) Test_AddMultipleRows_RowsExist() {
	return
}

func TestDBTestSuite(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
}
