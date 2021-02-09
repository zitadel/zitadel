package eventstore

import (
	"context"
	"github.com/caos/zitadel/internal/v2/domain"
	"strings"

	iam_es "github.com/caos/zitadel/internal/iam/repository/eventsourcing"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/sdk"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	iam_view_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	mgmt_view "github.com/caos/zitadel/internal/management/repository/eventsourcing/view"
	global_model "github.com/caos/zitadel/internal/model"
	org_model "github.com/caos/zitadel/internal/org/model"
	org_es "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/org/repository/view/model"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	usr_model "github.com/caos/zitadel/internal/user/model"
	usr_es "github.com/caos/zitadel/internal/user/repository/eventsourcing"
)

const (
	orgOwnerRole = "ORG_OWNER"
)

type OrgRepository struct {
	SearchLimit uint64
	*org_es.OrgEventstore
	UserEvents     *usr_es.UserEventstore
	IAMEventstore  *iam_es.IAMEventstore
	View           *mgmt_view.View
	Roles          []string
	SystemDefaults systemdefaults.SystemDefaults
}

func (repo *OrgRepository) OrgByID(ctx context.Context, id string) (*org_model.OrgView, error) {
	org, err := repo.View.OrgByID(id)
	if err != nil {
		return nil, err
	}
	return model.OrgToModel(org), nil
}

func (repo *OrgRepository) OrgByDomainGlobal(ctx context.Context, domain string) (*org_model.OrgView, error) {
	verifiedDomain, err := repo.View.VerifiedOrgDomain(domain)
	if err != nil {
		return nil, err
	}
	return repo.OrgByID(ctx, verifiedDomain.OrgID)
}

func (repo *OrgRepository) CreateOrg(ctx context.Context, name string) (*org_model.Org, error) {
	org, aggregates, err := repo.OrgEventstore.PrepareCreateOrg(ctx, &org_model.Org{Name: name}, nil)
	if err != nil {
		return nil, err
	}

	member := org_model.NewOrgMemberWithRoles(org.AggregateID, authz.GetCtxData(ctx).UserID, orgOwnerRole)
	_, memberAggregate, err := repo.OrgEventstore.PrepareAddOrgMember(ctx, member, org.AggregateID)
	if err != nil {
		return nil, err
	}
	aggregates = append(aggregates, memberAggregate)

	err = sdk.PushAggregates(ctx, repo.Eventstore.PushAggregates, org.AppendEvents, aggregates...)
	if err != nil {
		return nil, err
	}

	return org_es_model.OrgToModel(org), nil
}

func (repo *OrgRepository) UpdateOrg(ctx context.Context, org *org_model.Org) (*org_model.Org, error) {
	return nil, errors.ThrowUnimplemented(nil, "EVENT-RkurR", "not implemented")
}

func (repo *OrgRepository) DeactivateOrg(ctx context.Context, id string) (*org_model.Org, error) {
	return repo.OrgEventstore.DeactivateOrg(ctx, id)
}

func (repo *OrgRepository) ReactivateOrg(ctx context.Context, id string) (*org_model.Org, error) {
	return repo.OrgEventstore.ReactivateOrg(ctx, id)
}

func (repo *OrgRepository) GetMyOrgIamPolicy(ctx context.Context) (*iam_model.OrgIAMPolicyView, error) {
	policy, err := repo.View.OrgIAMPolicyByAggregateID(authz.GetCtxData(ctx).OrgID)
	if errors.IsNotFound(err) {
		policy, err = repo.View.OrgIAMPolicyByAggregateID(repo.SystemDefaults.IamID)
		if err != nil {
			return nil, err
		}
		policy.Default = true
	}
	if err != nil {
		return nil, err
	}
	return iam_es_model.OrgIAMViewToModel(policy), err
}

func (repo *OrgRepository) SearchMyOrgDomains(ctx context.Context, request *org_model.OrgDomainSearchRequest) (*org_model.OrgDomainSearchResponse, error) {
	request.EnsureLimit(repo.SearchLimit)
	request.Queries = append(request.Queries, &org_model.OrgDomainSearchQuery{Key: org_model.OrgDomainSearchKeyOrgID, Method: global_model.SearchMethodEquals, Value: authz.GetCtxData(ctx).OrgID})
	sequence, sequenceErr := repo.View.GetLatestOrgDomainSequence()
	logging.Log("EVENT-SLowp").OnError(sequenceErr).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Warn("could not read latest org domain sequence")
	domains, count, err := repo.View.SearchOrgDomains(request)
	if err != nil {
		return nil, err
	}
	result := &org_model.OrgDomainSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: uint64(count),
		Result:      model.OrgDomainsToModel(domains),
	}
	if sequenceErr == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.LastSuccessfulSpoolerRun
	}
	return result, nil
}

