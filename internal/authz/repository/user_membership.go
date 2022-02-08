package repository

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
)

type UserMembershipRepository interface {
	SearchMyMemberships(ctx context.Context) ([]*authz.Membership, error)
}
