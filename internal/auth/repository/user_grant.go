package repository

import (
	"context"
	"github.com/caos/zitadel/internal/usergrant/model"
)

type UserGrantRepository interface {
	SearchMyUserGrants(ctx context.Context, request *model.UserGrantSearchRequest) (*model.UserGrantSearchResponse, error)
	SearchMyProjectOrgs(ctx context.Context, request *model.UserGrantSearchRequest) (*model.ProjectOrgSearchResponse, error)
	SearchMyZitadelPermissions(ctx context.Context) ([]string, error)
}
