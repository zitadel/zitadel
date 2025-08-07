package webauthn

import (
	"context"
	"strings"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"

	"github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/domain"
)

func WebAuthNsToCredentials(ctx context.Context, webAuthNs []*domain.WebAuthNToken, rpID string) []webauthn.Credential {
	creds := make([]webauthn.Credential, 0)
	for _, webAuthN := range webAuthNs {
		// only add credentials that are ready and
		// either match the rpID or
		// if they were added through Console / old login UI, there is no stored rpID set;
		// then we check if the requested rpID matches the instance domain
		if webAuthN.State == domain.MFAStateReady &&
			(webAuthN.RPID == rpID ||
				(webAuthN.RPID == "" && rpID == strings.Split(http.DomainContext(ctx).InstanceHost, ":")[0])) {
			creds = append(creds, webauthn.Credential{
				ID:              webAuthN.KeyID,
				PublicKey:       webAuthN.PublicKey,
				AttestationType: webAuthN.AttestationType,
				Authenticator: webauthn.Authenticator{
					AAGUID:    webAuthN.AAGUID,
					SignCount: webAuthN.SignCount,
				},
			})
		}
	}
	return creds
}

func WebAuthNToSessionData(webAuthN *domain.WebAuthNToken) webauthn.SessionData {
	return webauthn.SessionData{
		Challenge:            webAuthN.Challenge,
		UserID:               []byte(webAuthN.AggregateID),
		AllowedCredentialIDs: webAuthN.AllowedCredentialIDs,
		UserVerification:     UserVerificationFromDomain(webAuthN.UserVerification),
	}
}

func WebAuthNLoginToSessionData(webAuthN *domain.WebAuthNLogin) webauthn.SessionData {
	return webauthn.SessionData{
		Challenge:            webAuthN.Challenge,
		UserID:               []byte(webAuthN.AggregateID),
		AllowedCredentialIDs: webAuthN.AllowedCredentialIDs,
		UserVerification:     UserVerificationFromDomain(webAuthN.UserVerification),
	}
}

func UserVerificationToDomain(verification protocol.UserVerificationRequirement) domain.UserVerificationRequirement {
	switch verification {
	case protocol.VerificationRequired:
		return domain.UserVerificationRequirementRequired
	case protocol.VerificationPreferred:
		return domain.UserVerificationRequirementPreferred
	case protocol.VerificationDiscouraged:
		return domain.UserVerificationRequirementDiscouraged
	default:
		return domain.UserVerificationRequirementUnspecified
	}
}

func UserVerificationFromDomain(verification domain.UserVerificationRequirement) protocol.UserVerificationRequirement {
	switch verification {
	case domain.UserVerificationRequirementRequired:
		return protocol.VerificationRequired
	case domain.UserVerificationRequirementPreferred:
		return protocol.VerificationPreferred
	case domain.UserVerificationRequirementDiscouraged:
		return protocol.VerificationDiscouraged
	default:
		return protocol.VerificationDiscouraged
	}
}

func AuthenticatorAttachmentFromDomain(authType domain.AuthenticatorAttachment) protocol.AuthenticatorAttachment {
	switch authType {
	case domain.AuthenticatorAttachmentPlattform:
		return protocol.Platform
	case domain.AuthenticatorAttachmentCrossPlattform:
		return protocol.CrossPlatform
	default:
		return ""
	}
}
