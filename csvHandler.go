package main

import (
	"encoding/csv"
	"go.uber.org/zap"
	"os"
	"strconv"
	"strings"
)

//CsvFile is an instance of a csv file.
type CsvFile struct {
	Reader  *csv.Reader
	Writer  *csv.Writer
	Content [][]string
}

// GetCsvFile verifies if the csv file exists on the project and opens it. If none, it creates one.
func GetCsvFile() CsvFile {
	fileName := "Vacuum-for-hire.csv"
	if os.Getenv("DEV_ENV") == "test" {
		fileName = "test-csv-file.csv"
	}
	_, err := os.Stat(fileName)
	if err != nil {
		Logger.Info("No existing CSV file found, creating a new one.")
		file, err := os.Create(fileName)
		if err != nil {
			Logger.Fatal("An error occurred while creating CSV file.", zap.Error(err))
		}
		err = file.Close()
		if err != nil {
			Logger.Fatal("An error occurred while closing CSV file.", zap.Error(err))
		}
	}

	Logger.Info("Opening csv file.")
	csvFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_RDWR, os.ModeAppend)
	if err != nil {
		Logger.Fatal("An error occurred while opening CSV file.", zap.Error(err))
	}
	csvReader := csv.NewReader(csvFile)
	csvWriter := csv.NewWriter(csvFile)

	records, err := csvReader.ReadAll()
	if err != nil {
		Logger.Fatal("An error occurred while reading csv file.", zap.Error(err))
	}

	return CsvFile{
		Reader:  csvReader,
		Writer:  csvWriter,
		Content: records,
	}
}

// IsSeeded verifies if the csv file contains data and fills it with the top bar.
func (file CsvFile) IsSeeded() bool {
	record := []string{"ID", "JobTitle", "CompanyName", "CompanyLocation", "JobSnippet", "Date", "URL"}
	if len(file.Content) == 0 {
		Logger.Info("Empty csv file detected. Seeding.")
		err := file.Writer.Write(record)
		if err != nil {
			Logger.Error("An error occurred while writing to csv.", zap.Error(err))
			return false
		}
		Logger.Debug("Flushing CSV writer buffered data.")
		file.Writer.Flush()
		if file.Writer.Error() != nil {
			Logger.Error("An error occurred while flushing writer.", zap.Error(file.Writer.Error()))
			return false
		}
	}
	return true
}

// getIdColumnIndex returns the index of the "date" column.
func (file CsvFile) getIDColumnIndex() (int, error) {
	Logger.Debug("Trying to retrieve ID column index.")
	if len(file.Content) > 0 {
		for i := 0; i < len(file.Content[0]); i++ {
			if strings.ToLower(file.Content[0][i]) == "id" {
				Logger.Debug("ID column index found.", zap.String("ID Index", strconv.Itoa(i)))
				return i, nil
			}
		}
	}
	Logger.Debug("Empty csv file. Couldn't get date column index.")
	return 0, ErrEmptyFile
}

//getLastImportID retrieves the most recent ID from the CSV file and returns it.
func (file CsvFile) getLastImportID() int {
	maxID := 0
	file.IsSeeded()
	columnIndex, err := file.getIDColumnIndex()
	if err != nil {
		return 0
	}
	for i := 0; i < len(file.Content); i++ {
		postID, _ := strconv.Atoi(file.Content[i][columnIndex])
		if postID > maxID {
			maxID = postID
		}
	}
	return maxID
}

//get json data

//importMissingData synchronises the DB and the csv file by adding missing data to it.
func (file CsvFile) importMissingData(content [][]string) error {
	Logger.Debug("Trying to write data in CSV file. ")
	err := file.Writer.WriteAll(content)
	if err != nil {
		Logger.Error("An error occurred while writing to csv.", zap.Error(err))
		return err
	}
	Logger.Info("Data imported to csv file.")
	return nil
}
