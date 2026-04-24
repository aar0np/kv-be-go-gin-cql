package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"killrvideo/go-backend-astra-cql/models"
	repo "killrvideo/go-backend-astra-cql/repository"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
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

var youTubeAPIKey = os.Getenv("YOUTUBE_API_KEY")
var hfAPIKey = os.Getenv("HF_API_KEY")

const youTubeAPI = "https://www.googleapis.com/youtube/v3/videos?part=snippet&id={YOUTUBE_ID}&key={API_KEY}"
const hFAploetzSpaceEndpoint = "https://aploetz-granite-embeddings.hf.space/embed"
const modelId = "ibm-granite/granite-embedding-30m-english"

type VideoController struct {
	videoDAL   repo.VideoDAL
	ratingsDAL repo.RatingsDAL
	commentDAL repo.CommentDAL
	authDAL    repo.AuthDAL
}

func NewVideoController(session *apachegocql.Session) *VideoController {
	return &VideoController{
		videoDAL:   *repo.NewVideoDAL(session),
		ratingsDAL: *repo.NewRatingsDAL(session),
		commentDAL: *repo.NewCommentDAL(session),
		authDAL:    *repo.NewAuthDAL(session),
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

func (vc *VideoController) SubmitVideo(c *gin.Context) {
	// get userid from auth
	userid, err1 := getUserIdFromToken(c)
	if err1 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err1.Error()})
		fmt.Println("count not parse userid from auth")
		return
	}

	// get request body
	var submitRequest models.VideoSubmitRequest
	if err3 := c.ShouldBindJSON(&submitRequest); err3 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err3.Error()})
		fmt.Println("count not bind from VideoSubmitRequest")
		return
	}

	videoid, err6 := apachegocql.RandomUUID()
	if err6 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err6.Error()})
		fmt.Println("count not gen a UUID for the videoid")
		return
	}

	// build local video object
	youtubeId := extractYouTubeId(submitRequest.YouTubeUrl)

	video := models.Video{
		Description:  submitRequest.Description,
		Tags:         submitRequest.Tags,
		Location:     submitRequest.YouTubeUrl,
		YouTubeId:    youtubeId,
		Userid:       userid,
		Videoid:      videoid,
		LocationType: 1,
		AddedDate:    time.Now().UTC(),
		Views:        0,
	}

	youtubeData, err4 := getYouTubeMetadata(youtubeId)
	if err4 != nil {
		fmt.Println("count not pull Youtube metadata")
		c.JSON(http.StatusBadRequest, gin.H{"error": err4})
		return
	}
	//fmt.Println(youtubeData.Title)

	video.Name = youtubeData.Title
	video.PreviewImageLocation = youtubeData.ThumbnailUrl
	video.Tags = youtubeData.Tags

	// get embedding from HuggingFace Space endpoint
	hfResponse, err5 := getHFEmbeddingData(video.Name)
	if err5 != nil {
		fmt.Println("count not generate embedding")
	} else {
		video.ContentFeatures = hfResponse.Embedding
	}

	// save video
	vc.videoDAL.SaveVideo(video)

	//save latest_video
	latest := models.LatestVideo{
		Day:                  time.Now().UTC(),
		AddedDate:            video.AddedDate,
		Videoid:              videoid,
		Name:                 video.Name,
		PreviewImageLocation: video.PreviewImageLocation,
		Userid:               userid,
		ContentRating:        video.ContentRating,
		Category:             video.Category,
	}
	vc.videoDAL.SaveLatestVideo(latest)

	// build and return response
	videoResponse := models.VideoResponse{
		Key:             videoid,
		Videoid:         videoid,
		Userid:          userid,
		Title:           video.Name,
		Description:     video.Description,
		Tags:            video.Tags,
		Location:        video.Location,
		ThumbnailUrl:    video.PreviewImageLocation,
		SubmittedAt:     video.AddedDate,
		UploadDate:      video.AddedDate,
		Creator:         userid.String(),
		CommentCount:    0,
		Views:           video.Views,
		AverageRating:   0.0,
		ContentFeatures: video.ContentFeatures,
		YouTubeId:       youtubeId,
		Category:        video.Category,
	}

	c.JSON(http.StatusOK, videoResponse)
}

