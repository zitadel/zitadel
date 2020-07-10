package repository

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
)

type UserGrantRepository interface {
	ResolveGrants(ctx context.Context) (*authz.Grant, error)
	SearchMyZitadelPermissions(ctx context.Context) ([]string, error)
}
