package controllers

import (
	repo "killrvideo/go-backend-astra-cql/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	gocql "github.com/gocql/gocql"
)

type VideoController struct {
	videoDAL repo.VideoDAL
}

func NewVideoController(session *gocql.Session) *VideoController {
	return &VideoController{
		videoDAL: *repo.NewVideoDAL(session),
	}
}

func (vc *VideoController) GetVideo(c *gin.Context) {
	id, err1 := gocql.ParseUUID(c.Param("id"))
	if err1 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err1.Error()})
	}

	video, err2 := vc.videoDAL.GetVideo(id)
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err2.Error()})
	}

	c.JSON(http.StatusOK, video)
}
