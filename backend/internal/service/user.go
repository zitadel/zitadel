package service

import (
	"context"

	"github.com/zitadel/zitadel/backend/internal/port"
)

type User struct {
	ID         string `consistent:"id,pk"`
	InstanceID string `consistent:"instance_id,pk"`

	Username string `consistent:"username"`
}

func (u User) Columns() []*port.Column {
	return []*port.Column{
		{Name: "id", Value: u.ID},
		{Name: "username", Value: u.Username},
	}
}

type CreateUserRequest struct {
	Username string
}

type UserRepository interface {
	// CreateUser creates a new user
	CreateUser(ctx context.Context, executor port.Executor[User], user *User) error
}
