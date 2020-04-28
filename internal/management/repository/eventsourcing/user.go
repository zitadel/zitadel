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

func (repo *UserRepo) CreateEmailVerificationCode(ctx context.Context, userID string) error {
	return repo.UserEvents.CreateEmailVerificationCode(ctx, userID)
}

func (repo *UserRepo) PhoneByID(ctx context.Context, userID string) (*usr_model.Phone, error) {
	return repo.UserEvents.PhoneByID(ctx, userID)
}

func (repo *UserRepo) ChangePhone(ctx context.Context, email *usr_model.Phone) (*usr_model.Phone, error) {
	return repo.UserEvents.ChangePhone(ctx, email)
}

func (repo *UserRepo) CreatePhoneVerificationCode(ctx context.Context, userID string) error {
	return repo.UserEvents.CreatePhoneVerificationCode(ctx, userID)
}

func (repo *UserRepo) AddressByID(ctx context.Context, userID string) (*usr_model.Address, error) {
	return repo.UserEvents.AddressByID(ctx, userID)
}

func (repo *UserRepo) ChangeAddress(ctx context.Context, address *usr_model.Address) (*usr_model.Address, error) {
	return repo.UserEvents.ChangeAddress(ctx, address)
}

func (repo *UserRepo) UserGrantByID(ctx context.Context, userID, grantID string) (*usr_model.UserGrant, error) {
	return repo.UserEvents.UserGrantByIDs(ctx, userID, grantID)
}

func (repo *UserRepo) AddUserGrant(ctx context.Context, grant *usr_model.UserGrant) (*usr_model.UserGrant, error) {
	return repo.UserEvents.AddUserGrant(ctx, grant)
}

func (repo *UserRepo) ChangeUserGrant(ctx context.Context, grant *usr_model.UserGrant) (*usr_model.UserGrant, error) {
	return repo.UserEvents.ChangeUserGrant(ctx, grant)
}

func (repo *UserRepo) DeactivateUserGrant(ctx context.Context, userID, grantID string) (*usr_model.UserGrant, error) {
	return repo.UserEvents.DeactivateUserGrant(ctx, userID, grantID)
}

func (repo *UserRepo) ReactivateUserGrant(ctx context.Context, userID, grantID string) (*usr_model.UserGrant, error) {
	return repo.UserEvents.ReactivateUserGrant(ctx, userID, grantID)
}

func (repo *UserRepo) RemoveUserGrant(ctx context.Context, userID, grantID string) error {
	user := usr_model.NewUserGrant(userID, grantID)
	return repo.UserEvents.RemoveUserGrant(ctx, user)
}
