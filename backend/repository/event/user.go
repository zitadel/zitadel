package event

import (
	"context"

	"github.com/zitadel/zitadel/backend/repository"
	"github.com/zitadel/zitadel/backend/storage/eventstore"
)

var _ repository.UserRepository = (*User)(nil)

type User struct {
	*eventstore.Eventstore

	next repository.UserRepository
}

func NewUser(eventstore *eventstore.Eventstore, next repository.UserRepository) *User {
	return &User{next: next, Eventstore: eventstore}
}

func (i *User) ByID(ctx context.Context, id string) (*repository.User, error) {
	return i.next.ByID(ctx, id)
}

func (i *User) Create(ctx context.Context, user *repository.User) error {
	err := i.next.Create(ctx, user)
	if err != nil {
		return err
	}

	return i.Push(ctx, user)
}
