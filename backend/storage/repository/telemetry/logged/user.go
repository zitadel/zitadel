package logged

import (
	"context"
	"log/slog"

	"github.com/zitadel/zitadel/backend/storage/repository"
)

type User struct {
	*slog.Logger

	next repository.UserRepository
}

func NewUser(logger *slog.Logger, next repository.UserRepository) *User {
	return &User{Logger: logger, next: next}
}

var _ repository.UserRepository = (*User)(nil)

func (i *User) ByID(ctx context.Context, id string) (*repository.User, error) {
	i.Logger.InfoContext(ctx, "By ID Query", slog.String("id", id))
	return i.next.ByID(ctx, id)
}

func (i *User) Create(ctx context.Context, user *repository.User) error {
	err := i.next.Create(ctx, user)
	if err != nil {
		i.Logger.ErrorContext(ctx, "Failed to create user", slog.Any("user", user), slog.Any("cause", err))
		return err
	}
	i.Logger.InfoContext(ctx, "User created successfully", slog.Any("user", user))
	return nil
}
