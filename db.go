package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"log"
	"os"
)

type Database struct {
	DB *sql.DB
	ID int64
}

var db Database

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

	return &db, sqliteDatabase
}

func (db Database) CreateTable() *error {
	createJobTableSQL := `CREATE TABLE JobList (
		"ID" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
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
	exec, err := statement.Exec() // Execute SQL Statements
	if err != nil {
		Logger.Fatal("Error while executing SQL statement.", zap.Error(err))
	}
	db.ID, err = exec.LastInsertId() // Get last insert ID
	if err != nil {
		Logger.Error("Error while getting LastInsertedId.", zap.Error(err))
	}
	Logger.Info("Table created.")

	return &err
}

func (db Database) InsertDataInTable(jobList []Post) *error {
	Logger.Info("Inserting jobs in database.")
	var err error
	for i := range jobList {
		insertJobSQL := `INSERT INTO JobList(ID, JobTitle, CompanyName, CompanyLocation, JobSnippet, Date, Url) VALUES (?, ?, ?, ?, ?, ?, ?)`
		statement, err := db.DB.Prepare(insertJobSQL) // Prepare statement.
		// This is good to avoid SQL injections
		if err != nil {
			Logger.Fatal("Error while preparing the SQL statement.", zap.Error(err))
		}
		_, err = statement.Exec(db.ID+int64(i), jobList[i].JobTitle, jobList[i].CompanyName, jobList[i].CompanyLocation, jobList[i].JobSnippet, jobList[i].Date, jobList[i].Url)
		if err != nil {
			Logger.Fatal("Error while executing SQL statement.", zap.Error(err))
		}
	}
	return &err
}

func (db Database) GetDataFromTable() {
	row, err := db.DB.Query("SELECT * FROM JobList") //Voir ce que je veux rechercher
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var ID int
		var JobTitle string
		var CompanyName string
		var CompanyLocation string
		var JobSnippet string
		var Date string
		var Url string
		row.Scan(&ID, &JobTitle, &CompanyName, &CompanyLocation, &JobSnippet, &Date, &Url)
		Logger.Debug("Jobs: " + JobTitle + " " + CompanyName + " " + CompanyLocation + " " + JobSnippet + " " + Date + " " + Url)
	}
}
