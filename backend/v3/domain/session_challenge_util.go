package domain

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type newOTPCodeFunc func(g crypto.Generator) (*crypto.CryptoValue, string, error)

type OTPType int64

const (
	OTPTypeSMS OTPType = iota
	OTPTypeEmail
)

func GetOTPCryptoGeneratorConfigWithDefault(ctx context.Context, instanceID string, opts *InvokeOpts, defaultConfig *crypto.GeneratorConfig, otpType OTPType) (*crypto.GeneratorConfig, error) {
	if defaultConfig == nil {
		return nil, zerrors.ThrowInternal(nil, "DOM-3AcM0U", "missing default config")
	}

	if opts == nil || opts.secretGeneratorSettingsRepo == nil {
		return defaultConfig, nil
	}

	settingsRepo := opts.secretGeneratorSettingsRepo
	cfg, err := settingsRepo.Get(
		ctx,
		opts.DB(),
		database.WithCondition(
			settingsRepo.UniqueCondition(instanceID, nil, SettingTypeSecretGenerator, SettingStateActive),
		),
	)
	if err := handleGetError(err, "DOM-x7Yd3E", "SecretGeneratorSettings"); err != nil {
		return nil, err // todo: or return defaultConfig?
	}

	if cfg.State != SettingStateActive {
		return defaultConfig, nil
	}

	var attrs SecretGeneratorAttrsWithExpiry
	switch otpType {
	case OTPTypeSMS:
		if cfg.OTPSMS == nil {
			return defaultConfig, nil
		}
		attrs = cfg.OTPSMS.SecretGeneratorAttrsWithExpiry
	case OTPTypeEmail:
		if cfg.OTPEmail == nil {
			return defaultConfig, nil
		}
		attrs = cfg.OTPEmail.SecretGeneratorAttrsWithExpiry
	default:
		return nil, zerrors.ThrowInternal(nil, "DOM-3AcM0U", "invalid otp type")
	}
	return &crypto.GeneratorConfig{
		Length:              getValueOrDefault(attrs.Length, defaultConfig.Length),
		Expiry:              getValueOrDefault(attrs.Expiry, defaultConfig.Expiry),
		IncludeLowerLetters: getValueOrDefault(attrs.IncludeLowerLetters, defaultConfig.IncludeLowerLetters),
		IncludeUpperLetters: getValueOrDefault(attrs.IncludeUpperLetters, defaultConfig.IncludeUpperLetters),
		IncludeDigits:       getValueOrDefault(attrs.IncludeDigits, defaultConfig.IncludeDigits),
		IncludeSymbols:      getValueOrDefault(attrs.IncludeSymbols, defaultConfig.IncludeSymbols),
	}, nil
}

// getValueOrDefault safely dereferences a pointer or returns a default value
func getValueOrDefault[T any](ptr *T, defaultVal T) T {
	if ptr != nil {
		return *ptr
	}
	return defaultVal
}
