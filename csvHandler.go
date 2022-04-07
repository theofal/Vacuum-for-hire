package main

import (
	"encoding/csv"
	"go.uber.org/zap"
	"os"
	"strconv"
	"strings"
)

type CsvFile struct {
	Reader  *csv.Reader
	Writer  *csv.Writer
	Content [][]string
}

// GetCsvFile verifies if the csv file exists on the project and opens it. If none, it creates one.
func GetCsvFile(fileName string) CsvFile {
	//fileName := "Vacuum-for-hire.csv"
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

// IsSeeded Seed le CSV avec la top barre.
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
func (file CsvFile) getIdColumnIndex() (int, error) {
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

//récupérer l'ID du dernier import
func (file CsvFile) getLastImportID() int {
	maxId := 0
	columnIndex, err := file.getIdColumnIndex()
	if err != nil {
		return 0
	}
	for i := 0; i < len(file.Content); i++ {
		postID, _ := strconv.Atoi(file.Content[i][columnIndex])
		if postID > maxId {
			maxId = postID
		}
	}
	return maxId
}

//get json data

//Insérer les données dans le csv
func (file CsvFile) importMissingData(content [][]string) error {
	Logger.Debug("Trying to write data in CSV file. ")
	err := file.Writer.WriteAll(content)
	if err != nil {
		Logger.Error("An error occurred while writing to csv.", zap.Error(err))
		return err
	}
	Logger.Debug("Data imported to csv file.")
	return nil
}
