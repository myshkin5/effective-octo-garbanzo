package persistence

import (
	"context"
	"database/sql"
)

type User struct {
	Id        int
	FirstName string
	LastName  string
}

func FetchAllUsers(ctx context.Context, database Database) ([]User, error) {
	query := "select id, first_name, last_name from garbanzo_user order by id"

	rows, err := database.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var id int
		var firstName, lastName string
		err = rows.Scan(&id, &firstName, &lastName)
		if err != nil {
			return nil, err
		}

		user := User{
			Id:        id,
			FirstName: firstName,
			LastName:  lastName,
		}
		users = append(users, user)
	}

	return users, nil
}

func FetchUserById(ctx context.Context, database Database, id int) (User, error) {
	query := "select first_name, last_name from garbanzo_user where id = $1"

	var firstName, lastName string
	err := database.QueryRowContext(ctx, query, id).Scan(&firstName, &lastName)
	if err == sql.ErrNoRows {
		return User{}, ErrNotFound
	} else if err != nil {
		return User{}, err
	}

	return User{
		Id:        id,
		FirstName: firstName,
		LastName:  lastName,
	}, nil
}

func CreateUser(ctx context.Context, database Database, user User) (int, error) {
	query := "insert into garbanzo_user (first_name, last_name) values ($1, $2) returning id"

	var userId int
	err := database.QueryRowContext(ctx, query, user.FirstName, user.LastName).Scan(&userId)
	if err != nil {
		return 0, err
	}

	return userId, nil
}

func DeleteUserById(ctx context.Context, database Database, id int) error {
	query := "delete from garbanzo_user where id = $1"

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
