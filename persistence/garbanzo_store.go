package persistence

import (
	"context"
	"database/sql"

	"github.com/satori/go.uuid"

	"github.com/myshkin5/effective-octo-garbanzo/logs"
	"github.com/myshkin5/effective-octo-garbanzo/persistence/data"
)

type GarbanzoStore struct{}

func (GarbanzoStore) FetchByOctoName(ctx context.Context, database Database, octoName string) ([]data.Garbanzo, error) {
	query := `select g.id, g.api_uuid, g.garbanzo_type_id, g.octo_id, g.diameter_mm from garbanzo g
		join octo o on g.octo_id = o.id
		join org on o.org_id = org.id
		where o.name = $1 and org.name = $2
		order by g.id`

	rows, err := database.Query(ctx, query, octoName, org(ctx))
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

func (GarbanzoStore) FetchByAPIUUIDAndOctoName(ctx context.Context, database Database, apiUUID uuid.UUID, octoName string) (data.Garbanzo, error) {
	query := `select g.id, g.garbanzo_type_id, g.octo_id, g.diameter_mm from garbanzo g
		join octo o on g.octo_id = o.id
		join org on o.org_id = org.id
		where g.api_uuid = $1 and o.name = $2 and org.name = $3`

	var id int
	var garbanzoType data.GarbanzoType
	var octoId int
	var diameterMM float32
	err := database.QueryRow(ctx, query, apiUUID, octoName, org(ctx)).Scan(&id, &garbanzoType, &octoId, &diameterMM)
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
	query := `insert into garbanzo (api_uuid, garbanzo_type_id, octo_id, diameter_mm)
		values (
			$1,
			$2,
			(select o.id from octo o join org on o.org_id = org.id where o.id = $3 and org.name = $4),
			$5) returning id`
	return ExecInsert(ctx, database, query, garbanzo.APIUUID, garbanzo.GarbanzoType, garbanzo.OctoId, org(ctx), garbanzo.DiameterMM)
}

func (GarbanzoStore) DeleteByAPIUUIDAndOctoName(ctx context.Context, database Database, apiUUID uuid.UUID, octoName string) error {
	query := `delete from garbanzo
		where api_uuid = $1 and octo_id = (
			select o.id from octo o
			join org on o.org_id = org.id
			where o.name = $2 and org.name = $3)`
	rowsAffected, err := ExecDelete(ctx, database, query, apiUUID, octoName, org(ctx))
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

func (GarbanzoStore) DeleteByOctoId(ctx context.Context, database Database, octoId int) error {
	query := `delete from garbanzo where octo_id = (
			select o.id from octo o
			join org on o.org_id = org.id
			where o.id = $1 and org.name = $2)`
	_, err := ExecDelete(ctx, database, query, octoId, org(ctx))
	if err != nil {
		return err
	}

	return nil
}
