package eventsourcing

import (
	"context"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	"github.com/caos/zitadel/internal/id"
	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
)

type OrgEventstore struct {
	eventstore.Eventstore
	IAMDomain           string
	idGenerator         id.Generator
	defaultOrgIamPolicy *org_model.OrgIamPolicy
}

type OrgConfig struct {
	eventstore.Eventstore
	IAMDomain string
}

func StartOrg(conf OrgConfig, defaults systemdefaults.SystemDefaults) *OrgEventstore {
	policy := defaults.DefaultPolicies.OrgIam
	policy.Default = true
	return &OrgEventstore{
		Eventstore:          conf.Eventstore,
		idGenerator:         id.SonyFlakeGenerator,
		IAMDomain:           conf.IAMDomain,
		defaultOrgIamPolicy: &policy,
	}
}

func (es *OrgEventstore) PrepareCreateOrg(ctx context.Context, orgModel *org_model.Org) (*model.Org, []*es_models.Aggregate, error) {
	if orgModel == nil || !orgModel.IsValid() {
		return nil, nil, errors.ThrowInvalidArgument(nil, "EVENT-OeLSk", "org not valid")
	}
	orgModel.AddIAMDomain(es.IAMDomain)

	id, err := es.idGenerator.Next()
	if err != nil {
		return nil, nil, errors.ThrowInternal(err, "EVENT-OwciI", "id gen failed")
	}
	orgModel.AggregateID = id
	org := model.OrgFromModel(orgModel)

	aggregates, err := orgCreatedAggregates(ctx, es.AggregateCreator(), org)

	return org, aggregates, err
}

func (es *OrgEventstore) CreateOrg(ctx context.Context, orgModel *org_model.Org) (*org_model.Org, error) {
	org, aggregates, err := es.PrepareCreateOrg(ctx, orgModel)
	err = es_sdk.PushAggregates(ctx, es.PushAggregates, org.AppendEvents, aggregates...)
	if err != nil {
		return nil, err
	}

	return model.OrgToModel(org), nil
}

func (es *OrgEventstore) OrgByID(ctx context.Context, org *org_model.Org) (*org_model.Org, error) {
	if org == nil {
		return nil, errors.ThrowInvalidArgument(nil, "EVENT-gQTYP", "org not set")
	}
	query, err := OrgByIDQuery(org.AggregateID, org.Sequence)
	if err != nil {
		return nil, err
	}

	esOrg := model.OrgFromModel(org)
	err = es_sdk.Filter(ctx, es.FilterEvents, esOrg.AppendEvents, query)
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}
	if esOrg.Sequence == 0 {
		return nil, errors.ThrowNotFound(nil, "EVENT-kVLb2", "org not found")
	}

	return model.OrgToModel(esOrg), nil
}

func (es *OrgEventstore) IsOrgUnique(ctx context.Context, name, domain string) (isUnique bool, err error) {
	var found bool
	err = es_sdk.Filter(ctx, es.FilterEvents, isUniqueValidation(&found), OrgNameUniqueQuery(name))
	if (err != nil && !errors.IsNotFound(err)) || found {
		return false, err
	}

	err = es_sdk.Filter(ctx, es.FilterEvents, isUniqueValidation(&found), OrgDomainUniqueQuery(domain))
	if err != nil && !errors.IsNotFound(err) {
		return false, err
	}

	return !found, nil
}

func isUniqueValidation(unique *bool) func(events ...*es_models.Event) error {
	return func(events ...*es_models.Event) error {
		if len(events) == 0 {
			return nil
		}
		*unique = *unique || events[0].Type == model.OrgDomainReserved || events[0].Type == model.OrgNameReserved

		return nil
	}
}

func (es *OrgEventstore) DeactivateOrg(ctx context.Context, orgID string) (*org_model.Org, error) {
	existingOrg, err := es.OrgByID(ctx, org_model.NewOrg(orgID))
	if err != nil {
		return nil, errors.ThrowInvalidArgument(nil, "EVENT-oL9nT", "org not found")
	}
	org := model.OrgFromModel(existingOrg)

	aggregate := orgDeactivateAggregate(es.AggregateCreator(), org)
	err = es_sdk.Push(ctx, es.PushAggregates, org.AppendEvents, aggregate)
	if err != nil {
		return nil, err
	}

	return model.OrgToModel(org), nil
}

