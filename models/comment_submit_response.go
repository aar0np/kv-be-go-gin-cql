package models

import (
	"time"
)

type CommentSubmitResponse struct {
	CommentId      string
	VideoId        string
	UserId         string
	Comment        string
	Timestamp      time.Time
	SentimentScore float32
	FirstName      string
	LastName       string
	UserName       string
}
