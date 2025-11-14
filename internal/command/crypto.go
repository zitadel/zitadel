package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
)

type encrypedCodeFunc func(ctx context.Context, filter preparation.FilterToQueryReducer, typ domain.SecretGeneratorType, alg crypto.EncryptionAlgorithm) (*EncryptedCode, error)

type encryptedCodeWithDefaultFunc func(ctx context.Context, filter preparation.FilterToQueryReducer, typ domain.SecretGeneratorType, alg crypto.EncryptionAlgorithm, defaultConfig *crypto.GeneratorConfig) (*EncryptedCode, error)

type encryptedCodeGeneratorWithDefaultFunc func(ctx context.Context, filter preparation.FilterToQueryReducer, typ domain.SecretGeneratorType, alg crypto.EncryptionAlgorithm, defaultConfig *crypto.GeneratorConfig) (*EncryptedCode, string, error)

var emptyConfig = &crypto.GeneratorConfig{}

type EncryptedCode struct {
	Crypted *crypto.CryptoValue
	Plain   string
	Expiry  time.Duration
}

func (e *EncryptedCode) CryptedCode() *crypto.CryptoValue {
	if e == nil {
		return nil
	}
	return e.Crypted
}

func (e *EncryptedCode) PlainCode() string {
	if e == nil {
		return ""
	}
	return e.Plain
}

func (e *EncryptedCode) CodeExpiry() time.Duration {
	if e == nil {
		return 0
	}
	return e.Expiry
}

func newEncryptedCode(ctx context.Context, filter preparation.FilterToQueryReducer, typ domain.SecretGeneratorType, alg crypto.EncryptionAlgorithm) (*EncryptedCode, error) {
	return newEncryptedCodeWithDefaultConfig(ctx, filter, typ, alg, emptyConfig)
}

func newEncryptedCodeWithDefaultConfig(ctx context.Context, filter preparation.FilterToQueryReducer, typ domain.SecretGeneratorType, alg crypto.EncryptionAlgorithm, defaultConfig *crypto.GeneratorConfig) (*EncryptedCode, error) {
	gen, config, err := encryptedCodeGenerator(ctx, filter, typ, alg, defaultConfig)
	if err != nil {
		return nil, err
	}
	crypted, plain, err := crypto.NewCode(gen)
	if err != nil {
		return nil, err
	}
	return &EncryptedCode{
		Crypted: crypted,
		Plain:   plain,
		Expiry:  config.Expiry,
	}, nil
}

func encryptedCodeGenerator(ctx context.Context, filter preparation.FilterToQueryReducer, typ domain.SecretGeneratorType, alg crypto.EncryptionAlgorithm, defaultConfig *crypto.GeneratorConfig) (crypto.Generator, *crypto.GeneratorConfig, error) {
	config, err := cryptoGeneratorConfigWithDefault(ctx, filter, typ, defaultConfig)
	if err != nil {
		return nil, nil, err
	}
	return crypto.NewEncryptionGenerator(*config, alg), config, nil
}

type hashedSecretFunc func(ctx context.Context, filter preparation.FilterToQueryReducer) (encodedHash, plain string, err error)

func newHashedSecretWithDefault(hasher *crypto.Hasher, defaultConfig *crypto.GeneratorConfig) hashedSecretFunc {
	return func(ctx context.Context, filter preparation.FilterToQueryReducer) (encodedHash string, plain string, err error) {
		config, err := cryptoGeneratorConfigWithDefault(ctx, filter, domain.SecretGeneratorTypeAppSecret, defaultConfig)
		if err != nil {
			return "", "", err
		}
		generator := crypto.NewHashGenerator(*config, hasher)
		return generator.NewCode()
	}
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
