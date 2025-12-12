package domain

import (
	"context"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"

	"github.com/zitadel/zitadel/internal/api/http"
	old_domain "github.com/zitadel/zitadel/internal/domain"
)

func PasskeysToCredentials(ctx context.Context, passkeys []*Passkey, rpID string) []webauthn.Credential {
	creds := make([]webauthn.Credential, 0)

	for _, pkey := range passkeys {
		// TODO(IAM-Marco): There is no state to check against, is this a problem?
		if pkey.RelyingPartyID == rpID ||
			(pkey.RelyingPartyID == "" && rpID == http.DomainContext(ctx).InstanceDomain()) {
			creds = append(creds, webauthn.Credential{
				ID:              pkey.KeyID,
				PublicKey:       pkey.PublicKey,
				AttestationType: pkey.AttestationType,
				Authenticator: webauthn.Authenticator{
					AAGUID:    pkey.AuthenticatorAttestationGUID,
					SignCount: pkey.SignCount,
				},
			})

		}
	}

	return creds
}

func UserVerificationFromDomain(verification old_domain.UserVerificationRequirement) protocol.UserVerificationRequirement {
	switch verification {
	case old_domain.UserVerificationRequirementRequired:
		return protocol.VerificationRequired
	case old_domain.UserVerificationRequirementPreferred:
		return protocol.VerificationPreferred
	case old_domain.UserVerificationRequirementDiscouraged, old_domain.UserVerificationRequirementUnspecified:
		fallthrough
	default:
		return protocol.VerificationDiscouraged
	}
}

type webAuthNUser struct {
	userID      string
	username    string
	displayName string
	creds       []webauthn.Credential
}

// WebAuthnCredentials implements [webauthn.User].
func (w *webAuthNUser) WebAuthnCredentials() []webauthn.Credential {
	return w.creds
}

// WebAuthnDisplayName implements [webauthn.User].
func (w *webAuthNUser) WebAuthnDisplayName() string {
	return w.displayName
}

// WebAuthnID implements [webauthn.User].
func (w *webAuthNUser) WebAuthnID() []byte {
	return []byte(w.userID)
}

// WebAuthnIcon implements [webauthn.User].
func (w *webAuthNUser) WebAuthnIcon() string {
	return ""
}

// WebAuthnName implements [webauthn.User].
func (w *webAuthNUser) WebAuthnName() string {
	return w.username
}

var _ webauthn.User = (*webAuthNUser)(nil)
