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
		Token:    os.Getenv("ASTRA_DB_APPLICATION_TOKEN"),
		Keyspace: os.Getenv("ASTRA_DB_KEYSPACE"),
		ScbDir:   os.Getenv("ASTRA_DB_SCB_DIR"),
		Hostname: os.Getenv("ASTRA_DB_HOSTNAME"),
	}

	session, err := repo.NewAstraSession(cfg)

	if err != nil {
		log.Fatalf("Failed to connect to Astra: %v", err)
	}
	defer session.Close()

	// controller definitions
	authController := controllers.NewAuthController(session)
	healthController := controllers.NewHealthController()
	ratingsController := controllers.NewRatingsController(session)
	videoController := controllers.NewVideoController(session)

	// route definitions
	router := gin.Default()
	api := router.Group("/api/v1")
	{
		health := api.Group("/health")
		{
			health.GET("", healthController.GetHealth)
		}
		auth := api.Group("/users")
		{
			auth.POST("/login", authController.Login)
			auth.GET("/me", authController.GetCurrentUser)
			auth.GET(":id", authController.GetUser)
		}
		videos := api.Group("/videos")
		{
			videos.GET("/id/:id", videoController.GetVideo)
			videos.GET("/latest", videoController.GetLatestVideos)
			videos.GET("/id/:id/related", videoController.GetSimilarVideos)
			videos.GET("/:id/ratings", ratingsController.GetRatingsByVideoId)
			videos.POST("/id/:id/view", videoController.RecordVideoView)
			videos.GET("/:id/comments", videoController.GetComments)
			videos.POST("/:id/comments", videoController.SubmitComment)
		}
	}

	router.RunTLS("localhost:8443", "localhost.pem", "localhost-key.pem")
}
