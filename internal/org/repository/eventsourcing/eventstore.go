package eventsourcing

import (
	"context"
	"encoding/json"

	"github.com/caos/logging"
	http_utils "github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/id"
	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	"github.com/golang/protobuf/ptypes"
)

type OrgEventstore struct {
	eventstore.Eventstore
	IAMDomain             string
	IamID                 string
	idGenerator           id.Generator
	verificationAlgorithm crypto.EncryptionAlgorithm
	verificationGenerator crypto.Generator
	verificationValidator func(domain string, token string, verifier string, checkType http_utils.CheckType) error
	secretCrypto          crypto.Crypto
}

type OrgConfig struct {
	eventstore.Eventstore
	IAMDomain          string
	VerificationConfig *crypto.KeyConfig
}

func StartOrg(conf OrgConfig, defaults systemdefaults.SystemDefaults) *OrgEventstore {
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
	existingOrg, err := es.OrgByID(ctx, org_model.NewOrg(domain.AggregateID))
	if err != nil {
		return nil, err
	}
	repoOrg := model.OrgFromModel(existingOrg)
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
	existingOrg, err := es.OrgByID(ctx, org_model.NewOrg(domain.AggregateID))
	if err != nil {
		return "", "", err
	}
	_, d := existingOrg.GetDomain(domain)
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

	repoOrg := model.OrgFromModel(existingOrg)
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
	existingOrg, err := es.OrgByID(ctx, org_model.NewOrg(domain.AggregateID))
	if err != nil {
		return err
	}
	_, existingDomain := existingOrg.GetDomain(domain)
	if existingDomain == nil {
		return errors.ThrowPreconditionFailed(nil, "EVENT-Sjdi3", "Errors.Org.DomainNotOnOrg")
	}
	if existingDomain.Verified {
		return errors.ThrowPreconditionFailed(nil, "EVENT-4gT342", "Errors.Org.DomainAlreadyVerified")
	}
	if existingDomain.ValidationCode == nil || existingDomain.ValidationType == org_model.OrgDomainValidationTypeUnspecified {
		return errors.ThrowPreconditionFailed(nil, "EVENT-SFBB3", "Errors.Org.DomainVerificationMissing")
	}
	validationCode, err := crypto.DecryptString(existingDomain.ValidationCode, es.verificationAlgorithm)
	if err != nil {
		return err
	}
	repoOrg := model.OrgFromModel(existingOrg)
	repoDomain := model.OrgDomainFromModel(domain)
	checkType, _ := existingDomain.ValidationType.CheckType()
	err = es.verificationValidator(existingDomain.Domain, validationCode, validationCode, checkType)
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
	existingOrg, err := es.OrgByID(ctx, org_model.NewOrg(domain.AggregateID))
	if err != nil {
		return err
	}
	_, existingDomain := existingOrg.GetDomain(domain)
	if existingDomain == nil {
		return errors.ThrowPreconditionFailed(nil, "EVENT-GDfA3", "Errors.Org.DomainNotOnOrg")
	}
	if !existingDomain.Verified {
		return errors.ThrowPreconditionFailed(nil, "EVENT-Ggd32", "Errors.Org.DomainNotVerified")
	}
	repoOrg := model.OrgFromModel(existingOrg)
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
	existingOrg, err := es.OrgByID(ctx, org_model.NewOrg(domain.AggregateID))
	if err != nil {
		return err
	}
	_, existingDomain := existingOrg.GetDomain(domain)
	if existingDomain == nil {
		return errors.ThrowPreconditionFailed(nil, "EVENT-Sjdi3", "Errors.Org.DomainNotOnOrg")
	}
	if existingDomain.Primary {
		return errors.ThrowPreconditionFailed(nil, "EVENT-Sjdi3", "Errors.Org.PrimaryDomainNotDeletable")
	}
	repoOrg := model.OrgFromModel(existingOrg)
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
		return nil, nil, errors.ThrowPreconditionFailed(nil, "EVENT-jRFLz", "Errors.Org.InvalidMember")
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
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-ara6l", "Errors.Org.InvalidMember")
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

func (es *OrgEventstore) GetOrgIAMPolicy(ctx context.Context, orgID string) (*iam_model.OrgIAMPolicy, error) {
	existingOrg, err := es.OrgByID(ctx, org_model.NewOrg(orgID))
	if err != nil {
		return nil, err
	}
	if existingOrg.OrgIamPolicy == nil {
		return nil, errors.ThrowNotFound(nil, "EVENT-3F9sf", "Errors.Org.OrgIAM.NotExisting")
	}
	return existingOrg.OrgIamPolicy, nil
}

func (es *OrgEventstore) AddOrgIAMPolicy(ctx context.Context, policy *iam_model.OrgIAMPolicy) (*iam_model.OrgIAMPolicy, error) {
	existingOrg, err := es.OrgByID(ctx, org_model.NewOrg(policy.AggregateID))
	if err != nil {
		return nil, err
	}
	if existingOrg.OrgIamPolicy != nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-7Usj3", "Errors.Org.PolicyAlreadyExists")
	}
	repoOrg := model.OrgFromModel(existingOrg)
	repoPolicy := iam_es_model.OrgIAMPolicyFromModel(policy)
	orgAggregate := OrgIAMPolicyAddedAggregate(es.Eventstore.AggregateCreator(), repoOrg, repoPolicy)
	if err != nil {
		return nil, err
	}
	err = es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, orgAggregate)
	if err != nil {
		return nil, err
	}

	return iam_es_model.OrgIAMPolicyToModel(repoOrg.OrgIAMPolicy), nil
}

