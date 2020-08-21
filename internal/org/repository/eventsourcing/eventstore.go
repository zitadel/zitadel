package eventsourcing

import (
	"context"
	"encoding/json"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"

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
	IamID                 string
	idGenerator           id.Generator
	verificationAlgorithm crypto.EncryptionAlgorithm
	verificationGenerator crypto.Generator
	defaultOrgIamPolicy   *org_model.OrgIamPolicy
	verificationValidator func(domain string, token string, verifier string, checkType http_utils.CheckType) error
	secretCrypto          crypto.Crypto
}

type OrgConfig struct {
	eventstore.Eventstore
	IAMDomain          string
	VerificationConfig *crypto.KeyConfig
}

func StartOrg(conf OrgConfig, defaults systemdefaults.SystemDefaults) *OrgEventstore {
	policy := defaults.DefaultPolicies.OrgIam
	policy.Default = true
	policy.IamDomain = conf.IAMDomain
	verificationAlg, err := crypto.NewAESCrypto(defaults.DomainVerification.VerificationKey)
	logging.Log("EVENT-aZ22d").OnError(err).Panic("cannot create verificationAlgorithm for domain verification")
	verificationGen := crypto.NewEncryptionGenerator(defaults.DomainVerification.VerificationGenerator, verificationAlg)

	aesCrypto, err := crypto.NewAESCrypto(defaults.IDPConfigVerificationKey)
	logging.Log("EVENT-Sn8du").OnError(err).Panic("cannot create verificationAlgorithm for idp config verification")

	return &OrgEventstore{
		Eventstore:            conf.Eventstore,
		idGenerator:           id.SonyFlakeGenerator,
		verificationAlgorithm: verificationAlg,
		verificationGenerator: verificationGen,
		verificationValidator: http_utils.ValidateDomain,
		IAMDomain:             conf.IAMDomain,
		IamID:                 defaults.IamID,
		defaultOrgIamPolicy:   &policy,
		secretCrypto:          aesCrypto,
	}
}

func (es *OrgEventstore) PrepareCreateOrg(ctx context.Context, orgModel *org_model.Org, users func(context.Context, string) ([]*es_models.Aggregate, error)) (*model.Org, []*es_models.Aggregate, error) {
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

	aggregates, err := orgCreatedAggregates(ctx, es.AggregateCreator(), org, users)

	return org, aggregates, err
}

func (es *OrgEventstore) CreateOrg(ctx context.Context, orgModel *org_model.Org, users func(context.Context, string) ([]*es_models.Aggregate, error)) (*org_model.Org, error) {
	org, aggregates, err := es.PrepareCreateOrg(ctx, orgModel, users)
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
	if domain == nil || !domain.IsValid() {
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
	if domain == nil || !domain.IsValid() {
		return "", "", errors.ThrowPreconditionFailed(nil, "EVENT-R24hb", "Errors.Org.InvalidDomain")
	}
	checkType, ok := domain.ValidationType.CheckType()
	if !ok {
		return "", "", errors.ThrowPreconditionFailed(nil, "EVENT-Gsw31", "Errors.Org.DomainVerificationTypeInvalid")
	}
	existing, err := es.OrgByID(ctx, org_model.NewOrg(domain.AggregateID))
	if err != nil {
		return "", "", err
	}
	_, d := existing.GetDomain(domain)
	if d == nil {
		return "", "", errors.ThrowPreconditionFailed(nil, "EVENT-AGD31", "Errors.Org.DomainNotOnOrg")
	}
	if d.Verified {
		return "", "", errors.ThrowPreconditionFailed(nil, "EVENT-HGw21", "Errors.Org.DomainAlreadyVerified")
	}
	token, err := domain.GenerateVerificationCode(es.verificationGenerator)
	if err != nil {
		return "", "", err
	}
	url, err := http_utils.TokenUrl(domain.Domain, token, checkType)
	if err != nil {
		return "", "", errors.ThrowPreconditionFailed(err, "EVENT-Bae21", "Errors.Org.DomainVerificationTypeInvalid")
	}

	repoOrg := model.OrgFromModel(existing)
	repoDomain := model.OrgDomainFromModel(domain)
	aggregate := OrgDomainValidationGeneratedAggregate(es.Eventstore.AggregateCreator(), repoOrg, repoDomain)

	err = es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, aggregate)
	if err != nil {
		return "", "", err
	}
	return token, url, err
}

