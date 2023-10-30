package eventstore

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

type UserMembershipRepo struct {
	Queries *query.Queries
}

func (repo *UserMembershipRepo) SearchMyMemberships(ctx context.Context, orgID string, shouldTriggerBulk bool) (_ []*authz.Membership, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	memberships, err := repo.searchUserMemberships(ctx, orgID, shouldTriggerBulk)
	if err != nil {
		return nil, err
	}
	return userMembershipsToMemberships(memberships), nil
}

func (repo *UserMembershipRepo) searchUserMemberships(ctx context.Context, orgID string, shouldTriggerBulk bool) (_ []*query.Membership, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	ctxData := authz.GetCtxData(ctx)
	userIDQuery, err := query.NewMembershipUserIDQuery(ctxData.UserID)
	if err != nil {
		return nil, err
	}
	orgIDsQuery, err := query.NewMembershipResourceOwnersSearchQuery(orgID, authz.GetInstance(ctx).InstanceID())
	if err != nil {
		return nil, err
	}
	grantedIDQuery, err := query.NewMembershipGrantedOrgIDSearchQuery(orgID)
	if err != nil {
		return nil, err
	}
	memberships, err := repo.Queries.Memberships(ctx, &query.MembershipSearchQuery{
		Queries: []query.SearchQuery{userIDQuery, query.Or(orgIDsQuery, grantedIDQuery)},
	}, false, shouldTriggerBulk)
	if err != nil {
		return nil, err
	}
	return memberships.Memberships, nil
}

func userMembershipToMembership(membership *query.Membership) *authz.Membership {
	if membership.IAM != nil {
		return &authz.Membership{
			MemberType:  authz.MemberTypeIAM,
			AggregateID: membership.IAM.IAMID,
			ObjectID:    membership.IAM.IAMID,
			Roles:       membership.Roles,
		}
	}
	if membership.Org != nil {
		return &authz.Membership{
			MemberType:  authz.MemberTypeOrganization,
			AggregateID: membership.Org.OrgID,
			ObjectID:    membership.Org.OrgID,
			Roles:       membership.Roles,
		}
	}
	if membership.Project != nil {
		return &authz.Membership{
			MemberType:  authz.MemberTypeProject,
			AggregateID: membership.Project.ProjectID,
			ObjectID:    membership.Project.ProjectID,
			Roles:       membership.Roles,
		}
	}
	return &authz.Membership{
		MemberType:  authz.MemberTypeProjectGrant,
		AggregateID: membership.ProjectGrant.ProjectID,
		ObjectID:    membership.ProjectGrant.GrantID,
		Roles:       membership.Roles,
	}
}

func userMembershipsToMemberships(memberships []*query.Membership) []*authz.Membership {
	result := make([]*authz.Membership, len(memberships))
	for i, m := range memberships {
		result[i] = userMembershipToMembership(m)
	}
	return result
}
