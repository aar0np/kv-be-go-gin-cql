package models

import (
	"time"

	apachegocql "github.com/apache/cassandra-gocql-driver/v2"
)

type User struct {
	Userid        apachegocql.UUID
	Email         string
	FirstName     string
	LastName      string
	AccountStatus string
	CreatedDate   time.Time
	LastLoginDate time.Time
}
