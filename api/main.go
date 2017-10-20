package main

import (
	"net/http"
	"os"

	gorilla_handlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"

	"fmt"

	"github.com/myshkin5/effective-octo-garbanzo/api/handlers"
	"github.com/myshkin5/effective-octo-garbanzo/api/handlers/middleware"
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

	headersHandler := middleware.StandardHeadersHandler

	middleware := alice.New(headersHandler)

	handlers.MapHealthRoutes(router, middleware)

	port := persistence.GetEnvWithDefault("PORT", "8080")
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = fmt.Sprintf("http://localhost:%v/", port)
	}

	handlers.MapGarbanzoCollectionRoutes(baseURL, router, middleware, garbanzoService)
	handlers.MapGarbanzoRoutes(baseURL, router, middleware, garbanzoService)

	// Must be last mapping
	handlers.MapCatchAllRoutes(baseURL, router, middleware)

	loggingHandler := gorilla_handlers.LoggingHandler(os.Stdout, router)

	serverAddr := persistence.GetEnvWithDefault("SERVER_ADDR", "localhost")
	logs.Logger.Infof("Listening on %s:%s...", serverAddr, port)
	err = http.ListenAndServe(serverAddr+":"+port, loggingHandler)
	if err != nil {
		logs.Logger.Critical("ListenAndServe:", err)
	}
}
