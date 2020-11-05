package model

import (
	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"

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

type WebauthNToken struct {
	es_models.ObjectRoot

	SessionID                    string
	CredentialCreationDataString string
	CredentialCreationData       *protocol.CredentialCreation
	State                        MfaState
	SessionData                  *webauthn.SessionData
	PublicKey                    []byte
	AttestationType              string
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