func (repo *OrgRepository) AddMyOrgDomain(ctx context.Context, domain *org_model.OrgDomain) (*org_model.OrgDomain, error) {
	domain.AggregateID = authz.GetCtxData(ctx).OrgID
	return repo.OrgEventstore.AddOrgDomain(ctx, domain)
}

func (repo *OrgRepository) GenerateMyOrgDomainValidation(ctx context.Context, domain *org_model.OrgDomain) (string, string, error) {
	domain.AggregateID = authz.GetCtxData(ctx).OrgID
	return repo.OrgEventstore.GenerateOrgDomainValidation(ctx, domain)
}

func (repo *OrgRepository) ValidateMyOrgDomain(ctx context.Context, domain *org_model.OrgDomain) error {
	domain.AggregateID = authz.GetCtxData(ctx).OrgID
	users := func(ctx context.Context, domain string) ([]*es_models.Aggregate, error) {
		userIDs, err := repo.View.UserIDsByDomain(domain)
		if err != nil {
			return nil, err
		}
		return repo.UserEvents.PrepareDomainClaimed(ctx, userIDs)
	}
	return repo.OrgEventstore.ValidateOrgDomain(ctx, domain, users)
}

func (repo *OrgRepository) SetMyPrimaryOrgDomain(ctx context.Context, domain *org_model.OrgDomain) error {
	domain.AggregateID = authz.GetCtxData(ctx).OrgID
	return repo.OrgEventstore.SetPrimaryOrgDomain(ctx, domain)
}

func (repo *OrgRepository) RemoveMyOrgDomain(ctx context.Context, domain string) error {
	d := org_model.NewOrgDomain(authz.GetCtxData(ctx).OrgID, domain)
	return repo.OrgEventstore.RemoveOrgDomain(ctx, d)
}

func (repo *OrgRepository) OrgChanges(ctx context.Context, id string, lastSequence uint64, limit uint64, sortAscending bool) (*org_model.OrgChanges, error) {
	changes, err := repo.OrgEventstore.OrgChanges(ctx, id, lastSequence, limit, sortAscending)
	if err != nil {
		return nil, err
	}
	for _, change := range changes.Changes {
		change.ModifierName = change.ModifierId
		user, _ := repo.UserEvents.UserByID(ctx, change.ModifierId)
		if user != nil {
			if user.Human != nil {
				change.ModifierName = user.DisplayName
			}
			if user.Machine != nil {
				change.ModifierName = user.Machine.Name
			}
		}
	}
	return changes, nil
}

func (repo *OrgRepository) OrgMemberByID(ctx context.Context, orgID, userID string) (*org_model.OrgMemberView, error) {
	member, err := repo.View.OrgMemberByIDs(orgID, userID)
	if err != nil {
		return nil, err
	}
	return model.OrgMemberToModel(member), nil
}

func (repo *OrgRepository) AddMyOrgMember(ctx context.Context, member *org_model.OrgMember) (*org_model.OrgMember, error) {
	member.AggregateID = authz.GetCtxData(ctx).OrgID
	return repo.OrgEventstore.AddOrgMember(ctx, member)
}

func (repo *OrgRepository) ChangeMyOrgMember(ctx context.Context, member *org_model.OrgMember) (*org_model.OrgMember, error) {
	member.AggregateID = authz.GetCtxData(ctx).OrgID
	return repo.OrgEventstore.ChangeOrgMember(ctx, member)
}

func (repo *OrgRepository) RemoveMyOrgMember(ctx context.Context, userID string) error {
	member := org_model.NewOrgMember(authz.GetCtxData(ctx).OrgID, userID)
	return repo.OrgEventstore.RemoveOrgMember(ctx, member)
}

func (repo *OrgRepository) SearchMyOrgMembers(ctx context.Context, request *org_model.OrgMemberSearchRequest) (*org_model.OrgMemberSearchResponse, error) {
	request.EnsureLimit(repo.SearchLimit)
	request.Queries[len(request.Queries)-1] = &org_model.OrgMemberSearchQuery{Key: org_model.OrgMemberSearchKeyOrgID, Method: global_model.SearchMethodEquals, Value: authz.GetCtxData(ctx).OrgID}
	sequence, sequenceErr := repo.View.GetLatestOrgMemberSequence()
	logging.Log("EVENT-Smu3d").OnError(sequenceErr).Warn("could not read latest org member sequence")
	members, count, err := repo.View.SearchOrgMembers(request)
	if err != nil {
		return nil, err
	}
	result := &org_model.OrgMemberSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: count,
		Result:      model.OrgMembersToModel(members),
	}
	if sequenceErr == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.LastSuccessfulSpoolerRun
	}
	return result, nil
}

