package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type cryptoCodeFunc func(ctx context.Context, filter preparation.FilterToQueryReducer, typ domain.SecretGeneratorType, alg crypto.Crypto) (*CryptoCode, error)

type cryptoCodeWithDefaultFunc func(ctx context.Context, filter preparation.FilterToQueryReducer, typ domain.SecretGeneratorType, alg crypto.Crypto, defaultConfig *crypto.GeneratorConfig) (*CryptoCode, error)

var emptyConfig = &crypto.GeneratorConfig{}

type CryptoCode struct {
	Crypted *crypto.CryptoValue
	Plain   string
	Expiry  time.Duration
}

func newCryptoCode(ctx context.Context, filter preparation.FilterToQueryReducer, typ domain.SecretGeneratorType, alg crypto.Crypto) (*CryptoCode, error) {
	return newCryptoCodeWithDefaultConfig(ctx, filter, typ, alg, emptyConfig)
}

func newCryptoCodeWithDefaultConfig(ctx context.Context, filter preparation.FilterToQueryReducer, typ domain.SecretGeneratorType, alg crypto.Crypto, defaultConfig *crypto.GeneratorConfig) (*CryptoCode, error) {
	gen, config, err := secretGenerator(ctx, filter, typ, alg, defaultConfig)
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

func verifyCryptoCode(ctx context.Context, filter preparation.FilterToQueryReducer, typ domain.SecretGeneratorType, alg crypto.Crypto, creation time.Time, expiry time.Duration, crypted *crypto.CryptoValue, plain string) error {
	gen, _, err := secretGenerator(ctx, filter, typ, alg, emptyConfig)
	if err != nil {
		return err
	}
	return crypto.VerifyCode(creation, expiry, crypted, plain, gen)
}

func secretGenerator(ctx context.Context, filter preparation.FilterToQueryReducer, typ domain.SecretGeneratorType, alg crypto.Crypto, defaultConfig *crypto.GeneratorConfig) (crypto.Generator, *crypto.GeneratorConfig, error) {
	config, err := secretGeneratorConfigWithDefault(ctx, filter, typ, defaultConfig)
	if err != nil {
		return nil, nil, err
	}
	switch a := alg.(type) {
	case crypto.HashAlgorithm:
		return crypto.NewHashGenerator(*config, a), config, nil
	case crypto.EncryptionAlgorithm:
		return crypto.NewEncryptionGenerator(*config, a), config, nil
	default:
		return nil, nil, zerrors.ThrowInternalf(nil, "COMMA-RreV6", "Errors.Internal unsupported crypto algorithm type %T", a)
	}
}

func secretGeneratorConfig(ctx context.Context, filter preparation.FilterToQueryReducer, typ domain.SecretGeneratorType) (*crypto.GeneratorConfig, error) {
	return secretGeneratorConfigWithDefault(ctx, filter, typ, emptyConfig)
}

func secretGeneratorConfigWithDefault(ctx context.Context, filter preparation.FilterToQueryReducer, typ domain.SecretGeneratorType, defaultConfig *crypto.GeneratorConfig) (*crypto.GeneratorConfig, error) {
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
