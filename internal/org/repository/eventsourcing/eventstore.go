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

func (es *OrgEventstore) OrgEventsByID(ctx context.Context, id string, sequence uint64) ([]*es_models.Event, error) {
	query, err := OrgByIDQuery(id, sequence)
	if err != nil {
		return nil, err
	}
	return es.FilterEvents(ctx, query)
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
