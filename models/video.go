package models

import (
	"time"

	apachegocql "github.com/apache/cassandra-gocql-driver/v2"
)

type Video struct {
	Videoid              apachegocql.UUID `json:"videoId"`
	Userid               apachegocql.UUID `json:"userId"`
	Name                 string           `json:"title"`
	Description          string           `json:"description"`
	Location             string           `json:"location"`
	PreviewImageLocation string           `json:"thumbnailUrl"`
	ContentFeatures      [384]float32     `json:"contentFeatures"`
	AddedDate            time.Time        `json:"submittedAt"`
	Views                int              `json:"views"`
	Score                float32          `json:"averageRating"`
	YouTubeId            string           `json:"youtubeVideoId"`
}

func NewVideo() *Video {
	return &Video{Views: 0}
}
