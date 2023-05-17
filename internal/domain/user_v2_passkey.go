package domain

import "io"

type PasskeyURLData struct {
	UserID        string
	ResourceOwner string
	CodeID        string
	Code          string
}

// RenderPasskeyURLTemplate parses and renders tmpl.
// userID, resourceOwner, codeID and code are passed into the [PasskeyURLData].
func RenderPasskeyURLTemplate(w io.Writer, tmpl, userID, resourceOwner, codeID, code string) error {
	return renderURLTemplate(w, tmpl, &PasskeyURLData{userID, resourceOwner, codeID, code})
}

type PasskeyCodeDetails struct {
	*ObjectDetails
	CodeID string
	Code   string
}

type PasskeyRegistrationDetails struct {
	*ObjectDetails

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

func (w *WebAuthNToken) PasskeyRegistrationDetails(details *ObjectDetails) *PasskeyRegistrationDetails {
	return &PasskeyRegistrationDetails{
		ObjectDetails:          details,
		WebAuthNTokenID:        w.WebAuthNTokenID,
		CredentialCreationData: w.CredentialCreationData,
		State:                  w.State,
		Challenge:              w.Challenge,
		AllowedCredentialIDs:   w.AllowedCredentialIDs,
		UserVerification:       w.UserVerification,
		KeyID:                  w.KeyID,
		PublicKey:              w.PublicKey,
		AttestationType:        w.AttestationType,
		AAGUID:                 w.AAGUID,
		SignCount:              w.SignCount,
		WebAuthNTokenName:      w.WebAuthNTokenName,
	}
}
