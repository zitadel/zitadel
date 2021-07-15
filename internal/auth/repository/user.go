package repository

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/domain"
	key_model "github.com/caos/zitadel/internal/key/model"

	"github.com/caos/zitadel/internal/user/model"
)

type UserRepository interface {
	myUserRepo

	UserSessionUserIDsByAgentID(ctx context.Context, agentID string) ([]string, error)

	UserByID(ctx context.Context, userID string) (*model.UserView, error)
	UserByLoginName(ctx context.Context, loginName string) (*model.UserView, error)

	MachineKeyByID(ctx context.Context, keyID string) (*key_model.AuthNKeyView, error)

	SearchUsers(ctx context.Context, request *model.UserSearchRequest) (*model.UserSearchResponse, error)
}

type myUserRepo interface {
	MyUser(ctx context.Context) (*model.UserView, error)

	MyProfile(ctx context.Context) (*model.Profile, error)

	MyEmail(ctx context.Context) (*model.Email, error)

	MyPhone(ctx context.Context) (*model.Phone, error)

	MyAddress(ctx context.Context) (*model.Address, error)

	SearchMyExternalIDPs(ctx context.Context, request *model.ExternalIDPSearchRequest) (*model.ExternalIDPSearchResponse, error)

	MyUserMFAs(ctx context.Context) ([]*model.MultiFactor, error)

	GetMyPasswordless(ctx context.Context) ([]*model.WebAuthNView, error)

	MyUserChanges(ctx context.Context, lastSequence uint64, limit uint64, sortAscending bool, retention time.Duration) (*model.UserChanges, error)

	SearchMyUserMemberships(ctx context.Context, request *model.UserMembershipSearchRequest) (*model.UserMembershipSearchResponse, error)

	GetMyMetaDataByKey(ctx context.Context, key string) (*domain.MetaData, error)
	SearchMyMetaData(ctx context.Context, req *domain.MetaDataSearchRequest) (*domain.MetaDataSearchResponse, error)
}
