package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func MapCatchAllRoutes(baseURL string, router *mux.Router, middleware alice.Chain) {
	router.PathPrefix("/").Handler(middleware.ThenFunc(catchAll(baseURL)))
}

func catchAll(baseURL string) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet || req.RequestURI != "/" {
			Error(w, "Not Found", http.StatusNotFound, nil, nil)
			return
		}

		Respond(w, http.StatusOK, JSONObject{
			"health": baseURL + "health",
			"octos":  baseURL + "octos",
		})
	}
}
