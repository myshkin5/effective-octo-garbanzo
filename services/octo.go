package services

import (
	"context"

	"github.com/myshkin5/effective-octo-garbanzo/persistence"
	"github.com/myshkin5/effective-octo-garbanzo/persistence/data"
)

type OctoStore interface {
	FetchAllOctos(ctx context.Context, database persistence.Database) (octos []data.Octo, err error)
	FetchOctoByName(ctx context.Context, database persistence.Database, name string) (octo data.Octo, err error)
	CreateOcto(ctx context.Context, database persistence.Database, octo data.Octo) (octoId int, err error)
	DeleteOctoByName(ctx context.Context, database persistence.Database, name string) (err error)
}

type OctoService struct {
	store    OctoStore
	database persistence.Database
}

func NewOctoService(store OctoStore, database persistence.Database) *OctoService {
	return &OctoService{
		store:    store,
		database: database,
	}
}

func (s *OctoService) FetchAllOctos(ctx context.Context) ([]data.Octo, error) {
	return s.store.FetchAllOctos(ctx, s.database)
}

func (s *OctoService) FetchOctoByName(ctx context.Context, name string) (data.Octo, error) {
	return s.store.FetchOctoByName(ctx, s.database, name)
}

func (s *OctoService) CreateOcto(ctx context.Context, octo data.Octo) (data.Octo, error) {
	if len(octo.Name) == 0 {
		return data.Octo{}, NewValidationError("'name' is required")
	}

	var err error
	octo.Id, err = s.store.CreateOcto(ctx, s.database, octo)
	if err != nil {
		return data.Octo{}, err
	}

	return octo, nil
}

func (s *OctoService) DeleteOctoByName(ctx context.Context, name string) error {
	return s.store.DeleteOctoByName(ctx, s.database, name)
}