package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func MapCatchAllRoutes(router *mux.Router, middleware alice.Chain) {
	router.PathPrefix("/").Handler(middleware.ThenFunc(catchAll))
}

func catchAll(w http.ResponseWriter, _ *http.Request) {
	Error(w, "Not Found", http.StatusNotFound)
}
