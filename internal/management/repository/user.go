package repository

import (
	"context"

	"github.com/caos/zitadel/internal/user/model"
)

type UserRepository interface {
	UserByID(ctx context.Context, id string) (*model.UserView, error)
	SearchUsers(ctx context.Context, request *model.UserSearchRequest) (*model.UserSearchResponse, error)

	GetUserByLoginNameGlobal(ctx context.Context, email string) (*model.UserView, error)
	IsUserUnique(ctx context.Context, userName, email string) (bool, error)

	UserChanges(ctx context.Context, id string, lastSequence uint64, limit uint64, sortAscending bool) (*model.UserChanges, error)

	ProfileByID(ctx context.Context, userID string) (*model.Profile, error)

	UserMFAs(ctx context.Context, userID string) ([]*model.MultiFactor, error)

	GetPasswordless(ctx context.Context, userID string) ([]*model.WebAuthNToken, error)

	SearchExternalIDPs(ctx context.Context, request *model.ExternalIDPSearchRequest) (*model.ExternalIDPSearchResponse, error)

	SearchMachineKeys(ctx context.Context, request *model.MachineKeySearchRequest) (*model.MachineKeySearchResponse, error)
	GetMachineKey(ctx context.Context, userID, keyID string) (*model.MachineKeyView, error)

	EmailByID(ctx context.Context, userID string) (*model.Email, error)

	PhoneByID(ctx context.Context, userID string) (*model.Phone, error)

	AddressByID(ctx context.Context, userID string) (*model.Address, error)

	SearchUserMemberships(ctx context.Context, request *model.UserMembershipSearchRequest) (*model.UserMembershipSearchResponse, error)
}
