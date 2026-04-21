package models

import (
	"time"

	apachegocql "github.com/apache/cassandra-gocql-driver/v2"
)

type LatestVideo struct {
	Key                  apachegocql.UUID `json:"key"`
	Videoid              apachegocql.UUID `json:"videoId"`
	Userid               apachegocql.UUID `json:"userId"`
	Name                 string           `json:"title"`
	PreviewImageLocation string           `json:"thumbnailUrl"`
	ContentRating        string           `json:"contentRating"`
	Category             string           `json:"category"`
	AddedDate            time.Time        `json:"submittedAt"`
	Day                  time.Time        `json:"day"`
	Score                float32          `json:"averageRating"`
	ViewCount            int              `json:"views"`
}

func NewLatestVideo() *LatestVideo {
	return &LatestVideo{ViewCount: 0, Score: 0.0}
}
