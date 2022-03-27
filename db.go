package main

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"os"
)

type Database struct {
	DB *sql.DB
}

var db Database

//TODO GERER LES ERREURS (ex: if err.is(blabla))

func CreateDbFile() (*Database, *sql.DB) {
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

func (db Database) CreateTable() *error {
	createJobTableSQL := `CREATE TABLE JobList (
		--"ID" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
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
	_, err = statement.Exec() // Execute SQL Statements
	if err != nil {
		Logger.Fatal("Error while executing SQL statement.", zap.Error(err))
	}

	Logger.Info("Table created.")

	return &err
}

func (db Database) InsertDataInTable(jobList []Post) *error {
	Logger.Info("Inserting jobs in database.")
	var err error
	for i := range jobList {
		insertJobSQL := `INSERT INTO JobList(SearchedTerm, JobTitle, CompanyName, CompanyLocation, JobSnippet, Date, Url) VALUES (?, ?, ?, ?, ?, ?, ?)`
		statement, err := db.DB.Prepare(insertJobSQL) // Prepare statement.
		// This is good to avoid SQL injections
		if err != nil {
			Logger.Fatal("Error while preparing the SQL statement.", zap.Error(err))
		}
		_, err = statement.Exec(TermToSearch, jobList[i].JobTitle, jobList[i].CompanyName, jobList[i].CompanyLocation, jobList[i].JobSnippet, jobList[i].Date, jobList[i].Url)
		if err != nil {
			Logger.Fatal("Error while executing SQL statement.", zap.Error(err))
		}
	}
	Logger.Info("Jobs inserted in database.")
	return &err
}

func (db Database) GetDataFromTable() *error {
	row, err := db.DB.Query("SELECT * FROM JobList") //Voir ce que je veux rechercher
	if err != nil {
		Logger.Error("Error while querying the database.", zap.Error(err))
		return &err
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
		var JobTitle string
		var CompanyName string
		var CompanyLocation string
		var JobSnippet string
		var Date string
		var Url string
		err := row.Scan(&JobTitle, &CompanyName, &CompanyLocation, &JobSnippet, &Date, &Url)
		if err != nil {
			return &err
		}
		Logger.Debug("Jobs: " + JobTitle + " " + CompanyName + " " + CompanyLocation + " " + JobSnippet + " " + Date + " " + Url)
	}

	return &err
}

func (db Database) GetTableLength() (int64, error) {
	var tableLength int64
	row, err := db.DB.Query("SELECT * FROM JobList") //Voir ce que je veux rechercher
	if err != nil {
		Logger.Warn("Error while querying the database.", zap.Error(err))
		return 0, ErrSeedNotFound
	}
	defer func(row *sql.Rows) {
		err := row.Close()
		if err != nil {
			Logger.Error("Error while closing DB row query.", zap.Error(err))
		}
	}(row)
	for row.Next() { // Iterate and fetch the records from result cursor
		tableLength += 1
	}
	return tableLength, err
}

func (db Database) IsSeeded() error {
	_, err := db.DB.Query("SELECT * FROM JobList")
	if err != nil {
		Logger.Warn("Error while querying the database.", zap.Error(err))
		return ErrSeedNotFound
	}
	return err
}