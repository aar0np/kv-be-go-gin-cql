package models

import (
	"time"

	apachegocql "github.com/apache/cassandra-gocql-driver/v2"
)

type User struct {
	Userid        apachegocql.UUID `json:"userId"`
	Email         string           `json:"email"`
	FirstName     string           `json:"firstName"`
	LastName      string           `json:"lastName"`
	AccountStatus string           `json:"accountStatus"`
	CreatedDate   time.Time        `json:"createdAt"`
	LastLoginDate time.Time        `json:"lastLoginDate"`
}
