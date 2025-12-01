package domain

import (
	"context"

	"github.com/go-webauthn/webauthn/webauthn"

	"github.com/zitadel/zitadel/internal/api/http"
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
