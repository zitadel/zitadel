package repository

import (
	"context"
	"github.com/caos/zitadel/internal/usergrant/model"
)

type UserGrantRepository interface {
	UserGrantByID(ctx context.Context, grantID string) (*model.UserGrant, error)
	AddUserGrant(ctx context.Context, grant *model.UserGrant) (*model.UserGrant, error)
	ChangeUserGrant(ctx context.Context, grant *model.UserGrant) (*model.UserGrant, error)
	DeactivateUserGrant(ctx context.Context, grantID string) (*model.UserGrant, error)
	ReactivateUserGrant(ctx context.Context, grantID string) (*model.UserGrant, error)
	RemoveUserGrant(ctx context.Context, grantID string) error
	SearchUserGrants(ctx context.Context, request *model.UserGrantSearchRequest) (*model.UserGrantSearchResponse, error)

	BulkAddUserGrant(ctx context.Context, grant ...*model.UserGrant) error
	BulkChangeUserGrant(ctx context.Context, grant ...*model.UserGrant) error
	BulkRemoveUserGrant(ctx context.Context, grantIDs ...string) error
}
