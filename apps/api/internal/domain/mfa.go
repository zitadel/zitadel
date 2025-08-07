package domain

import "github.com/zitadel/zitadel/internal/crypto"

type MFAState int32

const (
	MFAStateUnspecified MFAState = iota
	MFAStateNotReady
	MFAStateReady
	MFAStateRemoved

	stateCount
)

type MultifactorConfigs struct {
	OTP OTPConfig
}

type OTPConfig struct {
	Issuer    string
	CryptoMFA crypto.EncryptionAlgorithm
}
