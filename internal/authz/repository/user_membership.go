package repository

import (
	"context"

	"github.com/zitadel/zitadel/v2/internal/api/authz"
)

type UserMembershipRepository interface {
	SearchMyMemberships(ctx context.Context, orgID string, shouldTriggerBulk bool) ([]*authz.Membership, error)
}
