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
		"SELECT userid, name, description, location, preview_image_location, added_date, content_features, views, youtube_id FROM videos WHERE videoid = ?",
		id,
	).Scan(&video.Userid, &video.Name, &video.Description, &video.Location, &video.PreviewImageLocation, &video.AddedDate, &video.ContentFeatures, &video.Views, &video.YouTubeId)

	if err1 != nil {
		return nil, fmt.Errorf("query has failed: %w", err1)
	}

	return video, nil
}

func (r *VideoDAL) SaveVideo(video models.Video) {
	r.DB.Query(`INSERT INTO videos (videoid, userid, location, preview_image_location, content_features, added_date, youtube_id, content_rating, category, language, name, description, views, tags, location_type) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		video.Videoid, video.Userid, video.Location, video.PreviewImageLocation, video.ContentFeatures, video.AddedDate, video.YouTubeId, video.ContentRating, video.Category, video.Language, video.Name, video.Description, video.Views, video.Tags, video.LocationType,
	).Exec()
}

func (r *VideoDAL) SaveLatestVideo(video models.LatestVideo) {
	r.DB.Query(`INSERT INTO latest_videos (videoid, userid, preview_image_location, added_date, content_rating, category, name, day) VALUES(?,?,?,?,?,?,?,?)`,
		video.Videoid, video.Userid, video.PreviewImageLocation, video.AddedDate, video.ContentRating, video.Category, video.Name, video.Day,
	).Exec()
}

func (r *VideoDAL) UpdateYoutubeId(videoid apachegocql.UUID, youtubeId string) {
	r.DB.Query(
		"UPDATE videos SET youtube_id = ? WHERE videoid = ?",
		youtubeId,
		videoid,
	).Exec()
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
		"SELECT videoid as key, day, added_date, videoid, category, content_rating,name, preview_image_location, userid FROM latest_videos LIMIT ?", limit,
	).Iter()

	var latest models.LatestVideo
	var returnVal []models.LatestVideo

	for iter.Scan(&latest.Key, &latest.Day, &latest.AddedDate, &latest.Videoid, &latest.Category, &latest.ContentRating, &latest.Name, &latest.PreviewImageLocation, &latest.Userid) {
		returnVal = append(returnVal, latest)
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("query has failed: %w", err)
	}

	return &returnVal, nil
}

func (r *VideoDAL) GetVideosByVector(vector [384]float32, limit int) (*[]models.Video, error) {
	iter := r.DB.Query(
		"SELECT videoid, userid, name, description, location, preview_image_location, added_date, views, youtube_id FROM videos ORDER BY content_features ANN OF ? LIMIT ?", vector, limit,
	).Iter()
	var video models.Video
	var returnVal []models.Video
	for iter.Scan(&video.Videoid, &video.Userid, &video.Name, &video.Description, &video.Location, &video.PreviewImageLocation, &video.AddedDate, &video.Views, &video.YouTubeId) {
		returnVal = append(returnVal, video)
	}
	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("vector query has failed: %w", err)
	}

	return &returnVal, nil
}

func (r *VideoDAL) UpdateVideoView(videoid apachegocql.UUID, views int) {
	r.DB.Query(
		"UPDATE videos SET views = ? WHERE videoid = ?",
		views,
		videoid,
	).Exec()
}
