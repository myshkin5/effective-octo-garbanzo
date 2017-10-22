package services

import (
	"context"

	"github.com/satori/go.uuid"

	"github.com/myshkin5/effective-octo-garbanzo/persistence"
	"github.com/myshkin5/effective-octo-garbanzo/persistence/data"
)

type GarbanzoStore interface {
	FetchAllGarbanzos(ctx context.Context, database persistence.Database) (garbanzos []data.Garbanzo, err error)
	FetchGarbanzoByAPIUUID(ctx context.Context, database persistence.Database, apiUUID uuid.UUID) (garbanzo data.Garbanzo, err error)
	CreateGarbanzo(ctx context.Context, database persistence.Database, garbanzo data.Garbanzo) (garbanzoId int, err error)
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

func (s *GarbanzoService) FetchAllGarbanzos(ctx context.Context) ([]data.Garbanzo, error) {
	return s.store.FetchAllGarbanzos(ctx, s.database)
}

func (s *GarbanzoService) FetchGarbanzoByAPIUUID(ctx context.Context, apiUUID uuid.UUID) (data.Garbanzo, error) {
	return s.store.FetchGarbanzoByAPIUUID(ctx, s.database, apiUUID)
}

func (s *GarbanzoService) CreateGarbanzo(ctx context.Context, garbanzo data.Garbanzo) (data.Garbanzo, error) {
	garbanzo.APIUUID = uuid.NewV4()

	var err error
	garbanzo.Id, err = s.store.CreateGarbanzo(ctx, s.database, garbanzo)
	if err != nil {
		return data.Garbanzo{}, err
	}

	return garbanzo, nil
}

func (s *GarbanzoService) DeleteGarbanzoByAPIUUID(ctx context.Context, apiUUID uuid.UUID) error {
	return s.store.DeleteGarbanzoByAPIUUID(ctx, s.database, apiUUID)
}
