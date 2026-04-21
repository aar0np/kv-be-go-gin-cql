package models

import (
	"time"

	apachegocql "github.com/apache/cassandra-gocql-driver/v2"
)

type Rating struct {
	Videoid       apachegocql.UUID `json:"videoid"`
	Userid        apachegocql.UUID `json:"userid"`
	RatingCounter int              `json:"ratingCount"`
	RatingTotal   int              `json:"ratingTotal"`
	RatingDate    time.Time        `json:"ratingDate"`
	Score         float32          `json:"averageRating"`
}

func NewRating() *Rating {
	return &Rating{Score: 0.0, RatingCounter: 0, RatingTotal: 0}
}