func (repo *OrgRepository) GetOrgMemberRoles() []string {
	roles := make([]string, 0)
	for _, roleMap := range repo.Roles {
		if strings.HasPrefix(roleMap, "ORG") {
			roles = append(roles, roleMap)
		}
	}
	return roles
}

func (repo *OrgRepository) IDPConfigByID(ctx context.Context, idpConfigID string) (*iam_model.IDPConfigView, error) {
	idp, err := repo.View.IDPConfigByID(idpConfigID)
	if err != nil {
		return nil, err
	}
	return iam_view_model.IDPConfigViewToModel(idp), nil
}
func (repo *OrgRepository) AddOIDCIDPConfig(ctx context.Context, idp *iam_model.IDPConfig) (*iam_model.IDPConfig, error) {
	idp.AggregateID = authz.GetCtxData(ctx).OrgID
	return repo.OrgEventstore.AddIDPConfig(ctx, idp)
}

func (repo *OrgRepository) ChangeIDPConfig(ctx context.Context, idp *iam_model.IDPConfig) (*iam_model.IDPConfig, error) {
	idp.AggregateID = authz.GetCtxData(ctx).OrgID
	return repo.OrgEventstore.ChangeIDPConfig(ctx, idp)
}

func (repo *OrgRepository) DeactivateIDPConfig(ctx context.Context, idpConfigID string) (*iam_model.IDPConfig, error) {
	return repo.OrgEventstore.DeactivateIDPConfig(ctx, authz.GetCtxData(ctx).OrgID, idpConfigID)
}

func (repo *OrgRepository) ReactivateIDPConfig(ctx context.Context, idpConfigID string) (*iam_model.IDPConfig, error) {
	return repo.OrgEventstore.ReactivateIDPConfig(ctx, authz.GetCtxData(ctx).OrgID, idpConfigID)
}

func (repo *OrgRepository) RemoveIDPConfig(ctx context.Context, idpConfigID string) error {
	aggregates := make([]*es_models.Aggregate, 0)
	idp := iam_model.NewIDPConfig(authz.GetCtxData(ctx).OrgID, idpConfigID)
	_, agg, err := repo.OrgEventstore.PrepareRemoveIDPConfig(ctx, idp)
	if err != nil {

	}
	aggregates = append(aggregates, agg)
	externalIDPs, err := repo.View.ExternalIDPsByIDPConfigID(idpConfigID)
	if err != nil {
		return err
	}
	for _, externalIDP := range externalIDPs {
		idpRemove := &usr_model.ExternalIDP{ObjectRoot: es_models.ObjectRoot{AggregateID: externalIDP.UserID}, IDPConfigID: externalIDP.IDPConfigID, UserID: externalIDP.ExternalUserID}
		idpAgg := make([]*es_models.Aggregate, 0)
		_, idpAgg, err = repo.UserEvents.PrepareRemoveExternalIDP(ctx, idpRemove, true)
		if err != nil {
			return err
		}
		aggregates = append(aggregates, idpAgg...)
	}
	return sdk.PushAggregates(ctx, repo.Eventstore.PushAggregates, nil, aggregates...)
}

func (repo *OrgRepository) ChangeOIDCIDPConfig(ctx context.Context, oidcConfig *iam_model.OIDCIDPConfig) (*iam_model.OIDCIDPConfig, error) {
	oidcConfig.AggregateID = authz.GetCtxData(ctx).OrgID
	return repo.OrgEventstore.ChangeIDPOIDCConfig(ctx, oidcConfig)
}

func (repo *OrgRepository) SearchIDPConfigs(ctx context.Context, request *iam_model.IDPConfigSearchRequest) (*iam_model.IDPConfigSearchResponse, error) {
	request.EnsureLimit(repo.SearchLimit)
	request.AppendMyOrgQuery(authz.GetCtxData(ctx).OrgID, repo.SystemDefaults.IamID)

	sequence, sequenceErr := repo.View.GetLatestIDPConfigSequence()
	logging.Log("EVENT-Dk8si").OnError(sequenceErr).Warn("could not read latest idp config sequence")
	idps, count, err := repo.View.SearchIDPConfigs(request)
	if err != nil {
		return nil, err
	}
	result := &iam_model.IDPConfigSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: count,
		Result:      iam_view_model.IdpConfigViewsToModel(idps),
	}
	if sequenceErr == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.LastSuccessfulSpoolerRun
	}
	return result, nil
}

func (repo *OrgRepository) GetLabelPolicy(ctx context.Context) (*iam_model.LabelPolicyView, error) {
	policy, err := repo.View.LabelPolicyByAggregateID(authz.GetCtxData(ctx).OrgID)
	if errors.IsNotFound(err) {
		policy, err = repo.View.LabelPolicyByAggregateID(repo.SystemDefaults.IamID)
		if err != nil {
			return nil, err
		}
		policy.Default = true
	}
	if err != nil {
		return nil, err
	}
	return iam_es_model.LabelPolicyViewToModel(policy), err
}

