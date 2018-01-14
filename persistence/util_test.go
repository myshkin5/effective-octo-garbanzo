package persistence_test

import (
	"context"
	"fmt"
	"time"

	"github.com/myshkin5/effective-octo-garbanzo/logs"
	"github.com/myshkin5/effective-octo-garbanzo/persistence"
)

var ctx = context.Background()

func cleanDatabase(database persistence.Database) {
	execute("delete from garbanzo", database)
	execute("delete from octo", database)
	execute("delete from org where name like 'int_test_org_%'", database)
}

func execute(query string, database persistence.Database) {
	_, err := database.Exec(ctx, query)
	if err != nil {
		logs.Logger.Panicf("Error cleaning database with query %s: %v", query, err)
	}
}

func createOrg(suffix string, database persistence.Database) (int, string) {
	name := fmt.Sprintf("int_test_org_%d_%s", time.Now().Sub(start), suffix)
	id, err := persistence.ExecInsert(ctx, database, "insert into org (name) values ($1) returning id", name)
	if err != nil {
		logs.Logger.Panicf("Error creating org %s: %v", name, err)
	}

	return id, name
}
