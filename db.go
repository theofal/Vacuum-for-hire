package main

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"os"
)

//Database is an instance of the database.
type Database struct {
	DB *sql.DB
}

var db Database

// GetDbFile checks if a 'Vacuum-database.db' file is present on the project
// and instantiates a new seed.
func GetDbFile() (*Database, *sql.DB) {
	var emptyDb bool

	_, err := os.Stat("vacuum-database.db")
	if err != nil {
		emptyDb = true
		Logger.Info("No existing database found, creating a new one.")
		file, err := os.Create("Vacuum-database.db") // Create SQLite file
		if err != nil {
			Logger.Fatal("Couldn't create database file.", zap.Error(err))
		}
		err = file.Close()
		if err != nil {
			Logger.Fatal("Couldn't close database file.", zap.Error(err))
		}
	}

	Logger.Info("Opening database.")
	sqliteDatabase, _ := sql.Open("sqlite3", "./vacuum-database.db") // Open the created SQLite File
	Logger.Info("Database opened.")

	db.DB = sqliteDatabase

	if emptyDb {
		db.CreateTable()
	}

	err = db.IsSeeded()
	if err != nil {
		if errors.Is(ErrSeedNotFound, err) {
			db.CreateTable()
		} else {
			Logger.Error("Error while verifying db seed.", zap.Error(err))
		}
	}

	return &db, sqliteDatabase
}

// CreateTable creates a new table in a given database.
func (db Database) CreateTable() *error {
	createJobTableSQL := `CREATE TABLE JobList (
    	"ID" INTEGER PRIMARY KEY AUTOINCREMENT,
		"SearchedTerm" TEXT,
		"JobTitle" TEXT,
		"CompanyName" TEXT,
		"CompanyLocation" TEXT,
		"JobSnippet" TEXT,
		"Date" TEXT,
		"Url" TEXT
	  );`
	Logger.Info("Creating database.")

	statement, err := db.DB.Prepare(createJobTableSQL) // Prepare SQL Statement
	if err != nil {
		Logger.Fatal("Error while preparing the SQL statement.", zap.Error(err))
	}
	defer func(statement *sql.Stmt) {
		err := statement.Close()
		if err != nil {
			Logger.Error("Error while closing the SQL statement.", zap.Error(err))
		}
	}(statement)

	_, err = statement.Exec() // Execute SQL Statements
	if err != nil {
		Logger.Fatal("Error while executing SQL statement.", zap.Error(err))
	}

	Logger.Info("Table created.")

	return &err
}

// InsertDataInTable inserts data in a given database table.
func (db Database) InsertDataInTable(jobList []Post) error {
	Logger.Info("Inserting jobs in database.")
	var err error
	for i := range jobList {
		insertJobSQL := `INSERT INTO JobList(SearchedTerm, JobTitle, CompanyName, CompanyLocation, JobSnippet, Date, Url) VALUES (?, ?, ?, ?, ?, ?, ?)`
		statement, err := db.DB.Prepare(insertJobSQL) // Prepare statement.
		// This is good to avoid SQL injections
		if err != nil {
			Logger.Fatal("Error while preparing the SQL statement.", zap.Error(err))
			return err
		}
		_, err = statement.Exec(TermToSearch, jobList[i].JobTitle, jobList[i].CompanyName, jobList[i].CompanyLocation, jobList[i].JobSnippet, jobList[i].Date, jobList[i].URL)
		if err != nil {
			Logger.Fatal("Error while executing SQL statement.", zap.Error(err))
			return err
		}
	}
	Logger.Info("Jobs inserted in database.")
	return err
}

// GetDataSinceSpecificID retrieves data (posterior to an input date) from a given database table.
func (db Database) GetDataSinceSpecificID(ID int) ([]Post, error) {
	var allJobs []Post

	row, err := db.DB.Query("SELECT * FROM JobList WHERE ROWID > ?", ID)
	if err != nil {
		Logger.Error("Error while querying the database.", zap.Error(err))
		return nil, err
	}
	defer func(row *sql.Rows) *error {
		err := row.Close()
		if err != nil {
			Logger.Error("Error while closing sql query.", zap.Error(err))
			return &err
		}
		return &err
	}(row)

	for row.Next() { // Iterate and fetch the records from result cursor
		var ID string
		var SearchedTerm string
		var JobTitle string
		var CompanyName string
		var CompanyLocation string
		var JobSnippet string
		var Date string
		var URL string
		err := row.Scan(&ID, &SearchedTerm, &JobTitle, &CompanyName, &CompanyLocation, &JobSnippet, &Date, &URL)
		if err != nil {
			return nil, err
		}
		allJobs = append(allJobs,
			Post{
				ID:              ID,
				JobTitle:        JobTitle,
				Date:            Date,
				CompanyName:     CompanyName,
				CompanyLocation: CompanyLocation,
				URL:             URL,
			})
		//Logger.Debug("Jobs: " + SearchedTerm + JobTitle + " " + CompanyName + " " + CompanyLocation + " " + JobSnippet + " " + Date + " " + URL)
	}

	return allJobs, err
}

// IsSeeded checks if a given database is seeded or empty.
func (db Database) IsSeeded() error {
	row, err := db.DB.Query("SELECT * FROM JobList")
	defer func(row *sql.Rows) {
		err := row.Close()
		if err != nil {
			Logger.Error("Error while closing the SQL statement.", zap.Error(err))
		}
	}(row)
	if err != nil {
		Logger.Warn("Error while querying the database.", zap.Error(err))
		return ErrSeedNotFound
	}
	return nil
}
