package eventstore

import (
	"context"
	"github.com/caos/logging"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_view_model "github.com/caos/zitadel/internal/iam/repository/view/model"

	auth_model "github.com/caos/zitadel/internal/auth/model"
	auth_view "github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/sdk"
	org_model "github.com/caos/zitadel/internal/org/model"
	org_es "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	"github.com/caos/zitadel/internal/org/repository/view/model"
	policy_model "github.com/caos/zitadel/internal/policy/model"
	policy_es "github.com/caos/zitadel/internal/policy/repository/eventsourcing"
	usr_es "github.com/caos/zitadel/internal/user/repository/eventsourcing"
)

const (
	orgOwnerRole = "ORG_OWNER"
)

type OrgRepository struct {
	SearchLimit      uint64
	OrgEventstore    *org_es.OrgEventstore
	UserEventstore   *usr_es.UserEventstore
	PolicyEventstore *policy_es.PolicyEventstore

	View *auth_view.View
}

func (repo *OrgRepository) SearchOrgs(ctx context.Context, request *org_model.OrgSearchRequest) (*org_model.OrgSearchResult, error) {
	request.EnsureLimit(repo.SearchLimit)
	sequence, err := repo.View.GetLatestOrgSequence()
	logging.Log("EVENT-7Udhz").OnError(err).Warn("could not read latest org sequence")
	members, count, err := repo.View.SearchOrgs(request)
	if err != nil {
		return nil, err
	}
	result := &org_model.OrgSearchResult{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: count,
		Result:      model.OrgsToModel(members),
	}
	if err == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.CurrentTimestamp
	}
	return result, nil
}

func (repo *OrgRepository) RegisterOrg(ctx context.Context, register *auth_model.RegisterOrg) (*auth_model.RegisterOrg, error) {
	pwPolicy, err := repo.PolicyEventstore.GetPasswordComplexityPolicy(ctx, policy_model.DefaultPolicy)
	if err != nil {
		return nil, err
	}
	orgPolicy, err := repo.OrgEventstore.GetOrgIAMPolicy(ctx, policy_model.DefaultPolicy)
	if err != nil {
		return nil, err
	}
	users := func(ctx context.Context, domain string) ([]*es_models.Aggregate, error) {
		userIDs, err := repo.View.UserIDsByDomain(domain)
		if err != nil {
			return nil, err
		}
		return repo.UserEventstore.PrepareDomainClaimed(ctx, userIDs)
	}
	org, aggregates, err := repo.OrgEventstore.PrepareCreateOrg(ctx, register.Org, users)
	if err != nil {
		return nil, err
	}
	user, userAggregates, err := repo.UserEventstore.PrepareRegisterUser(ctx, register.User, nil, pwPolicy, orgPolicy, org.AggregateID)
	if err != nil {
		return nil, err
	}

	aggregates = append(aggregates, userAggregates...)
	registerModel := &Register{Org: org, User: user}

	member := org_model.NewOrgMemberWithRoles(org.AggregateID, user.AggregateID, orgOwnerRole)
	_, memberAggregate, err := repo.OrgEventstore.PrepareAddOrgMember(ctx, member, org.AggregateID)
	if err != nil {
		return nil, err
	}
	aggregates = append(aggregates, memberAggregate)

	err = sdk.PushAggregates(ctx, repo.OrgEventstore.PushAggregates, registerModel.AppendEvents, aggregates...)
	if err != nil {
		return nil, err
	}

	return RegisterToModel(registerModel), nil
}

func (repo *OrgRepository) GetDefaultOrgIamPolicy(ctx context.Context) *org_model.OrgIAMPolicy {
	return repo.OrgEventstore.GetDefaultOrgIAMPolicy(ctx)
}

func (repo *OrgRepository) GetOrgIamPolicy(ctx context.Context, orgID string) (*org_model.OrgIAMPolicy, error) {
	return repo.OrgEventstore.GetOrgIAMPolicy(ctx, orgID)
}

func (repo *OrgRepository) GetIDPConfigByID(ctx context.Context, idpConfigID string) (*iam_model.IDPConfigView, error) {
	idpConfig, err := repo.View.IDPConfigByID(idpConfigID)
	if err != nil {
		return nil, err
	}
	return iam_view_model.IDPConfigViewToModel(idpConfig), nil
}