func (es *OrgEventstore) ChangeOrgIAMPolicy(ctx context.Context, policy *iam_model.OrgIAMPolicy) (*iam_model.OrgIAMPolicy, error) {
	existingOrg, err := es.OrgByID(ctx, org_model.NewOrg(policy.AggregateID))
	if err != nil {
		return nil, err
	}
	if existingOrg.OrgIamPolicy == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-8juSd", "Errors.Org.PolicyNotExisting")
	}
	repoOrg := model.OrgFromModel(existingOrg)
	repoPolicy := iam_es_model.OrgIAMPolicyFromModel(policy)
	orgAggregate := OrgIAMPolicyChangedAggregate(es.Eventstore.AggregateCreator(), repoOrg, repoPolicy)
	if err != nil {
		return nil, err
	}
	err = es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, orgAggregate)
	if err != nil {
		return nil, err
	}

	return iam_es_model.OrgIAMPolicyToModel(repoOrg.OrgIAMPolicy), nil
}

func (es *OrgEventstore) RemoveOrgIAMPolicy(ctx context.Context, orgID string) error {
	existingOrg, err := es.OrgByID(ctx, org_model.NewOrg(orgID))
	if err != nil {
		return err
	}
	if existingOrg.OrgIamPolicy == nil {
		return errors.ThrowPreconditionFailed(nil, "EVENT-z6Dse", "Errors.Org.PolicyNotExisting")
	}
	repoOrg := model.OrgFromModel(existingOrg)
	orgAggregate := OrgIamPolicyRemovedAggregate(es.Eventstore.AggregateCreator(), repoOrg)
	if err != nil {
		return err
	}
	return es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, orgAggregate)
}

func (es *OrgEventstore) AddIDPConfig(ctx context.Context, idp *iam_model.IDPConfig) (*iam_model.IDPConfig, error) {
	if idp == nil || !idp.IsValid(true) {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Ms89d", "Errors.Org.IdpInvalid")
	}
	org, err := es.OrgByID(ctx, org_model.NewOrg(idp.AggregateID))
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
	repoOrg := model.OrgFromModel(org)
	repoIdp := iam_es_model.IDPConfigFromModel(idp)

	addAggregate := IDPConfigAddedAggregate(es.Eventstore.AggregateCreator(), repoOrg, repoIdp)
	err = es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	if _, idpConfig := iam_es_model.GetIDPConfig(repoOrg.IDPs, idp.IDPConfigID); idpConfig != nil {
		return iam_es_model.IDPConfigToModel(idpConfig), nil
	}
	return nil, errors.ThrowInternal(nil, "EVENT-Cmsj8d", "Errors.Internal")
}

