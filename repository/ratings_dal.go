package repository

import (
	"fmt"
	"killrvideo/go-backend-astra-cql/models"

	apachegocql "github.com/apache/cassandra-gocql-driver/v2"
)

type RatingsDAL struct {
	DB *apachegocql.Session
}

func NewRatingsDAL(session *apachegocql.Session) *RatingsDAL {
	return &RatingsDAL{
		DB: session,
	}
}

func (r *RatingsDAL) GetRatingsByVideoId(videoid apachegocql.UUID) (*[]models.Rating, error) {
	var ratings []models.Rating
	var rating models.Rating

	iter := r.DB.Query(
		"SELECT videoid, rating_counter, rating_total, CAST(rating_total AS float) / rating_counter AS score FROM video_ratings_no_counters WHERE videoid = ?",
		videoid,
	).Iter()

	for iter.Scan(&rating.Videoid, &rating.RatingCounter, &rating.RatingTotal, &rating.Score) {
		ratings = append(ratings, rating)
	}

	if err1 := iter.Close(); err1 != nil {
		fmt.Println(err1)
		return nil, fmt.Errorf("query has failed: %w", err1)
	}
	return &ratings, nil
}

func (r *RatingsDAL) GetSingleRating(videoid apachegocql.UUID) (*models.Rating, error) {
	var rating models.Rating

	err1 := r.DB.Query(
		"SELECT videoid, rating_counter, rating_total, CAST(rating_total AS float) / rating_counter AS score FROM video_ratings_no_counters WHERE videoid = ?",
		videoid).
		Scan(&rating.Videoid, &rating.RatingCounter, &rating.RatingTotal, &rating.Score)

	if err1 != nil {
		return nil, fmt.Errorf("query has failed: %w", err1)
	}
	return &rating, nil
}

//func (r *RatingsDAL) GetRatingByVideoIdAndUserId(videoid string, userid string) (*models.Rating, error) {
//
//}
