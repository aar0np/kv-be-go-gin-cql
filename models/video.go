package models

import (
	"time"

	gocql "github.com/gocql/gocql"
)

type Video struct {
	Videoid              gocql.UUID
	Userid               gocql.UUID
	Name                 string
	Description          string
	Location             string
	PreviewImageLocation string
	//ContentFeatures      [384]float32
	AddedDate time.Time
}
