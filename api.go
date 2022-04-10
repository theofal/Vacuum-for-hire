package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

//InitAPIServer initialises a new API server.
func InitAPIServer() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/api/")
	{
		v1.GET("posts/:id", getAllPostsSinceLastID)
	}
	err := router.Run(":8090")
	if err != nil {
		Logger.Error("An error occurred while running router", zap.Error(err))
		return nil
	}
	return router
}

//getAllPostsSinceLastID returns the list of posts where their ID is superior to the given ID.
func getAllPostsSinceLastID(context *gin.Context) {
	id := context.Param("id")
	idInt, _ := strconv.Atoi(id)
	db, err := sql.Open("sqlite3", "./vacuum-database.db")
	if err != nil {
		Logger.Error("An error occurred while opening the database.", zap.Error(err))
		panic(err)
	}
	// TODO : Ameliorer pour éviter la possibilité d'ouvrir la db à deux endroits en meme temps (https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f)
	listOfJobs, err := Database{DB: db}.GetDataSinceSpecificID(idInt)
	if err != nil {
		Logger.Error("An error occurred while querying the database.", zap.Error(err))
		return
	}

	allTheJobs := ParseToJson(listOfJobs)
	context.JSON(http.StatusOK, gin.H{"data": allTheJobs})
}

// TODO : channel (?) to make the router stop once data has been retrieved (https://github.com/gin-gonic/gin#graceful-shutdown-or-restart)

func main2() {
	_ = InitAPIServer()
}