func (es *OrgEventstore) ChangeIDPConfig(ctx context.Context, idp *iam_model.IDPConfig) (*iam_model.IDPConfig, error) {
	if idp == nil || !idp.IsValid(false) {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Mslo9", "Errors.Org.IdpInvalid")
	}
	org, err := es.OrgByID(ctx, org_model.NewOrg(idp.AggregateID))
	if err != nil {
		return nil, err
	}
	if _, i := org.GetIDP(idp.IDPConfigID); i == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Aji8e", "Errors.Org.IdpNotExisting")
	}
	repoOrg := model.OrgFromModel(org)
	repoIdp := iam_es_model.IDPConfigFromModel(idp)

	iamAggregate := IDPConfigChangedAggregate(es.Eventstore.AggregateCreator(), repoOrg, repoIdp)
	err = es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, iamAggregate)
	if err != nil {
		return nil, err
	}
	if _, idpConfig := iam_es_model.GetIDPConfig(repoOrg.IDPs, idp.IDPConfigID); idpConfig != nil {
		return iam_es_model.IDPConfigToModel(idpConfig), nil
	}
	return nil, errors.ThrowInternal(nil, "EVENT-Ml9xs", "Errors.Internal")
}

func (es *OrgEventstore) PrepareRemoveIDPConfig(ctx context.Context, idp *iam_model.IDPConfig) (*model.Org, *es_models.Aggregate, error) {
	if idp.IDPConfigID == "" {
		return nil, nil, errors.ThrowPreconditionFailed(nil, "EVENT-Wz7sD", "Errors.Org.IDMissing")
	}
	org, err := es.OrgByID(ctx, org_model.NewOrg(idp.AggregateID))
	if err != nil {
		return nil, nil, err
	}
	if _, i := org.GetIDP(idp.IDPConfigID); i == nil {
		return nil, nil, errors.ThrowPreconditionFailed(nil, "EVENT-Smiu8", "Errors.Org.IdpNotExisting")
	}
	repoOrg := model.OrgFromModel(org)
	repoIdp := iam_es_model.IDPConfigFromModel(idp)
	provider := new(iam_es_model.IDPProvider)
	if repoOrg.LoginPolicy != nil {
		_, provider = iam_es_model.GetIDPProvider(repoOrg.LoginPolicy.IDPProviders, idp.IDPConfigID)
	}
	agg, err := IDPConfigRemovedAggregate(ctx, es.Eventstore.AggregateCreator(), repoOrg, repoIdp, provider)
	if err != nil {
		return nil, nil, err
	}
	return repoOrg, agg, nil
}

func (es *OrgEventstore) RemoveIDPConfig(ctx context.Context, idp *iam_model.IDPConfig) error {
	repoOrg, agg, err := es.PrepareRemoveIDPConfig(ctx, idp)
	if err != nil {
		return err
	}
	return es_sdk.PushAggregates(ctx, es.PushAggregates, repoOrg.AppendEvents, agg)
}

func (es *OrgEventstore) DeactivateIDPConfig(ctx context.Context, orgID, idpID string) (*iam_model.IDPConfig, error) {
	if idpID == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Smk8d", "Errors.Org.IDMissing")
	}
	org, err := es.OrgByID(ctx, org_model.NewOrg(orgID))
	if err != nil {
		return nil, err
	}
	idp := &iam_model.IDPConfig{IDPConfigID: idpID}
	if _, app := org.GetIDP(idp.IDPConfigID); app == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Amk8d", "Errors.Org.IdpNotExisting")
	}
	repoOrg := model.OrgFromModel(org)
	repoIdp := iam_es_model.IDPConfigFromModel(idp)

	iamAggregate := IDPConfigDeactivatedAggregate(es.Eventstore.AggregateCreator(), repoOrg, repoIdp)
	err = es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, iamAggregate)
	if err != nil {
		return nil, err
	}
	if _, idpConfig := iam_es_model.GetIDPConfig(repoOrg.IDPs, idp.IDPConfigID); idpConfig != nil {
		return iam_es_model.IDPConfigToModel(idpConfig), nil
	}
	return nil, errors.ThrowInternal(nil, "EVENT-Amk9c", "Errors.Internal")
}