func (es *OrgEventstore) ValidateOrgDomain(ctx context.Context, domain *org_model.OrgDomain, users func(context.Context, string) ([]*es_models.Aggregate, error)) error {
	if domain == nil || !domain.IsValid() {
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
	if d.Verified {
		return errors.ThrowPreconditionFailed(nil, "EVENT-4gT342", "Errors.Org.DomainAlreadyVerified")
	}
	if d.ValidationCode == nil || d.ValidationType == org_model.OrgDomainValidationTypeUnspecified {
		return errors.ThrowPreconditionFailed(nil, "EVENT-SFBB3", "Errors.Org.DomainVerificationMissing")
	}
	validationCode, err := crypto.DecryptString(d.ValidationCode, es.verificationAlgorithm)
	if err != nil {
		return err
	}
	repoOrg := model.OrgFromModel(existing)
	repoDomain := model.OrgDomainFromModel(domain)
	checkType, _ := d.ValidationType.CheckType()
	err = es.verificationValidator(d.Domain, validationCode, validationCode, checkType)
	if err == nil {
		orgAggregates, err := OrgDomainVerifiedAggregate(ctx, es.Eventstore.AggregateCreator(), repoOrg, repoDomain, users)
		if err != nil {
			return err
		}
		return es_sdk.PushAggregates(ctx, es.PushAggregates, repoOrg.AppendEvents, orgAggregates...)
	}
	if err := es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, OrgDomainValidationFailedAggregate(es.Eventstore.AggregateCreator(), repoOrg, repoDomain)); err != nil {
		return err
	}
	return errors.ThrowInvalidArgument(err, "EVENT-GH3s", "Errors.Org.DomainVerificationFailed")
}

func (es *OrgEventstore) SetPrimaryOrgDomain(ctx context.Context, domain *org_model.OrgDomain) error {
	if domain == nil || !domain.IsValid() {
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

func (es *OrgEventstore) AddIdpConfiguration(ctx context.Context, idp *iam_model.IdpConfig) (*iam_model.IdpConfig, error) {
	if idp == nil || !idp.IsValid(true) {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Ms89d", "Errors.Org.IdpInvalid")
	}
	existing, err := es.OrgByID(ctx, org_model.NewOrg(idp.AggregateID))
	if err != nil {
		return nil, err
	}
	id, err := es.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	idp.IDPConfigID = id

	if idp.OIDCConfig != nil {
		idp.OIDCConfig.IDPConfigID = id
		err = idp.OIDCConfig.CryptSecret(es.secretCrypto)
	}
	repoOrg := model.OrgFromModel(existing)
	repoIdp := iam_es_model.IdpConfigFromModel(idp)

	addAggregate := IdpConfigurationAddedAggregate(es.Eventstore.AggregateCreator(), repoOrg, repoIdp)
	err = es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	if _, i := iam_es_model.GetIdpConfig(repoOrg.IDPs, idp.IDPConfigID); i != nil {
		return iam_es_model.IdpConfigToModel(i), nil
	}
	return nil, errors.ThrowInternal(nil, "EVENT-Cmsj8d", "Errors.Internal")
}

func (es *OrgEventstore) ChangeIdpConfiguration(ctx context.Context, idp *iam_model.IdpConfig) (*iam_model.IdpConfig, error) {
	if idp == nil || !idp.IsValid(false) {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Mslo9", "Errors.Org.IdpInvalid")
	}
	existing, err := es.OrgByID(ctx, org_model.NewOrg(idp.AggregateID))
	if err != nil {
		return nil, err
	}
	if _, i := existing.GetIDP(idp.IDPConfigID); i == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Aji8e", "Errors.Org.IdpNotExisting")
	}
	repoOrg := model.OrgFromModel(existing)
	repoIdp := iam_es_model.IdpConfigFromModel(idp)

	iamAggregate := IdpConfigurationChangedAggregate(es.Eventstore.AggregateCreator(), repoOrg, repoIdp)
	err = es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, iamAggregate)
	if err != nil {
		return nil, err
	}
	if _, i := iam_es_model.GetIdpConfig(repoOrg.IDPs, idp.IDPConfigID); i != nil {
		return iam_es_model.IdpConfigToModel(i), nil
	}
	return nil, errors.ThrowInternal(nil, "EVENT-Ml9xs", "Errors.Internal")
}

