package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

//InitAPIServer initialises a new API server.
func _() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/api/")
	{
		v1.GET("posts/:id", getAllPostsSinceLastID)
	}
	err := router.Run()
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
	// TODO : Ameliorer pour le rendre thread safe et éviter la var ouverte Db (https://medium.com/golang-issue/how-singleton-pattern-works-with-golang-2fdd61cd5a7f)
	listOfJobs, err := Db.GetDataSinceSpecificID(idInt)
	if err != nil {
		Logger.Error("An error occurred while querying the database.")
		return
	}
	context.JSON(http.StatusOK, gin.H{"data": listOfJobs})
}

// TODO : channel (?) to make the router stop once data has been retrieved (https://github.com/gin-gonic/gin#graceful-shutdown-or-restart)
