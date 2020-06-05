package repository

import (
	"context"
	"github.com/caos/zitadel/internal/api/auth"
)

type UserGrantRepository interface {
	ResolveGrants(ctx context.Context) (*auth.Grant, error)
	SearchMyZitadelPermissions(ctx context.Context) ([]string, error)
}
