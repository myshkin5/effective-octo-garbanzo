package garbanzo

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"

	"github.com/myshkin5/effective-octo-garbanzo/api/handlers"
	"github.com/myshkin5/effective-octo-garbanzo/persistence"
	"github.com/myshkin5/effective-octo-garbanzo/persistence/data"
)

type garbanzoCollection struct {
	garbanzoService GarbanzoService
	baseURL         string
}

func MapCollectionRoutes(baseURL string, router *mux.Router, middleware alice.Chain, garbanzoService GarbanzoService) {
	handler := &garbanzoCollection{
		garbanzoService: garbanzoService,
		baseURL:         baseURL,
	}
	methodHandler := make(handlers.MethodHandler)
	methodHandler[http.MethodGet] = http.HandlerFunc(handler.get)
	methodHandler[http.MethodPost] = http.HandlerFunc(handler.post)
	router.Handle("/octos/{octoName}/garbanzos", middleware.Then(methodHandler))
}

func (g *garbanzoCollection) get(w http.ResponseWriter, req *http.Request) {
	octoName := mux.Vars(req)["octoName"]
	garbanzos, err := g.garbanzoService.FetchByOctoName(req.Context(), octoName)
	if err != nil {
		handlers.Error(w, "Error fetching garbanzos", http.StatusInternalServerError, err, fieldMapping)
		return
	}

	list := []Garbanzo{}
	for _, garbanzo := range garbanzos {
		list = append(list, fromPersistence(garbanzo, g.baseURL, octoName))
	}

	handlers.Respond(w, http.StatusOK, list)
}

func (g *garbanzoCollection) post(w http.ResponseWriter, req *http.Request) {
	var dto Garbanzo
	err := json.NewDecoder(req.Body).Decode(&dto)
	if err != nil {
		handlers.Error(w, handlers.INVALID_JSON, http.StatusBadRequest, err, fieldMapping)
		return
	}

	garbanzoType, err := data.GarbanzoTypeFromString(dto.GarbanzoType)
	if err != nil {
		handlers.Error(w, err.Error(), http.StatusBadRequest, err, fieldMapping)
		return
	}

	octoName := mux.Vars(req)["octoName"]
	garbanzo, err := g.garbanzoService.Create(req.Context(), octoName, data.Garbanzo{
		GarbanzoType: garbanzoType,
		DiameterMM:   dto.DiameterMM,
	})
	if err == persistence.ErrNotFound {
		handlers.Error(w, fmt.Sprintf("Parent octo '%s' not found", octoName), http.StatusConflict, err, fieldMapping)
		return
	} else if err != nil {
		handlers.Error(w, "Error creating new garbanzo", http.StatusInternalServerError, err, fieldMapping)
		return
	}

	handlers.Respond(w, http.StatusCreated, fromPersistence(garbanzo, g.baseURL, octoName))
}