func (repo *OrgRepository) AddLabelPolicy(ctx context.Context, policy *iam_model.LabelPolicy) (*iam_model.LabelPolicy, error) {
	policy.AggregateID = authz.GetCtxData(ctx).OrgID
	return repo.OrgEventstore.AddLabelPolicy(ctx, policy)
}

func (repo *OrgRepository) ChangeLabelPolicy(ctx context.Context, policy *iam_model.LabelPolicy) (*iam_model.LabelPolicy, error) {
	policy.AggregateID = authz.GetCtxData(ctx).OrgID
	return repo.OrgEventstore.ChangeLabelPolicy(ctx, policy)
}

func (repo *OrgRepository) GetLoginPolicy(ctx context.Context) (*iam_model.LoginPolicyView, error) {
	policy, viewErr := repo.View.LoginPolicyByAggregateID(authz.GetCtxData(ctx).OrgID)
	if viewErr != nil && !errors.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if errors.IsNotFound(viewErr) {
		policy = new(iam_es_model.LoginPolicyView)
	}
	events, esErr := repo.OrgEventstore.OrgEventsByID(ctx, repo.SystemDefaults.IamID, policy.Sequence)
	if errors.IsNotFound(viewErr) && len(events) == 0 {
		return repo.GetDefaultLoginPolicy(ctx)
	}
	if esErr != nil {
		logging.Log("EVENT-38iTr").WithError(esErr).Debug("error retrieving new events")
		return iam_es_model.LoginPolicyViewToModel(policy), nil
	}
	policyCopy := *policy
	for _, event := range events {
		if err := policyCopy.AppendEvent(event); err != nil {
			return iam_es_model.LoginPolicyViewToModel(policy), nil
		}
	}
	return iam_es_model.LoginPolicyViewToModel(policy), nil
}

func (repo *OrgRepository) GetIDPProvidersByIDPConfigID(ctx context.Context, aggregateID, idpConfigID string) ([]*iam_model.IDPProviderView, error) {
	idpProviders, err := repo.View.IDPProvidersByIdpConfigID(aggregateID, idpConfigID)
	if err != nil {
		return nil, err
	}
	return iam_view_model.IDPProviderViewsToModel(idpProviders), err
}

func (repo *OrgRepository) GetDefaultLoginPolicy(ctx context.Context) (*iam_model.LoginPolicyView, error) {
	policy, viewErr := repo.View.LoginPolicyByAggregateID(domain.IAMID)
	if viewErr != nil && !errors.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if errors.IsNotFound(viewErr) {
		policy = new(iam_es_model.LoginPolicyView)
	}
	events, esErr := repo.IAMEventstore.IAMEventsByID(ctx, domain.IAMID, policy.Sequence)
	if errors.IsNotFound(viewErr) && len(events) == 0 {
		return nil, errors.ThrowNotFound(nil, "EVENT-cmO9s", "Errors.IAM.LoginPolicy.NotFound")
	}
	if esErr != nil {
		logging.Log("EVENT-28uLp").WithError(esErr).Debug("error retrieving new events")
		return iam_es_model.LoginPolicyViewToModel(policy), nil
	}
	policyCopy := *policy
	for _, event := range events {
		if err := policyCopy.AppendEvent(event); err != nil {
			return iam_es_model.LoginPolicyViewToModel(policy), nil
		}
	}
	policy.Default = true
	return iam_es_model.LoginPolicyViewToModel(policy), nil
}

func (repo *OrgRepository) AddLoginPolicy(ctx context.Context, policy *iam_model.LoginPolicy) (*iam_model.LoginPolicy, error) {
	policy.AggregateID = authz.GetCtxData(ctx).OrgID
	return repo.OrgEventstore.AddLoginPolicy(ctx, policy)
}

func (repo *OrgRepository) ChangeLoginPolicy(ctx context.Context, policy *iam_model.LoginPolicy) (*iam_model.LoginPolicy, error) {
	policy.AggregateID = authz.GetCtxData(ctx).OrgID
	return repo.OrgEventstore.ChangeLoginPolicy(ctx, policy)
}

func (repo *OrgRepository) RemoveLoginPolicy(ctx context.Context) error {
	policy := &iam_model.LoginPolicy{ObjectRoot: models.ObjectRoot{
		AggregateID: authz.GetCtxData(ctx).OrgID,
	}}
	return repo.OrgEventstore.RemoveLoginPolicy(ctx, policy)
}

