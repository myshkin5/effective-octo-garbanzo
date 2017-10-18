package main

import (
	"os"

	"github.com/mattes/migrate"
	_ "github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
	"github.com/myshkin5/effective-octo-garbanzo/logs"
)

func main() {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "INFO"
	}
	err := logs.Init(logLevel)
	if err != nil {
		panic(err)
	}

	migrator, err := migrate.New(os.Args[1], os.Args[2])
	if err != nil {
		logs.Logger.Panic(err)
	}

	err = migrator.Up()
	if err != nil {
		logs.Logger.Panic(err)
	}
}
