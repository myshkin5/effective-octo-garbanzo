package main

import (
	"github.com/myshkin5/effective-octo-garbanzo/logs"
	"github.com/myshkin5/effective-octo-garbanzo/persistence"
)

func main() {
	err := logs.Init()
	if err != nil {
		panic(err)
	}

	err = persistence.Migrate()
	if err != nil {
		logs.Logger.Panic("Could not migrate database", err)
	}
}
