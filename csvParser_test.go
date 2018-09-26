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

func (c *CSVTestSuite) Test_LineParse_MultipleLinesAndWorkers_NoEmptyTickets() {
	line1 := `51551278,2007-01-01 00:00:00,6014 W 64TH ST,90ad622c3274c9bdc9d8c812b79a01d0aaf7479f2bd7431f8935baa4048d0c86,IL,PAS,60638,0976160F,EXPIRED PLATES OR TEMPORARY REGISTRATION,8,CPD,CHEV,50,100,0,100,Paid,2007-05-21 00:00:00,SEIZ,"",5048648030,15227,"6000 w 64th st,
	chicago, il"`
	line2 := `51491256,2007-01-01 00:00:00,530 N MICHIGAN,bce4dc26b2c96965380cb2b838cdbb95632b7b5716061255c7ed9aa52b17163c,IL,PAS,606343801,0964150B,PARKING/STANDING PROHIBITED ANYTIME,18,CPD,CHRY,50,100,50,0,Define,2007-01-22 00:00:00,"","",0,18320,"500 n michigan, chicago, il"`
	line3 := `50433524,2007-01-01 00:01:00,4001 N LONG,44641e828f4d894c883c07c566063c2d99d08f2c03b3d41682d6d8201a0939bd,IL,PAS,60148,0976160F,EXPIRED PLATES OR TEMPORARY REGISTRATION,16,CPD,BUIC,50,100,0,50,Paid,2007-01-31 00:00:00,VIOL,"",5079875240,3207,"4000 n long, chicago, il"`
	line4 := `51430906,2007-01-01 00:01:00,303 E WACKER,eee50ca0d9be2debd0e7d45bad05b8674a6cf5b892230f54cf1923e36990ada9,IL,PAS,60601,0964110A,DOUBLE PARKING OR STANDING,152,CPD,NISS,100,200,0,100,Paid,2007-03-08 00:00:00,DETR,Liable,5023379950,19410,"300 e wacker, chicago, il"`
	line5 := `51501733,2007-01-01 00:04:00,2405 W 14TH ST,b27d76408581e0e3940aa2722fa87bd23cd5428be4c46c7c5d1682e10133ee58,IL,PAS,60651,0964110A,DOUBLE PARKING OR STANDING,10,CPD,DODG,100,200,244,0,Bankruptcy,2010-02-19 00:00:00,SEIZ,"",5038039180,08432,"2400 w 14th st, chicago, il"`

	ticketChan := make(chan *ticket, 100)
	lineChan := make(chan string, 0)

	// generate pool of line parsers
	workerCount := 2
	for i := 0; i < workerCount; i++ {
		go parseLine(lineChan, ticketChan)
	}

	lineChan <- line1
	lineChan <- line2
	lineChan <- line3
	lineChan <- line4
	lineChan <- line5
	close(lineChan)
	for i := 0; i < 5; i++ {
		testTicket := <-ticketChan
		c.NotEqual(ticket{}, testTicket)
	}
}

func TestCSVTestSuite(t *testing.T) {
	suite.Run(t, new(CSVTestSuite))
}
