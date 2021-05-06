package repository

import (
	"context"

	"github.com/caos/zitadel/internal/user/model"
)

type UserRepository interface {
	SearchUsers(ctx context.Context, request *model.UserSearchRequest) (*model.UserSearchResponse, error)
}