func (es *OrgEventstore) PrepareRemoveIdpConfiguration(ctx context.Context, idp *iam_model.IdpConfig) (*model.Org, *es_models.Aggregate, error) {
	if idp.IDPConfigID == "" {
		return nil, nil, errors.ThrowPreconditionFailed(nil, "EVENT-Wz7sD", "Errors.Org.IDMissing")
	}
	existing, err := es.OrgByID(ctx, org_model.NewOrg(idp.AggregateID))
	if err != nil {
		return nil, nil, err
	}
	if _, i := existing.GetIDP(idp.IDPConfigID); i == nil {
		return nil, nil, errors.ThrowPreconditionFailed(nil, "EVENT-Smiu8", "Errors.Org.IdpNotExisting")
	}
	repoOrg := model.OrgFromModel(existing)
	repoIdp := iam_es_model.IdpConfigFromModel(idp)
	provider := new(iam_es_model.IdpProvider)
	if repoOrg.LoginPolicy != nil {
		_, provider = iam_es_model.GetIdpProvider(repoOrg.LoginPolicy.IdpProviders, idp.IDPConfigID)
	}
	agg, err := IdpConfigurationRemovedAggregate(ctx, es.Eventstore.AggregateCreator(), repoOrg, repoIdp, provider)
	if err != nil {
		return nil, nil, err
	}
	return repoOrg, agg, nil
}

func (es *OrgEventstore) RemoveIdpConfiguration(ctx context.Context, idp *iam_model.IdpConfig) error {
	repoOrg, agg, err := es.PrepareRemoveIdpConfiguration(ctx, idp)
	if err != nil {
		return err
	}
	return es_sdk.PushAggregates(ctx, es.PushAggregates, repoOrg.AppendEvents, agg)
}

func (es *OrgEventstore) DeactivateIdpConfiguration(ctx context.Context, orgID, idpID string) (*iam_model.IdpConfig, error) {
	if idpID == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Smk8d", "Errors.Org.IDMissing")
	}
	existing, err := es.OrgByID(ctx, org_model.NewOrg(orgID))
	if err != nil {
		return nil, err
	}
	idp := &iam_model.IdpConfig{IDPConfigID: idpID}
	if _, app := existing.GetIDP(idp.IDPConfigID); app == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Amk8d", "Errors.Org.IdpNotExisting")
	}
	repoOrg := model.OrgFromModel(existing)
	repoIdp := iam_es_model.IdpConfigFromModel(idp)

	iamAggregate := IdpConfigurationDeactivatedAggregate(es.Eventstore.AggregateCreator(), repoOrg, repoIdp)
	err = es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, iamAggregate)
	if err != nil {
		return nil, err
	}
	if _, i := iam_es_model.GetIdpConfig(repoOrg.IDPs, idp.IDPConfigID); i != nil {
		return iam_es_model.IdpConfigToModel(i), nil
	}
	return nil, errors.ThrowInternal(nil, "EVENT-Amk9c", "Errors.Internal")
}

func (es *OrgEventstore) ReactivateIdpConfiguration(ctx context.Context, orgID, idpID string) (*iam_model.IdpConfig, error) {
	if idpID == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Xm8df", "Errors.Org.IDMissing")
	}
	existing, err := es.OrgByID(ctx, org_model.NewOrg(orgID))
	if err != nil {
		return nil, err
	}
	idp := &iam_model.IdpConfig{IDPConfigID: idpID}
	if _, i := existing.GetIDP(idp.IDPConfigID); i == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Qls0f", "Errors.Org.IdpNotExisting")
	}
	repoOrg := model.OrgFromModel(existing)
	repoIdp := iam_es_model.IdpConfigFromModel(idp)

	iamAggregate := IdpConfigurationReactivatedAggregate(es.Eventstore.AggregateCreator(), repoOrg, repoIdp)
	err = es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, iamAggregate)
	if err != nil {
		return nil, err
	}
	if _, i := iam_es_model.GetIdpConfig(repoOrg.IDPs, idp.IDPConfigID); i != nil {
		return iam_es_model.IdpConfigToModel(i), nil
	}
	return nil, errors.ThrowInternal(nil, "EVENT-Al90s", "Errors.Internal")
}

