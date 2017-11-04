package garbanzo

import (
	"encoding/json"
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
	garbanzos, err := g.garbanzoService.FetchAllGarbanzos(req.Context())
	if err != nil {
		handlers.Error(w, "Error fetching all garbanzos", http.StatusInternalServerError, err, fieldMapping)
		return
	}

	list := []Garbanzo{}
	for _, garbanzo := range garbanzos {
		list = append(list, fromPersistence(garbanzo, g.baseURL, mux.Vars(req)["octoName"]))
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

	garbanzoType, _ /* TODO */ := data.GarbanzoTypeFromString(dto.GarbanzoType)

	octoName := mux.Vars(req)["octoName"]
	garbanzo, err := g.garbanzoService.CreateGarbanzo(req.Context(), octoName, data.Garbanzo{
		GarbanzoType: garbanzoType,
		DiameterMM:   dto.DiameterMM,
	})
	if err == persistence.ErrNotFound {
		// TODO Use ValidationError but with conflict flag
		handlers.Error(w, "Error creating new garbanzo", http.StatusConflict, err, fieldMapping)
		return
	} else if err != nil {
		handlers.Error(w, "Error creating new garbanzo", http.StatusInternalServerError, err, fieldMapping)
		return
	}

	handlers.Respond(w, http.StatusCreated, fromPersistence(garbanzo, g.baseURL, octoName))
}
