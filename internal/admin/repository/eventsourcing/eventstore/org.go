package eventstore

import (
	"context"

	admin_model "github.com/caos/zitadel/internal/admin/model"
	admin_view "github.com/caos/zitadel/internal/admin/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/sdk"
	org_model "github.com/caos/zitadel/internal/org/model"
	org_es "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	org_view "github.com/caos/zitadel/internal/org/repository/view"
	policy_es "github.com/caos/zitadel/internal/policy/repository/eventsourcing"
	usr_es "github.com/caos/zitadel/internal/user/repository/eventsourcing"
)

type OrgRepo struct {
	Eventstore       eventstore.Eventstore
	OrgEventstore    *org_es.OrgEventstore
	UserEventstore   *usr_es.UserEventstore
	PolicyEventstore *policy_es.PolicyEventstore

	View *admin_view.View

	SearchLimit uint64
}

func (repo *OrgRepo) SetUpOrg(ctx context.Context, setUp *admin_model.SetupOrg) (*admin_model.SetupOrg, error) {
	policy, err := repo.PolicyEventstore.GetPasswordComplexityPolicy(ctx, "0")
	if err != nil {
		return nil, err
	}
	org, aggregates, err := repo.OrgEventstore.PrepareCreateOrg(ctx, setUp.Org)
	if err != nil {
		return nil, err
	}

	user, userAggregates, err := repo.UserEventstore.PrepareCreateUser(ctx, setUp.User, policy, org.AggregateID)
	if err != nil {
		return nil, err
	}

	aggregates = append(aggregates, userAggregates...)
	setupModel := &Setup{Org: org, User: user}

	member := org_model.NewOrgMemberWithRoles(org.AggregateID, user.AggregateID, "ORG_ADMIN") //TODO: role as const
	_, memberAggregate, err := repo.OrgEventstore.PrepareAddOrgMember(ctx, member)
	if err != nil {
		return nil, err
	}
	aggregates = append(aggregates, memberAggregate)

	err = sdk.PushAggregates(ctx, repo.Eventstore.PushAggregates, setupModel.AppendEvents, aggregates...)
	if err != nil {
		return nil, err
	}

	return SetupToModel(setupModel), nil
}

func (repo *OrgRepo) OrgByID(ctx context.Context, id string) (*org_model.Org, error) {
	return repo.OrgEventstore.OrgByID(ctx, org_model.NewOrg(id))
}

func (repo *OrgRepo) SearchOrgs(ctx context.Context, query *org_model.OrgSearchRequest) (*org_model.OrgSearchResult, error) {
	query.EnsureLimit(repo.SearchLimit)
	orgs, count, err := repo.View.SearchOrgs(query)
	if err != nil {
		return nil, err
	}
	return &org_model.OrgSearchResult{
		Offset:      query.Offset,
		Limit:       query.Limit,
		TotalResult: uint64(count),
		Result:      org_view.OrgsToModel(orgs),
	}, nil
}

func (repo *OrgRepo) IsOrgUnique(ctx context.Context, name, domain string) (isUnique bool, err error) {
	return repo.OrgEventstore.IsOrgUnique(ctx, name, domain)
}
