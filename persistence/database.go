package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

var (
	ErrNotFound = errors.New("identified data not found")
)

type Database interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

func Open() (Database, error) {
	server := GetEnvWithDefault("DB_SERVER", "localhost")
	port := GetEnvWithDefault("DB_PORT", "5432")
	username := GetEnvWithDefault("DB_USERNAME", "garbanzo")
	password := GetEnvWithDefault("DB_PASSWORD", "garbanzo-secret")

	database, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/garbanzo?sslmode=disable",
		username, password, server, port))
	if err != nil {
		return nil, err
	}

	maxOpenConns, err := strconv.Atoi(GetEnvWithDefault("DB_MAX_OPEN_CONNS", "10"))
	if err != nil {
		return nil, err
	}

	database.SetMaxOpenConns(maxOpenConns)

	return database, nil
}

func GetEnvWithDefault(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if ok {
		return value
	} else {
		return defaultValue
	}
}
