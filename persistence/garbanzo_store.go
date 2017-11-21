package persistence

import (
	"context"
	"database/sql"

	"github.com/satori/go.uuid"

	"github.com/myshkin5/effective-octo-garbanzo/persistence/data"
)

type GarbanzoStore struct{}

func (GarbanzoStore) FetchAll(ctx context.Context, database Database) ([]data.Garbanzo, error) {
	query := "select id, api_uuid, garbanzo_type_id, octo_id, diameter_mm from garbanzo order by id"

	rows, err := database.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var garbanzos []data.Garbanzo
	for rows.Next() {
		var id int
		var apiUUID uuid.UUID
		var garbanzoType data.GarbanzoType
		var octoId int
		var diameterMM float32
		err = rows.Scan(&id, &apiUUID, &garbanzoType, &octoId, &diameterMM)
		if err != nil {
			return nil, err
		}

		garbanzo := data.Garbanzo{
			Id:           id,
			APIUUID:      apiUUID,
			GarbanzoType: garbanzoType,
			OctoId:       octoId,
			DiameterMM:   diameterMM,
		}
		garbanzos = append(garbanzos, garbanzo)
	}

	return garbanzos, nil
}

func (GarbanzoStore) FetchByAPIUUID(ctx context.Context, database Database, apiUUID uuid.UUID) (data.Garbanzo, error) {
	query := "select id, garbanzo_type_id, octo_id, diameter_mm from garbanzo where api_uuid = $1"

	var id int
	var garbanzoType data.GarbanzoType
	var octoId int
	var diameterMM float32
	err := database.QueryRow(ctx, query, apiUUID).Scan(&id, &garbanzoType, &octoId, &diameterMM)
	if err == sql.ErrNoRows {
		return data.Garbanzo{}, ErrNotFound
	} else if err != nil {
		return data.Garbanzo{}, err
	}

	return data.Garbanzo{
		Id:           id,
		APIUUID:      apiUUID,
		GarbanzoType: garbanzoType,
		OctoId:       octoId,
		DiameterMM:   diameterMM,
	}, nil
}

func (GarbanzoStore) Create(ctx context.Context, database Database, garbanzo data.Garbanzo) (int, error) {
	query := "insert into garbanzo (api_uuid, garbanzo_type_id, octo_id, diameter_mm) " +
		"values ($1, $2, $3, $4) returning id"
	return ExecInsert(ctx, database, query, garbanzo.APIUUID, garbanzo.GarbanzoType, garbanzo.OctoId, garbanzo.DiameterMM)
}

func (GarbanzoStore) DeleteByAPIUUID(ctx context.Context, database Database, apiUUID uuid.UUID) error {
	query := "delete from garbanzo where api_uuid = $1"
	return ExecDelete(ctx, database, query, apiUUID)
}
