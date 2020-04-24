package eventsourcing

import (
	"context"
	usr_model "github.com/caos/zitadel/internal/user/model"
	usr_event "github.com/caos/zitadel/internal/user/repository/eventsourcing"
)

type UserRepo struct {
	UserEvents *usr_event.UserEventstore
}

func (repo *UserRepo) UserByID(ctx context.Context, id string) (project *usr_model.User, err error) {
	return repo.UserEvents.UserByID(ctx, id)
}

func (repo *UserRepo) CreateUser(ctx context.Context, user *usr_model.User) (*usr_model.User, error) {
	return repo.UserEvents.CreateUser(ctx, user)
}

func (repo *UserRepo) RegisterUser(ctx context.Context, user *usr_model.User, resourceOwner string) (*usr_model.User, error) {
	return repo.UserEvents.RegisterUser(ctx, user, resourceOwner)
}

func (repo *UserRepo) DeactivateUser(ctx context.Context, id string) (*usr_model.User, error) {
	return repo.UserEvents.DeactivateUser(ctx, id)
}

func (repo *UserRepo) ReactivateUser(ctx context.Context, id string) (*usr_model.User, error) {
	return repo.UserEvents.ReactivateUser(ctx, id)
}

func (repo *UserRepo) LockUser(ctx context.Context, id string) (*usr_model.User, error) {
	return repo.UserEvents.LockUser(ctx, id)
}

func (repo *UserRepo) UnlockUser(ctx context.Context, id string) (*usr_model.User, error) {
	return repo.UserEvents.UnlockUser(ctx, id)
}

func (repo *UserRepo) SetOneTimePassword(ctx context.Context, password *usr_model.Password) (*usr_model.Password, error) {
	return repo.UserEvents.SetOneTimePassword(ctx, password)
}

func (repo *UserRepo) RequestSetPassword(ctx context.Context, id string, notifyType usr_model.NotificationType) error {
	return repo.UserEvents.RequestSetPassword(ctx, id, notifyType)
}

func (repo *UserRepo) ProfileByID(ctx context.Context, userID string) (*usr_model.Profile, error) {
	return repo.UserEvents.ProfileByID(ctx, userID)
}

func (repo *UserRepo) ChangeProfile(ctx context.Context, profile *usr_model.Profile) (*usr_model.Profile, error) {
	return repo.UserEvents.ChangeProfile(ctx, profile)
}

func (repo *UserRepo) EmailByID(ctx context.Context, userID string) (*usr_model.Email, error) {
	return repo.UserEvents.EmailByID(ctx, userID)
}

func (repo *UserRepo) ChangeEmail(ctx context.Context, email *usr_model.Email) (*usr_model.Email, error) {
	return repo.UserEvents.ChangeEmail(ctx, email)
}

func (repo *UserRepo) VerifyEmail(ctx context.Context, userID, code string) error {
	return repo.UserEvents.VerifyEmail(ctx, userID, code)
}

func (repo *UserRepo) CreateEmailVerificationCode(ctx context.Context, userID string) error {
	return repo.UserEvents.CreateEmailVerificationCode(ctx, userID)
}
