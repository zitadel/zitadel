package repository

import (
	"context"

	"github.com/caos/zitadel/internal/user/model"
)

type UserRepository interface {
	UserByID(ctx context.Context, id string) (*model.UserView, error)
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	RegisterUser(ctx context.Context, user *model.User, resourceOwner string) (*model.User, error)
	DeactivateUser(ctx context.Context, id string) (*model.User, error)
	ReactivateUser(ctx context.Context, id string) (*model.User, error)
	LockUser(ctx context.Context, id string) (*model.User, error)
	UnlockUser(ctx context.Context, id string) (*model.User, error)
	RemoveUser(ctx context.Context, id string) error
	SearchUsers(ctx context.Context, request *model.UserSearchRequest) (*model.UserSearchResponse, error)

	GetUserByLoginNameGlobal(ctx context.Context, email string) (*model.UserView, error)
	IsUserUnique(ctx context.Context, userName, email string) (bool, error)

	UserChanges(ctx context.Context, id string, lastSequence uint64, limit uint64, sortAscending bool) (*model.UserChanges, error)

	ChangeUsername(ctx context.Context, id, username string) error

	SetOneTimePassword(ctx context.Context, password *model.Password) (*model.Password, error)
	RequestSetPassword(ctx context.Context, id string, notifyType model.NotificationType) error

	ProfileByID(ctx context.Context, userID string) (*model.Profile, error)
	ChangeProfile(ctx context.Context, profile *model.Profile) (*model.Profile, error)

	UserMFAs(ctx context.Context, userID string) ([]*model.MultiFactor, error)
	RemoveOTP(ctx context.Context, userID string) error

	SearchExternalIDPs(ctx context.Context, request *model.ExternalIDPSearchRequest) (*model.ExternalIDPSearchResponse, error)
	RemoveExternalIDP(ctx context.Context, externalIDP *model.ExternalIDP) error

	SearchMachineKeys(ctx context.Context, request *model.MachineKeySearchRequest) (*model.MachineKeySearchResponse, error)
	GetMachineKey(ctx context.Context, userID, keyID string) (*model.MachineKeyView, error)
	ChangeMachine(ctx context.Context, machine *model.Machine) (*model.Machine, error)
	AddMachineKey(ctx context.Context, key *model.MachineKey) (*model.MachineKey, error)
	RemoveMachineKey(ctx context.Context, userID, keyID string) error

	EmailByID(ctx context.Context, userID string) (*model.Email, error)
	ChangeEmail(ctx context.Context, email *model.Email) (*model.Email, error)
	CreateEmailVerificationCode(ctx context.Context, userID string) error

	PhoneByID(ctx context.Context, userID string) (*model.Phone, error)
	ChangePhone(ctx context.Context, email *model.Phone) (*model.Phone, error)
	RemovePhone(ctx context.Context, userID string) error
	CreatePhoneVerificationCode(ctx context.Context, userID string) error

	AddressByID(ctx context.Context, userID string) (*model.Address, error)
	ChangeAddress(ctx context.Context, address *model.Address) (*model.Address, error)

	SearchUserMemberships(ctx context.Context, request *model.UserMembershipSearchRequest) (*model.UserMembershipSearchResponse, error)

	ResendInitialMail(ctx context.Context, userID, email string) error
}
