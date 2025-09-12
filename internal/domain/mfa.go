package domain

import (
	"github.com/zitadel/zitadel/internal/crypto"
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
	RecoveryCodeFormatUUID      RecoveryCodeFormat = "uuid"
	RecoveryCodeFormatSonyFlake RecoveryCodeFormat = "sonyflake"
)

type RecoveryCodesConfig struct {
	MaxCount int
	Format   RecoveryCodeFormat
}

func (f RecoveryCodeFormat) Valid() bool {
	return f == RecoveryCodeFormatUUID || f == RecoveryCodeFormatSonyFlake
}
