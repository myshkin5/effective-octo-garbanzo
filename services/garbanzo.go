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
	octoStore     OctoStore
	garbanzoStore GarbanzoStore
	database      persistence.Database
}

func NewGarbanzoService(octoStore OctoStore, garbanzoStore GarbanzoStore, database persistence.Database) *GarbanzoService {
	return &GarbanzoService{
		octoStore:     octoStore,
		garbanzoStore: garbanzoStore,
		database:      database,
	}
}

func (s *GarbanzoService) FetchAllGarbanzos(ctx context.Context) ([]data.Garbanzo, error) {
	return s.garbanzoStore.FetchAllGarbanzos(ctx, s.database)
}

func (s *GarbanzoService) FetchGarbanzoByAPIUUID(ctx context.Context, apiUUID uuid.UUID) (data.Garbanzo, error) {
	return s.garbanzoStore.FetchGarbanzoByAPIUUID(ctx, s.database, apiUUID)
}

func (s *GarbanzoService) CreateGarbanzo(ctx context.Context, octoName string, garbanzo data.Garbanzo) (garbanzoOut data.Garbanzo, err error) {
	errors := validate(garbanzo)
	if len(errors) > 0 {
		return data.Garbanzo{}, NewValidationError(errors...)
	}

	garbanzo.APIUUID = uuid.NewV4()

	database, err := s.database.BeginTx(ctx)
	if err != nil {
		return data.Garbanzo{}, err
	}
	defer func() {
		if err != nil {
			database.Rollback()
			return
		}
		err = database.Commit()
	}()

	octo, err := s.octoStore.FetchOctoByName(ctx, database, octoName)
	if err != nil {
		return data.Garbanzo{}, err
	}
	garbanzo.OctoId = octo.Id

	garbanzo.Id, err = s.garbanzoStore.CreateGarbanzo(ctx, database, garbanzo)
	if err != nil {
		return data.Garbanzo{}, err
	}

	return garbanzo, nil
}

func validate(garbanzo data.Garbanzo) []string {
	var errors []string
	if garbanzo.GarbanzoType != data.DESI && garbanzo.GarbanzoType != data.KABULI {
		errors = append(errors, "'type' is required and must be either 'DESI' or 'KABULI'")
	}
	if garbanzo.DiameterMM <= 0.0 {
		errors = append(errors, "'diameter-mm' is required to be a positive decimal value")
	}

	return errors
}

func (s *GarbanzoService) DeleteGarbanzoByAPIUUID(ctx context.Context, apiUUID uuid.UUID) error {
	return s.garbanzoStore.DeleteGarbanzoByAPIUUID(ctx, s.database, apiUUID)
}
