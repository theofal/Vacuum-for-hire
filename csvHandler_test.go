package main

import (
	"encoding/csv"
	"os"
	"testing"
)

func TestGetIDColumnIndex(t *testing.T) {
	Logger.Debug("Creating test csv file.")
	file, err := os.Create("test-csv-file.csv")
	if err != nil {
		t.Error(err)
	}
	csvWriter := csv.NewWriter(file)

	idTest := []string{"id"}
	err = csvWriter.Write(idTest)
	if err != nil {
		t.Error(err)
	}

	csvFile := GetCsvFile()

	want := 0
	got, _ := csvFile.getIDColumnIndex()
	if got != want {
		t.Errorf("TestGetIDColumnIndex FAILED : want %v, got %v.\n", want, got)
	}

	err = file.Close()
	if err != nil {
		return
	}
	Logger.Debug("Deleting test csv file.")
	err = os.Remove("test-csv-file.csv")
	if err != nil {
		return
	}
}

func TestIsSeeded(t *testing.T) {
	Logger.Debug("Creating test csv file.")
	file, err := os.Create("test-csv-file.csv")
	if err != nil {
		t.Error(err)
	}
	csvWriter := csv.NewWriter(file)

	idTest := []string{"id"}
	err = csvWriter.Write(idTest)
	if err != nil {
		t.Error(err)
	}

	csvFile := GetCsvFile()

	got := csvFile.IsSeeded()
	if got != true {
		t.Errorf("TestGetIDColumnIndex FAILED : want %v, got %v.\n", true, got)
	}

	err = file.Close()
	if err != nil {
		return
	}
	Logger.Debug("Deleting test csv file.")
	err = os.Remove("test-csv-file.csv")
	if err != nil {
		return
	}
}

func TestGetLastImportID(t *testing.T) {
	Logger.Debug("Creating test csv file.")
	file, err := os.Create("test-csv-file.csv")
	if err != nil {
		t.Error(err)
	}
	csvWriter := csv.NewWriter(file)

	idTest := make([]string, 1)
	idTest2 := make([]string, 1)
	testid := make([][]string, 2)
	idTest[0] = "id"
	idTest2[0] = "1"
	testid[0] = idTest
	testid[1] = idTest2
	err = csvWriter.WriteAll(testid)
	if err != nil {
		t.Error(err)
	}

	csvFile := GetCsvFile()

	want := 1
	got := csvFile.getLastImportID()
	if got != want {
		t.Errorf("TestGetIDColumnIndex FAILED : want %v, got %v.\n", want, got)
	}

	err = file.Close()
	if err != nil {
		return
	}
	Logger.Debug("Deleting test csv file.")
	err = os.Remove("test-csv-file.csv")
	if err != nil {
		return
	}
}
