package eventstore

import (
	"context"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/eventstore/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	iam_view_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"strings"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/sdk"
	mgmt_view "github.com/caos/zitadel/internal/management/repository/eventsourcing/view"
	global_model "github.com/caos/zitadel/internal/model"
	org_model "github.com/caos/zitadel/internal/org/model"
	org_es "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/org/repository/view/model"
	usr_es "github.com/caos/zitadel/internal/user/repository/eventsourcing"
)

const (
	orgOwnerRole = "ORG_OWNER"
)

type OrgRepository struct {
	SearchLimit uint64
	*org_es.OrgEventstore
	UserEvents     *usr_es.UserEventstore
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

func (repo *OrgRepository) GetMyOrgIamPolicy(ctx context.Context) (*org_model.OrgIAMPolicy, error) {
	return repo.OrgEventstore.GetOrgIAMPolicy(ctx, authz.GetCtxData(ctx).OrgID)
}

func (repo *OrgRepository) SearchMyOrgDomains(ctx context.Context, request *org_model.OrgDomainSearchRequest) (*org_model.OrgDomainSearchResponse, error) {
	request.EnsureLimit(repo.SearchLimit)
	request.Queries = append(request.Queries, &org_model.OrgDomainSearchQuery{Key: org_model.OrgDomainSearchKeyOrgID, Method: global_model.SearchMethodEquals, Value: authz.GetCtxData(ctx).OrgID})
	sequence, sequenceErr := repo.View.GetLatestOrgDomainSequence()
	logging.Log("EVENT-SLowp").OnError(sequenceErr).Warn("could not read latest org domain sequence")
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
		result.Timestamp = sequence.CurrentTimestamp
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
			change.ModifierName = user.DisplayName
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
		result.Timestamp = sequence.CurrentTimestamp
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
		result.Timestamp = sequence.CurrentTimestamp
	}
	return result, nil
}

func (repo *OrgRepository) GetLoginPolicy(ctx context.Context) (*iam_model.LoginPolicyView, error) {
	policy, err := repo.View.LoginPolicyByAggregateID(authz.GetCtxData(ctx).OrgID)
	if errors.IsNotFound(err) {
		policy, err = repo.View.LoginPolicyByAggregateID(repo.SystemDefaults.IamID)
		if err != nil {
			return nil, err
		}
		policy.Default = true
	}
	if err != nil {
		return nil, err
	}
	return iam_es_model.LoginPolicyViewToModel(policy), err
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
	_, err := repo.View.LoginPolicyByAggregateID(authz.GetCtxData(ctx).OrgID)
	if err != nil {
		if errors.IsNotFound(err) {
			request.AppendAggregateIDQuery(repo.SystemDefaults.IamID)
		}
	} else {
		request.AppendAggregateIDQuery(authz.GetCtxData(ctx).OrgID)
	}
	request.EnsureLimit(repo.SearchLimit)
	sequence, sequenceErr := repo.View.GetLatestIdpProviderSequence()
	logging.Log("EVENT-Tuiks").OnError(sequenceErr).Warn("could not read latest iam sequence")
	providers, count, err := repo.View.SearchIdpProviders(request)
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
		result.Timestamp = sequence.CurrentTimestamp
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
	externalIDPs, err := repo.View.ExternalIDPsByIDPConfigID(provider.IdpConfigID)
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