func (es *OrgEventstore) ReactivateIDPConfig(ctx context.Context, orgID, idpID string) (*iam_model.IDPConfig, error) {
	if idpID == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Xm8df", "Errors.Org.IDMissing")
	}
	org, err := es.OrgByID(ctx, org_model.NewOrg(orgID))
	if err != nil {
		return nil, err
	}
	idp := &iam_model.IDPConfig{IDPConfigID: idpID}
	if _, i := org.GetIDP(idp.IDPConfigID); i == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Qls0f", "Errors.Org.IdpNotExisting")
	}
	repoOrg := model.OrgFromModel(org)
	repoIdp := iam_es_model.IDPConfigFromModel(idp)

	iamAggregate := IDPConfigReactivatedAggregate(es.Eventstore.AggregateCreator(), repoOrg, repoIdp)
	err = es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, iamAggregate)
	if err != nil {
		return nil, err
	}
	if _, idpConfig := iam_es_model.GetIDPConfig(repoOrg.IDPs, idp.IDPConfigID); idpConfig != nil {
		return iam_es_model.IDPConfigToModel(idpConfig), nil
	}
	return nil, errors.ThrowInternal(nil, "EVENT-Al90s", "Errors.Internal")
}

func (es *OrgEventstore) ChangeIDPOIDCConfig(ctx context.Context, config *iam_model.OIDCIDPConfig) (*iam_model.OIDCIDPConfig, error) {
	if config == nil || !config.IsValid(false) {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Qs789", "Errors.Org.OIDCConfigInvalid")
	}
	org, err := es.OrgByID(ctx, org_model.NewOrg(config.AggregateID))
	if err != nil {
		return nil, err
	}
	var idp *iam_model.IDPConfig
	if _, idp = org.GetIDP(config.IDPConfigID); idp == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-pso0s", "Errors.Org.IdpNoExisting")
	}
	if idp.Type != iam_model.IDPConfigTypeOIDC {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Fms8w", "Errors.IAM.IdpIsNotOIDC")
	}
	if config.ClientSecretString != "" {
		err = idp.OIDCConfig.CryptSecret(es.secretCrypto)
	} else {
		config.ClientSecret = nil
	}
	repoOrg := model.OrgFromModel(org)
	repoConfig := iam_es_model.OIDCIDPConfigFromModel(config)

	iamAggregate := OIDCIDPConfigChangedAggregate(es.Eventstore.AggregateCreator(), repoOrg, repoConfig)
	err = es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, iamAggregate)
	if err != nil {
		return nil, err
	}
	if _, idpConfig := iam_es_model.GetIDPConfig(repoOrg.IDPs, idp.IDPConfigID); idpConfig != nil {
		return iam_es_model.OIDCIDPConfigToModel(idpConfig.OIDCIDPConfig), nil
	}
	return nil, errors.ThrowInternal(nil, "EVENT-Sldk8", "Errors.Internal")
}

func (es *OrgEventstore) AddLabelPolicy(ctx context.Context, policy *iam_model.LabelPolicy) (*iam_model.LabelPolicy, error) {
	if policy == nil || !policy.IsValid() {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-37rSC", "Errors.Org.LabelPolicyInvalid")
	}
	org, err := es.OrgByID(ctx, org_model.NewOrg(policy.AggregateID))
	if err != nil {
		return nil, err
	}

	repoOrg := model.OrgFromModel(org)
	repoLabelPolicy := iam_es_model.LabelPolicyFromModel(policy)

	addAggregate := LabelPolicyAddedAggregate(es.Eventstore.AggregateCreator(), repoOrg, repoLabelPolicy)
	err = es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	return iam_es_model.LabelPolicyToModel(repoOrg.LabelPolicy), nil
}

func (es *OrgEventstore) ChangeLabelPolicy(ctx context.Context, policy *iam_model.LabelPolicy) (*iam_model.LabelPolicy, error) {
	if policy == nil || !policy.IsValid() {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-0NBIw", "Errors.Org.LabelPolicyInvalid")
	}
	org, err := es.OrgByID(ctx, org_model.NewOrg(policy.AggregateID))
	if err != nil {
		return nil, err
	}

	repoOrg := model.OrgFromModel(org)
	repoLabelPolicy := iam_es_model.LabelPolicyFromModel(policy)

	addAggregate := LabelPolicyChangedAggregate(es.Eventstore.AggregateCreator(), repoOrg, repoLabelPolicy)
	err = es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	return iam_es_model.LabelPolicyToModel(repoOrg.LabelPolicy), nil
}