func (es *OrgEventstore) ReactivateOrg(ctx context.Context, orgID string) (*org_model.Org, error) {
	existingOrg, err := es.OrgByID(ctx, org_model.NewOrg(orgID))
	if err != nil {
		return nil, errors.ThrowInvalidArgument(nil, "EVENT-oL9nT", "org not set")
	}
	org := model.OrgFromModel(existingOrg)

	aggregate := orgReactivateAggregate(es.AggregateCreator(), org)
	err = es_sdk.Push(ctx, es.PushAggregates, org.AppendEvents, aggregate)
	if err != nil {
		return nil, err

	}
	return model.OrgToModel(org), nil
}

func (es *OrgEventstore) AddOrgDomain(ctx context.Context, domain *org_model.OrgDomain) (*org_model.OrgDomain, error) {
	if !domain.IsValid() {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-8sFJW", "domain is invalid")
	}
	existing, err := es.OrgByID(ctx, org_model.NewOrg(domain.AggregateID))
	if err != nil {
		return nil, err
	}
	repoOrg := model.OrgFromModel(existing)
	repoDomain := model.OrgDomainFromModel(domain)
	orgAggregates, err := OrgDomainAddedAggregate(ctx, es.Eventstore.AggregateCreator(), repoOrg, repoDomain)
	if err != nil {
		return nil, err
	}
	err = es_sdk.PushAggregates(ctx, es.PushAggregates, repoOrg.AppendEvents, orgAggregates...)
	if err != nil {
		return nil, err
	}

	if _, d := model.GetDomain(repoOrg.Domains, domain.Domain); d != nil {
		return model.OrgDomainToModel(d), nil
	}
	return nil, errors.ThrowInternal(nil, "EVENT-ISOP0", "Could not find org in list")
}

func (es *OrgEventstore) RemoveOrgDomain(ctx context.Context, domain *org_model.OrgDomain) error {
	if domain.Domain == "" {
		return errors.ThrowPreconditionFailed(nil, "EVENT-SJsK3", "Domain is required")
	}
	existing, err := es.OrgByID(ctx, org_model.NewOrg(domain.AggregateID))
	if err != nil {
		return err
	}
	if !existing.ContainsDomain(domain) {
		return errors.ThrowPreconditionFailed(nil, "EVENT-Sjdi3", "Domain doesn't exist on project")
	}
	repoOrg := model.OrgFromModel(existing)
	repoDomain := model.OrgDomainFromModel(domain)
	orgAggregates, err := OrgDomainRemovedAggregate(ctx, es.Eventstore.AggregateCreator(), repoOrg, repoDomain)
	if err != nil {
		return err
	}
	err = es_sdk.PushAggregates(ctx, es.PushAggregates, repoOrg.AppendEvents, orgAggregates...)
	if err != nil {
		return err
	}
	return nil
}

func (es *OrgEventstore) OrgMemberByIDs(ctx context.Context, member *org_model.OrgMember) (*org_model.OrgMember, error) {
	if member == nil || member.UserID == "" || member.AggregateID == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-ld93d", "member not set")
	}

	org, err := es.OrgByID(ctx, &org_model.Org{ObjectRoot: member.ObjectRoot, Members: []*org_model.OrgMember{member}})
	if err != nil {
		return nil, err
	}

	for _, currentMember := range org.Members {
		if currentMember.UserID == member.UserID {
			return currentMember, nil
		}
	}

	return nil, errors.ThrowNotFound(nil, "EVENT-SXji6", "member not found")
}

func (es *OrgEventstore) PrepareAddOrgMember(ctx context.Context, member *org_model.OrgMember) (*model.OrgMember, *es_models.Aggregate, error) {
	if member == nil || !member.IsValid() {
		return nil, nil, errors.ThrowPreconditionFailed(nil, "EVENT-9dk45", "UserID and Roles are required")
	}

	repoMember := model.OrgMemberFromModel(member)
	addAggregate, err := orgMemberAddedAggregate(ctx, es.Eventstore.AggregateCreator(), repoMember)

	return repoMember, addAggregate, err
}

func (es *OrgEventstore) AddOrgMember(ctx context.Context, member *org_model.OrgMember) (*org_model.OrgMember, error) {
	repoMember, addAggregate, err := es.PrepareAddOrgMember(ctx, member)
	if err != nil {
		return nil, err
	}
	err = es_sdk.PushAggregates(ctx, es.PushAggregates, repoMember.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}

	return model.OrgMemberToModel(repoMember), nil
}

