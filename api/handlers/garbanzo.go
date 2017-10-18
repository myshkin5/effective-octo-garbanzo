package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/satori/go.uuid"

	"github.com/myshkin5/effective-octo-garbanzo/persistence"
)

type Garbanzo struct {
	Link      string `json:"link"`
	FirstName string `json:"first-name"`
	LastName  string `json:"last-name"`
}

type GarbanzoService interface {
	FetchAllGarbanzos(ctx context.Context) (garbanzos []persistence.Garbanzo, err error)
	FetchGarbanzoByAPIUUID(ctx context.Context, apiUUID uuid.UUID) (garbanzo persistence.Garbanzo, err error)
	CreateGarbanzo(ctx context.Context, garbanzoIn persistence.Garbanzo) (garbanzoOut persistence.Garbanzo, err error)
	DeleteGarbanzoByAPIUUID(ctx context.Context, apiUUID uuid.UUID) (err error)
}

type garbanzo struct {
	garbanzoService GarbanzoService
	baseURL         string
}

func MapGarbanzoRoutes(baseURL string, router *mux.Router, middleware alice.Chain, garbanzoService GarbanzoService) {
	handler := &garbanzo{
		garbanzoService: garbanzoService,
		baseURL:         baseURL + "garbanzos/",
	}
	methodHandler := make(MethodHandler)
	methodHandler[http.MethodGet] = http.HandlerFunc(handler.get)
	methodHandler[http.MethodDelete] = http.HandlerFunc(handler.delete)
	router.Handle("/garbanzos/{apiUUID}", middleware.Then(methodHandler))
}

func (g *garbanzo) get(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	apiUUID, err := uuid.FromString(vars["apiUUID"])
	if err != nil {
		Error(w, INVALID_UUID, http.StatusBadRequest, err)
		return
	}

	garbanzo, err := g.garbanzoService.FetchGarbanzoByAPIUUID(req.Context(), apiUUID)
	if err == persistence.ErrNotFound {
		Error(w, fmt.Sprintf("Garbanzo %s not found", apiUUID), http.StatusNotFound, err)
		return
	} else if err != nil {
		Error(w, "Error fetching garbanzo", http.StatusInternalServerError, err)
		return
	}

	Respond(w, http.StatusOK, JSONObject{
		"data": JSONObject{
			"garbanzo": fromPersistence(garbanzo, g.baseURL),
		},
	})
}

func (g *garbanzo) delete(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	apiUUID, err := uuid.FromString(vars["apiUUID"])
	if err != nil {
		Error(w, INVALID_UUID, http.StatusBadRequest, err)
		return
	}

	err = g.garbanzoService.DeleteGarbanzoByAPIUUID(req.Context(), apiUUID)
	if err == persistence.ErrNotFound {
		Error(w, fmt.Sprintf("Garbanzo %s not found", apiUUID), http.StatusNotFound, err)
		return
	} else if err != nil {
		Error(w, "Error fetching garbanzo", http.StatusInternalServerError, err)
		return
	}

	Respond(w, http.StatusNoContent, nil)
}

func fromPersistence(garbanzo persistence.Garbanzo, baseURL string) Garbanzo {
	return Garbanzo{
		Link:      baseURL + garbanzo.APIUUID.String(),
		FirstName: garbanzo.FirstName,
		LastName:  garbanzo.LastName,
	}
}