package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/myshkin5/effective-octo-garbanzo/logs"
	"github.com/satori/go.uuid"
)

//go:generate stringer -type=garbanzoType

type garbanzoType int

const (
	DESI   garbanzoType = 1001
	KABULI garbanzoType = 1002
)

func GarbanzoTypeFromString(gType string) (garbanzoType, error) {
	switch gType {
	case DESI.String():
		return DESI, nil
	case KABULI.String():
		return KABULI, nil
	default:
		return 0, fmt.Errorf("invalid garbanzo type: %s", gType)
	}
}

type Garbanzo struct {
	Id           int
	APIUUID      uuid.UUID
	GarbanzoType garbanzoType
	DiameterMM   float32
}

type GarbanzoStore struct{}

func (GarbanzoStore) FetchAllGarbanzos(ctx context.Context, database Database) ([]Garbanzo, error) {
	query := "select id, api_uuid, garbanzo_type_id, diameter_mm from garbanzo order by id"

	rows, err := database.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var garbanzos []Garbanzo
	for rows.Next() {
		var id int
		var apiUUID uuid.UUID
		var garbanzoTypeId garbanzoType
		var diameterMM float32
		err = rows.Scan(&id, &apiUUID, &garbanzoTypeId, &diameterMM)
		if err != nil {
			return nil, err
		}

		garbanzo := Garbanzo{
			Id:           id,
			APIUUID:      apiUUID,
			GarbanzoType: garbanzoTypeId,
			DiameterMM:   diameterMM,
		}
		garbanzos = append(garbanzos, garbanzo)
	}

	return garbanzos, nil
}

func (GarbanzoStore) FetchGarbanzoByAPIUUID(ctx context.Context, database Database, apiUUID uuid.UUID) (Garbanzo, error) {
	query := "select id, garbanzo_type_id, diameter_mm from garbanzo where api_uuid = $1"

	var id int
	var garbanzoTypeId garbanzoType
	var diameterMM float32
	err := database.QueryRowContext(ctx, query, apiUUID).Scan(&id, &garbanzoTypeId, &diameterMM)
	if err == sql.ErrNoRows {
		return Garbanzo{}, ErrNotFound
	} else if err != nil {
		return Garbanzo{}, err
	}

	return Garbanzo{
		Id:           id,
		APIUUID:      apiUUID,
		GarbanzoType: garbanzoTypeId,
		DiameterMM:   diameterMM,
	}, nil
}

func (GarbanzoStore) CreateGarbanzo(ctx context.Context, database Database, garbanzo Garbanzo) (int, error) {
	query := "insert into garbanzo (api_uuid, garbanzo_type_id, diameter_mm) values ($1, $2, $3) returning id"

	var garbanzoId int
	err := database.QueryRowContext(
		ctx,
		query,
		garbanzo.APIUUID,
		garbanzo.GarbanzoType,
		garbanzo.DiameterMM).Scan(&garbanzoId)
	if err != nil {
		return 0, err
	}

	return garbanzoId, nil
}

func (GarbanzoStore) DeleteGarbanzoByAPIUUID(ctx context.Context, database Database, apiUUID uuid.UUID) error {
	query := "delete from garbanzo where api_uuid = $1"

	result, err := database.ExecContext(ctx, query, apiUUID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	} else if rowsAffected > 1 {
		logs.Logger.Panic("Deleted multiple rows when expecting only one")
	}

	return nil
}
