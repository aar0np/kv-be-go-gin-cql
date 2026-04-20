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
	YouTubeId            string           `json:"youtubeVideoId"`
}
