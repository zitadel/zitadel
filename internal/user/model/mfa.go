package model

import (
	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
)

type MfaType interface {
	MfaState() MfaState
	MfaLevel() MfaLevel
}

func MfaIsReady(t MfaType) bool {
	return t.MfaState() == MFASTATE_READY
}

func MfaLevelSufficient(t MfaType, level MfaLevel) bool {
	return t.MfaLevel() >= level
}

type MfaLevel int

const (
	MfaLevelSoftware MfaLevel = iota
	MfaLevelHardware
	MfaLevelHardwareCertified
)

type OTP struct {
	es_models.ObjectRoot

	Secret       *crypto.CryptoValue
	SecretString string
	Url          string
	State        MfaState
}

func (o *OTP) MfaState() MfaState {
	return o.State
}
func (o *OTP) MfaLevel() MfaLevel {
	return MfaLevelSoftware
}

type MfaState int32

const (
	MFASTATE_UNSPECIFIED MfaState = iota
	MFASTATE_NOTREADY
	MFASTATE_READY
)
