package webauthn

import (
	"github.com/caos/zitadel/internal/user/model"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
)

func WebAuthNsToCredentials(webAuthNs []*model.WebAuthNToken) []webauthn.Credential {
	creds := make([]webauthn.Credential, len(webAuthNs))
	for i, webAuthN := range webAuthNs {
		creds[i] = webauthn.Credential{
			ID:              webAuthN.KeyID,
			PublicKey:       webAuthN.PublicKey,
			AttestationType: webAuthN.AttestationType,
			Authenticator: webauthn.Authenticator{
				AAGUID:    webAuthN.AAGUID,
				SignCount: webAuthN.SignCount,
			},
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
