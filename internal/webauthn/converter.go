package webauthn

import (
	"github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
)

func WebAuthNsToCredentials(webAuthNs []*domain.WebAuthNToken) []webauthn.Credential {
	creds := make([]webauthn.Credential, 0)
	for _, webAuthN := range webAuthNs {
		if webAuthN.State == domain.MFAStateReady {
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

func WebAuthNToSessionData(webAuthN *model.WebAuthNToken) webauthn.SessionData {
	return webauthn.SessionData{
		Challenge:            webAuthN.Challenge,
		UserID:               []byte(webAuthN.AggregateID),
		AllowedCredentialIDs: webAuthN.AllowedCredentialIDs,
		UserVerification:     UserVerificationFromDomain(webAuthN.UserVerification),
	}
}

func WebAuthNLoginToSessionData(webAuthN *model.WebAuthNLogin) webauthn.SessionData {
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
