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
	State        MFAState
}

type MFAState int32

const (
	MFAStateUnspecified MFAState = iota
	MFAStateNotReady
	MFAStateReady
)

type MultiFactor struct {
	Type      MFAType
	State     MFAState
	Attribute string
	ID        string
}

type MFAType int32

const (
	MFATypeUnspecified MFAType = iota
	MFATypeOTP
	MFATypeU2F
)
