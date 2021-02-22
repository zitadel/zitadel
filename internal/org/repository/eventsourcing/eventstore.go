package eventsourcing

import (
	"github.com/caos/logging"
	http_utils "github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/eventstore"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/id"
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
