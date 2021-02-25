package domain

import "github.com/caos/zitadel/internal/crypto"

type MFAState int32

const (
	MFAStateUnspecified MFAState = iota
	MFAStateNotReady
	MFAStateReady
	MFAStateRemoved

	stateCount
)

func (f MFAState) Valid() bool {
	return f >= 0 && f < stateCount
}

type MultifactorConfigs struct {
	OTP OTPConfig
}

type OTPConfig struct {
	Issuer    string
	CryptoMFA crypto.EncryptionAlgorithm
}
