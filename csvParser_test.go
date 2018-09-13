package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type CSVTestSuite struct {
	suite.Suite
}

func (c *CSVTestSuite) Test_GenerateTicket_SupplyValidInfo_TicketCorrect() {
	testTicket := &ticket{}
	ticketNum := "51551278"
	officerNum := "15227"
	violationLocation := "6014 W 64TH ST"
	fineAmt := "50"
	violationDescription := "EXPIRED PLATES OR TEMPORARY REGISTRATION"
	zipcode := "60638"
	violationCode := "0976160F"
	ticketDate := "2007-01-01 00:03:00"
	correctDate := time.Date(2007, time.January, 1, 0, 3, 0, 0, time.UTC)
	err := testTicket.addValue("ticket_number", ticketNum)
	c.NoError(err)
	err = testTicket.addValue("issue_date", ticketDate)
	c.NoError(err)
	err = testTicket.addValue("officer", officerNum)
	c.NoError(err)
	err = testTicket.addValue("violation_location", violationLocation)
	c.NoError(err)
	err = testTicket.addValue("fine_level1_amount", fineAmt)
	c.NoError(err)
	err = testTicket.addValue("violation_description", violationDescription)
	c.NoError(err)
	err = testTicket.addValue("zipcode", zipcode)
	c.NoError(err)
	err = testTicket.addValue("violation_code", violationCode)
	c.NoError(err)
	c.Equal(51551278, testTicket.ticketNumber)
	c.Equal(15227, testTicket.officer)
	c.Equal(float64(50), testTicket.fineAmt)
	c.Equal(violationLocation, testTicket.violationLocation)
	c.Equal(violationDescription, testTicket.violationDescription)
	c.Equal(60638, testTicket.zipcode)
	c.Equal(violationCode, testTicket.violationCode)
	c.Equal(correctDate, testTicket.issueDate)
}

func (c *CSVTestSuite) Test_GenerateTicket_WrongDataType_Error() {
	ticket := &ticket{}
	err := ticket.addValue("ticket_number", "this isn't a number")
	c.Error(err)
}

func (c *CSVTestSuite) Test_GenerateTicket_ColumnNameIncoorect_Error() {
	ticket := &ticket{}
	err := ticket.addValue("not_a_columnn", "0")
	c.Error(err)
}

func (c *CSVTestSuite) Test_ParseFullLine_CreatesFullTicket() {
	// first line from file
	line := `51551278,2007-01-01 00:00:00,6014 W 64TH ST,`
	line += `90ad622c3274c9bdc9d8c812b79a01d0aaf7479f2bd7431f8935baa4048d0c86,`
	line += `IL,PAS,60638,0976160F,EXPIRED PLATES OR TEMPORARY REGISTRATION,8,CPD,CHEV,50,`
	line += `100,0,100,Paid,2007-05-21 00:00:00,SEIZ,"",5048648030,15227,"6000 w 64th st,"chicago, il`

	// go through remaining rows
	ticketChan := make(chan *ticket, 0)
	lineChan := make(chan string, 0)

	// generate pool of line parsers
	workerCount := 1
	for i := 0; i < workerCount; i++ {
		go parseLine(lineChan, ticketChan)
	}

	lineChan <- line
	close(lineChan)
	testTicket := <-ticketChan
	goodTicket := ticket{
		ticketNumber:         51551278,
		zipcode:              60638,
		officer:              15227,
		issueDate:            time.Date(2007, time.January, 1, 0, 0, 0, 0, time.UTC),
		violationDescription: "EXPIRED PLATES OR TEMPORARY REGISTRATION",
		violationCode:        "0976160F",
		violationLocation:    "6014 W 64TH ST",
		fineAmt:              float64(50),
	}
	c.Equal(goodTicket, *testTicket)
}

func TestCSVTestSuite(t *testing.T) {
	suite.Run(t, new(CSVTestSuite))
}
