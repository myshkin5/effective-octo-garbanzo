package octo

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"

	"github.com/myshkin5/effective-octo-garbanzo/api/handlers"
	"github.com/myshkin5/effective-octo-garbanzo/persistence"
	"github.com/myshkin5/effective-octo-garbanzo/persistence/data"
)

type Octo struct {
	Link      string `json:"link"`
	Name      string `json:"name"`
	Garbanzos string `json:"garbanzos"`
}

var fieldMapping = map[string]string{
	"Link":      "link",
	"Name":      "name",
	"Garbanzos": "garbanzos",
}

type OctoService interface {
	FetchAll(ctx context.Context) (octos []data.Octo, err error)
	FetchByName(ctx context.Context, name string) (octo data.Octo, err error)
	Create(ctx context.Context, octoIn data.Octo) (octoOut data.Octo, err error)
	DeleteByName(ctx context.Context, name string) (err error)
}

type octo struct {
	octoService OctoService
	baseURL     string
}

func MapRoutes(baseURL string, router *mux.Router, middleware alice.Chain, octoService OctoService) {
	handler := &octo{
		octoService: octoService,
		baseURL:     baseURL + "octos/",
	}
	methodHandler := make(handlers.MethodHandler)
	methodHandler[http.MethodGet] = http.HandlerFunc(handler.get)
	methodHandler[http.MethodDelete] = http.HandlerFunc(handler.delete)
	router.Handle("/octos/{name}", middleware.Then(methodHandler))
}

func (g *octo) get(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	name := vars["name"]

	octo, err := g.octoService.FetchByName(req.Context(), name)
	if err == persistence.ErrNotFound {
		handlers.Error(w, fmt.Sprintf("Octo %s not found", name), http.StatusNotFound, err, fieldMapping)
		return
	} else if err != nil {
		handlers.Error(w, "Error fetching octo", http.StatusInternalServerError, err, fieldMapping)
		return
	}

	handlers.Respond(w, http.StatusOK, fromPersistence(octo, g.baseURL))
}

func (g *octo) delete(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	name := vars["name"]

	err := g.octoService.DeleteByName(req.Context(), name)
	if err == persistence.ErrNotFound {
		handlers.Error(w, fmt.Sprintf("Octo %s not found", name), http.StatusNotFound, err, fieldMapping)
		return
	} else if err != nil {
		handlers.Error(w, "Error fetching octo", http.StatusInternalServerError, err, fieldMapping)
		return
	}

	handlers.Respond(w, http.StatusNoContent, nil)
}

func fromPersistence(octo data.Octo, baseURL string) Octo {
	link := baseURL + octo.Name
	return Octo{
		Link:      link,
		Name:      octo.Name,
		Garbanzos: link + "/garbanzos",
	}
}
