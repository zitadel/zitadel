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

type RecoveryCodesConfig struct {
	Count int
}
