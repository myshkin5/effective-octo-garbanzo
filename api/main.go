package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"

	"github.com/myshkin5/effective-octo-garbanzo/api/handlers"
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

	router := mux.NewRouter()

	middleware := alice.New()

	handlers.MapHealthRoutes(router.Path("/health"), middleware)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	logs.Logger.Infof("Listening on %s...", port)
	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		logs.Logger.Critical("ListenAndServe:", err)
	}
}
