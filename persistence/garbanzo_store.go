package persistence

import (
	"context"
	"database/sql"

	"github.com/satori/go.uuid"
)

type Garbanzo struct {
	Id        int
	APIUUID   uuid.UUID
	FirstName string
	LastName  string
}

func FetchAllGarbanzos(ctx context.Context, database Database) ([]Garbanzo, error) {
	query := "select id, api_uuid, first_name, last_name from garbanzo order by id"

	rows, err := database.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var garbanzos []Garbanzo
	for rows.Next() {
		var id int
		var apiUUID uuid.UUID
		var firstName, lastName string
		err = rows.Scan(&id, &apiUUID, &firstName, &lastName)
		if err != nil {
			return nil, err
		}

		garbanzo := Garbanzo{
			Id:        id,
			APIUUID:   apiUUID,
			FirstName: firstName,
			LastName:  lastName,
		}
		garbanzos = append(garbanzos, garbanzo)
	}

	return garbanzos, nil
}

func FetchGarbanzoById(ctx context.Context, database Database, id int) (Garbanzo, error) {
	query := "select api_uuid, first_name, last_name from garbanzo where id = $1"

	var apiUUID uuid.UUID
	var firstName, lastName string
	err := database.QueryRowContext(ctx, query, id).Scan(&apiUUID, &firstName, &lastName)
	if err == sql.ErrNoRows {
		return Garbanzo{}, ErrNotFound
	} else if err != nil {
		return Garbanzo{}, err
	}

	return Garbanzo{
		Id:        id,
		APIUUID:   apiUUID,
		FirstName: firstName,
		LastName:  lastName,
	}, nil
}

func CreateGarbanzo(ctx context.Context, database Database, garbanzo Garbanzo) (int, error) {
	query := "insert into garbanzo (api_uuid, first_name, last_name) values ($1, $2, $3) returning id"

	var garbanzoId int
	err := database.QueryRowContext(
		ctx,
		query,
		garbanzo.APIUUID,
		garbanzo.FirstName,
		garbanzo.LastName).Scan(&garbanzoId)
	if err != nil {
		return 0, err
	}

	return garbanzoId, nil
}

func DeleteGarbanzoById(ctx context.Context, database Database, id int) error {
	query := "delete from garbanzo where id = $1"

	result, err := database.ExecContext(ctx, query, id)
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
		panic("Deleted multiple rows when expecting only one")
	}

	return nil
}
