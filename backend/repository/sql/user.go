package sql

import (
	"context"

	"github.com/zitadel/zitadel/backend/repository"
)

const userByIDQuery = `SELECT id, username FROM users WHERE id = $1`

func (q *querier[C]) UserByID(ctx context.Context, id string) (res *repository.User, err error) {
	row := q.client.QueryRow(ctx, userByIDQuery, id)
	var user repository.User
	if err := row.Scan(&user.ID, &user.Username); err != nil {
		return nil, err
	}
	return &user, nil
}

func (e *executor[C]) CreateUser(ctx context.Context, user *repository.User) (res *repository.User, err error) {
	err = e.client.Exec(ctx, "INSERT INTO users (id, username) VALUES ($1, $2)", user.ID, user.Username)
	if err != nil {
		return nil, err
	}
	return user, nil
}
