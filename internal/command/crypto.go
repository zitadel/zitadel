package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
)

type cryptoCodeFunc func(ctx context.Context, filter preparation.FilterToQueryReducer, typ domain.SecretGeneratorType, alg crypto.EncryptionAlgorithm) (*CryptoCode, error)

type cryptoCodeWithDefaultFunc func(ctx context.Context, filter preparation.FilterToQueryReducer, typ domain.SecretGeneratorType, alg crypto.EncryptionAlgorithm, defaultConfig *crypto.GeneratorConfig) (*CryptoCode, error)

var emptyConfig = &crypto.GeneratorConfig{}

type CryptoCode struct {
	Crypted *crypto.CryptoValue
	Plain   string
	Expiry  time.Duration
}

func newCryptoCode(ctx context.Context, filter preparation.FilterToQueryReducer, typ domain.SecretGeneratorType, alg crypto.EncryptionAlgorithm) (*CryptoCode, error) {
	return newCryptoCodeWithDefaultConfig(ctx, filter, typ, alg, emptyConfig)
}

func newCryptoCodeWithDefaultConfig(ctx context.Context, filter preparation.FilterToQueryReducer, typ domain.SecretGeneratorType, alg crypto.EncryptionAlgorithm, defaultConfig *crypto.GeneratorConfig) (*CryptoCode, error) {
	gen, config, err := cryptoCodeGenerator(ctx, filter, typ, alg, defaultConfig)
	if err != nil {
		return nil, err
	}
	crypted, plain, err := crypto.NewCode(gen)
	if err != nil {
		return nil, err
	}
	return &CryptoCode{
		Crypted: crypted,
		Plain:   plain,
		Expiry:  config.Expiry,
	}, nil
}

func verifyCryptoCode(ctx context.Context, filter preparation.FilterToQueryReducer, typ domain.SecretGeneratorType, alg crypto.EncryptionAlgorithm, creation time.Time, expiry time.Duration, crypted *crypto.CryptoValue, plain string) error {
	gen, _, err := cryptoCodeGenerator(ctx, filter, typ, alg, emptyConfig)
	if err != nil {
		return err
	}
	return crypto.VerifyCode(creation, expiry, crypted, plain, gen.Alg())
}

func cryptoCodeGenerator(ctx context.Context, filter preparation.FilterToQueryReducer, typ domain.SecretGeneratorType, alg crypto.EncryptionAlgorithm, defaultConfig *crypto.GeneratorConfig) (crypto.Generator, *crypto.GeneratorConfig, error) {
	config, err := cryptoGeneratorConfigWithDefault(ctx, filter, typ, defaultConfig)
	if err != nil {
		return nil, nil, err
	}
	return crypto.NewEncryptionGenerator(*config, alg), config, nil
}

func newHashedSecret(ctx context.Context, filter preparation.FilterToQueryReducer, hasher *crypto.PasswordHasher) (encodedHash, plain string, err error) {
	return newHashedSecretWithDefaultConfig(ctx, filter, hasher, emptyConfig)
}

func newHashedSecretWithDefaultConfig(ctx context.Context, filter preparation.FilterToQueryReducer, hasher *crypto.PasswordHasher, defaultConfig *crypto.GeneratorConfig) (encodedHash, plain string, err error) {
	generator, err := hashedSecretGeneratorWithDefaultConfig(ctx, filter, hasher, defaultConfig)
	if err != nil {
		return "", "", err
	}
	return generator.NewCode()
}

func hashedSecretGenerator(ctx context.Context, filter preparation.FilterToQueryReducer, hasher *crypto.PasswordHasher) (*crypto.HashGenerator, error) {
	return hashedSecretGeneratorWithDefaultConfig(ctx, filter, hasher, emptyConfig)
}

func hashedSecretGeneratorWithDefaultConfig(ctx context.Context, filter preparation.FilterToQueryReducer, hasher *crypto.PasswordHasher, defaultConfig *crypto.GeneratorConfig) (*crypto.HashGenerator, error) {
	config, err := cryptoGeneratorConfigWithDefault(ctx, filter, domain.SecretGeneratorTypeAppSecret, defaultConfig)
	if err != nil {
		return nil, err
	}
	return crypto.NewHashGenerator(*config, hasher), nil
}

func cryptoGeneratorConfig(ctx context.Context, filter preparation.FilterToQueryReducer, typ domain.SecretGeneratorType) (*crypto.GeneratorConfig, error) {
	return cryptoGeneratorConfigWithDefault(ctx, filter, typ, emptyConfig)
}

func cryptoGeneratorConfigWithDefault(ctx context.Context, filter preparation.FilterToQueryReducer, typ domain.SecretGeneratorType, defaultConfig *crypto.GeneratorConfig) (*crypto.GeneratorConfig, error) {
	wm := NewInstanceSecretGeneratorConfigWriteModel(ctx, typ)
	events, err := filter(ctx, wm.Query())
	if err != nil {
		return nil, err
	}
	wm.AppendEvents(events...)
	if err := wm.Reduce(); err != nil {
		return nil, err
	}
	if wm.State != domain.SecretGeneratorStateActive {
		return defaultConfig, nil
	}
	return &crypto.GeneratorConfig{
		Length:              wm.Length,
		Expiry:              wm.Expiry,
		IncludeLowerLetters: wm.IncludeLowerLetters,
		IncludeUpperLetters: wm.IncludeUpperLetters,
		IncludeDigits:       wm.IncludeDigits,
		IncludeSymbols:      wm.IncludeSymbols,
	}, nil
}
