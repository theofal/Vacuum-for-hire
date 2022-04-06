package main

import (
	"encoding/csv"
	"go.uber.org/zap"
	"os"
)

type CsvFile struct {
	Reader  *csv.Reader
	Content [][]string
}

// GetCsvFile verifies if the csv file exists on the project and opens it. If none, it creates one.
func GetCsvFile() CsvFile {
	fileName := "Vacuum-for-hire.csv"
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

	csvFile, err := os.Open(fileName)
	if err != nil {
		Logger.Fatal("An error occurred while opening CSV file", zap.Error(err))
	}
	csvReader := csv.NewReader(csvFile)

	records, err := csvReader.ReadAll()
	if err != nil {
		Logger.Fatal("An error occurred while reading csv file", zap.Error(err))
	}

	return CsvFile{
		Reader:  csvReader,
		Content: records,
	}
}

// getDateColumnIndex returns the index of the "date" column.
func (file CsvFile) getDateColumnIndex() (int, error) {

	if len(file.Content) < 0 {
		for i := 0; i < len(file.Content[0]); i++ {
			if file.Content[0][i] == "Date" {
				return i, nil
			}
		}
	}
	return 0, ErrEmptyFile
}
