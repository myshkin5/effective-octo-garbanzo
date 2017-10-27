package persistence_test

import (
	"context"

	"github.com/myshkin5/effective-octo-garbanzo/persistence"
)

func cleanDatabase(database persistence.Database) error {
	ctx := context.Background()
	err := cleanTable(ctx, "garbanzo", database)
	if err != nil {
		return err
	}
	err = cleanTable(ctx, "octo", database)
	if err != nil {
		return err
	}
	return nil
}

func cleanTable(ctx context.Context, table string, database persistence.Database) error {
	query := "delete from " + table
	_, err := database.Exec(ctx, query)
	return err
}
