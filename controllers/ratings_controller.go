package controllers

import (
	"killrvideo/go-backend-astra-cql/models"
	repo "killrvideo/go-backend-astra-cql/repository"
	"net/http"

	apachegocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/gin-gonic/gin"
)

type RatingsController struct {
	ratingsDAL repo.RatingsDAL
}

func NewRatingsController(session *apachegocql.Session) *RatingsController {
	return &RatingsController{
		ratingsDAL: *repo.NewRatingsDAL(session),
	}
}

func (rc *RatingsController) GetRatingsByVideoId(c *gin.Context) {
	id, err1 := apachegocql.ParseUUID(c.Param("id"))
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err1})
		return
	}

	ratings, err2 := rc.ratingsDAL.GetRatingsByVideoId(id)
	var avgRating float32 = 0.0

	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err2})
		return
	}

	for _, rating := range *ratings {
		avgRating = rating.Score
		break
	}

	returnVal := models.RatingsResponse{Data: *ratings, AverageRating: avgRating}
	c.JSON(http.StatusOK, returnVal)
}