func (es *OrgEventstore) AddLoginPolicy(ctx context.Context, policy *iam_model.LoginPolicy) (*iam_model.LoginPolicy, error) {
	if policy == nil || !policy.IsValid() {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Sjkl9", "Errors.Org.LoginPolicy.Invalid")
	}
	org, err := es.OrgByID(ctx, org_model.NewOrg(policy.AggregateID))
	if err != nil {
		return nil, err
	}

	repoOrg := model.OrgFromModel(org)
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
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Lso02", "Errors.Org.LoginPolicy.Invalid")
	}
	org, err := es.OrgByID(ctx, org_model.NewOrg(policy.AggregateID))
	if err != nil {
		return nil, err
	}

	if org.LoginPolicy == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Lso02", "Errors.Org.LoginPolicy.NotExisting")
	}

	repoOrg := model.OrgFromModel(org)
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
		return errors.ThrowPreconditionFailed(nil, "EVENT-O0s9e", "Errors.Org.LoginPolicy.Invalid")
	}
	org, err := es.OrgByID(ctx, org_model.NewOrg(policy.AggregateID))
	if err != nil {
		return err
	}
	repoOrg := model.OrgFromModel(org)

	addAggregate := LoginPolicyRemovedAggregate(es.Eventstore.AggregateCreator(), repoOrg)
	return es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, addAggregate)
}

func (es *OrgEventstore) GetIDPConfig(ctx context.Context, aggregateID, idpConfigID string) (*iam_model.IDPConfig, error) {
	existing, err := es.OrgByID(ctx, org_model.NewOrg(aggregateID))
	if err != nil {
		return nil, err
	}
	if _, i := existing.GetIDP(idpConfigID); i != nil {
		return i, nil
	}
	return nil, errors.ThrowNotFound(nil, "EVENT-Qlo0d", "Errors.Org.IdpNotExisting")
}

func (es *OrgEventstore) AddIDPProviderToLoginPolicy(ctx context.Context, provider *iam_model.IDPProvider) (*iam_model.IDPProvider, error) {
	if provider == nil || !provider.IsValid() {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Sjd8e", "Errors.Org.IdpProviderInvalid")
	}
	org, err := es.OrgByID(ctx, org_model.NewOrg(provider.AggregateID))
	if err != nil {
		return nil, err
	}
	if org.LoginPolicy == nil {
		return nil, errors.ThrowAlreadyExists(nil, "EVENT-sk9fW", "Errors.Org.LoginPolicy.NotExisting")
	}
	if _, m := org.LoginPolicy.GetIdpProvider(provider.IdpConfigID); m != nil {
		return nil, errors.ThrowAlreadyExists(nil, "EVENT-Lso9f", "Errors.Org.LoginPolicy.IdpProviderAlreadyExisting")
	}
	repoOrg := model.OrgFromModel(org)
	repoProvider := iam_es_model.IDPProviderFromModel(provider)

	addAggregate := LoginPolicyIDPProviderAddedAggregate(es.Eventstore.AggregateCreator(), repoOrg, repoProvider, es.IamID)
	err = es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	if _, m := iam_es_model.GetIDPProvider(repoOrg.LoginPolicy.IDPProviders, provider.IdpConfigID); m != nil {
		return iam_es_model.IDPProviderToModel(m), nil
	}
	return nil, errors.ThrowInternal(nil, "EVENT-Slf9s", "Errors.Internal")
}

func (es *OrgEventstore) PrepareRemoveIDPProviderFromLoginPolicy(ctx context.Context, provider *iam_model.IDPProvider, cascade bool) (*model.Org, *es_models.Aggregate, error) {
	if provider == nil || !provider.IsValid() {
		return nil, nil, errors.ThrowPreconditionFailed(nil, "EVENT-Esi8c", "Errors.IdpProviderInvalid")
	}
	org, err := es.OrgByID(ctx, org_model.NewOrg(provider.AggregateID))
	if err != nil {
		return nil, nil, err
	}
	if _, m := org.LoginPolicy.GetIdpProvider(provider.IdpConfigID); m == nil {
		return nil, nil, errors.ThrowPreconditionFailed(nil, "EVENT-29skr", "Errors.IAM.LoginPolicy.IdpProviderNotExisting")
	}
	repoOrg := model.OrgFromModel(org)
	providerID := &iam_es_model.IDPProviderID{provider.IdpConfigID}
	providerAggregates, err := LoginPolicyIDPProviderRemovedAggregate(ctx, es.Eventstore.AggregateCreator(), repoOrg, providerID, cascade)
	if err != nil {
		return nil, nil, err
	}
	return repoOrg, providerAggregates, nil
}

