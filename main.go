package main

import (
	"database/sql"
	"fmt"
	"github.com/braintree/manners"
	"github.com/joho/godotenv"
	_ "github.com/robfig/cron/v3"
	"github.com/theofal/Vacuum-for-hire/services/app"
	"go.uber.org/zap"
	"os"
	"sync"
)

// getDotEnvVar returns a specific variable in the .env file.
func getDotEnvVar(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		app.Logger.Error("Couldn't find .env find", zap.Error(err))
	}
	return os.Getenv(key)
}

func main() {

	//Logger initialisation
	if os.Args[1] == "" {
		app.TermToSearch = "Golang"
	}
	app.Logger = app.InitLogger()
	defer func(Logger *zap.Logger) {
		_ = Logger.Sync()
	}(app.Logger)

	//DB initialisation.
	db, sqlDb := app.GetDbFile()
	defer func(sqlDb *sql.DB) {
		err := sqlDb.Close()
		if err != nil {
			app.Logger.Error("Error while closing sql database.", zap.Error(err))
		}
	}(sqlDb)

	//Selenium webdriver instantiation + google search.
	allJobs, err := app.Webdriver().SearchGoogle(app.TermToSearch)
	// TODO : If err == ErrTimedOut -> flush ? puis relancer le code
	if err != nil {
		os.Exit(1)
	}

	//Data insertion in database.
	err = db.InsertDataInTable(allJobs)
	if err != nil {
		app.Logger.Error("Error while inserting data in table.", zap.Error(err))
	}

	//Csv file creation/retrieve.
	csvFile := app.GetCsvFile()

	//API server implementation and fetching data
	var wg sync.WaitGroup
	c := make(chan []app.Post)
	go app.InitAPIServer(c, csvFile.GetLastImportID())
	wg.Add(1)
	listOfJobs := <-c
	wg.Done()
	fmt.Println(listOfJobs)
	app.Logger.Info("Job done, closing router.")
	manners.Close()
	close(c)

	//Way to fetch data directly from db - Not used if the API is used.
	/*listOfJobs, err := db.GetDataSinceSpecificID(csvFile.getLastImportID())
	if err != nil {
		fmt.Println(err)
	}*/

	//Transforming Post struct into string array of arrays.
	allTheJobs := make([][]string, 0)
	for i := 0; i < len(listOfJobs); i++ {
		arrayOfJobs := app.ParseStructToArray(listOfJobs[i])
		allTheJobs = append(allTheJobs, arrayOfJobs)
	}

	//Uploading data that are not present in the csv file.
	err = csvFile.ImportMissingData(allTheJobs)
	if err != nil {
		app.Logger.Error("Error while importing data", zap.Error(err))
	}
}
