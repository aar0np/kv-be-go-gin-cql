package models

type JwtResponse struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	Token  string `json:"token"`
}
