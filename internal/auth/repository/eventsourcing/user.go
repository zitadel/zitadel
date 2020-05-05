package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/api/auth"
	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/model"
	user_event "github.com/caos/zitadel/internal/user/repository/eventsourcing"
)

type UserRepo struct {
	UserEvents *user_event.UserEventstore
}

func (repo *UserRepo) MyProfile(ctx context.Context) (*model.Profile, error) {
	return repo.UserEvents.ProfileByID(ctx, auth.GetCtxData(ctx).UserID)
}

func (repo *UserRepo) ChangeMyProfile(ctx context.Context, profile *model.Profile) (*model.Profile, error) {
	if err := checkIDs(ctx, profile.ObjectRoot); err != nil {
		return nil, err
	}
	return repo.UserEvents.ChangeProfile(ctx, profile)
}

func (repo *UserRepo) MyEmail(ctx context.Context) (*model.Email, error) {
	return repo.UserEvents.EmailByID(ctx, auth.GetCtxData(ctx).UserID)
}

func (repo *UserRepo) ChangeMyEmail(ctx context.Context, email *model.Email) (*model.Email, error) {
	if err := checkIDs(ctx, email.ObjectRoot); err != nil {
		return nil, err
	}
	return repo.UserEvents.ChangeEmail(ctx, email)
}

//func (repo *UserRepo) CreateEmailVerificationCode(ctx context.Context) error

func (repo *UserRepo) MyPhone(ctx context.Context) (*model.Phone, error) {
	return repo.UserEvents.PhoneByID(ctx, auth.GetCtxData(ctx).UserID)
}

func (repo *UserRepo) ChangeMyPhone(ctx context.Context, phone *model.Phone) (*model.Phone, error) {
	if err := checkIDs(ctx, phone.ObjectRoot); err != nil {
		return nil, err
	}
	return repo.UserEvents.ChangePhone(ctx, phone)
}

//func (repo *UserRepo) CreatePhoneVerificationCode(ctx context.Context) error

func (repo *UserRepo) MyAddress(ctx context.Context) (*model.Address, error) {
	return repo.UserEvents.AddressByID(ctx, auth.GetCtxData(ctx).UserID)
}

func (repo *UserRepo) ChangeMyAddress(ctx context.Context, address *model.Address) (*model.Address, error) {
	if err := checkIDs(ctx, address.ObjectRoot); err != nil {
		return nil, err
	}
	return repo.UserEvents.ChangeAddress(ctx, address)
}

func (repo *UserRepo) AddMfaOTP(ctx context.Context) (*model.OTP, error) {
	return repo.UserEvents.AddOTP(ctx, auth.GetCtxData(ctx).UserID)
}

func (repo *UserRepo) VerifyMfaOTP(ctx context.Context, code string) (*model.OTP, error) {
	return nil, repo.UserEvents.CheckMfaOTP(ctx, auth.GetCtxData(ctx).UserID, code) //TODO:
}

func (repo *UserRepo) RemoveMyMfaOTP(ctx context.Context) error {
	return repo.UserEvents.RemoveOTP(ctx, auth.GetCtxData(ctx).UserID)
}

func checkIDs(ctx context.Context, obj es_models.ObjectRoot) error {
	if obj.AggregateID != auth.GetCtxData(ctx).UserID {
		return errors.ThrowPermissionDenied(nil, "EVENT-kFi9w", "object does not belong to user")
	}
	return nil
}

//
//func (repo *UserRepo) UserByID(ctx context.Context, id string) (*usr_model.User, error) {
//	return repo.UserEvents.UserByID(ctx, id)
//}
//
////func (repo *UserRepo) UserByUsername(ctx context.Context, username string) (*usr_model.User, error) {
////	return nil, errors.ThrowUnimplemented(nil, "EVENT-asjod", "user by username not yet implemented")
////}
//
//func (repo *UserRepo) RegisterUser(ctx context.Context, user *usr_model.User, resourceOwner string) (*usr_model.User, error) {
//	return repo.UserEvents.RegisterUser(ctx, user, resourceOwner)
//}
//
////func (repo *UserRepo) CheckUserPassword(ctx context.Context, id, password string) error {
////	return repo.UserEvents.CheckPassword(ctx, id, password)
////}
//
//func (repo *UserRepo) RequestSetPassword(ctx context.Context, id string, notifyType usr_model.NotificationType) error {
//	return repo.UserEvents.RequestSetPassword(ctx, id, notifyType)
//}
//
//func (repo *UserRepo) SkipMfaInit(ctx context.Context, id, password string) error {
//	return repo.UserEvents.SkipMfaInit(ctx, id)
//}
//
//func (repo *UserRepo) AddOTP(ctx context.Context, id string) (*usr_model.OTP, error) {
//	return repo.UserEvents.AddOTP(ctx, id)
//}
//
//func (repo *UserRepo) CheckMfaOtp(ctx context.Context, id, code string) error {
//	return repo.UserEvents.CheckMfaOTP(ctx, id, code)
//}
