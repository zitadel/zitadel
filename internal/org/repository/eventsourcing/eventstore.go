package eventsourcing

import (
	"context"
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
