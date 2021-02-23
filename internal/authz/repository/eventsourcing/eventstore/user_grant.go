package eventstore

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v1"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	es_sdk "github.com/caos/zitadel/internal/eventstore/v1/sdk"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	iam_view "github.com/caos/zitadel/internal/iam/repository/view"
	"k8s.io/apimachinery/pkg/api/errors"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/authz/repository/eventsourcing/view"
	caos_errs "github.com/caos/zitadel/internal/errors"
	global_model "github.com/caos/zitadel/internal/model"
	user_model "github.com/caos/zitadel/internal/user/model"
	user_view_model "github.com/caos/zitadel/internal/user/repository/view/model"
	grant_model "github.com/caos/zitadel/internal/usergrant/model"
	"github.com/caos/zitadel/internal/v2/domain"
)

type UserGrantRepo struct {
	View         *view.View
	IamID        string
	IamProjectID string
	Auth         authz.Config
	Eventstore   v1.Eventstore
}

func (repo *UserGrantRepo) Health() error {
	return repo.View.Health()
}

func (repo *UserGrantRepo) SearchMyMemberships(ctx context.Context) ([]*authz.Membership, error) {
	memberships, err := repo.searchUserMemberships(ctx)
	if err != nil {
		return nil, err
	}
	return userMembershipsToMemberships(memberships), nil
}

func (repo *UserGrantRepo) SearchMyZitadelPermissions(ctx context.Context) ([]string, error) {
	memberships, err := repo.searchUserMemberships(ctx)
	if err != nil {
		return nil, err
	}
	permissions := &grant_model.Permissions{Permissions: []string{}}
	for _, membership := range memberships {
		for _, role := range membership.Roles {
			permissions = repo.mapRoleToPermission(permissions, membership, role)
		}
	}
	return permissions.Permissions, nil
}

func (repo *UserGrantRepo) searchUserMemberships(ctx context.Context) ([]*user_view_model.UserMembershipView, error) {
	ctxData := authz.GetCtxData(ctx)
	orgMemberships, orgCount, err := repo.View.SearchUserMemberships(&user_model.UserMembershipSearchRequest{
		Queries: []*user_model.UserMembershipSearchQuery{
			{
				Key:    user_model.UserMembershipSearchKeyUserID,
				Method: global_model.SearchMethodEquals,
				Value:  ctxData.UserID,
			},
			{
				Key:    user_model.UserMembershipSearchKeyResourceOwner,
				Method: global_model.SearchMethodEquals,
				Value:  ctxData.OrgID,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	iamMemberships, iamCount, err := repo.View.SearchUserMemberships(&user_model.UserMembershipSearchRequest{
		Queries: []*user_model.UserMembershipSearchQuery{
			{
				Key:    user_model.UserMembershipSearchKeyUserID,
				Method: global_model.SearchMethodEquals,
				Value:  ctxData.UserID,
			},
			{
				Key:    user_model.UserMembershipSearchKeyAggregateID,
				Method: global_model.SearchMethodEquals,
				Value:  repo.IamID,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	if orgCount == 0 && iamCount == 0 {
		return []*user_view_model.UserMembershipView{}, nil
	}
	return append(orgMemberships, iamMemberships...), nil
}

func (repo *UserGrantRepo) FillIamProjectID(ctx context.Context) error {
	if repo.IamProjectID != "" {
		return nil
	}
	iam, err := repo.getIAMByID(ctx)
	if err != nil {
		return err
	}
	if iam.SetUpDone < domain.StepCount-1 {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-skiwS", "Setup not done")
	}
	repo.IamProjectID = iam.IAMProjectID
	return nil
}

func (repo *UserGrantRepo) mapRoleToPermission(permissions *grant_model.Permissions, membership *user_view_model.UserMembershipView, role string) *grant_model.Permissions {
	for _, mapping := range repo.Auth.RolePermissionMappings {
		if mapping.Role == role {
			ctxID := ""
			if membership.MemberType == int32(user_model.MemberTypeProject) || membership.MemberType == int32(user_model.MemberTypeProjectGrant) {
				ctxID = membership.ObjectID
			}
			permissions.AppendPermissions(ctxID, mapping.Permissions...)
		}
	}
	return permissions
}

func (u *UserGrantRepo) getIAMByID(ctx context.Context) (*iam_model.IAM, error) {
	query, err := iam_view.IAMByIDQuery(domain.IAMID, 0)
	if err != nil {
		return nil, err
	}
	iam := &iam_es_model.IAM{
		ObjectRoot: models.ObjectRoot{
			AggregateID: domain.IAMID,
		},
	}
	err = es_sdk.Filter(ctx, u.Eventstore.FilterEvents, iam.AppendEvents, query)
	if err != nil && errors.IsNotFound(err) && iam.Sequence == 0 {
		return nil, err
	}
	return iam_es_model.IAMToModel(iam), nil
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
