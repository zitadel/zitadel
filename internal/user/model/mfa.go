package model

import (
	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
)

type OTP struct {
	es_models.ObjectRoot

	Secret       *crypto.CryptoValue
	SecretString string
	Url          string
	State        MfaState
}

type MfaState int32

const (
	MFASTATE_UNSPECIFIED MfaState = iota
	MFASTATE_NOTREADY
	MFASTATE_READY
)

type MultiFactor struct {
	Type  MFAType
	State MfaState
}

type MFAType int32

const (
	MFATYPE_UNSPECIFIED MFAType = iota
	MFATYPE_OTP
	MFATYPE_SMS
)
