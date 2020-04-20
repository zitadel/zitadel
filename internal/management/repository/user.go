package repository

import (
	"context"
	"github.com/caos/zitadel/internal/user/model"
)

type UserRepository interface {
	UserByID(ctx context.Context, id string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	RegisterUser(ctx context.Context, user *model.User, resourceOwner string) (*model.User, error)
	DeactivateUser(ctx context.Context, id string) (*model.User, error)
	ReactivateUser(ctx context.Context, id string) (*model.User, error)
	LockUser(ctx context.Context, id string) (*model.User, error)
	UnlockUser(ctx context.Context, id string) (*model.User, error)
}
