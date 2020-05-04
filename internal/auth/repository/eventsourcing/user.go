package eventsourcing

import (
	"context"

	usr_model "github.com/caos/zitadel/internal/user/model"
	user_event "github.com/caos/zitadel/internal/user/repository/eventsourcing"
)

type UserRepo struct {
	UserEvents *user_event.UserEventstore
}

func (repo *UserRepo) UserByID(ctx context.Context, id string) (*usr_model.User, error) {
	return repo.UserEvents.UserByID(ctx, id)
}

//func (repo *UserRepo) UserByUsername(ctx context.Context, username string) (*usr_model.User, error) {
//	return nil, errors.ThrowUnimplemented(nil, "EVENT-asjod", "user by username not yet implemented")
//}

func (repo *UserRepo) RegisterUser(ctx context.Context, user *usr_model.User, resourceOwner string) (*usr_model.User, error) {
	return repo.UserEvents.RegisterUser(ctx, user, resourceOwner)
}

//func (repo *UserRepo) CheckUserPassword(ctx context.Context, id, password string) error {
//	return repo.UserEvents.CheckPassword(ctx, id, password)
//}

func (repo *UserRepo) RequestSetPassword(ctx context.Context, id string, notifyType usr_model.NotificationType) error {
	return repo.UserEvents.RequestSetPassword(ctx, id, notifyType)
}

func (repo *UserRepo) SkipMfaInit(ctx context.Context, id, password string) error {
	return repo.UserEvents.SkipMfaInit(ctx, id)
}

func (repo *UserRepo) AddOTP(ctx context.Context, id string) (*usr_model.OTP, error) {
	return repo.UserEvents.AddOTP(ctx, id)
}

func (repo *UserRepo) CheckMfaOtp(ctx context.Context, id, code string) error {
	return repo.UserEvents.CheckMfaOTP(ctx, id, code)
}
