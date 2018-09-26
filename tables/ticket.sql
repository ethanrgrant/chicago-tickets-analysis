CREATE TABLE ticket (
    ticketNum             integer,
    zipcode               int,
    officer               int,
    issueDate             date,
    violationCode         text,
    fineAmt               real,
    PRIMARY KEY(ticketNum)
);