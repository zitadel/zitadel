package repository

import (
	"context"

	key_model "github.com/caos/zitadel/internal/key/model"
	org_model "github.com/caos/zitadel/internal/org/model"

	"github.com/caos/zitadel/internal/user/model"
)

type UserRepository interface {
	myUserRepo

	GetPasswordless(ctx context.Context, id string) ([]*model.WebAuthNToken, error)

	UserSessionUserIDsByAgentID(ctx context.Context, agentID string) ([]string, error)

	UserByID(ctx context.Context, userID string) (*model.UserView, error)
	UserByLoginName(ctx context.Context, loginName string) (*model.UserView, error)

	MachineKeyByID(ctx context.Context, keyID string) (*key_model.AuthNKeyView, error)
}

type myUserRepo interface {
	MyUser(ctx context.Context) (*model.UserView, error)

	MyProfile(ctx context.Context) (*model.Profile, error)

	MyEmail(ctx context.Context) (*model.Email, error)

	MyPhone(ctx context.Context) (*model.Phone, error)

	MyAddress(ctx context.Context) (*model.Address, error)

	SearchMyExternalIDPs(ctx context.Context, request *model.ExternalIDPSearchRequest) (*model.ExternalIDPSearchResponse, error)

	MyUserMFAs(ctx context.Context) ([]*model.MultiFactor, error)

	GetMyPasswordless(ctx context.Context) ([]*model.WebAuthNToken, error)

	MyUserChanges(ctx context.Context, lastSequence uint64, limit uint64, sortAscending bool) (*model.UserChanges, error)

	SearchMyUserMemberships(ctx context.Context, request *model.UserMembershipSearchRequest) (*model.UserMembershipSearchResponse, error)
}
