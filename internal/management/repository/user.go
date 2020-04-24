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

	SetOneTimePassword(ctx context.Context, password *model.Password) (*model.Password, error)
	RequestSetPassword(ctx context.Context, id string, notifyType model.NotificationType) error

	ProfileByID(ctx context.Context, userID string) (*model.Profile, error)
	ChangeProfile(ctx context.Context, profile *model.Profile) (*model.Profile, error)

	EmailByID(ctx context.Context, userID string) (*model.Email, error)
	ChangeEmail(ctx context.Context, email *model.Email) (*model.Email, error)
	VerifyEmail(ctx context.Context, userID, code string) error
	CreateEmailVerificationCode(ctx context.Context, userID string) error
}
