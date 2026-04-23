package repository

import (
	"fmt"
	"killrvideo/go-backend-astra-cql/models"

	apachegocql "github.com/apache/cassandra-gocql-driver/v2"
)

type CommentDAL struct {
	DB *apachegocql.Session
}

func NewCommentDAL(session *apachegocql.Session) *CommentDAL {
	return &CommentDAL{
		DB: session,
	}
}

func (r *CommentDAL) GetCommentsByVideoId(videoid apachegocql.UUID, limit int) (*[]models.Comment, error) {
	// Initialize with empty slice to ensure we never return nil
	comments := make([]models.Comment, 0)
	var comment models.Comment

	iter := r.DB.Query(
		"SELECT videoid, userid, commentid, comment, sentiment_score FROM comments WHERE videoid = ?",
		videoid,
	).Iter()

	for iter.Scan(&comment.Videoid, &comment.Userid, &comment.Commentid, &comment.CommentText, &comment.SentimentScore) {
		comments = append(comments, comment)
	}

	if err1 := iter.Close(); err1 != nil {
		fmt.Println(err1)
		return nil, fmt.Errorf("query has failed: %w", err1)
	}
	return &comments, nil
}

func (r *CommentDAL) SaveComment(comment models.Comment) {
	r.DB.Query(`INSERT INTO comments (videoid, userid, commentid, comment, sentiment_score) VALUES(?,?,?,?,?)`,
		comment.Videoid, comment.Userid, comment.Commentid, comment.CommentText, comment.SentimentScore,
	).Exec()
}

func (r *CommentDAL) SaveCommentByUser(comment models.Comment) {
	r.DB.Query(`INSERT INTO comments_by_user (videoid, userid, commentid, comment, sentiment_score) VALUES(?,?,?,?,?)`,
		comment.Videoid, comment.Userid, comment.Commentid, comment.CommentText, comment.SentimentScore,
	).Exec()
}
