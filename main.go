package main

import (
	"killrvideo/go-backend-astra-cql/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	healthController := controllers.NewHealthController()

	router.GET("/health", healthController.GetHealth)

	router.Run("localhost:8080")
}

//func getHealth(c *gin.Context) {
//	c.IndentedJSON(http.StatusOK, "Service is up and running!")
//}