func (vc *VideoController) GetLatestVideos(c *gin.Context) {
	page, err1 := strconv.Atoi(c.Query("page"))
	pageSize, err2 := strconv.Atoi(c.Query("pageSize"))

	if err1 != nil {
		fmt.Println("Could not convert page string to int")
		page = 0
	}

	if err2 != nil {
		fmt.Println("Could not convert pageSize string to int")
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

	//for _, latestVideo := range *latestVideos {
	for i := range *latestVideos {
		lVideo := &(*latestVideos)[i]

		// get ratings
		rating, err5 := vc.ratingsDAL.GetSingleRating(lVideo.Videoid)
		if err5 != nil {
			fmt.Println(err5)
		}

		if rating == nil {
			lVideo.Score = 0
		} else {
			lVideo.Score = rating.Score
		}

		// get views
		video, err6 := vc.videoDAL.GetVideo(lVideo.Videoid)
		if err6 != nil {
			fmt.Println(err6)
		}

		lVideo.ViewCount = video.Views
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

		// get ratings
		rating, err5 := vc.ratingsDAL.GetSingleRating(video.Videoid)
		if err5 != nil {
			fmt.Println(err5)
		}

		if rating == nil {
			video.Score = 0
		} else {
			video.Score = rating.Score
		}

		returnVal = append(returnVal, video)
		uniqueVideoIDs[video.Name] = struct{}{}

		if len(returnVal) >= limit {
			break
		}
	}

	c.JSON(http.StatusOK, returnVal)
}

func (vc *VideoController) RecordVideoView(c *gin.Context) {
	videoid, err1 := apachegocql.ParseUUID(c.Param("id"))
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"viewError1": err1})
		return
	}

	video, err2 := vc.videoDAL.GetVideo(videoid)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"viewError2": err2})
		return
	}

	views := video.Views + 1

	vc.videoDAL.UpdateVideoView(videoid, views)
}

func (vc *VideoController) GetComments(c *gin.Context) {
	videoid, err1 := apachegocql.ParseUUID(c.Param("id"))
	page, err2 := strconv.Atoi(c.Query("page"))
	pageSize, err3 := strconv.Atoi(c.Query("pageSize"))

	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"commentError1": err1})
		return
	}
	if err2 != nil {
		fmt.Println("Could not convert page string to integer")
		page = 0
	}

	if err3 != nil {
		fmt.Println("Could not convert pageSize string to integer")
		pageSize = 0
	}

	if page <= 0 || pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	comments, err4 := vc.commentDAL.GetCommentsByVideoId(videoid, pageSize)
	if err4 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"commentError3": err4})
		return
	}

	// Initialize response with empty slice to ensure JSON returns [] instead of null
	var commentData []models.Comment
	if comments != nil && len(*comments) > 0 {
		commentData = *comments
	} else {
		commentData = make([]models.Comment, 0)
	}

	pagination := models.Pagination{
		TotalPages:  1,
		PageSize:    pageSize,
		TotalItems:  len(commentData),
		CurrentPage: page,
	}

	returnVal := models.CommentResponse{Data: commentData, Pagination: pagination}

	c.JSON(http.StatusOK, returnVal)
}

