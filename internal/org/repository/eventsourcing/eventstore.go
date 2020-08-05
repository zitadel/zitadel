package eventsourcing

import (
	"context"
	"encoding/json"

	"github.com/caos/logging"
	"github.com/golang/protobuf/ptypes"

	http_utils "github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
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
	IAMDomain             string
	idGenerator           id.Generator
	verificationAlgorithm crypto.EncryptionAlgorithm
	verificationGenerator crypto.Generator
	defaultOrgIamPolicy   *org_model.OrgIamPolicy
}

type OrgConfig struct {
	eventstore.Eventstore
	IAMDomain          string
	VerificationConfig *crypto.KeyConfig
}

func StartOrg(conf OrgConfig, defaults systemdefaults.SystemDefaults) *OrgEventstore {
	policy := defaults.DefaultPolicies.OrgIam
	policy.Default = true
	verificationAlg, err := crypto.NewAESCrypto(defaults.DomainVerification.VerificationKey)
	logging.Log("EVENT-aZ22d").OnError(err).Panic("cannot create verificationAlgorithm for domain verification")
	verificationGen := crypto.NewEncryptionGenerator(defaults.DomainVerification.VerificationGenerator, verificationAlg)
	return &OrgEventstore{
		Eventstore:            conf.Eventstore,
		idGenerator:           id.SonyFlakeGenerator,
		verificationAlgorithm: verificationAlg,
		verificationGenerator: verificationGen,
		IAMDomain:             conf.IAMDomain,
		defaultOrgIamPolicy:   &policy,
	}
}

func (es *OrgEventstore) PrepareCreateOrg(ctx context.Context, orgModel *org_model.Org) (*model.Org, []*es_models.Aggregate, error) {
	if orgModel == nil || !orgModel.IsValid() {
		return nil, nil, errors.ThrowInvalidArgument(nil, "EVENT-OeLSk", "Errors.Org.Invalid")
	}
	orgModel.AddIAMDomain(es.IAMDomain)

	id, err := es.idGenerator.Next()
	if err != nil {
		return nil, nil, errors.ThrowInternal(err, "EVENT-OwciI", "Errors.Internal")
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
		return nil, errors.ThrowInvalidArgument(nil, "EVENT-gQTYP", "Errors.Org.Empty")
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
		return nil, errors.ThrowNotFound(nil, "EVENT-kVLb2", "Errors.Org.NotFound")
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
		return nil, errors.ThrowInvalidArgument(nil, "EVENT-oL9nT", "Errors.Org.NotFound")
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
		return nil, errors.ThrowInvalidArgument(nil, "EVENT-oL9nT", "Errors.Org.Empty")
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
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-8sFJW", "Errors.Org.InvalidDomain")
	}
	existing, err := es.OrgByID(ctx, org_model.NewOrg(domain.AggregateID))
	if err != nil {
		return nil, err
	}
	repoOrg := model.OrgFromModel(existing)
	repoDomain := model.OrgDomainFromModel(domain)
	aggregate := OrgDomainAddedAggregate(es.Eventstore.AggregateCreator(), repoOrg, repoDomain)

	err = es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, aggregate)
	if err != nil {
		return nil, err
	}

	if _, d := model.GetDomain(repoOrg.Domains, domain.Domain); d != nil {
		return model.OrgDomainToModel(d), nil
	}
	return nil, errors.ThrowInternal(nil, "EVENT-ISOP0", "Errors.Internal")
}

func (es *OrgEventstore) GenerateOrgDomainValidation(ctx context.Context, domain *org_model.OrgDomain) (string, string, error) {
	if !domain.IsValid() {
		return "", "", errors.ThrowPreconditionFailed(nil, "EVENT-R24hb", "Errors.Org.InvalidDomain")
	}
	existing, err := es.OrgByID(ctx, org_model.NewOrg(domain.AggregateID))
	if err != nil {
		return "", "", err
	}
	_, d := existing.GetDomain(domain)
	if d == nil {
		return "", "", errors.ThrowPreconditionFailed(nil, "EVENT-AGD31", "Errors.Org.DomainNotOnOrg")
	}
	token, err := domain.GenerateVerificationCode(es.verificationGenerator)
	if err != nil {
		return "", "", err
	}

	repoOrg := model.OrgFromModel(existing)
	repoDomain := model.OrgDomainFromModel(domain)
	aggregate := OrgDomainValidationGeneratedAggregate(es.Eventstore.AggregateCreator(), repoOrg, repoDomain)

	err = es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, aggregate)
	if err != nil {
		return "", "", err
	}

	url, err := http_utils.TokenUrl(domain.Domain, token, domain.ValidationType.CheckType())
	return token, url, err
}

