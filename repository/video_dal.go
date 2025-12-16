package repository

import (
	"fmt"
	"killrvideo/go-backend-astra-cql/models"
	"time"

	apachegocql "github.com/apache/cassandra-gocql-driver/v2"
)

type VideoDAL struct {
	DB *apachegocql.Session
}

func NewVideoDAL(session *apachegocql.Session) *VideoDAL {
	return &VideoDAL{
		DB: session,
	}
}

func (r *VideoDAL) GetVideo(id apachegocql.UUID) (*models.Video, error) {
	video := &models.Video{Videoid: id}
	//var vector []float32

	err1 := r.DB.Query(
		"SELECT userid, name, description, location, preview_image_location, added_date, content_features FROM videos WHERE videoid = ?",
		id,
	).Scan(&video.Userid, &video.Name, &video.Description, &video.Location, &video.PreviewImageLocation, &video.AddedDate, &video.ContentFeatures)

	if err1 != nil {
		return nil, fmt.Errorf("query has failed: %w", err1)
	}

	return video, nil
}

func (r *VideoDAL) GetLatestVideosToday(day time.Time, limit int) (*[]models.LatestVideo, error) {
	iter := r.DB.Query(
		"SELECT day, added_date, videoid, category, content_rating,name, preview_image_location, userid FROM latest_videos WHERE day=? LIMIT ?", day, limit,
	).Iter()

	var latest models.LatestVideo
	var returnVal []models.LatestVideo

	for iter.Scan(&latest.Day, &latest.AddedDate, &latest.Videoid, &latest.Category, &latest.ContentRating, &latest.Name, &latest.PreviewImageLocation, &latest.Userid) {
		returnVal = append(returnVal, latest)
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("query has failed: %w", err)
	}

	return &returnVal, nil
}

func (r *VideoDAL) GetLatestVideos(limit int) (*[]models.LatestVideo, error) {
	iter := r.DB.Query(
		"SELECT day, added_date, videoid, category, content_rating,name, preview_image_location, userid FROM latest_videos LIMIT ?", limit,
	).Iter()

	var latest models.LatestVideo
	var returnVal []models.LatestVideo

	for iter.Scan(&latest.Day, &latest.AddedDate, &latest.Videoid, &latest.Category, &latest.ContentRating, &latest.Name, &latest.PreviewImageLocation, &latest.Userid) {
		returnVal = append(returnVal, latest)
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("query has failed: %w", err)
	}

	return &returnVal, nil
}
