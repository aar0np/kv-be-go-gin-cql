package models

type RatingsResponse struct {
	Data          []Rating `json:"data"`
	AverageRating float32  `json:"averageRating"`
}
