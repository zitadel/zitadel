package repository

import "context"

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	ByID(ctx context.Context, id string) (*User, error)
}

type User struct {
	ID       string
	Username string
}
