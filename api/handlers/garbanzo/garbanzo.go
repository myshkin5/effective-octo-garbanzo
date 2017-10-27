package garbanzo

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/satori/go.uuid"

	"github.com/myshkin5/effective-octo-garbanzo/api/handlers"
	"github.com/myshkin5/effective-octo-garbanzo/persistence"
	"github.com/myshkin5/effective-octo-garbanzo/persistence/data"
)

type Garbanzo struct {
	Link         string  `json:"link"`
	GarbanzoType string  `json:"type"`
	DiameterMM   float32 `json:"diameter-mm"`
}

type GarbanzoService interface {
	FetchAllGarbanzos(ctx context.Context) (garbanzos []data.Garbanzo, err error)
	FetchGarbanzoByAPIUUID(ctx context.Context, apiUUID uuid.UUID) (garbanzo data.Garbanzo, err error)
	CreateGarbanzo(ctx context.Context, octoName string, garbanzoIn data.Garbanzo) (garbanzoOut data.Garbanzo, err error)
	DeleteGarbanzoByAPIUUID(ctx context.Context, apiUUID uuid.UUID) (err error)
}

type garbanzo struct {
	garbanzoService GarbanzoService
	baseURL         string
}

func MapRoutes(baseURL string, router *mux.Router, middleware alice.Chain, garbanzoService GarbanzoService) {
	handler := &garbanzo{
		garbanzoService: garbanzoService,
		baseURL:         baseURL,
	}
	methodHandler := make(handlers.MethodHandler)
	methodHandler[http.MethodGet] = http.HandlerFunc(handler.get)
	methodHandler[http.MethodDelete] = http.HandlerFunc(handler.delete)
	router.Handle("/octos/{octoName}/garbanzos/{apiUUID}", middleware.Then(methodHandler))
}

func (g *garbanzo) get(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	apiUUID, err := uuid.FromString(vars["apiUUID"])
	if err != nil {
		handlers.Error(w, handlers.INVALID_UUID, http.StatusBadRequest, err)
		return
	}

	garbanzo, err := g.garbanzoService.FetchGarbanzoByAPIUUID(req.Context(), apiUUID)
	if err == persistence.ErrNotFound {
		handlers.Error(w, fmt.Sprintf("Garbanzo %s not found", apiUUID), http.StatusNotFound, err)
		return
	} else if err != nil {
		handlers.Error(w, "Error fetching garbanzo", http.StatusInternalServerError, err)
		return
	}

	handlers.Respond(w, http.StatusOK, fromPersistence(garbanzo, g.baseURL, vars["octoName"]))
}

func (g *garbanzo) delete(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	apiUUID, err := uuid.FromString(vars["apiUUID"])
	if err != nil {
		handlers.Error(w, handlers.INVALID_UUID, http.StatusBadRequest, err)
		return
	}

	err = g.garbanzoService.DeleteGarbanzoByAPIUUID(req.Context(), apiUUID)
	if err == persistence.ErrNotFound {
		handlers.Error(w, fmt.Sprintf("Garbanzo %s not found", apiUUID), http.StatusNotFound, err)
		return
	} else if err != nil {
		handlers.Error(w, "Error fetching garbanzo", http.StatusInternalServerError, err)
		return
	}

	handlers.Respond(w, http.StatusNoContent, nil)
}

func fromPersistence(garbanzo data.Garbanzo, baseURL, octoName string) Garbanzo {
	return Garbanzo{
		Link:         fmt.Sprintf("%soctos/%s/garbanzos/%s", baseURL, octoName, garbanzo.APIUUID.String()),
		GarbanzoType: garbanzo.GarbanzoType.String(),
		DiameterMM:   garbanzo.DiameterMM,
	}
}
