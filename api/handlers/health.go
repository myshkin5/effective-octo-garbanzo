package handlers

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func MapHealthRoutes(router *mux.Router, middleware alice.Chain) {
	methodHandler := make(handlers.MethodHandler)
	methodHandler[http.MethodGet] = http.HandlerFunc(get)
	router.PathPrefix("/health").Subrouter().Handle("", middleware.Then(methodHandler))
}

func get(w http.ResponseWriter, _ *http.Request) {
	Respond(w, http.StatusOK, JSONObject{
		"health": "GOOD",
	})
}
