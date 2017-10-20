package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func MapHealthRoutes(router *mux.Router, middleware alice.Chain) {
	methodHandler := make(MethodHandler)
	methodHandler[http.MethodGet] = http.HandlerFunc(getHealth)
	router.PathPrefix("/health").Handler(middleware.Then(methodHandler))
}

func getHealth(w http.ResponseWriter, _ *http.Request) {
	Respond(w, http.StatusOK, JSONObject{
		"health": "GOOD",
	})
}
