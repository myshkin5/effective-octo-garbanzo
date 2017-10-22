package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"

	"github.com/myshkin5/effective-octo-garbanzo/persistence/data"
)

type garbanzoCollection struct {
	garbanzoService GarbanzoService
	baseURL         string
}

func MapGarbanzoCollectionRoutes(baseURL string, router *mux.Router, middleware alice.Chain, garbanzoService GarbanzoService) {
	handler := &garbanzoCollection{
		garbanzoService: garbanzoService,
		baseURL:         baseURL + "garbanzos/",
	}
	methodHandler := make(MethodHandler)
	methodHandler[http.MethodGet] = http.HandlerFunc(handler.get)
	methodHandler[http.MethodPost] = http.HandlerFunc(handler.post)
	router.Handle("/garbanzos", middleware.Then(methodHandler))
}

func (g *garbanzoCollection) get(w http.ResponseWriter, req *http.Request) {
	garbanzos, err := g.garbanzoService.FetchAllGarbanzos(req.Context())
	if err != nil {
		Error(w, "Error fetching all garbanzos", http.StatusInternalServerError, err)
		return
	}

	list := []Garbanzo{}
	for _, garbanzo := range garbanzos {
		list = append(list, fromPersistence(garbanzo, g.baseURL))
	}

	Respond(w, http.StatusOK, JSONObject{
		"data": JSONObject{
			"garbanzos": list,
		},
	})
}

func (g *garbanzoCollection) post(w http.ResponseWriter, req *http.Request) {
	var dto Garbanzo
	err := json.NewDecoder(req.Body).Decode(&dto)
	if err != nil {
		Error(w, INVALID_JSON, http.StatusBadRequest, err)
		return
	}

	garbanzoType, err := data.GarbanzoTypeFromString(dto.GarbanzoType)
	if err != nil {
		Error(w, fmt.Sprintf("Unknown garbanzo type: %s", dto.GarbanzoType), http.StatusBadRequest, err)
		return
	}

	garbanzo, err := g.garbanzoService.CreateGarbanzo(req.Context(), data.Garbanzo{
		GarbanzoType: garbanzoType,
		DiameterMM:   dto.DiameterMM,
	})
	if err != nil {
		// TODO: Separate bad request issues from internal errors
		Error(w, "Error creating new garbanzo", http.StatusInternalServerError, err)
		return
	}

	Respond(w, http.StatusCreated, JSONObject{
		"data": JSONObject{
			"garbanzo": fromPersistence(garbanzo, g.baseURL),
		},
	})
}
