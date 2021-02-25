package command

import (
	"context"
	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"time"

	"github.com/caos/zitadel/internal/api/http"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/id"
	iam_repo "github.com/caos/zitadel/internal/repository/iam"
	keypair "github.com/caos/zitadel/internal/repository/keypair"
	"github.com/caos/zitadel/internal/repository/org"
	proj_repo "github.com/caos/zitadel/internal/repository/project"
	usr_repo "github.com/caos/zitadel/internal/repository/user"
	usr_grant_repo "github.com/caos/zitadel/internal/repository/usergrant"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	webauthn_helper "github.com/caos/zitadel/internal/webauthn"
)

type Commands struct {
	eventstore  *eventstore.Eventstore
	idGenerator id.Generator
	iamDomain   string

	idpConfigSecretCrypto crypto.Crypto

	userPasswordAlg             crypto.HashAlgorithm
	initializeUserCode          crypto.Generator
	emailVerificationCode       crypto.Generator
	phoneVerificationCode       crypto.Generator
	passwordVerificationCode    crypto.Generator
	machineKeyAlg               crypto.EncryptionAlgorithm
	machineKeySize              int
	applicationKeySize          int
	applicationSecretGenerator  crypto.Generator
	domainVerificationAlg       *crypto.AESCrypto
	domainVerificationGenerator crypto.Generator
	domainVerificationValidator func(domain, token, verifier string, checkType http.CheckType) error
	multifactors                domain.MultifactorConfigs

	webauthn           *webauthn_helper.WebAuthN
	keySize            int
	keyAlgorithm       crypto.EncryptionAlgorithm
	privateKeyLifetime time.Duration
	publicKeyLifetime  time.Duration
}

type Config struct {
	Eventstore types.SQLUser
}

func StartCommands(eventstore *eventstore.Eventstore, defaults sd.SystemDefaults) (repo *Commands, err error) {
	repo = &Commands{
		eventstore:         eventstore,
		idGenerator:        id.SonyFlakeGenerator,
		iamDomain:          defaults.Domain,
		keySize:            defaults.KeyConfig.Size,
		privateKeyLifetime: defaults.KeyConfig.PrivateKeyLifetime.Duration,
		publicKeyLifetime:  defaults.KeyConfig.PublicKeyLifetime.Duration,
	}
	iam_repo.RegisterEventMappers(repo.eventstore)
	org.RegisterEventMappers(repo.eventstore)
	usr_repo.RegisterEventMappers(repo.eventstore)
	usr_grant_repo.RegisterEventMappers(repo.eventstore)
	proj_repo.RegisterEventMappers(repo.eventstore)
	keypair.RegisterEventMappers(repo.eventstore)

	repo.idpConfigSecretCrypto, err = crypto.NewAESCrypto(defaults.IDPConfigVerificationKey)
	if err != nil {
		return nil, err
	}
	userEncryptionAlgorithm, err := crypto.NewAESCrypto(defaults.UserVerificationKey)
	if err != nil {
		return nil, err
	}
	repo.initializeUserCode = crypto.NewEncryptionGenerator(defaults.SecretGenerators.InitializeUserCode, userEncryptionAlgorithm)
	repo.emailVerificationCode = crypto.NewEncryptionGenerator(defaults.SecretGenerators.EmailVerificationCode, userEncryptionAlgorithm)
	repo.phoneVerificationCode = crypto.NewEncryptionGenerator(defaults.SecretGenerators.PhoneVerificationCode, userEncryptionAlgorithm)
	repo.passwordVerificationCode = crypto.NewEncryptionGenerator(defaults.SecretGenerators.PasswordVerificationCode, userEncryptionAlgorithm)
	repo.userPasswordAlg = crypto.NewBCrypt(defaults.SecretGenerators.PasswordSaltCost)
	repo.machineKeyAlg = userEncryptionAlgorithm
	repo.machineKeySize = int(defaults.SecretGenerators.MachineKeySize)
	repo.applicationKeySize = int(defaults.SecretGenerators.ApplicationKeySize)

	aesOTPCrypto, err := crypto.NewAESCrypto(defaults.Multifactors.OTP.VerificationKey)
	if err != nil {
		return nil, err
	}
	repo.multifactors = domain.MultifactorConfigs{
		OTP: domain.OTPConfig{
			CryptoMFA: aesOTPCrypto,
			Issuer:    defaults.Multifactors.OTP.Issuer,
		},
	}
	passwordAlg := crypto.NewBCrypt(defaults.SecretGenerators.PasswordSaltCost)
	repo.applicationSecretGenerator = crypto.NewHashGenerator(defaults.SecretGenerators.ClientSecretGenerator, passwordAlg)

	repo.domainVerificationAlg, err = crypto.NewAESCrypto(defaults.DomainVerification.VerificationKey)
	if err != nil {
		return nil, err
	}
	repo.domainVerificationGenerator = crypto.NewEncryptionGenerator(defaults.DomainVerification.VerificationGenerator, repo.domainVerificationAlg)
	repo.domainVerificationValidator = http.ValidateDomain
	web, err := webauthn_helper.StartServer(defaults.WebAuthN)
	if err != nil {
		return nil, err
	}
	repo.webauthn = web

	keyAlgorithm, err := crypto.NewAESCrypto(defaults.KeyConfig.EncryptionConfig)
	if err != nil {
		return nil, err
	}
	repo.keyAlgorithm = keyAlgorithm
	return repo, nil
}

func (c *Commands) getIAMWriteModel(ctx context.Context) (_ *IAMWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewIAMWriteModel()
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	return writeModel, nil
}

func AppendAndReduce(object interface {
	AppendEvents(...eventstore.EventReader)
	Reduce() error
}, events ...eventstore.EventReader) error {
	object.AppendEvents(events...)
	return object.Reduce()
}
