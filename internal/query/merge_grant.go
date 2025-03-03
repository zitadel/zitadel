package query

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
)

type Grant struct {
	ID               string
	ResourceOwner    string
	ProjectID        string
	UserID           string
	GroupID          string
	OrgPrimaryDomain string
	Roles            []string
}

func (q *Queries) MergeUserAndGroupGrants(ctx context.Context, userID string, groupIDs []string) ([]Grant, error) {
	userIdQuery, err := NewUserGrantUserIDSearchQuery(userID)
	if err != nil {
		return nil, err
	}
	ownerQuery, err := NewUserGrantResourceOwnerSearchQuery(authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	userGrant, err := q.UserGrant(ctx, true, userIdQuery, ownerQuery)
	if err != nil {
		return nil, err
	}
	var groupGrants []*GroupGrant
	for _, groupID := range groupIDs {
		groupIdQuery, err := NewGroupGrantGroupIDSearchQuery(groupID)
		if err != nil {
			return nil, err
		}
		groupGrant, err := q.GroupGrant(ctx, true, groupIdQuery, ownerQuery)
		if err != nil {
			return nil, err
		}
		groupGrants = append(groupGrants, groupGrant)
	}

	mergedGrants := make([]Grant, 0, 1+len(groupGrants))
	mergedGrants = append(mergedGrants, Grant{
		ID:               userGrant.ID,
		ResourceOwner:    userGrant.ResourceOwner,
		ProjectID:        userGrant.ProjectID,
		UserID:           userGrant.UserID,
		OrgPrimaryDomain: userGrant.OrgPrimaryDomain,
		Roles:            userGrant.Roles,
	})

	for _, groupGrant := range groupGrants {
		mergedGrants = append(mergedGrants, Grant{
			ID:               groupGrant.ID,
			ResourceOwner:    groupGrant.ResourceOwner,
			ProjectID:        groupGrant.ProjectID,
			GroupID:          groupGrant.GroupID,
			OrgPrimaryDomain: groupGrant.OrgPrimaryDomain,
			Roles:            groupGrant.Roles,
		})
	}

	return mergedGrants, nil
}

// func mergeRoles(userGrant *Grant, groupGrants []*Grant) []*Grant {
// 	roleMap := make(map[string]struct{})
// 	merged := make([]*Grant, 0, 1+len(groupGrants))

// 	for _, role := range userGrant.Roles {
// 		roleMap[role] = struct{}{}
// 	}
// 	merged = append(merged, userGrant)

// 	for _, grant := range groupGrants {
// 		for _, role := range grant.Roles {
// 			if _, exists := roleMap[role]; !exists {
// 				roleMap[role] = struct{}{}
// 			}
// 		}
// 		merged = append(merged, grant)
// 	}

// 	return merged
// }
