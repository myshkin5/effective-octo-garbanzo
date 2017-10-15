package main

import (
	"net/http"
	"os"

	gorilla_handlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"

	"github.com/myshkin5/effective-octo-garbanzo/api/handlers"
	"github.com/myshkin5/effective-octo-garbanzo/logs"
	"github.com/myshkin5/effective-octo-garbanzo/persistence"
	"github.com/myshkin5/effective-octo-garbanzo/services"
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

	database, err := persistence.Open()
	if err != nil {
		logs.Logger.Panic("Could not open database", err)
	}

	garbanzoService := services.NewGarbanzoService(persistence.GarbanzoStore{}, database)

	router := mux.NewRouter()

	middleware := alice.New()

	handlers.MapHealthRoutes(router, middleware)

	// TODO: Get proper base URL for absolute URIs
	baseURL := "./"

	handlers.MapGarbanzoCollectionRoutes(baseURL, router, middleware, garbanzoService)
	handlers.MapGarbanzoRoutes(baseURL, router, middleware, garbanzoService)

	loggingHandler := gorilla_handlers.LoggingHandler(os.Stdout, router)

	serverAddr := persistence.GetEnvWithDefault("SERVER_ADDR", "localhost")
	port := persistence.GetEnvWithDefault("PORT", "8080")
	logs.Logger.Infof("Listening on %s:%s...", serverAddr, port)
	err = http.ListenAndServe(serverAddr+":"+port, loggingHandler)
	if err != nil {
		logs.Logger.Critical("ListenAndServe:", err)
	}
}
