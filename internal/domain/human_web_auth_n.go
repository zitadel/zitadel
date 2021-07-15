package domain

import (
	"bytes"
	"time"

	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
)

type WebAuthNToken struct {
	es_models.ObjectRoot

	WebAuthNTokenID        string
	CredentialCreationData []byte
	State                  MFAState
	Challenge              string
	AllowedCredentialIDs   [][]byte
	UserVerification       UserVerificationRequirement
	KeyID                  []byte
	PublicKey              []byte
	AttestationType        string
	AAGUID                 []byte
	SignCount              uint32
	WebAuthNTokenName      string
}

type WebAuthNLogin struct {
	es_models.ObjectRoot

	CredentialAssertionData []byte
	Challenge               string
	AllowedCredentialIDs    [][]byte
	UserVerification        UserVerificationRequirement
}

type UserVerificationRequirement int32

const (
	UserVerificationRequirementUnspecified UserVerificationRequirement = iota
	UserVerificationRequirementRequired
	UserVerificationRequirementPreferred
	UserVerificationRequirementDiscouraged
)

type AuthenticatorAttachment int32

const (
	AuthenticatorAttachmentUnspecified AuthenticatorAttachment = iota
	AuthenticatorAttachmentPlattform
	AuthenticatorAttachmentCrossPlattform
)

func GetTokenToVerify(tokens []*WebAuthNToken) (int, *WebAuthNToken) {
	for i, u2f := range tokens {
		if u2f.State == MFAStateNotReady {
			return i, u2f
		}
	}
	return -1, nil
}

func GetTokenByKeyID(tokens []*WebAuthNToken, keyID []byte) (int, *WebAuthNToken) {
	for i, token := range tokens {
		if bytes.Compare(token.KeyID, keyID) == 0 {
			return i, token
		}
	}
	return -1, nil
}

type PasswordlessInitCodeState int32

const (
	PasswordlessInitCodeStateUnspecified PasswordlessInitCodeState = iota
	PasswordlessInitCodeStateActive
	PasswordlessInitCodeStateRemoved
)

type PasswordlessInitCode struct {
	es_models.ObjectRoot

	CodeID     string
	Code       string
	Expiration time.Duration
	Active     bool
}

//func NewPasswordlessInitCode(generator crypto.Generator) (*PasswordlessInitCode, string, error) {
//	initCodeCrypto, code, err := crypto.NewCode(generator)
//	if err != nil {
//		return nil, "", err
//	}
//	return &PasswordlessInitCode{
//		Code:   initCodeCrypto,
//		Expiry: generator.Expiry(),
//	}, code, nil
//}
