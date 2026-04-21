package controllers

import (
	"fmt"
	"killrvideo/go-backend-astra-cql/models"
	repo "killrvideo/go-backend-astra-cql/repository"
	"net/http"
	"regexp"
	"strconv"
	"time"

	apachegocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/gin-gonic/gin"
)

var youTubePatterns = [4]string{
	"(?:https?://)?(?:www\\.)?youtu\\.be/(?<id>[A-Za-z0-9_-]{11})",
	"(?:https?://)?(?:www\\.)?youtube\\.com/watch\\?v=(?<id>[A-Za-z0-9_-]{11})",
	"(?:https?://)?(?:www\\.)?youtube\\.com/embed/(?<id>[A-Za-z0-9_-]{11})",
	"(?:https?://)?(?:www\\.)?youtube\\.com/v/(?<id>[A-Za-z0-9_-]{11})",
}

type VideoController struct {
	videoDAL repo.VideoDAL
}

func NewVideoController(session *apachegocql.Session) *VideoController {
	return &VideoController{
		videoDAL: *repo.NewVideoDAL(session),
	}
}

func (vc *VideoController) GetVideo(c *gin.Context) {
	id, err1 := apachegocql.ParseUUID(c.Param("id"))
	if err1 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err1.Error()})
	}

	video, err2 := vc.videoDAL.GetVideo(id)
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err2.Error()})
	}

	// make sure that we have a YouTubeID
	if video.YouTubeId == "" {
		video.YouTubeId = extractYouTubeId(video.Location)
		vc.videoDAL.UpdateYoutubeId(id, video.YouTubeId)
	}

	c.JSON(http.StatusOK, video)
}

func (vc *VideoController) GetLatestVideos(c *gin.Context) {
	page, err1 := strconv.Atoi(c.Query("page"))
	pageSize, err2 := strconv.Atoi(c.Query("page_size"))

	if err1 != nil {
		page = 0
	}

	if err2 != nil {
		pageSize = 0
	}

	//today := time.Now().Format("2001-01-01")
	today := time.Now()

	if page <= 0 || pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	latestVideos, err3 := vc.videoDAL.GetLatestVideosToday(today, pageSize)

	if err3 != nil {
		fmt.Println(err3)
	}

	if latestVideos != nil && len(*latestVideos) < pageSize {
		newLimit := pageSize - len(*latestVideos)
		additionalVideos, err4 := vc.videoDAL.GetLatestVideos(newLimit)

		if err4 != nil {
			fmt.Println(err4)
		}

		*latestVideos = append(*latestVideos, *additionalVideos...)
	}

	returnVal := models.LatestVideoResponse{Data: *latestVideos}

	c.JSON(http.StatusOK, returnVal)
}

func (vc *VideoController) GetSimilarVideos(c *gin.Context) {
	id, err1 := apachegocql.ParseUUID(c.Param("id"))
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err1})
		return
	}

	limit, err2 := strconv.Atoi(c.Query("limit"))
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err2})
		return
	}
	// make sure limit behaves
	if limit < 1 || limit > 20 {
		// default to 5
		limit = 5
	}

	// get original video so we can use its vector
	originalVideo, err3 := vc.videoDAL.GetVideo(id)
	if err3 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err3})
		return
	}

	similarVideos, err4 := vc.videoDAL.GetVideosByVector(originalVideo.ContentFeatures, (limit+1)*2)
	if err4 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err4})
		return
	}

	var returnVal []models.Video
	uniqueVideoIDs := make(map[string]struct{})

	for _, video := range *similarVideos {
		if video.Name == originalVideo.Name {
			continue
		}

		if _, exists := uniqueVideoIDs[video.Name]; exists {
			continue
		}

		returnVal = append(returnVal, video)
		uniqueVideoIDs[video.Name] = struct{}{}

		if len(returnVal) >= limit {
			break
		}
	}

	c.JSON(http.StatusOK, returnVal)
}

func extractYouTubeId(location string) string {
	for _, pattern := range youTubePatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(location)

		if len(matches) > 1 {
			return matches[1]
		}
	}
	return ""
}
