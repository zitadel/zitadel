package sql

import (
	"context"

	"github.com/zitadel/zitadel/backend/storage/database"
	"github.com/zitadel/zitadel/backend/storage/repository"
)

func NewUser(client database.QueryExecutor) repository.UserRepository {
	return &User{client: client}
}

type User struct {
	client database.QueryExecutor
}

const userByIDQuery = `SELECT id, username FROM users WHERE id = $1`

// ByID implements [UserRepository].
func (r *User) ByID(ctx context.Context, id string) (*repository.User, error) {
	row := r.client.QueryRow(ctx, userByIDQuery, id)
	var user repository.User
	if err := row.Scan(&user.ID, &user.Username); err != nil {
		return nil, err
	}
	return &user, nil
}

// Create implements [UserRepository].
func (r *User) Create(ctx context.Context, user *repository.User) error {
	return r.client.Exec(ctx, "INSERT INTO users (id, username) VALUES ($1, $2)", user.ID, user.Username)
}
