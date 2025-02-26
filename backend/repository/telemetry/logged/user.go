package logged

import (
	"context"
	"log/slog"

	"github.com/zitadel/zitadel/backend/repository"
	"github.com/zitadel/zitadel/backend/telemetry/logging"
)

type User struct {
	logger *logging.Logger

	next repository.UserRepository
}

func NewUser(logger *logging.Logger, next repository.UserRepository) *User {
	return &User{logger: logger, next: next}
}

var _ repository.UserRepository = (*User)(nil)

func (i *User) ByID(ctx context.Context, id string) (*repository.User, error) {
	i.logger.InfoContext(ctx, "By ID Query", slog.String("id", id))
	return i.next.ByID(ctx, id)
}

func (i *User) Create(ctx context.Context, user *repository.User) error {
	err := i.next.Create(ctx, user)
	if err != nil {
		i.logger.ErrorContext(ctx, "Failed to create user", slog.Any("user", user), slog.Any("cause", err))
		return err
	}
	i.logger.InfoContext(ctx, "User created successfully", slog.Any("user", user))
	return nil
}
