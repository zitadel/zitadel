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
	RecoveryCodeFormatUUID         RecoveryCodeFormat = "uuid"
	RecoveryCodeFormatAlphanumeric RecoveryCodeFormat = "alphanumeric"
)

type RecoveryCodesConfig struct {
	MaxCount   int
	Format     RecoveryCodeFormat
	Length     int
	WithHyphen bool
}