func (repo *OrgRepository) SearchIDPProviders(ctx context.Context, request *iam_model.IDPProviderSearchRequest) (*iam_model.IDPProviderSearchResponse, error) {
	policy, err := repo.View.LoginPolicyByAggregateID(authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	if policy.Default {
		request.AppendAggregateIDQuery(domain.IAMID)
	} else {
		request.AppendAggregateIDQuery(authz.GetCtxData(ctx).OrgID)
	}
	request.EnsureLimit(repo.SearchLimit)
	sequence, sequenceErr := repo.View.GetLatestIDPProviderSequence()
	logging.Log("EVENT-Tuiks").OnError(sequenceErr).Warn("could not read latest iam sequence")
	providers, count, err := repo.View.SearchIDPProviders(request)
	if err != nil {
		return nil, err
	}
	result := &iam_model.IDPProviderSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: count,
		Result:      iam_es_model.IDPProviderViewsToModel(providers),
	}
	if sequenceErr == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.LastSuccessfulSpoolerRun
	}
	return result, nil
}

func (repo *OrgRepository) AddIDPProviderToLoginPolicy(ctx context.Context, provider *iam_model.IDPProvider) (*iam_model.IDPProvider, error) {
	provider.AggregateID = authz.GetCtxData(ctx).OrgID
	return repo.OrgEventstore.AddIDPProviderToLoginPolicy(ctx, provider)
}

func (repo *OrgRepository) RemoveIDPProviderFromIdpProvider(ctx context.Context, provider *iam_model.IDPProvider) error {
	aggregates := make([]*es_models.Aggregate, 0)
	provider.AggregateID = authz.GetCtxData(ctx).OrgID
	_, agg, err := repo.OrgEventstore.PrepareRemoveIDPProviderFromLoginPolicy(ctx, provider, false)
	if err != nil {
		return err
	}
	aggregates = append(aggregates, agg)
	externalIDPs, err := repo.View.ExternalIDPsByIDPConfigID(provider.IDPConfigID)
	if err != nil {
		return err
	}
	for _, externalIDP := range externalIDPs {
		idpRemove := &usr_model.ExternalIDP{ObjectRoot: es_models.ObjectRoot{AggregateID: externalIDP.UserID}, IDPConfigID: externalIDP.IDPConfigID, UserID: externalIDP.ExternalUserID}
		idpAgg := make([]*es_models.Aggregate, 0)
		_, idpAgg, err = repo.UserEvents.PrepareRemoveExternalIDP(ctx, idpRemove, true)
		if err != nil {
			return err
		}
		aggregates = append(aggregates, idpAgg...)
	}
	return sdk.PushAggregates(ctx, repo.Eventstore.PushAggregates, nil, aggregates...)
}

func (repo *OrgRepository) SearchSecondFactors(ctx context.Context) (*iam_model.SecondFactorsSearchResponse, error) {
	policy, err := repo.GetLoginPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return &iam_model.SecondFactorsSearchResponse{
		TotalResult: uint64(len(policy.SecondFactors)),
		Result:      policy.SecondFactors,
	}, nil
}

func (repo *OrgRepository) AddSecondFactorToLoginPolicy(ctx context.Context, mfa iam_model.SecondFactorType) (iam_model.SecondFactorType, error) {
	return repo.OrgEventstore.AddSecondFactorToLoginPolicy(ctx, authz.GetCtxData(ctx).OrgID, mfa)
}

func (repo *OrgRepository) RemoveSecondFactorFromLoginPolicy(ctx context.Context, mfa iam_model.SecondFactorType) error {
	return repo.OrgEventstore.RemoveSecondFactorFromLoginPolicy(ctx, authz.GetCtxData(ctx).OrgID, mfa)
}

func (repo *OrgRepository) SearchMultiFactors(ctx context.Context) (*iam_model.MultiFactorsSearchResponse, error) {
	policy, err := repo.GetLoginPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return &iam_model.MultiFactorsSearchResponse{
		TotalResult: uint64(len(policy.MultiFactors)),
		Result:      policy.MultiFactors,
	}, nil
}

func (repo *OrgRepository) AddMultiFactorToLoginPolicy(ctx context.Context, mfa iam_model.MultiFactorType) (iam_model.MultiFactorType, error) {
	return repo.OrgEventstore.AddMultiFactorToLoginPolicy(ctx, authz.GetCtxData(ctx).OrgID, mfa)
}

func (repo *OrgRepository) RemoveMultiFactorFromLoginPolicy(ctx context.Context, mfa iam_model.MultiFactorType) error {
	return repo.OrgEventstore.RemoveMultiFactorFromLoginPolicy(ctx, authz.GetCtxData(ctx).OrgID, mfa)
}

