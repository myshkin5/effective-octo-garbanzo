package services

import (
	"context"

	"github.com/satori/go.uuid"

	"github.com/myshkin5/effective-octo-garbanzo/persistence"
	"github.com/myshkin5/effective-octo-garbanzo/persistence/data"
)

type GarbanzoStore interface {
	FetchByOctoName(ctx context.Context, database persistence.Database, octoName string) (garbanzos []data.Garbanzo, err error)
	FetchByAPIUUIDAndOctoName(ctx context.Context, database persistence.Database, apiUUID uuid.UUID, octoName string) (garbanzo data.Garbanzo, err error)
	Create(ctx context.Context, database persistence.Database, garbanzo data.Garbanzo) (garbanzoId int, err error)
	DeleteByAPIUUIDAndOctoName(ctx context.Context, database persistence.Database, apiUUID uuid.UUID, octoName string) (err error)
	DeleteByOctoId(ctx context.Context, database persistence.Database, octoId int) (err error)
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

func (s *GarbanzoService) FetchByOctoName(ctx context.Context, octoName string) ([]data.Garbanzo, error) {
	return s.garbanzoStore.FetchByOctoName(ctx, s.database, octoName)
}

func (s *GarbanzoService) FetchByAPIUUIDAndOctoName(ctx context.Context, apiUUID uuid.UUID, octoName string) (data.Garbanzo, error) {
	return s.garbanzoStore.FetchByAPIUUIDAndOctoName(ctx, s.database, apiUUID, octoName)
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

	octo, err := s.octoStore.FetchByName(ctx, database, octoName, true)
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

func (s *GarbanzoService) DeleteByAPIUUIDAndOctoName(ctx context.Context, apiUUID uuid.UUID, octoName string) error {
	return s.garbanzoStore.DeleteByAPIUUIDAndOctoName(ctx, s.database, apiUUID, octoName)
}
