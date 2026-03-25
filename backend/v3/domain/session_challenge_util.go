package domain

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type newOTPCodeFunc func(g crypto.Generator) (*crypto.CryptoValue, string, error)

type OTPRequestType int64

const (
	OTPSMSRequestType OTPRequestType = iota
	OTPEmailRequestType
)

func GetOTPCryptoGeneratorConfigWithDefault(ctx context.Context, instanceID string, opts *InvokeOpts, defaultConfig *crypto.GeneratorConfig, otpType OTPRequestType) (*crypto.GeneratorConfig, error) {
	settingsRepo := opts.secretGeneratorSettingsRepo
	cfg, err := settingsRepo.Get(
		ctx,
		opts.DB(),
		database.WithCondition(
			database.And(
				settingsRepo.InstanceIDCondition(instanceID),
				settingsRepo.TypeCondition(SettingTypeSecretGenerator),
			),
		),
	)
	if err := handleGetError(err, "DOM-x7Yd3E", "SecretGeneratorSettings"); err != nil {
		return nil, err
	}

	if cfg.State != SettingStateActive {
		return defaultConfig, nil
	}

	var attrs SecretGeneratorAttrsWithExpiry
	switch otpType {
	case OTPSMSRequestType:
		if cfg.OTPSMS == nil {
			return defaultConfig, nil
		}
		attrs = cfg.OTPSMS.SecretGeneratorAttrsWithExpiry
	case OTPEmailRequestType:
		if cfg.OTPEmail == nil {
			return defaultConfig, nil
		}
		attrs = cfg.OTPEmail.SecretGeneratorAttrsWithExpiry
	default:
		return nil, zerrors.ThrowInternal(nil, "DOM-3AcM0U", "Errors.Invalid.OTP.Type")
	}
	return &crypto.GeneratorConfig{
		Length:              *attrs.Length,
		Expiry:              *attrs.Expiry,
		IncludeLowerLetters: *attrs.IncludeLowerLetters,
		IncludeUpperLetters: *attrs.IncludeUpperLetters,
		IncludeDigits:       *attrs.IncludeDigits,
		IncludeSymbols:      *attrs.IncludeSymbols,
	}, nil
}
