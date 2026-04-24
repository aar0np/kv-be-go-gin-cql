package models

import apachegocql "github.com/apache/cassandra-gocql-driver/v2"

type VideoSubmitRequest struct {
	YouTubeUrl  string           `json:"youtubeUrl"`
	Description string           `json:"description"`
	Tags        []string         `json:"tags"`
	UserId      apachegocql.UUID `json:"userid"`
}
