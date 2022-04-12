package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/braintree/manners"
	"github.com/gin-gonic/gin"
	_ "github.com/go-ping/ping"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

//InitAPIServer initialises a new API server.
func InitAPIServer(c chan []Post, index int) {
	Logger.Info("Initialising API server.")
	router := gin.Default()

	v1 := router.Group("/api/")
	{
		v1.GET("posts/:id", getAllPostsSinceLastID)
	}

	go requestPostsFromAPI(c, index)

	Logger.Info("Starting server.")
	err := manners.ListenAndServe(":8090", router)
	if err != nil {
		Logger.Error("An error occurred while running router", zap.Error(err))
		os.Exit(1)
	}

}

//requestPostsFromAPI goroutine called in InitAPIServer. This function requests data from API server and exports data
//as []Post to main function via a channel.
func requestPostsFromAPI(c chan []Post, index int) {
	var resp *http.Response

	for i := 1; i <= 5; i++ {
		Logger.Debug("Trying to get data from API server.", zap.String("try number", strconv.Itoa(i)))
		resp, _ = http.Get("http://localhost:8090/api/posts/" + strconv.Itoa(index))
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			time.Sleep(1 * time.Second)
			continue
		} else {
			Logger.Debug("Got all the data I needed!?")
			break
		}
	}
	body, _ := ioutil.ReadAll(resp.Body)

	var jobs []Post
	err := json.Unmarshal(body, &jobs)
	if err != nil {
		Logger.Error("An error occurred while unmarshalling json.", zap.Error(err))
		os.Exit(1)
	}
	Logger.Info("Done getting data from API server.", zap.String("Number of jobs", strconv.Itoa(len(jobs))))
	c <- jobs
}

//getAllPostsSinceLastID returns the list of posts where their ID is superior to the given ID.
func getAllPostsSinceLastID(context *gin.Context) {
	id := context.Param("id")
	idInt, _ := strconv.Atoi(id)
	db, err := sql.Open("sqlite3", "./vacuum-database.db")
	if err != nil {
		Logger.Error("An error occurred while opening the database.", zap.Error(err))
		os.Exit(1)
	}
	// TODO : Ameliorer pour éviter la possibilité d'ouvrir la db à deux endroits en meme temps (https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f)
	// TODO : remove the database instantiation via a getter ?
	jobList, err := database{DB: db}.GetDataSinceSpecificID(idInt)
	if err != nil {
		Logger.Error("An error occurred while querying the database.", zap.Error(err))
		os.Exit(1)
	}

	allTheJobs := ParseToJson(jobList)
	context.JSON(http.StatusOK, allTheJobs)
}

// TODO : channel (?) to make the router stop once data has been retrieved (https://github.com/gin-gonic/gin#graceful-shutdown-or-restart)

func main2() {
	var wg sync.WaitGroup
	c := make(chan []Post)

	Logger = InitLogger()
	go InitAPIServer(c, 49)
	wg.Add(1)
	jobsList := <-c
	wg.Done()
	fmt.Println(jobsList)
	fmt.Println(len(jobsList))
	Logger.Info("Job done, closing router.")
	manners.Close()
	close(c)
}
