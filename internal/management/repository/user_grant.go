package repository

import (
	"context"

	"github.com/caos/zitadel/internal/usergrant/model"
)

type UserGrantRepository interface {
	SearchUserGrants(ctx context.Context, request *model.UserGrantSearchRequest) (*model.UserGrantSearchResponse, error)
}
