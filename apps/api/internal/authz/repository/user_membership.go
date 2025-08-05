package repository

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
)

type UserMembershipRepository interface {
	SearchMyMemberships(ctx context.Context, orgID string, shouldTriggerBulk bool) ([]*authz.Membership, error)
}
