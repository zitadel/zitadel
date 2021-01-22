package repository

import (
	"context"
	"github.com/caos/zitadel/internal/usergrant/model"
)

type UserGrantRepository interface {
	UserGrantByID(ctx context.Context, grantID string) (*model.UserGrantView, error)
	SearchUserGrants(ctx context.Context, request *model.UserGrantSearchRequest) (*model.UserGrantSearchResponse, error)
}
