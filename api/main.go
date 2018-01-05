package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/myshkin5/effective-octo-garbanzo/identity"

	"github.com/myshkin5/effective-octo-garbanzo/api/handlers"
	"github.com/myshkin5/effective-octo-garbanzo/api/handlers/garbanzo"
	"github.com/myshkin5/effective-octo-garbanzo/api/handlers/octo"
	apiMiddleware "github.com/myshkin5/effective-octo-garbanzo/api/middleware"
	"github.com/myshkin5/effective-octo-garbanzo/logs"
	"github.com/myshkin5/effective-octo-garbanzo/persistence"
	"github.com/myshkin5/effective-octo-garbanzo/services"
	"github.com/myshkin5/effective-octo-garbanzo/utils"
)

func main() {
	initLogging()

	utils.InitStackTracer()

	database := initDatabase()

	garbanzoService := services.NewGarbanzoService(persistence.OctoStore{}, persistence.GarbanzoStore{}, database)
	octoService := services.NewOctoService(persistence.OctoStore{}, persistence.GarbanzoStore{}, database)

	port := persistence.GetEnvWithDefault("PORT", "8080")
	router := initRoutes(port, octoService, garbanzoService)

	serverAddr := persistence.GetEnvWithDefault("SERVER_ADDR", "localhost")

	initPProf(serverAddr)

	listenAndServe(serverAddr, port, router)
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
		logs.Logger.Panic("Could not open database: ", err)
	}

	err = persistence.Migrate()
	if err != nil {
		logs.Logger.Panic("Could not migrate database: ", err)
	}

	return database
}

func initRoutes(port string, octoService *services.OctoService, garbanzoService *services.GarbanzoService) *mux.Router {
	router := mux.NewRouter()

	headersHandler := apiMiddleware.StandardHeadersHandler

	client := &http.Client{}
	verifierKeyInsecure := os.Getenv("VERIFIER_KEY_INSECURE")
	if verifierKeyInsecure == "true" {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	verifierKeyURI := os.Getenv("VERIFIER_KEY_URI")
	publicKeys := identity.MustFetchKeys(verifierKeyURI, client)
	validator := identity.NewValidator(publicKeys)

	loginURI := os.Getenv("LOGIN_URI")
	authHandler := func(h http.Handler) http.Handler {
		return apiMiddleware.AuthenticatedHandler(h, loginURI, validator)
	}

	middleware := alice.New(handlers.LoggingHandler, headersHandler, authHandler)

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

func initPProf(serverAddr string) {
	// Typically 6060
	pprofPort, ok := os.LookupEnv("PPROF_PORT")
	if ok {
		logs.Logger.Infof("PProf listening on %s:%s...", serverAddr, pprofPort)
		go func() {
			logs.Logger.Info(http.ListenAndServe(serverAddr+":"+pprofPort, nil))
		}()
	}
}

func listenAndServe(serverAddr, port string, router *mux.Router) {
	logs.Logger.Infof("Listening on %s:%s...", serverAddr, port)
	err := http.ListenAndServe(serverAddr+":"+port, router)
	if err != nil {
		logs.Logger.Panic("ListenAndServe: ", err)
	}
}
