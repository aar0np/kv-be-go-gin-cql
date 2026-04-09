package models

import (
	apachegocql "github.com/apache/cassandra-gocql-driver/v2"
)

type UserCredentials struct {
	Email         string
	Password      string
	Userid        apachegocql.UUID
	AccountLocked bool
}
