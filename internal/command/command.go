package command

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/http"
	authz_repo "github.com/caos/zitadel/internal/authz/repository"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/id"
	"github.com/caos/zitadel/internal/repository/action"
	iam_repo "github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/keypair"
	"github.com/caos/zitadel/internal/repository/org"
	proj_repo "github.com/caos/zitadel/internal/repository/project"
	usr_repo "github.com/caos/zitadel/internal/repository/user"
	usr_grant_repo "github.com/caos/zitadel/internal/repository/usergrant"
	"github.com/caos/zitadel/internal/static"
	webauthn_helper "github.com/caos/zitadel/internal/webauthn"
)

type Commands struct {
	eventstore   *eventstore.Eventstore
	static       static.Storage
	idGenerator  id.Generator
	iamDomain    string
	zitadelRoles []authz.RoleMapping

	idpConfigSecretCrypto crypto.EncryptionAlgorithm
	smtpPasswordCrypto    crypto.EncryptionAlgorithm
	smsCrypto             crypto.EncryptionAlgorithm

	userPasswordAlg             crypto.HashAlgorithm
	machineKeySize              int
	applicationKeySize          int
	domainVerificationAlg       crypto.EncryptionAlgorithm
	domainVerificationGenerator crypto.Generator
	domainVerificationValidator func(domain, token, verifier string, checkType http.CheckType) error
	multifactors                domain.MultifactorConfigs

	webauthn           *webauthn_helper.WebAuthN
	keySize            int
	keyAlgorithm       crypto.EncryptionAlgorithm
	privateKeyLifetime time.Duration
	publicKeyLifetime  time.Duration
	tokenVerifier      orgFeatureChecker
}

type orgFeatureChecker interface {
	CheckOrgFeatures(ctx context.Context, orgID string, requiredFeatures ...string) error
}

func StartCommands(es *eventstore.Eventstore, defaults sd.SystemDefaults, authZConfig authz.Config, staticStore static.Storage, authZRepo authz_repo.Repository, encryptionKeyConfig *crypto.EncryptionKeys, keyStorage crypto.KeyStorage, webAuthN webauthn_helper.Config) (repo *Commands, err error) {
	repo = &Commands{
		eventstore:         es,
		static:             staticStore,
		idGenerator:        id.SonyFlakeGenerator,
		iamDomain:          defaults.Domain,
		zitadelRoles:       authZConfig.RolePermissionMappings,
		keySize:            defaults.KeyConfig.Size,
		privateKeyLifetime: defaults.KeyConfig.PrivateKeyLifetime,
		publicKeyLifetime:  defaults.KeyConfig.PublicKeyLifetime,
	}
	iam_repo.RegisterEventMappers(repo.eventstore)
	org.RegisterEventMappers(repo.eventstore)
	usr_repo.RegisterEventMappers(repo.eventstore)
	usr_grant_repo.RegisterEventMappers(repo.eventstore)
	proj_repo.RegisterEventMappers(repo.eventstore)
	keypair.RegisterEventMappers(repo.eventstore)
	action.RegisterEventMappers(repo.eventstore)

	repo.idpConfigSecretCrypto, err = crypto.NewAESCrypto(encryptionKeyConfig.IDPConfig, keyStorage)
	if err != nil {
		return nil, err
	}

	repo.userPasswordAlg = crypto.NewBCrypt(defaults.SecretGenerators.PasswordSaltCost)
	repo.machineKeySize = int(defaults.SecretGenerators.MachineKeySize)
	repo.applicationKeySize = int(defaults.SecretGenerators.ApplicationKeySize)

	aesOTPCrypto, err := crypto.NewAESCrypto(encryptionKeyConfig.OTP, keyStorage)
	if err != nil {
		return nil, err
	}
	repo.multifactors = domain.MultifactorConfigs{
		OTP: domain.OTPConfig{
			CryptoMFA: aesOTPCrypto,
			Issuer:    defaults.Multifactors.OTP.Issuer,
		},
	}

	repo.smtpPasswordCrypto, err = crypto.NewAESCrypto(encryptionKeyConfig.SMTP, keyStorage)
	if err != nil {
		return nil, err
	}
	repo.smsCrypto, err = crypto.NewAESCrypto(encryptionKeyConfig.SMS, keyStorage)
	if err != nil {
		return nil, err
	}

	repo.domainVerificationAlg, err = crypto.NewAESCrypto(encryptionKeyConfig.DomainVerification, keyStorage)
	if err != nil {
		return nil, err
	}
	repo.domainVerificationGenerator = crypto.NewEncryptionGenerator(defaults.DomainVerification.VerificationGenerator, repo.domainVerificationAlg)
	repo.domainVerificationValidator = http.ValidateDomain
	web, err := webauthn_helper.StartServer(webAuthN)
	if err != nil {
		return nil, err
	}
	repo.webauthn = web

	keyAlgorithm, err := crypto.NewAESCrypto(encryptionKeyConfig.OIDC, keyStorage)
	if err != nil {
		return nil, err
	}
	repo.keyAlgorithm = keyAlgorithm

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
