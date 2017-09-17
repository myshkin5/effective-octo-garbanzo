package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

type health struct {
	subroute *mux.Route
}

func MapHealthRoutes(subroute *mux.Route, middleware alice.Chain) {
	handler := &health{
		subroute: subroute,
	}
	methodHandler := make(handlers.MethodHandler)
	methodHandler["GET"] = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		handler.get(w, req)
	})
	subroute.Handler(middleware.Then(methodHandler))
}

func (h *health) get(w http.ResponseWriter, req *http.Request) {
	health := map[string]interface{}{
		"health": "GOOD",
	}
	bytes, _ := json.Marshal(health)

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
