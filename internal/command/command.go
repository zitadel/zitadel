package command

import (
	"context"
	"net/http"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	api_http "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/command/preparation"
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
	"github.com/zitadel/zitadel/internal/repository/quota"
	usr_repo "github.com/zitadel/zitadel/internal/repository/user"
	usr_grant_repo "github.com/zitadel/zitadel/internal/repository/usergrant"
	"github.com/zitadel/zitadel/internal/static"
	webauthn_helper "github.com/zitadel/zitadel/internal/webauthn"
)

type Commands struct {
	httpClient *http.Client

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
	domainVerificationValidator func(domain, token, verifier string, checkType api_http.CheckType) error

	multifactors         domain.MultifactorConfigs
	webauthnConfig       *webauthn_helper.Config
	keySize              int
	keyAlgorithm         crypto.EncryptionAlgorithm
	certificateAlgorithm crypto.EncryptionAlgorithm
	certKeySize          int
	privateKeyLifetime   time.Duration
	publicKeyLifetime    time.Duration
	certificateLifetime  time.Duration
}

func StartCommands(es *eventstore.Eventstore,
	defaults sd.SystemDefaults,
	zitadelRoles []authz.RoleMapping,
	staticStore static.Storage,
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
	httpClient *http.Client,
) (repo *Commands, err error) {
	if externalDomain == "" {
		return nil, errors.ThrowInvalidArgument(nil, "COMMAND-Df21s", "no external domain specified")
	}
	repo = &Commands{
		eventstore:            es,
		static:                staticStore,
		idGenerator:           id.SonyFlakeGenerator(),
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
		httpClient:            httpClient,
	}

	instance_repo.RegisterEventMappers(repo.eventstore)
	org.RegisterEventMappers(repo.eventstore)
	usr_repo.RegisterEventMappers(repo.eventstore)
	usr_grant_repo.RegisterEventMappers(repo.eventstore)
	proj_repo.RegisterEventMappers(repo.eventstore)
	keypair.RegisterEventMappers(repo.eventstore)
	action.RegisterEventMappers(repo.eventstore)
	quota.RegisterEventMappers(repo.eventstore)

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
	repo.domainVerificationValidator = api_http.ValidateDomain
	return repo, nil
}

func AppendAndReduce(object interface {
	AppendEvents(...eventstore.Event)
	// TODO: Why is it allowed to return an error here?
	Reduce() error
}, events ...eventstore.Event) error {
	object.AppendEvents(events...)
	return object.Reduce()
}

func queryAndReduce(ctx context.Context, filter preparation.FilterToQueryReducer, wm eventstore.QueryReducer) error {
	events, err := filter(ctx, wm.Query())
	if err != nil {
		return err
	}
	if len(events) == 0 {
		return nil
	}
	wm.AppendEvents(events...)
	return wm.Reduce()
}

type existsWriteModel interface {
	Exists() bool
	eventstore.QueryReducer
}

func exists(ctx context.Context, filter preparation.FilterToQueryReducer, wm existsWriteModel) (bool, error) {
	err := queryAndReduce(ctx, filter, wm)
	if err != nil {
		return false, err
	}
	return wm.Exists(), nil
}
