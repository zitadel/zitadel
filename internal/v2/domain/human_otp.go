package domain

import (
	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
)

type OTP struct {
	es_models.ObjectRoot

	Secret       *crypto.CryptoValue
	SecretString string
	Url          string
	State        MFAState
}

type OTPState int32

const (
	OTPStateUnspecified OTPState = iota
	OTPStateActive
	OTPStateRemoved

	otpStateCount
)

func (s OTPState) Valid() bool {
	return s >= 0 && s < otpStateCount
}
