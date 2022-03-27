package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type URL struct {
	Base     string
	Term     string
	Endpoint string
}

type Post struct {
	JobTitle        string
	CompanyName     string
	CompanyLocation string
	JobSnippet      string
	Date            string
	URL             string
}

var (
	AllJobs      []Post
	Logger       *zap.Logger
	TermToSearch string
)

// getDotEnvVar returns a specific variable in the .env file.
func getDotEnvVar(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		Logger.Error("Couldn't find .env find", zap.Error(err))
	}
	return os.Getenv(key)
}

// ParseDate returns a date in a string format to when the job post was uploaded.
func ParseDate(date string) string {
	var amount string
	//Logger.Debug("Parsing date.", zap.String("Date", date))
	for _, v := range date {
		if unicode.IsDigit(v) {
			amount += string(v)
			//Logger.Debug("Working.. ", zap.String("Value of V", string(v)))
		} /* else {
			Logger.Debug("Working.. ", zap.String("Value of V", string(v)))
		}*/
	}
	intAmount, _ := strconv.Atoi(amount)
	timeMinusMinutes := time.Now().Add(-time.Minute * time.Duration(intAmount))
	timeMinusHours := time.Now().Add(-time.Hour * time.Duration(intAmount))
	timeMinusDays := time.Now().AddDate(0, 0, -2)
	switch {
	case date == "PostedPubliée à l'instant" || date == "PostedAujourd'hui":
		return fmt.Sprintf("%d/%d/%d", time.Now().Day(), time.Now().Month(), time.Now().Year())
	case strings.Contains(strings.ToLower(date), "minute"):
		return fmt.Sprintf("%v/%v/%v at %v:%v\n", timeMinusMinutes.Day(), timeMinusMinutes.Month(), timeMinusMinutes.Year(), timeMinusMinutes.Hour(), timeMinusMinutes.Minute())
	case strings.Contains(strings.ToLower(date), "heure"):
		return fmt.Sprintf("%v/%v/%v at %v:%v\n", timeMinusHours.Day(), timeMinusHours.Month(), timeMinusHours.Year(), timeMinusHours.Hour(), timeMinusHours.Minute())
	case strings.Contains(strings.ToLower(date), "jour"):
		return fmt.Sprintf("%v/%v/%v at %v:%v\n", timeMinusDays.Day(), timeMinusDays.Month(), timeMinusDays.Year(), timeMinusDays.Hour(), timeMinusDays.Minute())

	}
	return fmt.Sprintf("Couldn't parse time \"%v\".", date)
}

func main() {
	TermToSearch = "Golang"
	Logger = InitLogger()
	defer func(Logger *zap.Logger) {
		err := Logger.Sync()
		if err != nil {
			Logger.Error("Error while syncing logger.", zap.Error(err))
		}
	}(Logger)

	db, sqlDb := CreateDbFile()
	defer func(sqlDb *sql.DB) {
		err := sqlDb.Close()
		if err != nil {
			Logger.Error("Error while closing sql database.", zap.Error(err))
		}
	}(sqlDb)

	Webdriver().SearchGoogle(TermToSearch)

	db.InsertDataInTable(AllJobs)
}
