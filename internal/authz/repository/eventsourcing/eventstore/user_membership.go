package eventstore

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/authz/repository/eventsourcing/view"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	user_model "github.com/zitadel/zitadel/internal/user/model"
	user_view_model "github.com/zitadel/zitadel/internal/user/repository/view/model"
)

type UserMembershipRepo struct {
	View *view.View
}

func (repo *UserMembershipRepo) Health() error {
	return repo.View.Health()
}

func (repo *UserMembershipRepo) SearchMyMemberships(ctx context.Context) (_ []*authz.Membership, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	memberships, err := repo.searchUserMemberships(ctx)
	if err != nil {
		return nil, err
	}
	return userMembershipsToMemberships(memberships), nil
}

func (repo *UserMembershipRepo) searchUserMemberships(ctx context.Context) (_ []*user_view_model.UserMembershipView, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	ctxData := authz.GetCtxData(ctx)
	instance := authz.GetInstance(ctx)
	ctx, orgSpan := tracing.NewSpan(ctx)
	orgMemberships, orgCount, err := repo.View.SearchUserMemberships(&user_model.UserMembershipSearchRequest{
		Queries: []*user_model.UserMembershipSearchQuery{
			{
				Key:    user_model.UserMembershipSearchKeyUserID,
				Method: domain.SearchMethodEquals,
				Value:  ctxData.UserID,
			},
			{
				Key:    user_model.UserMembershipSearchKeyResourceOwner,
				Method: domain.SearchMethodEquals,
				Value:  ctxData.OrgID,
			},
			{
				Key:    user_model.UserMembershipSearchKeyInstanceID,
				Method: domain.SearchMethodEquals,
				Value:  instance.InstanceID(),
			},
		},
	})
	orgSpan.EndWithError(err)
	if err != nil {
		return nil, err
	}
	ctx, iamSpan := tracing.NewSpan(ctx)
	iamMemberships, iamCount, err := repo.View.SearchUserMemberships(&user_model.UserMembershipSearchRequest{
		Queries: []*user_model.UserMembershipSearchQuery{
			{
				Key:    user_model.UserMembershipSearchKeyUserID,
				Method: domain.SearchMethodEquals,
				Value:  ctxData.UserID,
			},
			{
				Key:    user_model.UserMembershipSearchKeyAggregateID,
				Method: domain.SearchMethodEquals,
				Value:  instance.InstanceID(),
			},
			{
				Key:    user_model.UserMembershipSearchKeyInstanceID,
				Method: domain.SearchMethodEquals,
				Value:  instance.InstanceID(),
			},
		},
	})
	iamSpan.EndWithError(err)
	if err != nil {
		return nil, err
	}
	if orgCount == 0 && iamCount == 0 {
		return []*user_view_model.UserMembershipView{}, nil
	}
	return append(orgMemberships, iamMemberships...), nil
}

func userMembershipToMembership(membership *user_view_model.UserMembershipView) *authz.Membership {
	return &authz.Membership{
		MemberType:  authz.MemberType(membership.MemberType),
		AggregateID: membership.AggregateID,
		ObjectID:    membership.ObjectID,
		Roles:       membership.Roles,
	}
}

func userMembershipsToMemberships(memberships []*user_view_model.UserMembershipView) []*authz.Membership {
	result := make([]*authz.Membership, len(memberships))
	for i, m := range memberships {
		result[i] = userMembershipToMembership(m)
	}
	return result
}