func (es *OrgEventstore) ValidateOrgDomain(ctx context.Context, domain *org_model.OrgDomain) error {
	if !domain.IsValid() {
		return errors.ThrowPreconditionFailed(nil, "EVENT-R24hb", "Errors.Org.InvalidDomain")
	}
	existing, err := es.OrgByID(ctx, org_model.NewOrg(domain.AggregateID))
	if err != nil {
		return err
	}
	_, d := existing.GetDomain(domain)
	if d == nil {
		return errors.ThrowPreconditionFailed(nil, "EVENT-Sjdi3", "Errors.Org.DomainNotOnOrg")
	}
	validationCode, err := crypto.DecryptString(d.ValidationCode, es.verificationAlgorithm)
	if err != nil {
		return err
	}
	repoOrg := model.OrgFromModel(existing)
	repoDomain := model.OrgDomainFromModel(domain)
	err = http_utils.ValidateDomain(d.Domain, validationCode, validationCode, d.ValidationType.CheckType())
	if err == nil {
		orgAggregates, err := OrgDomainVerifiedAggregate(ctx, es.Eventstore.AggregateCreator(), repoOrg, repoDomain)
		if err != nil {
			return err
		}
		return es_sdk.PushAggregates(ctx, es.PushAggregates, repoOrg.AppendEvents, orgAggregates...)
	}
	if err := es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, OrgDomainValidationFailedAggregate(es.Eventstore.AggregateCreator(), repoOrg, repoDomain)); err != nil {
		return err
	}
	return errors.ThrowInvalidArgument(err, "EVENT-GH3s", "Errors.Org.Code.Invalid")
}

func (es *OrgEventstore) SetPrimaryOrgDomain(ctx context.Context, domain *org_model.OrgDomain) error {
	if !domain.IsValid() {
		return errors.ThrowPreconditionFailed(nil, "EVENT-SsDG2", "Errors.Org.InvalidDomain")
	}
	existing, err := es.OrgByID(ctx, org_model.NewOrg(domain.AggregateID))
	if err != nil {
		return err
	}
	_, d := existing.GetDomain(domain)
	if d == nil {
		return errors.ThrowPreconditionFailed(nil, "EVENT-GDfA3", "Errors.Org.DomainNotOnOrg")
	}
	if !d.Verified {
		return errors.ThrowPreconditionFailed(nil, "EVENT-Ggd32", "Errors.Org.DomainNotVerified")
	}
	repoOrg := model.OrgFromModel(existing)
	repoDomain := model.OrgDomainFromModel(domain)
	if err := es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, OrgDomainSetPrimaryAggregate(es.Eventstore.AggregateCreator(), repoOrg, repoDomain)); err != nil {
		return err
	}
	return nil
}

