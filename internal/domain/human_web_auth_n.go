package domain

import (
	"bytes"
	"fmt"
	"time"

	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
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
	RPID                   string
}

type WebAuthNLogin struct {
	es_models.ObjectRoot

	CredentialAssertionData []byte
	Challenge               string
	AllowedCredentialIDs    [][]byte
	UserVerification        UserVerificationRequirement
	RPID                    string
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

func GetTokenByKeyID(tokens []*WebAuthNToken, keyID []byte) *WebAuthNToken {
	for _, token := range tokens {
		if bytes.Equal(token.KeyID, keyID) {
			return token
		}
	}
	return nil
}

type PasswordlessInitCodeState int32

const (
	PasswordlessInitCodeStateUnspecified PasswordlessInitCodeState = iota
	PasswordlessInitCodeStateRequested
	PasswordlessInitCodeStateActive
	PasswordlessInitCodeStateRemoved
)

type PasswordlessInitCode struct {
	es_models.ObjectRoot

	CodeID     string
	Code       string
	Expiration time.Duration
	State      PasswordlessInitCodeState
}

func (p *PasswordlessInitCode) Link(baseURL string) string {
	return PasswordlessInitCodeLink(baseURL, p.AggregateID, p.ResourceOwner, p.CodeID, p.Code)
}

func PasswordlessInitCodeLink(baseURL, userID, resourceOwner, codeID, code string) string {
	return fmt.Sprintf("%s?userID=%s&orgID=%s&codeID=%s&code=%s", baseURL, userID, resourceOwner, codeID, code)
}

func PasswordlessInitCodeLinkTemplate(baseURL, userID, resourceOwner, codeID string) string {
	return PasswordlessInitCodeLink(baseURL, userID, resourceOwner, codeID, "{{.Code}}")
}
