package models

import (
	"time"

	apachegocql "github.com/apache/cassandra-gocql-driver/v2"
)

type LatestVideo struct {
	Videoid              apachegocql.UUID
	Userid               apachegocql.UUID
	Name                 string
	PreviewImageLocation string
	ContentRating        string
	Category             string
	AddedDate            time.Time
	Day                  time.Time
}
