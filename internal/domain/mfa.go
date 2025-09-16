package domain

import (
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type MFAState int32

const (
	MFAStateUnspecified MFAState = iota
	MFAStateNotReady
	MFAStateReady
	MFAStateRemoved

	stateCount
)

type MultifactorConfigs struct {
	OTP           OTPConfig
	RecoveryCodes RecoveryCodesConfig
}

type OTPConfig struct {
	Issuer    string
	CryptoMFA crypto.EncryptionAlgorithm
}

type RecoveryCodeFormat string

const (
	RecoveryCodeFormatUUID         RecoveryCodeFormat = "uuid"
	RecoveryCodeFormatAlphanumeric RecoveryCodeFormat = "alphanumeric"
)

type RecoveryCodesConfig struct {
	MaxCount   int
	Format     RecoveryCodeFormat
	Length     int
	WithHyphen bool
}

func (f RecoveryCodeFormat) Valid() bool {
	return f == RecoveryCodeFormatUUID || f == RecoveryCodeFormatAlphanumeric
}

func (c RecoveryCodesConfig) Valid() error {
	if !c.Format.Valid() {
		return zerrors.ThrowInvalidArgument(nil, "DOMAIN-8xke2", "Errors.User.MFA.RecoveryCodes.FormatInvalid")
	}

	if c.MaxCount <= 0 || c.MaxCount > 100 {
		return zerrors.ThrowInvalidArgument(nil, "DOMAIN-9xl3w", "Errors.User.MFA.RecoveryCodes.MaxCountInvalid")
	}

	if c.Format == RecoveryCodeFormatAlphanumeric {
		if c.Length <= 0 || c.Length > 50 {
			return zerrors.ThrowInvalidArgument(nil, "DOMAIN-7xm4k", "Errors.User.MFA.RecoveryCodes.LengthInvalid")
		}
	}

	return nil
}