func (repo *OrgRepository) GetPasswordComplexityPolicy(ctx context.Context) (*iam_model.PasswordComplexityPolicyView, error) {
	policy, viewErr := repo.View.PasswordComplexityPolicyByAggregateID(authz.GetCtxData(ctx).OrgID)
	if viewErr != nil && !errors.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if errors.IsNotFound(viewErr) {
		policy = new(iam_es_model.PasswordComplexityPolicyView)
	}
	events, esErr := repo.OrgEventstore.OrgEventsByID(ctx, repo.SystemDefaults.IamID, policy.Sequence)
	if errors.IsNotFound(viewErr) && len(events) == 0 {
		return repo.GetDefaultPasswordComplexityPolicy(ctx)
	}
	if esErr != nil {
		logging.Log("EVENT-1Bx8s").WithError(esErr).Debug("error retrieving new events")
		return iam_es_model.PasswordComplexityViewToModel(policy), nil
	}
	policyCopy := *policy
	for _, event := range events {
		if err := policyCopy.AppendEvent(event); err != nil {
			return iam_es_model.PasswordComplexityViewToModel(policy), nil
		}
	}
	return iam_es_model.PasswordComplexityViewToModel(policy), nil
}

func (repo *OrgRepository) GetDefaultPasswordComplexityPolicy(ctx context.Context) (*iam_model.PasswordComplexityPolicyView, error) {
	policy, viewErr := repo.View.PasswordComplexityPolicyByAggregateID(repo.SystemDefaults.IamID)
	if viewErr != nil && !errors.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if errors.IsNotFound(viewErr) {
		policy = new(iam_es_model.PasswordComplexityPolicyView)
	}
	events, esErr := repo.IAMEventstore.IAMEventsByID(ctx, repo.SystemDefaults.IamID, policy.Sequence)
	if errors.IsNotFound(viewErr) && len(events) == 0 {
		return nil, errors.ThrowNotFound(nil, "EVENT-cmO9s", "Errors.IAM.PasswordComplexityPolicy.NotFound")
	}
	if esErr != nil {
		logging.Log("EVENT-pL9sw").WithError(esErr).Debug("error retrieving new events")
		return iam_es_model.PasswordComplexityViewToModel(policy), nil
	}
	policyCopy := *policy
	for _, event := range events {
		if err := policyCopy.AppendEvent(event); err != nil {
			return iam_es_model.PasswordComplexityViewToModel(policy), nil
		}
	}
	policy.Default = true
	return iam_es_model.PasswordComplexityViewToModel(policy), nil
}

func (repo *OrgRepository) AddPasswordComplexityPolicy(ctx context.Context, policy *iam_model.PasswordComplexityPolicy) (*iam_model.PasswordComplexityPolicy, error) {
	policy.AggregateID = authz.GetCtxData(ctx).OrgID
	return repo.OrgEventstore.AddPasswordComplexityPolicy(ctx, policy)
}

func (repo *OrgRepository) ChangePasswordComplexityPolicy(ctx context.Context, policy *iam_model.PasswordComplexityPolicy) (*iam_model.PasswordComplexityPolicy, error) {
	policy.AggregateID = authz.GetCtxData(ctx).OrgID
	return repo.OrgEventstore.ChangePasswordComplexityPolicy(ctx, policy)
}

func (repo *OrgRepository) RemovePasswordComplexityPolicy(ctx context.Context) error {
	policy := &iam_model.PasswordComplexityPolicy{ObjectRoot: models.ObjectRoot{
		AggregateID: authz.GetCtxData(ctx).OrgID,
	}}
	return repo.OrgEventstore.RemovePasswordComplexityPolicy(ctx, policy)
}

func (repo *OrgRepository) GetPasswordAgePolicy(ctx context.Context) (*iam_model.PasswordAgePolicyView, error) {
	policy, viewErr := repo.View.PasswordAgePolicyByAggregateID(authz.GetCtxData(ctx).OrgID)
	if viewErr != nil && !errors.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if errors.IsNotFound(viewErr) {
		policy = new(iam_es_model.PasswordAgePolicyView)
	}
	events, esErr := repo.OrgEventstore.OrgEventsByID(ctx, repo.SystemDefaults.IamID, policy.Sequence)
	if errors.IsNotFound(viewErr) && len(events) == 0 {
		return repo.GetDefaultPasswordAgePolicy(ctx)
	}
	if esErr != nil {
		logging.Log("EVENT-5Mx7s").WithError(esErr).Debug("error retrieving new events")
		return iam_es_model.PasswordAgeViewToModel(policy), nil
	}
	policyCopy := *policy
	for _, event := range events {
		if err := policyCopy.AppendEvent(event); err != nil {
			return iam_es_model.PasswordAgeViewToModel(policy), nil
		}
	}
	return iam_es_model.PasswordAgeViewToModel(policy), nil
}

