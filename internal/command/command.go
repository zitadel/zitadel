package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/http"
	authz_repo "github.com/zitadel/zitadel/internal/authz/repository"
	sd "github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/repository/action"
	instance_repo "github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/keypair"
	"github.com/zitadel/zitadel/internal/repository/org"
	proj_repo "github.com/zitadel/zitadel/internal/repository/project"
	usr_repo "github.com/zitadel/zitadel/internal/repository/user"
	usr_grant_repo "github.com/zitadel/zitadel/internal/repository/usergrant"
	"github.com/zitadel/zitadel/internal/static"
	webauthn_helper "github.com/zitadel/zitadel/internal/webauthn"
)

type Commands struct {
	eventstore     *eventstore.Eventstore
	static         static.Storage
	idGenerator    id.Generator
	zitadelRoles   []authz.RoleMapping
	externalDomain string
	externalSecure bool
	externalPort   uint16

	idpConfigEncryption         crypto.EncryptionAlgorithm
	smtpEncryption              crypto.EncryptionAlgorithm
	smsEncryption               crypto.EncryptionAlgorithm
	userEncryption              crypto.EncryptionAlgorithm
	userPasswordAlg             crypto.HashAlgorithm
	machineKeySize              int
	applicationKeySize          int
	domainVerificationAlg       crypto.EncryptionAlgorithm
	domainVerificationGenerator crypto.Generator
	domainVerificationValidator func(domain, token, verifier string, checkType http.CheckType) error

	multifactors         domain.MultifactorConfigs
	webauthnConfig       *webauthn_helper.Config
	keySize              int
	keyAlgorithm         crypto.EncryptionAlgorithm
	certificateAlgorithm crypto.EncryptionAlgorithm
	certKeySize          int
	privateKeyLifetime   time.Duration
	publicKeyLifetime    time.Duration
	certificateLifetime  time.Duration

	tokenVerifier orgFeatureChecker
}

type orgFeatureChecker interface {
	CheckOrgFeatures(ctx context.Context, orgID string, requiredFeatures ...string) error
}

func StartCommands(es *eventstore.Eventstore,
	defaults sd.SystemDefaults,
	zitadelRoles []authz.RoleMapping,
	staticStore static.Storage,
	authZRepo authz_repo.Repository,
	webAuthN *webauthn_helper.Config,
	externalDomain string,
	externalSecure bool,
	externalPort uint16,
	idpConfigEncryption,
	otpEncryption,
	smtpEncryption,
	smsEncryption,
	userEncryption,
	domainVerificationEncryption,
	oidcEncryption,
	samlEncryption crypto.EncryptionAlgorithm,
) (repo *Commands, err error) {
	if externalDomain == "" {
		return nil, errors.ThrowInvalidArgument(nil, "COMMAND-Df21s", "not external domain specified")
	}
	repo = &Commands{
		eventstore:            es,
		static:                staticStore,
		idGenerator:           id.SonyFlakeGenerator,
		zitadelRoles:          zitadelRoles,
		externalDomain:        externalDomain,
		externalSecure:        externalSecure,
		externalPort:          externalPort,
		keySize:               defaults.KeyConfig.Size,
		certKeySize:           defaults.KeyConfig.CertificateSize,
		privateKeyLifetime:    defaults.KeyConfig.PrivateKeyLifetime,
		publicKeyLifetime:     defaults.KeyConfig.PublicKeyLifetime,
		certificateLifetime:   defaults.KeyConfig.CertificateLifetime,
		idpConfigEncryption:   idpConfigEncryption,
		smtpEncryption:        smtpEncryption,
		smsEncryption:         smsEncryption,
		userEncryption:        userEncryption,
		domainVerificationAlg: domainVerificationEncryption,
		keyAlgorithm:          oidcEncryption,
		certificateAlgorithm:  samlEncryption,
		webauthnConfig:        webAuthN,
	}

	instance_repo.RegisterEventMappers(repo.eventstore)
	org.RegisterEventMappers(repo.eventstore)
	usr_repo.RegisterEventMappers(repo.eventstore)
	usr_grant_repo.RegisterEventMappers(repo.eventstore)
	proj_repo.RegisterEventMappers(repo.eventstore)
	keypair.RegisterEventMappers(repo.eventstore)
	action.RegisterEventMappers(repo.eventstore)

	repo.userPasswordAlg = crypto.NewBCrypt(defaults.SecretGenerators.PasswordSaltCost)
	repo.machineKeySize = int(defaults.SecretGenerators.MachineKeySize)
	repo.applicationKeySize = int(defaults.SecretGenerators.ApplicationKeySize)

	repo.multifactors = domain.MultifactorConfigs{
		OTP: domain.OTPConfig{
			CryptoMFA: otpEncryption,
			Issuer:    defaults.Multifactors.OTP.Issuer,
		},
	}

	repo.domainVerificationGenerator = crypto.NewEncryptionGenerator(defaults.DomainVerification.VerificationGenerator, repo.domainVerificationAlg)
	repo.domainVerificationValidator = http.ValidateDomain

	repo.tokenVerifier = authZRepo
	return repo, nil
}

func AppendAndReduce(object interface {
	AppendEvents(...eventstore.Event)
	Reduce() error
}, events ...eventstore.Event) error {
	object.AppendEvents(events...)
	return object.Reduce()
}
