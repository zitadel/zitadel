package webauthn

import (
	"github.com/caos/zitadel/internal/user/model"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
)

func WebAuthNsToCredentials(webAuthNs []*model.WebAuthNToken) []webauthn.Credential {
	creds := make([]webauthn.Credential, 0)
	for _, webAuthN := range webAuthNs {
		if webAuthN.State == model.MFAStateReady {
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
		UserVerification:     UserVerificationFromModel(webAuthN.UserVerification),
	}
}

func WebAuthNLoginToSessionData(webAuthN *model.WebAuthNLogin) webauthn.SessionData {
	return webauthn.SessionData{
		Challenge:            webAuthN.Challenge,
		UserID:               []byte(webAuthN.AggregateID),
		AllowedCredentialIDs: webAuthN.AllowedCredentialIDs,
		UserVerification:     UserVerificationFromModel(webAuthN.UserVerification),
	}
}

func UserVerificationToModel(verification protocol.UserVerificationRequirement) model.UserVerificationRequirement {
	switch verification {
	case protocol.VerificationRequired:
		return model.UserVerificationRequirementRequired
	case protocol.VerificationPreferred:
		return model.UserVerificationRequirementPreferred
	case protocol.VerificationDiscouraged:
		return model.UserVerificationRequirementDiscouraged
	default:
		return model.UserVerificationRequirementUnspecified
	}
}

func UserVerificationFromModel(verification model.UserVerificationRequirement) protocol.UserVerificationRequirement {
	switch verification {
	case model.UserVerificationRequirementRequired:
		return protocol.VerificationRequired
	case model.UserVerificationRequirementPreferred:
		return protocol.VerificationPreferred
	case model.UserVerificationRequirementDiscouraged:
		return protocol.VerificationDiscouraged
	default:
		return protocol.VerificationDiscouraged
	}
}

func AuthenticatorAttachmentFromModel(authType model.AuthenticatorAttachment) protocol.AuthenticatorAttachment {
	switch authType {
	case model.AuthenticatorAttachmentPlattform:
		return protocol.Platform
	case model.AuthenticatorAttachmentCrossPlattform:
		return protocol.CrossPlatform
	default:
		return ""
	}
}
