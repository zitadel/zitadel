package repository

import (
	"context"
	"log"
)

const userByIDQuery = `SELECT id, username FROM users WHERE id = $1`

func (q *querier) UserByID(ctx context.Context, id string) (res *User, err error) {
	log.Println("sql.user.byID")
	row := q.client.QueryRow(ctx, userByIDQuery, id)
	var user User
	if err := row.Scan(&user.ID, &user.Username); err != nil {
		return nil, err
	}
	return &user, nil
}

func (e *executor) CreateUser(ctx context.Context, user *User) (res *User, err error) {
	log.Println("sql.user.create")
	err = e.client.Exec(ctx, "INSERT INTO users (id, username) VALUES ($1, $2)", user.ID, user.Username)
	if err != nil {
		return nil, err
	}
	return user, nil
}
