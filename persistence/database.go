package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	// Used by main.go and tests to import the proper database driver
	_ "github.com/lib/pq"
	"github.com/mattes/migrate"
	// Used by main.go and tests to import the proper migration drivers
	_ "github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"

	"github.com/myshkin5/effective-octo-garbanzo/logs"
)

var (
	ErrNotFound = errors.New("identified data not found")
)

type Database interface {
	Exec(ctx context.Context, query string, args ...interface{}) (result sql.Result, err error)
	Query(ctx context.Context, query string, args ...interface{}) (rows *sql.Rows, err error)
	QueryRow(ctx context.Context, query string, args ...interface{}) (row *sql.Row)
	BeginTx(ctx context.Context) (database Database, err error)
	Commit() (err error)
	Rollback() (err error)
}

const OrgContextKey = "org"

func org(ctx context.Context) string {
	return ctx.Value(OrgContextKey).(string)
}

func Open() (Database, error) {
	db, err := sql.Open("postgres", getDatabaseURL())
	if err != nil {
		return nil, err
	}

	maxOpenConns, err := strconv.Atoi(GetEnvWithDefault("DB_MAX_OPEN_CONNS", "10"))
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)

	verifyConnection(db)

	return &database{
		internalDB: db,
		internalTx: nil,
	}, nil
}

func ExecInsert(ctx context.Context, database Database, query string, args ...interface{}) (int, error) {
	var id int
	err := database.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func ExecDelete(ctx context.Context, database Database, query string, args ...interface{}) (int64, error) {
	result, err := database.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func verifyConnection(db *sql.DB) {
	query := "select 1"
	for {
		for i := 0; i < 20; i++ {
			_, err := db.Query(query)
			if err == nil {
				return
			}

			time.Sleep(1 * time.Second)
		}
		logs.Logger.Warn("Could not connect to database. Continuing to try...")
	}
}

type migrateLogger struct{}

func (l migrateLogger) Printf(format string, v ...interface{}) {
	logs.Logger.Infof(format, v...)
}

func (l migrateLogger) Verbose() bool {
	return false
}

func Migrate() error {
	sourceURL := GetEnvWithDefault("DB_SOURCE_URL", "file://./persistence/ddl")
	databaseURL := getDatabaseURL()

	// Open() verifies that the database is up and running
	_, err := Open()
	if err != nil {
		return err
	}

	migrator, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		return err
	}

	migrator.Log = migrateLogger{}

	err = migrator.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func getDatabaseURL() string {
	server := GetEnvWithDefault("DB_SERVER", "localhost")
	port := GetEnvWithDefault("DB_PORT", "5432")
	username := GetEnvWithDefault("DB_USERNAME", "garbanzo")
	password := GetEnvWithDefault("DB_PASSWORD", "garbanzo-secret")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/garbanzo?sslmode=disable", username, password, server, port)
}

func GetEnvWithDefault(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if ok {
		return value
	}

	return defaultValue
}

type internalDBAndTx interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type internalTx interface {
	internalDBAndTx
	Commit() error
	Rollback() error
}

type internalDB interface {
	internalDBAndTx
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type database struct {
	internalDB internalDB
	internalTx internalTx
}

func (d *database) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if d.internalTx == nil {
		return d.internalDB.ExecContext(ctx, query, args...)
	}

	return d.internalTx.ExecContext(ctx, query, args...)
}

func (d *database) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if d.internalTx == nil {
		return d.internalDB.QueryContext(ctx, query, args...)
	}

	return d.internalTx.QueryContext(ctx, query, args...)
}

func (d *database) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	if d.internalTx == nil {
		return d.internalDB.QueryRowContext(ctx, query, args...)
	}

	return d.internalTx.QueryRowContext(ctx, query, args...)
}

func (d *database) BeginTx(ctx context.Context) (Database, error) {
	tx, err := d.internalDB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &database{
		internalDB: nil,
		internalTx: tx,
	}, nil
}

func (d *database) Commit() error {
	return d.internalTx.Commit()
}

func (d *database) Rollback() error {
	return d.internalTx.Rollback()
}
