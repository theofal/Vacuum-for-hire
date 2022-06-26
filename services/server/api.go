package server

import (
	"database/sql"
	"encoding/json"
	"github.com/braintree/manners"
	"github.com/gin-gonic/gin"
	_ "github.com/go-ping/ping"
	"github.com/theofal/Vacuum-for-hire/services"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

//InitAPIServer initialises a new API server.
func InitAPIServer(c chan []services.Post, index int) {
	services.Logger.Info("Initialising API server.")
	router := gin.Default()

	v1 := router.Group("/server/")
	{
		v1.GET("posts/:id", getAllPostsSinceLastID)
	}

	go requestPostsFromAPI(c, index)

	services.Logger.Info("Starting server.")
	err := manners.ListenAndServe(":8090", router)
	if err != nil {
		services.Logger.Error("An error occurred while running router", zap.Error(err))
		os.Exit(1)
	}

}

//requestPostsFromAPI goroutine called in InitAPIServer. This function requests data from API server and exports data
//as []Post to main function via a channel.
func requestPostsFromAPI(c chan []services.Post, index int) {
	var resp *http.Response

	for i := 1; i <= 5; i++ {
		services.Logger.Debug("Trying to get data from API server.", zap.String("try number", strconv.Itoa(i)))
		resp, _ = http.Get("http://localhost:8090/server/posts/" + strconv.Itoa(index))
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			time.Sleep(1 * time.Second)
			continue
		} else {
			services.Logger.Debug("Got all the data I needed!?")
			break
		}
	}
	body, _ := ioutil.ReadAll(resp.Body)

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			os.Exit(1)
		}
	}(resp.Body)

	var jobs []services.Post
	err := json.Unmarshal(body, &jobs)
	if err != nil {
		services.Logger.Error("An error occurred while unmarshalling json.", zap.Error(err))
		os.Exit(1)
	}
	services.Logger.Info("Done getting data from API server.", zap.String("Number of jobs", strconv.Itoa(len(jobs))))
	c <- jobs
}

//getAllPostsSinceLastID returns the list of posts where their ID is superior to the given ID.
func getAllPostsSinceLastID(context *gin.Context) {
	id := context.Param("id")
	idInt, _ := strconv.Atoi(id)
	db, err := sql.Open("sqlite3", "./vacuum-database.db")
	if err != nil {
		services.Logger.Error("An error occurred while opening the database.", zap.Error(err))
		os.Exit(1)
	}
	// TODO : Ameliorer pour éviter la possibilité d'ouvrir la db à deux endroits en meme temps (https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f)
	// TODO : remove the database instantiation via a getter ?
	jobList, err := Database{DB: db}.GetDataSinceSpecificID(idInt)
	if err != nil {
		services.Logger.Error("An error occurred while querying the database.", zap.Error(err))
		os.Exit(1)
	}

	allTheJobs := services.ParseToJSON(jobList)
	context.JSON(http.StatusOK, allTheJobs)
}
