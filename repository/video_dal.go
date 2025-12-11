package repository

import (
	"fmt"
	"killrvideo/go-backend-astra-cql/models"

	gocql "github.com/gocql/gocql"
)

type VideoDAL struct {
	DB *gocql.Session
}

func NewVideoDAL(session *gocql.Session) *VideoDAL {
	return &VideoDAL{
		DB: session,
	}
}

func (r *VideoDAL) GetVideo(id gocql.UUID) (*models.Video, error) {
	video := &models.Video{Videoid: id}

	err1 := r.DB.Query(
		"SELECT userid, name, description, location, preview_image_location, added_date FROM videos WHERE videoid = ?",
		id,
	).Scan(&video.Userid, &video.Name, &video.Description, &video.Location, &video.PreviewImageLocation, &video.AddedDate)

	if err1 != nil {
		return nil, fmt.Errorf("query has failed: %w", err1)
	}

	return video, nil
}