func (es *OrgEventstore) RemoveOrgDomain(ctx context.Context, domain *org_model.OrgDomain) error {
	if domain.Domain == "" {
		return errors.ThrowPreconditionFailed(nil, "EVENT-SJsK3", "Errors.Org.DomainMissing")
	}
	existing, err := es.OrgByID(ctx, org_model.NewOrg(domain.AggregateID))
	if err != nil {
		return err
	}
	_, d := existing.GetDomain(domain)
	if d == nil {
		return errors.ThrowPreconditionFailed(nil, "EVENT-Sjdi3", "Errors.Org.DomainNotOnOrg")
	}
	if d.Primary {
		return errors.ThrowPreconditionFailed(nil, "EVENT-Sjdi3", "Errors.Org.PrimaryDomainNotDeletable")
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

func (es *OrgEventstore) OrgChanges(ctx context.Context, id string, lastSequence uint64, limit uint64, sortAscending bool) (*org_model.OrgChanges, error) {
	query := ChangesQuery(id, lastSequence, limit, sortAscending)

	events, err := es.Eventstore.FilterEvents(context.Background(), query)
	if err != nil {
		logging.Log("EVENT-ZRffs").WithError(err).Warn("eventstore unavailable")
		return nil, errors.ThrowInternal(err, "EVENT-328b1", "Errors.Org.NotFound")
	}
	if len(events) == 0 {
		return nil, errors.ThrowNotFound(nil, "EVENT-FpQqK", "Errors.Changes.NotFound")
	}

	changes := make([]*org_model.OrgChange, len(events))

	for i, event := range events {
		creationDate, err := ptypes.TimestampProto(event.CreationDate)
		logging.Log("EVENT-qxIR7").OnError(err).Debug("unable to parse timestamp")
		change := &org_model.OrgChange{
			ChangeDate: creationDate,
			EventType:  event.Type.String(),
			ModifierId: event.EditorUser,
			Sequence:   event.Sequence,
		}

		if event.Data != nil {
			org := new(model.Org)
			err := json.Unmarshal(event.Data, org)
			logging.Log("EVENT-XCLEm").OnError(err).Debug("unable to unmarshal data")
			change.Data = org
		}

		changes[i] = change
		if lastSequence < event.Sequence {
			lastSequence = event.Sequence
		}
	}

	return &org_model.OrgChanges{
		Changes:      changes,
		LastSequence: lastSequence,
	}, nil
}

func ChangesQuery(orgID string, latestSequence, limit uint64, sortAscending bool) *es_models.SearchQuery {
	query := es_models.NewSearchQuery().
		AggregateTypeFilter(model.OrgAggregate)

	if !sortAscending {
		query.OrderDesc()
	}

	query.LatestSequenceFilter(latestSequence).
		AggregateIDFilter(orgID).
		SetLimit(limit)
	return query
}

func (es *OrgEventstore) OrgMemberByIDs(ctx context.Context, member *org_model.OrgMember) (*org_model.OrgMember, error) {
	if member == nil || member.UserID == "" || member.AggregateID == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-ld93d", "Errors.Org.MemberIDMissing")
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

	return nil, errors.ThrowNotFound(nil, "EVENT-SXji6", "Errors.Org.MemberNotFound")
}

func (es *OrgEventstore) PrepareAddOrgMember(ctx context.Context, member *org_model.OrgMember, resourceOwner string) (*model.OrgMember, *es_models.Aggregate, error) {
	if member == nil || !member.IsValid() {
		return nil, nil, errors.ThrowPreconditionFailed(nil, "EVENT-9dk45", "Errors.Org.InvalidMember")
	}

	repoMember := model.OrgMemberFromModel(member)
	addAggregate, err := orgMemberAddedAggregate(ctx, es.Eventstore.AggregateCreator(), repoMember, resourceOwner)

	return repoMember, addAggregate, err
}

func (es *OrgEventstore) AddOrgMember(ctx context.Context, member *org_model.OrgMember) (*org_model.OrgMember, error) {
	repoMember, addAggregate, err := es.PrepareAddOrgMember(ctx, member, "")
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
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-9dk45", "Errors.Org.InvalidMember")
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
		return errors.ThrowInvalidArgument(nil, "EVENT-d43fs", "Errors.Org.UserIDMissing")
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
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}
	if existing != nil && existing.OrgIamPolicy != nil {
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
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-7Usj3", "Errors.Org.PolicyAlreadyExists")
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
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-8juSd", "Errors.Org.PolicyNotExisting")
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
		return errors.ThrowPreconditionFailed(nil, "EVENT-z6Dse", "Errors.Org.PolicyNotExisting")
	}
	repoOrg := model.OrgFromModel(existing)
	orgAggregate := OrgIamPolicyRemovedAggregate(es.Eventstore.AggregateCreator(), repoOrg)
	if err != nil {
		return err
	}
	return es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, orgAggregate)
}
