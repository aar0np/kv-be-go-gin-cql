package models

import (
	"time"

	apachegocql "github.com/apache/cassandra-gocql-driver/v2"
)

type VideoResponse struct {
	Key              apachegocql.UUID `json:"key"`
	Videoid          apachegocql.UUID `json:"videoId"`
	Userid           apachegocql.UUID `json:"userId"`
	Title            string           `json:"title"`
	Description      string           `json:"description"`
	Tags             []string         `json:"tags"`
	Location         string           `json:"location"`
	ThumbnailUrl     string           `json:"thumbnailUrl"`
	SubmittedAt      time.Time        `json:"submittedAt"`
	UploadDate       time.Time        `json:"uploadDate"`
	Creator          string           `json:"creator"`
	CommentCount     int              `json:"commentCount"`
	Views            int              `json:"views"`
	ProcessingStatus string           `json:"processingStatus"`
	AverageRating    float32          `json:"averageRating"`
	ContentFeatures  [384]float32     `json:"content_features"`
	YouTubeId        string           `json:"youtubeVideoId"`
	ContentRating    string           `json:"contentRating"`
	Category         string           `json:"category"`
}
