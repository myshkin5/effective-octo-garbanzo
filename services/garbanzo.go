package services

import (
	"context"

	"github.com/myshkin5/effective-octo-garbanzo/persistence"
	"github.com/satori/go.uuid"
)

type GarbanzoStore interface {
	FetchAllGarbanzos(ctx context.Context, database persistence.Database) (garbanzos []persistence.Garbanzo, err error)
	FetchGarbanzoByAPIUUID(ctx context.Context, database persistence.Database, apiUUID uuid.UUID) (garbanzo persistence.Garbanzo, err error)
	CreateGarbanzo(ctx context.Context, database persistence.Database, garbanzo persistence.Garbanzo) (garbanzoId int, err error)
	DeleteGarbanzoByAPIUUID(ctx context.Context, database persistence.Database, apiUUID uuid.UUID) (err error)
}

type GarbanzoService struct {
	store    GarbanzoStore
	database persistence.Database
}

func NewGarbanzoService(store GarbanzoStore, database persistence.Database) *GarbanzoService {
	return &GarbanzoService{
		store:    store,
		database: database,
	}
}

func (s *GarbanzoService) FetchAllGarbanzos(ctx context.Context) ([]persistence.Garbanzo, error) {
	return s.store.FetchAllGarbanzos(ctx, s.database)
}

func (s *GarbanzoService) FetchGarbanzoByAPIUUID(ctx context.Context, apiUUID uuid.UUID) (persistence.Garbanzo, error) {
	return s.store.FetchGarbanzoByAPIUUID(ctx, s.database, apiUUID)
}

func (s *GarbanzoService) CreateGarbanzo(ctx context.Context, garbanzo persistence.Garbanzo) (int, error) {
	return s.store.CreateGarbanzo(ctx, s.database, garbanzo)
}

func (s *GarbanzoService) DeleteGarbanzoByAPIUUID(ctx context.Context, apiUUID uuid.UUID) error {
	return s.store.DeleteGarbanzoByAPIUUID(ctx, s.database, apiUUID)
}