func (es *OrgEventstore) RemoveIDPProviderFromLoginPolicy(ctx context.Context, provider *iam_model.IDPProvider) error {
	repoOrg, agg, err := es.PrepareRemoveIDPProviderFromLoginPolicy(ctx, provider, false)
	if err != nil {
		return err
	}
	return es_sdk.PushAggregates(ctx, es.PushAggregates, repoOrg.AppendEvents, agg)
}

func (es *OrgEventstore) AddPasswordComplexityPolicy(ctx context.Context, policy *iam_model.PasswordComplexityPolicy) (*iam_model.PasswordComplexityPolicy, error) {
	if policy == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Sjkl9", "Errors.Org.PasswordComplexityPolicy.Invalid")
	}

	if err := policy.IsValid(); err != nil {
		return nil, err
	}
	org, err := es.OrgByID(ctx, org_model.NewOrg(policy.AggregateID))
	if err != nil {
		return nil, err
	}

	repoOrg := model.OrgFromModel(org)
	repoPasswordComplexityPolicy := iam_es_model.PasswordComplexityPolicyFromModel(policy)

	addAggregate := PasswordComplexityPolicyAddedAggregate(es.Eventstore.AggregateCreator(), repoOrg, repoPasswordComplexityPolicy)
	err = es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	return iam_es_model.PasswordComplexityPolicyToModel(repoOrg.PasswordComplexityPolicy), nil
}

func (es *OrgEventstore) ChangePasswordComplexityPolicy(ctx context.Context, policy *iam_model.PasswordComplexityPolicy) (*iam_model.PasswordComplexityPolicy, error) {
	if policy == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-r5Hd", "Errors.Org.PasswordComplexityPolicy.Empty")
	}
	if err := policy.IsValid(); err != nil {
		return nil, err
	}
	org, err := es.OrgByID(ctx, org_model.NewOrg(policy.AggregateID))
	if err != nil {
		return nil, err
	}

	if org.PasswordComplexityPolicy == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-v6Hdr", "Errors.Org.PasswordComplexityPolicy.NotExisting")
	}

	repoOrg := model.OrgFromModel(org)
	repoPasswordComplexityPolicy := iam_es_model.PasswordComplexityPolicyFromModel(policy)

	addAggregate := PasswordComplexityPolicyChangedAggregate(es.Eventstore.AggregateCreator(), repoOrg, repoPasswordComplexityPolicy)
	err = es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	return iam_es_model.PasswordComplexityPolicyToModel(repoOrg.PasswordComplexityPolicy), nil
}

func (es *OrgEventstore) RemovePasswordComplexityPolicy(ctx context.Context, policy *iam_model.PasswordComplexityPolicy) error {
	if policy == nil || policy.AggregateID == "" {
		return errors.ThrowPreconditionFailed(nil, "EVENT-3Ghs8", "Errors.Org.PasswordComplexityPolicy.Invalid")
	}
	org, err := es.OrgByID(ctx, org_model.NewOrg(policy.AggregateID))
	if err != nil {
		return err
	}
	repoOrg := model.OrgFromModel(org)

	addAggregate := PasswordComplexityPolicyRemovedAggregate(es.Eventstore.AggregateCreator(), repoOrg)
	return es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, addAggregate)
}

func (es *OrgEventstore) AddPasswordAgePolicy(ctx context.Context, policy *iam_model.PasswordAgePolicy) (*iam_model.PasswordAgePolicy, error) {
	if policy == nil || policy.AggregateID == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Sjkl9", "Errors.Org.PasswordAgePolicy.Invalid")
	}
	org, err := es.OrgByID(ctx, org_model.NewOrg(policy.AggregateID))
	if err != nil {
		return nil, err
	}

	repoOrg := model.OrgFromModel(org)
	repoPasswordAgePolicy := iam_es_model.PasswordAgePolicyFromModel(policy)

	addAggregate := PasswordAgePolicyAddedAggregate(es.Eventstore.AggregateCreator(), repoOrg, repoPasswordAgePolicy)
	err = es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	return iam_es_model.PasswordAgePolicyToModel(repoOrg.PasswordAgePolicy), nil
}