func (es *OrgEventstore) ChangeIdpOidcConfiguration(ctx context.Context, config *iam_model.OidcIdpConfig) (*iam_model.OidcIdpConfig, error) {
	if config == nil || !config.IsValid(false) {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Qs789", "Errors.Org.OIDCConfigInvalid")
	}
	existing, err := es.OrgByID(ctx, org_model.NewOrg(config.AggregateID))
	if err != nil {
		return nil, err
	}
	var idp *iam_model.IdpConfig
	if _, idp = existing.GetIDP(config.IDPConfigID); idp == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-pso0s", "Errors.Org.IdpNoExisting")
	}
	if idp.Type != iam_model.IDPConfigTypeOIDC {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Fms8w", "Errors.Iam.IdpIsNotOIDC")
	}
	if config.ClientSecretString != "" {
		err = idp.OIDCConfig.CryptSecret(es.secretCrypto)
	} else {
		config.ClientSecret = nil
	}
	repoOrg := model.OrgFromModel(existing)
	repoConfig := iam_es_model.OidcIdpConfigFromModel(config)

	iamAggregate := OIDCIdpConfigurationChangedAggregate(es.Eventstore.AggregateCreator(), repoOrg, repoConfig)
	err = es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, iamAggregate)
	if err != nil {
		return nil, err
	}
	if _, a := iam_es_model.GetIdpConfig(repoOrg.IDPs, idp.IDPConfigID); a != nil {
		return iam_es_model.OidcIdpConfigToModel(a.OIDCIDPConfig), nil
	}
	return nil, errors.ThrowInternal(nil, "EVENT-Sldk8", "Errors.Internal")
}

func (es *OrgEventstore) AddLoginPolicy(ctx context.Context, policy *iam_model.LoginPolicy) (*iam_model.LoginPolicy, error) {
	if policy == nil || !policy.IsValid() {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Sjkl9", "Errors.Org.LoginPolicyInvalid")
	}
	existing, err := es.OrgByID(ctx, org_model.NewOrg(policy.AggregateID))
	if err != nil {
		return nil, err
	}

	repoOrg := model.OrgFromModel(existing)
	repoLoginPolicy := iam_es_model.LoginPolicyFromModel(policy)

	addAggregate := LoginPolicyAddedAggregate(es.Eventstore.AggregateCreator(), repoOrg, repoLoginPolicy)
	err = es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	return iam_es_model.LoginPolicyToModel(repoOrg.LoginPolicy), nil
}

func (es *OrgEventstore) ChangeLoginPolicy(ctx context.Context, policy *iam_model.LoginPolicy) (*iam_model.LoginPolicy, error) {
	if policy == nil || !policy.IsValid() {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Lso02", "Errors.Org.LoginPolicyInvalid")
	}
	existing, err := es.OrgByID(ctx, org_model.NewOrg(policy.AggregateID))
	if err != nil {
		return nil, err
	}

	repoOrg := model.OrgFromModel(existing)
	repoLoginPolicy := iam_es_model.LoginPolicyFromModel(policy)

	addAggregate := LoginPolicyChangedAggregate(es.Eventstore.AggregateCreator(), repoOrg, repoLoginPolicy)
	err = es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	return iam_es_model.LoginPolicyToModel(repoOrg.LoginPolicy), nil
}

