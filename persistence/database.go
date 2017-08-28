package persistence

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
)

var (
	ErrNotFound = errors.New("Identified data not found")
)

type Database interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}
