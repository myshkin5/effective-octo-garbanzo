package persistence

import (
	"context"
	"database/sql"

	"github.com/satori/go.uuid"

	"github.com/myshkin5/effective-octo-garbanzo/persistence/data"
)

type GarbanzoStore struct{}

func (GarbanzoStore) FetchAllGarbanzos(ctx context.Context, database Database) ([]data.Garbanzo, error) {
	query := "select id, api_uuid, garbanzo_type_id, diameter_mm from garbanzo order by id"

	rows, err := database.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var garbanzos []data.Garbanzo
	for rows.Next() {
		var id int
		var apiUUID uuid.UUID
		var garbanzoTypeId data.GarbanzoType
		var diameterMM float32
		err = rows.Scan(&id, &apiUUID, &garbanzoTypeId, &diameterMM)
		if err != nil {
			return nil, err
		}

		garbanzo := data.Garbanzo{
			Id:           id,
			APIUUID:      apiUUID,
			GarbanzoType: garbanzoTypeId,
			DiameterMM:   diameterMM,
		}
		garbanzos = append(garbanzos, garbanzo)
	}

	return garbanzos, nil
}

func (GarbanzoStore) FetchGarbanzoByAPIUUID(ctx context.Context, database Database, apiUUID uuid.UUID) (data.Garbanzo, error) {
	query := "select id, garbanzo_type_id, diameter_mm from garbanzo where api_uuid = $1"

	var id int
	var garbanzoTypeId data.GarbanzoType
	var diameterMM float32
	err := database.QueryRowContext(ctx, query, apiUUID).Scan(&id, &garbanzoTypeId, &diameterMM)
	if err == sql.ErrNoRows {
		return data.Garbanzo{}, ErrNotFound
	} else if err != nil {
		return data.Garbanzo{}, err
	}

	return data.Garbanzo{
		Id:           id,
		APIUUID:      apiUUID,
		GarbanzoType: garbanzoTypeId,
		DiameterMM:   diameterMM,
	}, nil
}

func (GarbanzoStore) CreateGarbanzo(ctx context.Context, database Database, garbanzo data.Garbanzo) (int, error) {
	query := "insert into garbanzo (api_uuid, garbanzo_type_id, diameter_mm) values ($1, $2, $3) returning id"
	return ExecInsert(ctx, database, query, garbanzo.APIUUID, garbanzo.GarbanzoType, garbanzo.DiameterMM)
}

func (GarbanzoStore) DeleteGarbanzoByAPIUUID(ctx context.Context, database Database, apiUUID uuid.UUID) error {
	query := "delete from garbanzo where api_uuid = $1"
	return ExecDelete(ctx, database, query, apiUUID)
}
