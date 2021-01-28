package eventstore

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_view_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	"github.com/caos/zitadel/internal/telemetry/tracing"

	auth_model "github.com/caos/zitadel/internal/auth/model"
	auth_view "github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/sdk"
	org_model "github.com/caos/zitadel/internal/org/model"
	org_es "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	"github.com/caos/zitadel/internal/org/repository/view/model"
	usr_es "github.com/caos/zitadel/internal/user/repository/eventsourcing"
)

const (
	orgOwnerRole = "ORG_OWNER"
)

type OrgRepository struct {
	SearchLimit    uint64
	OrgEventstore  *org_es.OrgEventstore
	UserEventstore *usr_es.UserEventstore

	View           *auth_view.View
	SystemDefaults systemdefaults.SystemDefaults
}

func (repo *OrgRepository) SearchOrgs(ctx context.Context, request *org_model.OrgSearchRequest) (*org_model.OrgSearchResult, error) {
	request.EnsureLimit(repo.SearchLimit)
	sequence, err := repo.View.GetLatestOrgSequence()
	logging.Log("EVENT-7Udhz").OnError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Warn("could not read latest org sequence")
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
		result.Timestamp = sequence.LastSuccessfulSpoolerRun
	}
	return result, nil
}

func (repo *OrgRepository) OrgByPrimaryDomain(primaryDomain string) (*org_model.OrgView, error) {
	org, err := repo.View.OrgByPrimaryDomain(primaryDomain)
	if err != nil {
		return nil, err
	}
	return model.OrgToModel(org), nil
}

func (repo *OrgRepository) RegisterOrg(ctx context.Context, register *auth_model.RegisterOrg) (*auth_model.RegisterOrg, error) {
	pwPolicy, err := repo.View.PasswordComplexityPolicyByAggregateID(repo.SystemDefaults.IamID)
	if err != nil {
		return nil, err
	}
	pwPolicyView := iam_view_model.PasswordComplexityViewToModel(pwPolicy)
	orgPolicy, err := repo.View.OrgIAMPolicyByAggregateID(repo.SystemDefaults.IamID)
	if err != nil {
		return nil, err
	}
	orgPolicyView := iam_view_model.OrgIAMViewToModel(orgPolicy)
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
	user, userAggregates, err := repo.UserEventstore.PrepareRegisterUser(ctx, register.User, nil, pwPolicyView, orgPolicyView, org.AggregateID)
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

func (repo *OrgRepository) GetDefaultOrgIAMPolicy(ctx context.Context) (*iam_model.OrgIAMPolicyView, error) {
	orgPolicy, err := repo.View.OrgIAMPolicyByAggregateID(repo.SystemDefaults.IamID)
	if err != nil {
		return nil, err
	}
	policy := iam_view_model.OrgIAMViewToModel(orgPolicy)
	policy.IAMDomain = repo.SystemDefaults.Domain
	return policy, err
}

func (repo *OrgRepository) GetOrgIAMPolicy(ctx context.Context, orgID string) (*iam_model.OrgIAMPolicyView, error) {
	orgPolicy, err := repo.View.OrgIAMPolicyByAggregateID(orgID)
	if errors.IsNotFound(err) {
		orgPolicy, err = repo.View.OrgIAMPolicyByAggregateID(repo.SystemDefaults.IamID)
	}
	if err != nil {
		return nil, err
	}
	return iam_view_model.OrgIAMViewToModel(orgPolicy), nil
}

func (repo *OrgRepository) GetIDPConfigByID(ctx context.Context, idpConfigID string) (*iam_model.IDPConfigView, error) {
	idpConfig, err := repo.View.IDPConfigByID(idpConfigID)
	if err != nil {
		return nil, err
	}
	return iam_view_model.IDPConfigViewToModel(idpConfig), nil
}

func (repo *OrgRepository) GetMyPasswordComplexityPolicy(ctx context.Context) (*iam_model.PasswordComplexityPolicyView, error) {
	policy, err := repo.View.PasswordComplexityPolicyByAggregateID(authz.GetCtxData(ctx).OrgID)
	if errors.IsNotFound(err) {
		policy, err = repo.View.PasswordComplexityPolicyByAggregateID(repo.SystemDefaults.IamID)
		if err != nil {
			return nil, err
		}
		policy.Default = true
	}
	if err != nil {
		return nil, err
	}
	return iam_view_model.PasswordComplexityViewToModel(policy), err
}
