package repository

import (
	"fmt"
	"time"

	gocqlastra "github.com/datastax/gocql-astra"
	gocql "github.com/gocql/gocql"
)

type AstraConfig struct {
	AstraDBID string
	Region    string
	Token     string
	Keyspace  string
}

func NewAstraSession(cfg AstraConfig) (*gocql.Session, error) {
	cluster, err1 := gocqlastra.NewClusterFromURL("https://api.astra.datastax.com", cfg.AstraDBID, cfg.Token, 10*time.Second)

	if err1 != nil {
		return nil, fmt.Errorf("unable to connect to Astra cluster: %w", err1)
	}

	cluster.Keyspace = cfg.Keyspace
	cluster.ProtoVersion = 4

	session, err2 := cluster.CreateSession()

	if err2 != nil {
		return nil, fmt.Errorf("unable to create Astra session: %w", err2)
	}

	return session, nil
}
