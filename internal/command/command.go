package command

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore"
	"time"

	"github.com/caos/zitadel/internal/api/http"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/id"
	global_model "github.com/caos/zitadel/internal/model"
	iam_repo "github.com/caos/zitadel/internal/repository/iam"
	keypair "github.com/caos/zitadel/internal/repository/keypair"
	"github.com/caos/zitadel/internal/repository/org"
	proj_repo "github.com/caos/zitadel/internal/repository/project"
	usr_repo "github.com/caos/zitadel/internal/repository/user"
	usr_grant_repo "github.com/caos/zitadel/internal/repository/usergrant"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	webauthn_helper "github.com/caos/zitadel/internal/webauthn"
)

type CommandSide struct {
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
	//TODO: remove global model, or move to domain
	multifactors global_model.Multifactors

	webauthn           *webauthn_helper.WebAuthN
	keySize            int
	keyAlgorithm       crypto.EncryptionAlgorithm
	privateKeyLifetime time.Duration
	publicKeyLifetime  time.Duration
}

type Config struct {
	Eventstore     *eventstore.Eventstore
	SystemDefaults sd.SystemDefaults
}

func StartCommandSide(config *Config) (repo *CommandSide, err error) {
	repo = &CommandSide{
		eventstore:         config.Eventstore,
		idGenerator:        id.SonyFlakeGenerator,
		iamDomain:          config.SystemDefaults.Domain,
		keySize:            config.SystemDefaults.KeyConfig.Size,
		privateKeyLifetime: config.SystemDefaults.KeyConfig.PrivateKeyLifetime.Duration,
		publicKeyLifetime:  config.SystemDefaults.KeyConfig.PublicKeyLifetime.Duration,
	}
	iam_repo.RegisterEventMappers(repo.eventstore)
	org.RegisterEventMappers(repo.eventstore)
	usr_repo.RegisterEventMappers(repo.eventstore)
	usr_grant_repo.RegisterEventMappers(repo.eventstore)
	proj_repo.RegisterEventMappers(repo.eventstore)
	keypair.RegisterEventMappers(repo.eventstore)

	//TODO: simplify!!!!
	repo.idpConfigSecretCrypto, err = crypto.NewAESCrypto(config.SystemDefaults.IDPConfigVerificationKey)
	if err != nil {
		return nil, err
	}
	userEncryptionAlgorithm, err := crypto.NewAESCrypto(config.SystemDefaults.UserVerificationKey)
	if err != nil {
		return nil, err
	}
	repo.initializeUserCode = crypto.NewEncryptionGenerator(config.SystemDefaults.SecretGenerators.InitializeUserCode, userEncryptionAlgorithm)
	repo.emailVerificationCode = crypto.NewEncryptionGenerator(config.SystemDefaults.SecretGenerators.EmailVerificationCode, userEncryptionAlgorithm)
	repo.phoneVerificationCode = crypto.NewEncryptionGenerator(config.SystemDefaults.SecretGenerators.PhoneVerificationCode, userEncryptionAlgorithm)
	repo.passwordVerificationCode = crypto.NewEncryptionGenerator(config.SystemDefaults.SecretGenerators.PasswordVerificationCode, userEncryptionAlgorithm)
	repo.userPasswordAlg = crypto.NewBCrypt(config.SystemDefaults.SecretGenerators.PasswordSaltCost)
	repo.machineKeyAlg = userEncryptionAlgorithm
	repo.machineKeySize = int(config.SystemDefaults.SecretGenerators.MachineKeySize)
	repo.applicationKeySize = int(config.SystemDefaults.SecretGenerators.ApplicationKeySize)

	aesOTPCrypto, err := crypto.NewAESCrypto(config.SystemDefaults.Multifactors.OTP.VerificationKey)
	if err != nil {
		return nil, err
	}
	repo.multifactors = global_model.Multifactors{
		OTP: global_model.OTP{
			CryptoMFA: aesOTPCrypto,
			Issuer:    config.SystemDefaults.Multifactors.OTP.Issuer,
		},
	}
	passwordAlg := crypto.NewBCrypt(config.SystemDefaults.SecretGenerators.PasswordSaltCost)
	repo.applicationSecretGenerator = crypto.NewHashGenerator(config.SystemDefaults.SecretGenerators.ClientSecretGenerator, passwordAlg)

	repo.domainVerificationAlg, err = crypto.NewAESCrypto(config.SystemDefaults.DomainVerification.VerificationKey)
	if err != nil {
		return nil, err
	}
	repo.domainVerificationGenerator = crypto.NewEncryptionGenerator(config.SystemDefaults.DomainVerification.VerificationGenerator, repo.domainVerificationAlg)
	repo.domainVerificationValidator = http.ValidateDomain
	web, err := webauthn_helper.StartServer(config.SystemDefaults.WebAuthN)
	if err != nil {
		return nil, err
	}
	repo.webauthn = web

	keyAlgorithm, err := crypto.NewAESCrypto(config.SystemDefaults.KeyConfig.EncryptionConfig)
	if err != nil {
		return nil, err
	}
	repo.keyAlgorithm = keyAlgorithm
	return repo, nil
}

func (r *CommandSide) getIAMWriteModel(ctx context.Context) (_ *IAMWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewIAMWriteModel()
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
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
