package persistence

import (
	"context"
	"database/sql"

	"github.com/myshkin5/effective-octo-garbanzo/persistence/data"
)

type OctoStore struct{}

func (OctoStore) FetchAllOctos(ctx context.Context, database Database) ([]data.Octo, error) {
	query := "select id, name from octo order by id"

	rows, err := database.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var octos []data.Octo
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}

		octo := data.Octo{
			Id:   id,
			Name: name,
		}
		octos = append(octos, octo)
	}

	return octos, nil
}

func (OctoStore) FetchOctoByName(ctx context.Context, database Database, name string) (data.Octo, error) {
	query := "select id from octo where name = $1"

	var id int
	err := database.QueryRow(ctx, query, name).Scan(&id)
	if err == sql.ErrNoRows {
		return data.Octo{}, ErrNotFound
	} else if err != nil {
		return data.Octo{}, err
	}

	return data.Octo{
		Id:   id,
		Name: name,
	}, nil
}

func (OctoStore) CreateOcto(ctx context.Context, database Database, octo data.Octo) (int, error) {
	query := "insert into octo (name) values ($1) returning id"
	return ExecInsert(ctx, database, query, octo.Name)
}

func (OctoStore) DeleteOctoByName(ctx context.Context, database Database, name string) error {
	query := "delete from octo where name = $1"
	return ExecDelete(ctx, database, query, name)
}