func (es *OrgEventstore) RemoveLoginPolicy(ctx context.Context, policy *iam_model.LoginPolicy) error {
	if policy == nil || !policy.IsValid() {
		return errors.ThrowPreconditionFailed(nil, "EVENT-O0s9e", "Errors.Org.LoginPolicyInvalid")
	}
	existing, err := es.OrgByID(ctx, org_model.NewOrg(policy.AggregateID))
	if err != nil {
		return err
	}
	repoOrg := model.OrgFromModel(existing)

	addAggregate := LoginPolicyRemovedAggregate(es.Eventstore.AggregateCreator(), repoOrg)
	return es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, addAggregate)
}

func (es *OrgEventstore) GetIdpConfiguration(ctx context.Context, aggregateID, idpConfigID string) (*iam_model.IdpConfig, error) {
	existing, err := es.OrgByID(ctx, org_model.NewOrg(aggregateID))
	if err != nil {
		return nil, err
	}
	if _, i := existing.GetIDP(idpConfigID); i != nil {
		return i, nil
	}
	return nil, errors.ThrowNotFound(nil, "EVENT-Qlo0d", "Errors.Org.IdpNotExisting")
}

func (es *OrgEventstore) AddIdpProviderToLoginPolicy(ctx context.Context, provider *iam_model.IdpProvider) (*iam_model.IdpProvider, error) {
	if provider == nil || !provider.IsValid() {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Sjd8e", "Errors.Org.IdpProviderInvalid")
	}
	existing, err := es.OrgByID(ctx, org_model.NewOrg(provider.AggregateID))
	if err != nil {
		return nil, err
	}
	if existing.LoginPolicy == nil {
		return nil, errors.ThrowAlreadyExists(nil, "EVENT-sk9fW", "Errors.Org.LoginPolicy.NotExisting")
	}
	if _, m := existing.LoginPolicy.GetIdpProvider(provider.IdpConfigID); m != nil {
		return nil, errors.ThrowAlreadyExists(nil, "EVENT-Lso9f", "Errors.Org.LoginPolicy.IdpProviderAlreadyExisting")
	}
	repoOrg := model.OrgFromModel(existing)
	repoProvider := iam_es_model.IdpProviderFromModel(provider)

	addAggregate := LoginPolicyIdpProviderAddedAggregate(es.Eventstore.AggregateCreator(), repoOrg, repoProvider, es.IamID)
	err = es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	if _, m := iam_es_model.GetIdpProvider(repoOrg.LoginPolicy.IdpProviders, provider.IdpConfigID); m != nil {
		return iam_es_model.IdpProviderToModel(m), nil
	}
	return nil, errors.ThrowInternal(nil, "EVENT-Slf9s", "Errors.Internal")
}

func (es *OrgEventstore) PrepareRemoveIdpProviderFromLoginPolicy(ctx context.Context, provider *iam_model.IdpProvider, cascade bool) (*model.Org, *es_models.Aggregate, error) {
	if provider == nil || !provider.IsValid() {
		return nil, nil, errors.ThrowPreconditionFailed(nil, "EVENT-Esi8c", "Errors.IdpProviderInvalid")
	}
	existing, err := es.OrgByID(ctx, org_model.NewOrg(provider.AggregateID))
	if err != nil {
		return nil, nil, err
	}
	if _, m := existing.LoginPolicy.GetIdpProvider(provider.IdpConfigID); m == nil {
		return nil, nil, errors.ThrowPreconditionFailed(nil, "EVENT-29skr", "Errors.Iam.LoginPolicy.IdpProviderNotExisting")
	}
	repoOrg := model.OrgFromModel(existing)
	providerID := &iam_es_model.IdpProviderID{provider.IdpConfigID}
	providerAggregates, err := LoginPolicyIdpProviderRemovedAggregate(ctx, es.Eventstore.AggregateCreator(), repoOrg, providerID, cascade)
	if err != nil {
		return nil, nil, err
	}
	return repoOrg, providerAggregates, nil
}

func (es *OrgEventstore) RemoveIdpProviderFromLoginPolicy(ctx context.Context, provider *iam_model.IdpProvider) error {
	repoOrg, agg, err := es.PrepareRemoveIdpProviderFromLoginPolicy(ctx, provider, false)
	if err != nil {
		return err
	}
	return es_sdk.PushAggregates(ctx, es.PushAggregates, repoOrg.AppendEvents, agg)
}
