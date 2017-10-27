package octo

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"

	"github.com/myshkin5/effective-octo-garbanzo/api/handlers"
	"github.com/myshkin5/effective-octo-garbanzo/persistence/data"
)

type octoCollection struct {
	octoService OctoService
	baseURL     string
}

func MapCollectionRoutes(baseURL string, router *mux.Router, middleware alice.Chain, octoService OctoService) {
	handler := &octoCollection{
		octoService: octoService,
		baseURL:     baseURL + "octos/",
	}
	methodHandler := make(handlers.MethodHandler)
	methodHandler[http.MethodGet] = http.HandlerFunc(handler.get)
	methodHandler[http.MethodPost] = http.HandlerFunc(handler.post)
	router.Handle("/octos", middleware.Then(methodHandler))
}

func (g *octoCollection) get(w http.ResponseWriter, req *http.Request) {
	octos, err := g.octoService.FetchAllOctos(req.Context())
	if err != nil {
		handlers.Error(w, "Error fetching all octos", http.StatusInternalServerError, err)
		return
	}

	list := []Octo{}
	for _, octo := range octos {
		list = append(list, fromPersistence(octo, g.baseURL))
	}

	handlers.Respond(w, http.StatusOK, list)
}

func (g *octoCollection) post(w http.ResponseWriter, req *http.Request) {
	var dto Octo
	err := json.NewDecoder(req.Body).Decode(&dto)
	if err != nil {
		handlers.Error(w, handlers.INVALID_JSON, http.StatusBadRequest, err)
		return
	}

	octo, err := g.octoService.CreateOcto(req.Context(), data.Octo{
		Name: dto.Name,
	})
	if err != nil {
		handlers.Error(w, "Error creating new octo", http.StatusInternalServerError, err)
		return
	}

	handlers.Respond(w, http.StatusCreated, fromPersistence(octo, g.baseURL))
}