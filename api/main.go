package main

import (
	"fmt"
	"net/http"
	"os"

	gorilla_handlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"

	"github.com/myshkin5/effective-octo-garbanzo/api/handlers"
	"github.com/myshkin5/effective-octo-garbanzo/api/handlers/garbanzo"
	"github.com/myshkin5/effective-octo-garbanzo/api/handlers/octo"
	api_middleware "github.com/myshkin5/effective-octo-garbanzo/api/middleware"
	"github.com/myshkin5/effective-octo-garbanzo/logs"
	"github.com/myshkin5/effective-octo-garbanzo/persistence"
	"github.com/myshkin5/effective-octo-garbanzo/services"
)

func main() {
	initLogging()

	database := initDatabase()

	garbanzoService := services.NewGarbanzoService(persistence.OctoStore{}, persistence.GarbanzoStore{}, database)
	octoService := services.NewOctoService(persistence.OctoStore{}, persistence.GarbanzoStore{}, database)

	port := persistence.GetEnvWithDefault("PORT", "8080")
	router := initRoutes(port, octoService, garbanzoService)

	listenAndServe(port, router)
}

func initLogging() {
	err := logs.Init()
	if err != nil {
		panic(err)
	}
}

func initDatabase() persistence.Database {
	database, err := persistence.Open()
	if err != nil {
		logs.Logger.Panic("Could not open database", err)
	}

	err = persistence.Migrate()
	if err != nil {
		logs.Logger.Panic("Could not migrate database", err)
	}

	return database
}

type loggingWriter struct{}

func (w loggingWriter) Write(p []byte) (int, error) {
	n := len(p)
	logs.Logger.Info(string(p[:n-1]))
	return n, nil
}

func initRoutes(port string, octoService *services.OctoService, garbanzoService *services.GarbanzoService) *mux.Router {
	router := mux.NewRouter()

	loggingHandler := func(handler http.Handler) http.Handler {
		return gorilla_handlers.LoggingHandler(loggingWriter{}, handler)
	}
	headersHandler := api_middleware.StandardHeadersHandler

	middleware := alice.New(loggingHandler, headersHandler)

	handlers.MapHealthRoutes(router, middleware)

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = fmt.Sprintf("http://localhost:%v/", port)
	}

	octo.MapCollectionRoutes(baseURL, router, middleware, octoService)
	octo.MapRoutes(baseURL, router, middleware, octoService)

	garbanzo.MapCollectionRoutes(baseURL, router, middleware, garbanzoService)
	garbanzo.MapRoutes(baseURL, router, middleware, garbanzoService)

	// Must be last mapping
	handlers.MapCatchAllRoutes(baseURL, router, middleware)

	return router
}

func listenAndServe(port string, router *mux.Router) {
	serverAddr := persistence.GetEnvWithDefault("SERVER_ADDR", "localhost")
	logs.Logger.Infof("Listening on %s:%s...", serverAddr, port)
	err := http.ListenAndServe(serverAddr+":"+port, router)
	if err != nil {
		logs.Logger.Panic("ListenAndServe:", err)
	}
}
