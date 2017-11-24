package persistence

import (
	"context"
	"database/sql"

	"github.com/myshkin5/effective-octo-garbanzo/logs"
	"github.com/myshkin5/effective-octo-garbanzo/persistence/data"
)

type OctoStore struct{}

func (OctoStore) FetchAll(ctx context.Context, database Database) ([]data.Octo, error) {
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

func (OctoStore) FetchByName(ctx context.Context, database Database, name string, selectForUpdate bool) (data.Octo, error) {
	query := "select id from octo where name = $1"
	if selectForUpdate {
		query += " for update"
	}

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

func (OctoStore) Create(ctx context.Context, database Database, octo data.Octo) (int, error) {
	query := "insert into octo (name) values ($1) returning id"
	return ExecInsert(ctx, database, query, octo.Name)
}

func (OctoStore) DeleteById(ctx context.Context, database Database, id int) error {
	query := "delete from octo where id = $1"
	rowsAffected, err := ExecDelete(ctx, database, query, id)
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