func (es *OrgEventstore) ChangePasswordAgePolicy(ctx context.Context, policy *iam_model.PasswordAgePolicy) (*iam_model.PasswordAgePolicy, error) {
	if policy == nil || policy.AggregateID == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-r5Hd", "Errors.Org.PasswordAgePolicy.Empty")
	}
	org, err := es.OrgByID(ctx, org_model.NewOrg(policy.AggregateID))
	if err != nil {
		return nil, err
	}

	if org.PasswordAgePolicy == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-v6Hdr", "Errors.Org.PasswordAgePolicy.NotExisting")
	}

	repoOrg := model.OrgFromModel(org)
	repoPasswordAgePolicy := iam_es_model.PasswordAgePolicyFromModel(policy)

	addAggregate := PasswordAgePolicyChangedAggregate(es.Eventstore.AggregateCreator(), repoOrg, repoPasswordAgePolicy)
	err = es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	return iam_es_model.PasswordAgePolicyToModel(repoOrg.PasswordAgePolicy), nil
}

func (es *OrgEventstore) RemovePasswordAgePolicy(ctx context.Context, policy *iam_model.PasswordAgePolicy) error {
	if policy == nil || policy.AggregateID == "" {
		return errors.ThrowPreconditionFailed(nil, "EVENT-3Ghs8", "Errors.Org.PasswordAgePolicy.Invalid")
	}
	org, err := es.OrgByID(ctx, org_model.NewOrg(policy.AggregateID))
	if err != nil {
		return err
	}
	repoOrg := model.OrgFromModel(org)

	addAggregate := PasswordAgePolicyRemovedAggregate(es.Eventstore.AggregateCreator(), repoOrg)
	return es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, addAggregate)
}

func (es *OrgEventstore) AddPasswordLockoutPolicy(ctx context.Context, policy *iam_model.PasswordLockoutPolicy) (*iam_model.PasswordLockoutPolicy, error) {
	if policy == nil || policy.AggregateID == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-6Zdk9", "Errors.Org.PasswordLockoutPolicy.Invalid")
	}
	org, err := es.OrgByID(ctx, org_model.NewOrg(policy.AggregateID))
	if err != nil {
		return nil, err
	}

	repoOrg := model.OrgFromModel(org)
	repoPasswordLockoutPolicy := iam_es_model.PasswordLockoutPolicyFromModel(policy)

	addAggregate := PasswordLockoutPolicyAddedAggregate(es.Eventstore.AggregateCreator(), repoOrg, repoPasswordLockoutPolicy)
	err = es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	return iam_es_model.PasswordLockoutPolicyToModel(repoOrg.PasswordLockoutPolicy), nil
}

func (es *OrgEventstore) ChangePasswordLockoutPolicy(ctx context.Context, policy *iam_model.PasswordLockoutPolicy) (*iam_model.PasswordLockoutPolicy, error) {
	if policy == nil || policy.AggregateID == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-lp0Sf", "Errors.Org.PasswordLockoutPolicy.Empty")
	}
	org, err := es.OrgByID(ctx, org_model.NewOrg(policy.AggregateID))
	if err != nil {
		return nil, err
	}

	if org.PasswordLockoutPolicy == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-3Fks9", "Errors.Org.PasswordLockoutPolicy.NotExisting")
	}

	repoOrg := model.OrgFromModel(org)
	repoPasswordLockoutPolicy := iam_es_model.PasswordLockoutPolicyFromModel(policy)

	addAggregate := PasswordLockoutPolicyChangedAggregate(es.Eventstore.AggregateCreator(), repoOrg, repoPasswordLockoutPolicy)
	err = es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	return iam_es_model.PasswordLockoutPolicyToModel(repoOrg.PasswordLockoutPolicy), nil
}

func (es *OrgEventstore) RemovePasswordLockoutPolicy(ctx context.Context, policy *iam_model.PasswordLockoutPolicy) error {
	if policy == nil || policy.AggregateID == "" {
		return errors.ThrowPreconditionFailed(nil, "EVENT-6Hls0", "Errors.Org.PasswordLockoutPolicy.Invalid")
	}
	org, err := es.OrgByID(ctx, org_model.NewOrg(policy.AggregateID))
	if err != nil {
		return err
	}
	repoOrg := model.OrgFromModel(org)

	addAggregate := PasswordLockoutPolicyRemovedAggregate(es.Eventstore.AggregateCreator(), repoOrg)
	return es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, addAggregate)
}

