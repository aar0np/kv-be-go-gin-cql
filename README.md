# kv-be-go-gin-cql
A reference backend for the KilrlVideo sample application rebuilt for 2026 using **Golang**, and **DataStax Astra DB**.

## Notice
Work-in-progress!

## Overview
This repo demonstrates modern API best-practices with:

* Restful, typed request/response models
* Role-based JWT auth
* Micro-service friendly layout – or run everything as a monolith

## Prerequisites
1. **Go 1.25.5** or later.
2. A **DataStax Astra DB** serverless database – [grab a free account](https://astra.datastax.com).

## Setup & Configuration
```bash
# clone
git clone git@github.com:KillrVideo/kv-be-go-gin-cql.git
cd kv-be-go-gin-cql

# build and install dependencies
go get github.com/golang-jwt/jwt/v5
go get github.com/apache/cassandra-gocql-driver/v2@latest
go get -u github.com/gin-gonic/gin
```

### Database schema:
1. Create a new keyspace named `killrvideo`.
2. Create the tables from the CQL file in the killrvideo-data repository: <https://github.com/KillrVideo/killrvideo-data/blob/master/schema-astra.cql>

### Environment variables (via `export`):
| Variable | Description |
|----------|-------------|
| `ASTRA_DB_APPLICATION_TOKEN` | The token created in the Astra UI |
| `ASTRA_DB_HOSTNAME` | The hostname for your Astra database |
| `ASTRA_DB_KEYSPACE` | `killrvideo` |
| `ASTRA_DB_SCB_DIR` | Path to the SCB file downloaded from the Astra UI |
| `JWT_KEY` | A random, 64-byte secret key used to sign the JSON Web Token |
| `YOUTUBE_API_KEY` | Required for pulling new video info from YouTube API |
| `HF_API_KEY` | HuggingFace key used to hit a HuggingFace Space to create an embedding |

## Running the Application
```bash
go run main.go
```

## Test the health check service
```bash
curl -X GET "https://localhost:8443/api/v1/health" \
--header "Content-Type: application/json" \
```
"Service is up and running!"