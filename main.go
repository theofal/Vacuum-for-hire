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

var (
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
	Logger.Debug("Parsing date.", zap.String("Date", date))
	for _, v := range date {
		if unicode.IsDigit(v) {
			amount += string(v)
		}
	}
	intAmount, _ := strconv.Atoi(amount)
	timeMinusMinutes := time.Now().Add(-time.Minute * time.Duration(intAmount)).Format("02/01/2006 15:04")
	timeMinusHours := time.Now().Add(-time.Hour * time.Duration(intAmount)).Format("02/01/2006 15:04")
	timeMinusDays := time.Now().AddDate(0, 0, -2).Format("02/01/2006 15:04")
	switch {
	case date == "PostedPubliée à l'instant" || date == "PostedAujourd'hui":
		return time.Now().Format("02/01/2006 15:04")
	case strings.Contains(strings.ToLower(date), "minute"):
		return timeMinusMinutes
	case strings.Contains(strings.ToLower(date), "heure"):
		return timeMinusHours
	case strings.Contains(strings.ToLower(date), "jour"):
		return timeMinusDays
	}
	return fmt.Sprintf("Couldn't parse time \"%v\".", date)
}

func main() {

	//Logger initialisation
	TermToSearch = "Golang"
	Logger = InitLogger()
	defer func(Logger *zap.Logger) {
		_ = Logger.Sync()
	}(Logger)

	//DB initialisation
	db, sqlDb := GetDbFile()
	defer func(sqlDb *sql.DB) {
		err := sqlDb.Close()
		if err != nil {
			Logger.Error("Error while closing sql database.", zap.Error(err))
		}
	}(sqlDb)

	/*	//Selenium instantiation + google search
		allJobs, err := Webdriver().SearchGoogle(TermToSearch)
		// TODO : If err == ErrTimedOut -> flush ? puis relancer le code
		if err != nil {
			os.Exit(1)
		}

		//Data insertion in database
		err = db.InsertDataInTable(allJobs)
		if err != nil {
			Logger.Error("Error while inserting data in table.", zap.Error(err))
		}*/

	//Retrieving data from DB
	listOfJobs, err := db.GetDataSinceSpecificID(4)
	if err != nil {
		fmt.Println(err)
	}

	// Transforming Post struct into array of array
	allJobs := make([]interface{}, 0)
	for i := 0; i < len(listOfJobs); i++ {
		arrayOfJobs := ParseStructToArray(listOfJobs[i])
		allJobs = append(allJobs, arrayOfJobs)
	}
	fmt.Println(allJobs)

	//CSV data insertion
}

// API ?
//1er cron qui lance le job et qui declenche un webhook
//webhook pour executer un cron qui va récupérer les données de la db et la mettre en ligne
