package command

import (
	"context"

	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/id"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

type CommandSide struct {
	eventstore  *eventstore.Eventstore
	idGenerator id.Generator
	iamID       string
	iamDomain   string

	idpConfigSecretCrypto crypto.Crypto

	userPasswordAlg          crypto.HashAlgorithm
	initializeUserCode       crypto.Generator
	emailVerificationCode    crypto.Generator
	phoneVerificationCode    crypto.Generator
	passwordVerificationCode crypto.Generator
	machineKeyAlg            crypto.EncryptionAlgorithm
	machineKeySize           int
}

type Config struct {
	Eventstore     *eventstore.Eventstore
	SystemDefaults sd.SystemDefaults
}

func StartCommandSide(config *Config) (repo *CommandSide, err error) {
	repo = &CommandSide{
		eventstore:  config.Eventstore,
		idGenerator: id.SonyFlakeGenerator,
		iamID:       config.SystemDefaults.IamID,
		iamDomain:   config.SystemDefaults.Domain,
	}
	iam_repo.RegisterEventMappers(repo.eventstore)

	repo.idpConfigSecretCrypto, err = crypto.NewAESCrypto(config.SystemDefaults.IDPConfigVerificationKey)
	if err != nil {
		return nil, err
	}
	aesCrypto, err := crypto.NewAESCrypto(config.SystemDefaults.UserVerificationKey)
	if err != nil {
		return nil, err
	}
	repo.initializeUserCode = crypto.NewEncryptionGenerator(config.SystemDefaults.SecretGenerators.InitializeUserCode, aesCrypto)
	repo.emailVerificationCode = crypto.NewEncryptionGenerator(config.SystemDefaults.SecretGenerators.EmailVerificationCode, aesCrypto)
	repo.phoneVerificationCode = crypto.NewEncryptionGenerator(config.SystemDefaults.SecretGenerators.PhoneVerificationCode, aesCrypto)
	repo.passwordVerificationCode = crypto.NewEncryptionGenerator(config.SystemDefaults.SecretGenerators.PasswordVerificationCode, aesCrypto)
	repo.userPasswordAlg = crypto.NewBCrypt(config.SystemDefaults.SecretGenerators.PasswordSaltCost)
	repo.machineKeyAlg = aesCrypto
	repo.machineKeySize = int(config.SystemDefaults.SecretGenerators.MachineKeySize)
	return repo, nil
}

func (r *CommandSide) iamByID(ctx context.Context, id string) (_ *IAMWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewIAMWriteModel(id)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	return writeModel, nil
}