func (vc *VideoController) SubmitComment(c *gin.Context) {
	// get userid from auth
	userid, err1 := getUserIdFromToken(c)
	if err1 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err1.Error()})
		fmt.Println("count not parse userid from auth")
		return
	}

	// get videoid from url
	videoid, err2 := apachegocql.ParseUUID(c.Param("id"))
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err2.Error()})
		fmt.Println("count not parse videoid from url")
		return
	}

	// get comment from request body
	var commentReq models.CommentSubmitRequest
	if err3 := c.ShouldBindJSON(&commentReq); err3 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err3.Error()})
		fmt.Println("count not commentReq videoid from CommentSubmitRequest")
		return
	}

	// generate commentid TimeUUID and sentiment score
	commentid := apachegocql.TimeUUID()

	// save comment to database
	comment := models.Comment{
		Videoid:        videoid,
		Commentid:      commentid,
		Userid:         userid,
		CommentText:    commentReq.Text,
		SentimentScore: 0,
	}
	vc.commentDAL.SaveComment(comment)
	vc.commentDAL.SaveCommentByUser(comment)

	firstName := ""
	lastName := ""
	userName := ""

	user, err4 := vc.authDAL.GetUserById(userid)
	if err4 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err4.Error()})
		fmt.Println("count not get user from authDAL")
	} else {
		firstName = user.FirstName
		lastName = user.LastName
		userName = user.FirstName + " " + user.LastName
	}

	response := models.CommentSubmitResponse{
		CommentId:      commentid.String(),
		VideoId:        videoid.String(),
		UserId:         userid.String(),
		Comment:        commentReq.Text,
		Timestamp:      time.Now(),
		SentimentScore: 0,
		FirstName:      firstName,
		LastName:       lastName,
		UserName:       userName,
	}

	c.JSON(http.StatusCreated, response)
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

func getYouTubeMetadata(youTubeId string) (models.YouTubeMetadata, error) {
	// "https://www.googleapis.com/youtube/v3/videos?part=snippet&id={YOUTUBE_ID}&key={API_KEY}"
	url := strings.Replace(youTubeAPI, "{API_KEY}", youTubeAPIKey, 1)
	url = strings.Replace(url, "{YOUTUBE_ID}", youTubeId, 1)

	youTubeClient := &http.Client{}
	var response models.YouTubeMetadata

	req, _ := http.NewRequest("GET", url, nil)

	resp, err1 := youTubeClient.Do(req)
	if err1 != nil {
		fmt.Println("Error calling the YouTube API")
		return response, err1
	}
	defer resp.Body.Close()

	var youTubeResponse models.YouTubeResponse

	if err2 := json.NewDecoder(resp.Body).Decode(&youTubeResponse); err2 != nil {
		fmt.Println("Error processing the YouTube response")
		return response, err2
	}

	if len(youTubeResponse.Items) > 0 {
		snippet := youTubeResponse.Items[0].Snippet
		response.Title = snippet.Title
		response.Description = snippet.Description
		response.Tags = snippet.Tags

		if snippet.Thumbnails.High.URL != "" {
			response.ThumbnailUrl = snippet.Thumbnails.High.URL
		} else if snippet.Thumbnails.Medium.URL != "" {
			response.ThumbnailUrl = snippet.Thumbnails.Medium.URL
		} else if snippet.Thumbnails.Default.URL != "" {
			response.ThumbnailUrl = snippet.Thumbnails.Default.URL
		}
	}

	return response, nil
}

func getHFEmbeddingData(text string) (models.HuggingFaceResponse, error) {
	hfClient := &http.Client{}

	//hfRequest := models.HuggingFaceRequest{
	//	Text:  text,
	//	Model: modelId,
	//}
	hfRequestString := "{\"text\": \"" + text + "\", \"model\": \"" + modelId + "\"}"
	jsonBody := []byte(hfRequestString)
	bodyReader := bytes.NewReader(jsonBody)

	req, _ := http.NewRequest(http.MethodPost, hFAploetzSpaceEndpoint, bodyReader)
	var response models.HuggingFaceResponse

	resp, err1 := hfClient.Do(req)
	if err1 != nil {
		fmt.Println("Error calling the HuggingFace Spaces endpoint")
		return response, err1
	}
	defer resp.Body.Close()

	if err2 := json.NewDecoder(resp.Body).Decode(&response); err2 != nil {
		fmt.Println("Error processing the HuggingFace Spaces response")
		return response, err2
	}

	return response, nil
}
