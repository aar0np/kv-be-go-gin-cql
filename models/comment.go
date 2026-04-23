package models

import (
	"time"

	apachegocql "github.com/apache/cassandra-gocql-driver/v2"
)

type Comment struct {
	Commentid      apachegocql.UUID `json:"commentid"`
	Videoid        apachegocql.UUID `json:"videoid"`
	Userid         apachegocql.UUID `json:"userid"`
	UserName       string           `json:"user_name"`
	CommentText    string           `json:"comment"`
	SentimentScore float32          `json:"sentiment_score"`
	Timestamp      time.Time        `json:"timestamp"`
}

func NewComment() *Comment {
	return &Comment{SentimentScore: 0.0}
}
