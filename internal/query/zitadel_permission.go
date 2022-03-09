package query

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
)

func (q *Queries) MyZitadelPermissions(ctx context.Context, orgID, userID string) (*domain.Permissions, error) {
	userIDQuery, err := NewMembershipUserIDQuery(userID)
	if err != nil {
		return nil, err
	}
	orgIDsQuery, err := NewMembershipResourceOwnersSearchQuery(orgID, domain.IAMID)
	if err != nil {
		return nil, err
	}
	memberships, err := q.Memberships(ctx, &MembershipSearchQuery{
		Queries: []SearchQuery{userIDQuery, orgIDsQuery},
	})
	if err != nil {
		return nil, err
	}
	permissions := &domain.Permissions{Permissions: []string{}}
	for _, membership := range memberships.Memberships {
		for _, role := range membership.Roles {
			permissions = q.mapRoleToPermission(permissions, membership, role)
		}
	}
	return permissions, nil
}

func (q *Queries) mapRoleToPermission(permissions *domain.Permissions, membership *Membership, role string) *domain.Permissions {
	for _, mapping := range q.zitadelRoles {
		if mapping.Role == role {
			ctxID := ""
			if membership.Project != nil {
				ctxID = membership.Project.ProjectID
			} else if membership.ProjectGrant != nil {
				ctxID = membership.ProjectGrant.GrantID
			}
			permissions.AppendPermissions(ctxID, mapping.Permissions...)
		}
	}
	return permissions
}
