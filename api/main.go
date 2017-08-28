package main

import (
	"os"

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

	logs.Logger.Info("Hi logger")
}