func (repo *OrgRepository) GetDefaultPasswordAgePolicy(ctx context.Context) (*iam_model.PasswordAgePolicyView, error) {
	policy, viewErr := repo.View.PasswordAgePolicyByAggregateID(repo.SystemDefaults.IamID)
	if viewErr != nil && !errors.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if errors.IsNotFound(viewErr) {
		policy = new(iam_es_model.PasswordAgePolicyView)
	}
	events, esErr := repo.IAMEventstore.IAMEventsByID(ctx, repo.SystemDefaults.IamID, policy.Sequence)
	if errors.IsNotFound(viewErr) && len(events) == 0 {
		return nil, errors.ThrowNotFound(nil, "EVENT-cmO9s", "Errors.IAM.PasswordAgePolicy.NotFound")
	}
	if esErr != nil {
		logging.Log("EVENT-3I90s").WithError(esErr).Debug("error retrieving new events")
		return iam_es_model.PasswordAgeViewToModel(policy), nil
	}
	policyCopy := *policy
	for _, event := range events {
		if err := policyCopy.AppendEvent(event); err != nil {
			return iam_es_model.PasswordAgeViewToModel(policy), nil
		}
	}
	policy.Default = true
	return iam_es_model.PasswordAgeViewToModel(policy), nil
}

func (repo *OrgRepository) AddPasswordAgePolicy(ctx context.Context, policy *iam_model.PasswordAgePolicy) (*iam_model.PasswordAgePolicy, error) {
	policy.AggregateID = authz.GetCtxData(ctx).OrgID
	return repo.OrgEventstore.AddPasswordAgePolicy(ctx, policy)
}

func (repo *OrgRepository) ChangePasswordAgePolicy(ctx context.Context, policy *iam_model.PasswordAgePolicy) (*iam_model.PasswordAgePolicy, error) {
	policy.AggregateID = authz.GetCtxData(ctx).OrgID
	return repo.OrgEventstore.ChangePasswordAgePolicy(ctx, policy)
}

func (repo *OrgRepository) RemovePasswordAgePolicy(ctx context.Context) error {
	policy := &iam_model.PasswordAgePolicy{ObjectRoot: models.ObjectRoot{
		AggregateID: authz.GetCtxData(ctx).OrgID,
	}}
	return repo.OrgEventstore.RemovePasswordAgePolicy(ctx, policy)
}

func (repo *OrgRepository) GetPasswordLockoutPolicy(ctx context.Context) (*iam_model.PasswordLockoutPolicyView, error) {
	policy, viewErr := repo.View.PasswordLockoutPolicyByAggregateID(authz.GetCtxData(ctx).OrgID)
	if viewErr != nil && !errors.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if errors.IsNotFound(viewErr) {
		policy = new(iam_es_model.PasswordLockoutPolicyView)
	}
	events, esErr := repo.OrgEventstore.OrgEventsByID(ctx, repo.SystemDefaults.IamID, policy.Sequence)
	if errors.IsNotFound(viewErr) && len(events) == 0 {
		return repo.GetDefaultPasswordLockoutPolicy(ctx)
	}
	if esErr != nil {
		logging.Log("EVENT-mS9od").WithError(esErr).Debug("error retrieving new events")
		return iam_es_model.PasswordLockoutViewToModel(policy), nil
	}
	policyCopy := *policy
	for _, event := range events {
		if err := policyCopy.AppendEvent(event); err != nil {
			return iam_es_model.PasswordLockoutViewToModel(policy), nil
		}
	}
	return iam_es_model.PasswordLockoutViewToModel(policy), nil
}

func (repo *OrgRepository) GetDefaultPasswordLockoutPolicy(ctx context.Context) (*iam_model.PasswordLockoutPolicyView, error) {
	policy, viewErr := repo.View.PasswordLockoutPolicyByAggregateID(repo.SystemDefaults.IamID)
	if viewErr != nil && !errors.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if errors.IsNotFound(viewErr) {
		policy = new(iam_es_model.PasswordLockoutPolicyView)
	}
	events, esErr := repo.IAMEventstore.IAMEventsByID(ctx, repo.SystemDefaults.IamID, policy.Sequence)
	if errors.IsNotFound(viewErr) && len(events) == 0 {
		return nil, errors.ThrowNotFound(nil, "EVENT-cmO9s", "Errors.IAM.PasswordLockoutPolicy.NotFound")
	}
	if esErr != nil {
		logging.Log("EVENT-2Ms9f").WithError(esErr).Debug("error retrieving new events")
		return iam_es_model.PasswordLockoutViewToModel(policy), nil
	}
	policyCopy := *policy
	for _, event := range events {
		if err := policyCopy.AppendEvent(event); err != nil {
			return iam_es_model.PasswordLockoutViewToModel(policy), nil
		}
	}
	policy.Default = true
	return iam_es_model.PasswordLockoutViewToModel(policy), nil
}

