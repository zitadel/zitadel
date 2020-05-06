package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/api/auth"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/model"
	user_event "github.com/caos/zitadel/internal/user/repository/eventsourcing"
)

type UserRepo struct {
	UserEvents  *user_event.UserEventstore
	PasswordAlg crypto.HashAlgorithm
}

func (repo *UserRepo) Register(ctx context.Context, user *model.User, resourceOwner string) (*model.User, error) {
	return repo.UserEvents.RegisterUser(ctx, user, resourceOwner)
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

func (repo *UserRepo) VerifyMyEmail(ctx context.Context, code string) error {
	return repo.UserEvents.VerifyEmail(ctx, auth.GetCtxData(ctx).UserID, code)
}

func (repo *UserRepo) ResendMyEmailVerificationMail(ctx context.Context) error {
	return repo.UserEvents.CreateEmailVerificationCode(ctx, auth.GetCtxData(ctx).UserID)
}

func (repo *UserRepo) MyPhone(ctx context.Context) (*model.Phone, error) {
	return repo.UserEvents.PhoneByID(ctx, auth.GetCtxData(ctx).UserID)
}

func (repo *UserRepo) ChangeMyPhone(ctx context.Context, phone *model.Phone) (*model.Phone, error) {
	if err := checkIDs(ctx, phone.ObjectRoot); err != nil {
		return nil, err
	}
	return repo.UserEvents.ChangePhone(ctx, phone)
}

func (repo *UserRepo) VerifyMyPhone(ctx context.Context, code string) error {
	return repo.UserEvents.VerifyPhone(ctx, auth.GetCtxData(ctx).UserID, code)
}

func (repo *UserRepo) ResendMyPhoneVerificationCode(ctx context.Context) error {
	return repo.UserEvents.CreatePhoneVerificationCode(ctx, auth.GetCtxData(ctx).UserID)
}

func (repo *UserRepo) MyAddress(ctx context.Context) (*model.Address, error) {
	return repo.UserEvents.AddressByID(ctx, auth.GetCtxData(ctx).UserID)
}

func (repo *UserRepo) ChangeMyAddress(ctx context.Context, address *model.Address) (*model.Address, error) {
	if err := checkIDs(ctx, address.ObjectRoot); err != nil {
		return nil, err
	}
	return repo.UserEvents.ChangeAddress(ctx, address)
}

func (repo *UserRepo) ChangeMyPassword(ctx context.Context, old, new string) error {
	_, err := repo.UserEvents.ChangePassword(ctx, auth.GetCtxData(ctx).UserID, old, new)
	return err
}

func (repo *UserRepo) AddMyMfaOTP(ctx context.Context) (*model.OTP, error) {
	return repo.UserEvents.AddOTP(ctx, auth.GetCtxData(ctx).UserID)
}

func (repo *UserRepo) VerifyMyMfaOTP(ctx context.Context, code string) (*model.OTP, error) {
	return nil, repo.UserEvents.CheckMfaOTP(ctx, auth.GetCtxData(ctx).UserID, code) //TODO: return?
}

func (repo *UserRepo) RemoveMyMfaOTP(ctx context.Context) error {
	return repo.UserEvents.RemoveOTP(ctx, auth.GetCtxData(ctx).UserID)
}

func (repo *UserRepo) RequestPasswordReset(ctx context.Context, username string) error {
	//TODO: get id from view
	var userID string
	return repo.UserEvents.RequestSetPassword(ctx, userID, model.NOTIFICATIONTYPE_EMAIL) //TODO: ?
}

func (repo *UserRepo) SetPassword(ctx context.Context, userID, code, password string) error {
	return repo.UserEvents.SetPassword(ctx, userID, code, password)
}

func checkIDs(ctx context.Context, obj es_models.ObjectRoot) error {
	if obj.AggregateID != auth.GetCtxData(ctx).UserID {
		return errors.ThrowPermissionDenied(nil, "EVENT-kFi9w", "object does not belong to user")
	}
	return nil
}