func (es *OrgEventstore) ChangeOrgMember(ctx context.Context, member *org_model.OrgMember) (*org_model.OrgMember, error) {
	if member == nil || !member.IsValid() {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-9dk45", "UserID and Roles are required")
	}

	existingMember, err := es.OrgMemberByIDs(ctx, member)
	if err != nil {
		return nil, err
	}

	member.ObjectRoot = existingMember.ObjectRoot
	repoMember := model.OrgMemberFromModel(member)
	repoExistingMember := model.OrgMemberFromModel(existingMember)

	orgAggregate := orgMemberChangedAggregate(es.Eventstore.AggregateCreator(), repoExistingMember, repoMember)
	err = es_sdk.Push(ctx, es.PushAggregates, repoMember.AppendEvents, orgAggregate)
	if err != nil {
		return nil, err
	}

	return model.OrgMemberToModel(repoMember), nil
}

func (es *OrgEventstore) RemoveOrgMember(ctx context.Context, member *org_model.OrgMember) error {
	if member == nil || member.UserID == "" {
		return errors.ThrowInvalidArgument(nil, "EVENT-d43fs", "UserID is required")
	}

	existingMember, err := es.OrgMemberByIDs(ctx, member)
	if errors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}

	member.ObjectRoot = existingMember.ObjectRoot
	repoMember := model.OrgMemberFromModel(member)

	orgAggregate := orgMemberRemovedAggregate(es.Eventstore.AggregateCreator(), repoMember)
	return es_sdk.Push(ctx, es.PushAggregates, repoMember.AppendEvents, orgAggregate)
}

func (es *OrgEventstore) GetOrgIamPolicy(ctx context.Context, orgID string) (*org_model.OrgIamPolicy, error) {
	existing, err := es.OrgByID(ctx, org_model.NewOrg(orgID))
	if err != nil {
		return nil, err
	}
	if existing.OrgIamPolicy != nil {
		return existing.OrgIamPolicy, nil
	}
	return es.defaultOrgIamPolicy, nil
}

func (es *OrgEventstore) AddOrgIamPolicy(ctx context.Context, policy *org_model.OrgIamPolicy) (*org_model.OrgIamPolicy, error) {
	existing, err := es.OrgByID(ctx, org_model.NewOrg(policy.AggregateID))
	if err != nil {
		return nil, err
	}
	if existing.OrgIamPolicy != nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-7Usj3", "Policy already exists")
	}
	repoOrg := model.OrgFromModel(existing)
	repoPolicy := model.OrgIamPolicyFromModel(policy)
	orgAggregate := OrgIamPolicyAddedAggregate(es.Eventstore.AggregateCreator(), repoOrg, repoPolicy)
	if err != nil {
		return nil, err
	}
	err = es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, orgAggregate)
	if err != nil {
		return nil, err
	}

	return model.OrgIamPolicyToModel(repoOrg.OrgIamPolicy), nil
}

func (es *OrgEventstore) ChangeOrgIamPolicy(ctx context.Context, policy *org_model.OrgIamPolicy) (*org_model.OrgIamPolicy, error) {
	existing, err := es.OrgByID(ctx, org_model.NewOrg(policy.AggregateID))
	if err != nil {
		return nil, err
	}
	if existing.OrgIamPolicy == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-8juSd", "Policy doesnt exist")
	}
	repoOrg := model.OrgFromModel(existing)
	repoPolicy := model.OrgIamPolicyFromModel(policy)
	orgAggregate := OrgIamPolicyChangedAggregate(es.Eventstore.AggregateCreator(), repoOrg, repoPolicy)
	if err != nil {
		return nil, err
	}
	err = es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, orgAggregate)
	if err != nil {
		return nil, err
	}

	return model.OrgIamPolicyToModel(repoOrg.OrgIamPolicy), nil
}

func (es *OrgEventstore) RemoveOrgIamPolicy(ctx context.Context, orgID string) error {
	existing, err := es.OrgByID(ctx, org_model.NewOrg(orgID))
	if err != nil {
		return err
	}
	if existing.OrgIamPolicy == nil {
		return errors.ThrowPreconditionFailed(nil, "EVENT-z6Dse", "Policy doesnt exist")
	}
	repoOrg := model.OrgFromModel(existing)
	orgAggregate := OrgIamPolicyRemovedAggregate(es.Eventstore.AggregateCreator(), repoOrg)
	if err != nil {
		return err
	}
	return es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, orgAggregate)
}
