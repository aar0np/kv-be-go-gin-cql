package main

import (
	"killrvideo/go-backend-astra-cql/controllers"
	repo "killrvideo/go-backend-astra-cql/repository"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	// define DB connection
	cfg := repo.AstraConfig{
		AstraDBID: os.Getenv("ASTRA_DB_ID"),
		Region:    "",
		Token:     os.Getenv("ASTRA_DB_APPLICATION_TOKEN"),
		Keyspace:  os.Getenv("ASTRA_DB_KEYSPACE"),
	}

	session, err := repo.NewAstraSession(cfg)

	if err != nil {
		log.Fatalf("Failed to connect to Astra: %v", err)
	}
	defer session.Close()

	// controller definitions
	healthController := controllers.NewHealthController()
	videoController := controllers.NewVideoController(session)

	// route definitions
	router := gin.Default()
	api := router.Group("/api/v1")
	{
		health := api.Group("/health")
		{
			health.GET("", healthController.GetHealth)
		}
		videos := api.Group("/videos")
		{
			videos.GET("/id/:id", videoController.GetVideo)
		}
	}

	//router.GET("/health", healthController.GetHealth)

	router.Run("localhost:8080")
}
