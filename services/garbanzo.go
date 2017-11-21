package services

import (
	"context"

	"github.com/satori/go.uuid"

	"github.com/myshkin5/effective-octo-garbanzo/persistence"
	"github.com/myshkin5/effective-octo-garbanzo/persistence/data"
)

type GarbanzoStore interface {
	FetchAll(ctx context.Context, database persistence.Database) (garbanzos []data.Garbanzo, err error)
	FetchByAPIUUID(ctx context.Context, database persistence.Database, apiUUID uuid.UUID) (garbanzo data.Garbanzo, err error)
	Create(ctx context.Context, database persistence.Database, garbanzo data.Garbanzo) (garbanzoId int, err error)
	DeleteByAPIUUID(ctx context.Context, database persistence.Database, apiUUID uuid.UUID) (err error)
	DeleteByOctoName(ctx context.Context, database persistence.Database, octoName string) (err error)
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

func (s *GarbanzoService) FetchAll(ctx context.Context) ([]data.Garbanzo, error) {
	return s.garbanzoStore.FetchAll(ctx, s.database)
}

func (s *GarbanzoService) FetchByAPIUUID(ctx context.Context, apiUUID uuid.UUID) (data.Garbanzo, error) {
	return s.garbanzoStore.FetchByAPIUUID(ctx, s.database, apiUUID)
}

func (s *GarbanzoService) Create(ctx context.Context, octoName string, garbanzo data.Garbanzo) (garbanzoOut data.Garbanzo, err error) {
	err = validate(garbanzo)
	if err != nil {
		return data.Garbanzo{}, err
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

	octo, err := s.octoStore.FetchByName(ctx, database, octoName)
	if err != nil {
		return data.Garbanzo{}, err
	}
	garbanzo.OctoId = octo.Id

	garbanzo.Id, err = s.garbanzoStore.Create(ctx, database, garbanzo)
	if err != nil {
		return data.Garbanzo{}, err
	}

	return garbanzo, nil
}

func validate(garbanzo data.Garbanzo) error {
	errors := make(map[string][]string)
	if garbanzo.GarbanzoType == 0 {
		errors["GarbanzoType"] = append(errors["GarbanzoType"], "must be present")
	}
	if garbanzo.GarbanzoType != data.DESI && garbanzo.GarbanzoType != data.KABULI {
		errors["GarbanzoType"] = append(errors["GarbanzoType"], "must be either 'DESI' or 'KABULI'")
	}
	if garbanzo.DiameterMM == 0.0 {
		errors["DiameterMM"] = append(errors["DiameterMM"], "must be present")
	}
	if garbanzo.DiameterMM <= 0.0 {
		errors["DiameterMM"] = append(errors["DiameterMM"], "must be a positive decimal value")
	}

	if len(errors) > 0 {
		return NewValidationError(errors)
	}

	return nil
}

func (s *GarbanzoService) DeleteByAPIUUID(ctx context.Context, apiUUID uuid.UUID) error {
	return s.garbanzoStore.DeleteByAPIUUID(ctx, s.database, apiUUID)
}