func (es *OrgEventstore) AddMailTemplate(ctx context.Context, template *iam_model.MailTemplate) (*iam_model.MailTemplate, error) {
	if template == nil || !template.IsValid() {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-nb66d", "Errors.Org.MailTemplateInvalid")
	}
	org, err := es.OrgByID(ctx, org_model.NewOrg(template.AggregateID))
	if err != nil {
		return nil, err
	}

	repoOrg := model.OrgFromModel(org)
	repoMailTemplate := iam_es_model.MailTemplateFromModel(template)

	addAggregate := MailTemplateAddedAggregate(es.Eventstore.AggregateCreator(), repoOrg, repoMailTemplate)
	err = es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	return iam_es_model.MailTemplateToModel(repoOrg.MailTemplate), nil
}

func (es *OrgEventstore) ChangeMailTemplate(ctx context.Context, template *iam_model.MailTemplate) (*iam_model.MailTemplate, error) {
	if template == nil || !template.IsValid() {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-FV2qE", "Errors.Org.MailTemplateInvalid")
	}
	org, err := es.OrgByID(ctx, org_model.NewOrg(template.AggregateID))
	if err != nil {
		return nil, err
	}

	repoOrg := model.OrgFromModel(org)
	repoMailTemplate := iam_es_model.MailTemplateFromModel(template)

	addAggregate := MailTemplateChangedAggregate(es.Eventstore.AggregateCreator(), repoOrg, repoMailTemplate)
	err = es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	return iam_es_model.MailTemplateToModel(repoOrg.MailTemplate), nil
}

// ToDo Michi
func (es *OrgEventstore) AddMailText(ctx context.Context, mailtext *iam_model.MailText) (*iam_model.MailText, error) {
	if mailtext == nil || !mailtext.IsValid() {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-108Iz", "Errors.Org.MailTextInvalid")
	}
	org, err := es.OrgByID(ctx, org_model.NewOrg(mailtext.AggregateID))
	if err != nil {
		return nil, err
	}

	repoOrg := model.OrgFromModel(org)
	repoMailText := iam_es_model.MailTextFromModel(mailtext)

	addAggregate := MailTextAddedAggregate(es.Eventstore.AggregateCreator(), repoOrg, repoMailText)
	err = es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	// ToDo Michi
	return iam_es_model.MailTextToModel(repoOrg.MailTexts[0]), nil
}

func (es *OrgEventstore) ChangeMailText(ctx context.Context, mailtext *iam_model.MailText) (*iam_model.MailText, error) {
	if mailtext == nil || !mailtext.IsValid() {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-fdbqE", "Errors.Org.MailTextInvalid")
	}
	org, err := es.OrgByID(ctx, org_model.NewOrg(mailtext.AggregateID))
	if err != nil {
		return nil, err
	}

	repoOrg := model.OrgFromModel(org)
	repoMailText := iam_es_model.MailTextFromModel(mailtext)

	addAggregate := MailTextChangedAggregate(es.Eventstore.AggregateCreator(), repoOrg, repoMailText)
	err = es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, addAggregate)
	if err != nil {
		return nil, err
	}
	// ToDo Michi
	return iam_es_model.MailTextToModel(repoOrg.MailTexts[0]), nil
}

func (es *OrgEventstore) RemoveMailTemplate(ctx context.Context, policy *iam_model.MailTemplate) error {
	if policy == nil || !policy.IsValid() {
		return errors.ThrowPreconditionFailed(nil, "EVENT-LulaW", "Errors.Org.MailTemplate.Invalid")
	}
	org, err := es.OrgByID(ctx, org_model.NewOrg(policy.AggregateID))
	if err != nil {
		return err
	}
	repoOrg := model.OrgFromModel(org)

	addAggregate := MailTemplateRemovedAggregate(es.Eventstore.AggregateCreator(), repoOrg)
	return es_sdk.Push(ctx, es.PushAggregates, repoOrg.AppendEvents, addAggregate)
}
