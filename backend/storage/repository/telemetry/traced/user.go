package traced

import (
	"context"

	"github.com/zitadel/zitadel/backend/storage/repository"
	"github.com/zitadel/zitadel/backend/telemetry/tracing"
)

var _ repository.UserRepository = (*User)(nil)

type User struct {
	*tracing.Tracer

	next repository.UserRepository
}

func NewUser(tracer *tracing.Tracer, next repository.UserRepository) *User {
	return &User{Tracer: tracer, next: next}
}

func (i *User) SetNext(next repository.UserRepository) *User {
	return &User{Tracer: i.Tracer, next: next}
}

// ByID implements [repository.UserRepository].
func (i *User) ByID(ctx context.Context, id string) (user *repository.User, err error) {
	i.Tracer.Decorate(ctx, func(ctx context.Context) error {
		user, err = i.next.ByID(ctx, id)
		return err
	})

	return user, err
}

// Create implements [repository.UserRepository].
func (i *User) Create(ctx context.Context, user *repository.User) (err error) {
	i.Tracer.Decorate(ctx, func(ctx context.Context) error {
		err = i.next.Create(ctx, user)
		return err
	})

	return err
}
