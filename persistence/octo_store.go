package persistence

import (
	"context"
	"database/sql"

	"github.com/myshkin5/effective-octo-garbanzo/logs"
	"github.com/myshkin5/effective-octo-garbanzo/persistence/data"
)

type OctoStore struct{}

func (OctoStore) FetchAll(ctx context.Context, database Database) ([]data.Octo, error) {
	query := `select o.id, o.name from octo o
		join org on o.org_id = org.id
		where org.name = $1
		order by o.id`

	rows, err := database.Query(ctx, query, org(ctx))
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
	query := `select o.id from octo o
		join org on o.org_id = org.id
		where o.name = $1 and org.name = $2`
	if selectForUpdate {
		query += " for update"
	}

	var id int
	err := database.QueryRow(ctx, query, name, org(ctx)).Scan(&id)
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
	query := "insert into octo (name, org_id) values ($1, (select id from org where name = $2)) returning id"
	return ExecInsert(ctx, database, query, octo.Name, org(ctx))
}

func (OctoStore) DeleteById(ctx context.Context, database Database, id int) error {
	query := "delete from octo where id = $1 and org_id = (select id from org where name = $2)"
	rowsAffected, err := ExecDelete(ctx, database, query, id, org(ctx))
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
