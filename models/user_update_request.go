package models

type UserUpdateRequest struct {
	Email     string
	FirstName string
	LastName  string
	Password  string
}
