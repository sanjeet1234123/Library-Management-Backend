package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func initDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/postgres")
	if err != nil {
		return nil, err
	}
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS Library (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) UNIQUE
      )
    `

	_, err = db.Exec(createTableQuery)
	if err != nil {
		db.Close()
		return nil, err
	}

	createTableQuery1 := `
	CREATE TABLE IF NOT EXISTS Users (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        email VARCHAR(255) UNIQUE NOT NULL,
        contact_number VARCHAR(255) NOT NULL,
        role VARCHAR(255) NOT NULL,
        lib_id INT NOT NULL,
        FOREIGN KEY (lib_id) REFERENCES Library (id)
      );
    `

	_, err = db.Exec(createTableQuery1)
	if err != nil {
		db.Close()
		return nil, err
	}

	createTableQuery2 := `
	CREATE TABLE IF NOT EXISTS BookInventory (
        ISBN INT PRIMARY KEY,
        LibID INT NOT NULL,
        Title VARCHAR(255) NOT NULL,
        Authors VARCHAR(255) NOT NULL,
        Publisher VARCHAR(255) NOT NULL,
        Version VARCHAR(255) NOT NULL,
        TotalCopies INT NOT NULL,
        AvailableCopies INT NOT NULL,
        FOREIGN KEY (LibID) REFERENCES Library (ID)
      );
    `

	_, err = db.Exec(createTableQuery2)
	if err != nil {
		db.Close()
		return nil, err
	}

	createTableQuery3 := `
	CREATE TABLE IF NOT EXISTS RequestEvents (
        ReqID SERIAL PRIMARY KEY,
        BookID INT NOT NULL,
        ReaderID INT NOT NULL,
        RequestDate date NOT NULL,
        ApprovalDate date,
        ApproverID INT,
        RequestType VARCHAR(255) NOT NULL,
        FOREIGN KEY (BookID) REFERENCES BookInventory (ISBN),
        FOREIGN KEY (ReaderID) REFERENCES Users (ID),
        FOREIGN KEY (ApproverID) REFERENCES Users (ID)
      );

    `

	_, err = db.Exec(createTableQuery3)
	if err != nil {
		db.Close()
		return nil, err
	}

	createTableQuery4 := `
	CREATE TABLE IF NOT EXISTS IssueRegistry (
        IssueID SERIAL PRIMARY KEY,
        ISBN INT NOT NULL,
        ReaderID INT NOT NULL,
        IssueApproverID INT NOT NULL,
        IssueStatus VARCHAR(255) NOT NULL,
        IssueDate date NOT NULL,
        ExpectedReturnDate date NOT NULL,
        ReturnDate date,
        ReturnApproverID INT,
        FOREIGN KEY (ISBN) REFERENCES BookInventory (ISBN),
        FOREIGN KEY (ReaderID) REFERENCES Users (ID),
        FOREIGN KEY (IssueApproverID) REFERENCES Users (ID),
        FOREIGN KEY (ReturnApproverID) REFERENCES Users (ID)
      );

    `

	_, err = db.Exec(createTableQuery4)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
