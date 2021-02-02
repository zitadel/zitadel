package repository

import (
	"context"
	"github.com/caos/zitadel/internal/usergrant/model"
)

type UserGrantRepository interface {
	UserGrantByID(ctx context.Context, grantID string) (*model.UserGrantView, error)
	SearchUserGrants(ctx context.Context, request *model.UserGrantSearchRequest) (*model.UserGrantSearchResponse, error)
	UserGrantsByProjectID(ctx context.Context, projectID string) ([]*model.UserGrantView, error)
	UserGrantsByUserID(ctx context.Context, userID string) ([]*model.UserGrantView, error)
}
