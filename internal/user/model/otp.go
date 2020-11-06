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
	MfaStateUnspecified MfaState = iota
	MfaStateNotReady
	MfaStateReady
)

type MultiFactor struct {
	Type  MfaType
	State MfaState
}

type MfaType int32

const (
	MfaTypeUnspecified MfaType = iota
	MfaTypeOTP
	MfaTypeSMS
)