func (repo *OrgRepository) AddPasswordLockoutPolicy(ctx context.Context, policy *iam_model.PasswordLockoutPolicy) (*iam_model.PasswordLockoutPolicy, error) {
	policy.AggregateID = authz.GetCtxData(ctx).OrgID
	return repo.OrgEventstore.AddPasswordLockoutPolicy(ctx, policy)
}

func (repo *OrgRepository) ChangePasswordLockoutPolicy(ctx context.Context, policy *iam_model.PasswordLockoutPolicy) (*iam_model.PasswordLockoutPolicy, error) {
	policy.AggregateID = authz.GetCtxData(ctx).OrgID
	return repo.OrgEventstore.ChangePasswordLockoutPolicy(ctx, policy)
}

func (repo *OrgRepository) RemovePasswordLockoutPolicy(ctx context.Context) error {
	policy := &iam_model.PasswordLockoutPolicy{ObjectRoot: models.ObjectRoot{
		AggregateID: authz.GetCtxData(ctx).OrgID,
	}}
	return repo.OrgEventstore.RemovePasswordLockoutPolicy(ctx, policy)
}

func (repo *OrgRepository) GetDefaultMailTemplate(ctx context.Context) (*iam_model.MailTemplateView, error) {
	template, err := repo.View.MailTemplateByAggregateID(repo.SystemDefaults.IamID)
	if err != nil {
		return nil, err
	}
	template.Default = true
	return iam_es_model.MailTemplateViewToModel(template), err
}

func (repo *OrgRepository) GetMailTemplate(ctx context.Context) (*iam_model.MailTemplateView, error) {
	template, err := repo.View.MailTemplateByAggregateID(authz.GetCtxData(ctx).OrgID)
	if errors.IsNotFound(err) {
		template, err = repo.View.MailTemplateByAggregateID(repo.SystemDefaults.IamID)
		if err != nil {
			return nil, err
		}
		template.Default = true
	}
	if err != nil {
		return nil, err
	}
	return iam_es_model.MailTemplateViewToModel(template), err
}

func (repo *OrgRepository) AddMailTemplate(ctx context.Context, template *iam_model.MailTemplate) (*iam_model.MailTemplate, error) {
	template.AggregateID = authz.GetCtxData(ctx).OrgID
	return repo.OrgEventstore.AddMailTemplate(ctx, template)
}

func (repo *OrgRepository) ChangeMailTemplate(ctx context.Context, template *iam_model.MailTemplate) (*iam_model.MailTemplate, error) {
	template.AggregateID = authz.GetCtxData(ctx).OrgID
	return repo.OrgEventstore.ChangeMailTemplate(ctx, template)
}

func (repo *OrgRepository) RemoveMailTemplate(ctx context.Context) error {
	template := &iam_model.MailTemplate{ObjectRoot: models.ObjectRoot{
		AggregateID: authz.GetCtxData(ctx).OrgID,
	}}
	return repo.OrgEventstore.RemoveMailTemplate(ctx, template)
}

func (repo *OrgRepository) GetDefaultMailTexts(ctx context.Context) (*iam_model.MailTextsView, error) {
	texts, err := repo.View.MailTextsByAggregateID(repo.SystemDefaults.IamID)
	if err != nil {
		return nil, err
	}
	return iam_es_model.MailTextsViewToModel(texts, true), err
}

func (repo *OrgRepository) GetMailTexts(ctx context.Context) (*iam_model.MailTextsView, error) {
	defaultIn := false
	texts, err := repo.View.MailTextsByAggregateID(authz.GetCtxData(ctx).OrgID)
	if errors.IsNotFound(err) || len(texts) == 0 {
		texts, err = repo.View.MailTextsByAggregateID(repo.SystemDefaults.IamID)
		if err != nil {
			return nil, err
		}
		defaultIn = true
	}
	if err != nil {
		return nil, err
	}
	return iam_es_model.MailTextsViewToModel(texts, defaultIn), err
}

func (repo *OrgRepository) AddMailText(ctx context.Context, text *iam_model.MailText) (*iam_model.MailText, error) {
	text.AggregateID = authz.GetCtxData(ctx).OrgID
	return repo.OrgEventstore.AddMailText(ctx, text)
}

func (repo *OrgRepository) ChangeMailText(ctx context.Context, text *iam_model.MailText) (*iam_model.MailText, error) {
	text.AggregateID = authz.GetCtxData(ctx).OrgID
	return repo.OrgEventstore.ChangeMailText(ctx, text)
}

func (repo *OrgRepository) RemoveMailText(ctx context.Context, text *iam_model.MailText) error {
	text.AggregateID = authz.GetCtxData(ctx).OrgID
	return repo.OrgEventstore.RemoveMailText(ctx, text)
}
