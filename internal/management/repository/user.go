package repository

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/domain"
	key_model "github.com/caos/zitadel/internal/key/model"
	"github.com/caos/zitadel/internal/user/model"
)

type UserRepository interface {
	UserByID(ctx context.Context, id string) (*model.UserView, error)
	SearchUsers(ctx context.Context, request *model.UserSearchRequest, ensureLimit bool) (*model.UserSearchResponse, error)
	UserIDsByDomain(ctx context.Context, domain string) ([]string, error)

	GetUserByLoginNameGlobal(ctx context.Context, email string) (*model.UserView, error)
	IsUserUnique(ctx context.Context, userName, email string) (bool, error)

	GetMetadataByKey(ctx context.Context, userID, resourceOwner, key string) (*domain.Metadata, error)
	SearchMetadata(ctx context.Context, userID, resourceOwner string, req *domain.MetadataSearchRequest) (*domain.MetadataSearchResponse, error)

	UserChanges(ctx context.Context, id string, lastSequence uint64, limit uint64, sortAscending bool, retention time.Duration) (*model.UserChanges, error)

	ProfileByID(ctx context.Context, userID string) (*model.Profile, error)

	UserMFAs(ctx context.Context, userID string) ([]*model.MultiFactor, error)

	GetPasswordless(ctx context.Context, userID string) ([]*model.WebAuthNView, error)

	SearchExternalIDPs(ctx context.Context, request *model.ExternalIDPSearchRequest) (*model.ExternalIDPSearchResponse, error)
	ExternalIDPsByIDPConfigID(ctx context.Context, idpConfigID string) ([]*model.ExternalIDPView, error)
	ExternalIDPsByIDPConfigIDAndResourceOwner(ctx context.Context, idpConfigID, resourceOwner string) ([]*model.ExternalIDPView, error)

	SearchMachineKeys(ctx context.Context, request *key_model.AuthNKeySearchRequest) (*key_model.AuthNKeySearchResponse, error)
	GetMachineKey(ctx context.Context, userID, keyID string) (*key_model.AuthNKeyView, error)

	EmailByID(ctx context.Context, userID string) (*model.Email, error)

	PhoneByID(ctx context.Context, userID string) (*model.Phone, error)

	AddressByID(ctx context.Context, userID string) (*model.Address, error)

	SearchUserMemberships(ctx context.Context, request *model.UserMembershipSearchRequest) (*model.UserMembershipSearchResponse, error)
	UserMembershipsByUserID(ctx context.Context, userID string) ([]*model.UserMembershipView, error)
}
