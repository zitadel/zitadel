package domain

import "io"

type PasskeyURLData struct {
	UserID string
	OrgID  string
	CodeID string
	Code   string
}

// RenderPasskeyURLTemplate parses and renders tmpl.
// userID, orgID, codeID and code are passed into the [PasskeyURLData].
func RenderPasskeyURLTemplate(w io.Writer, tmpl, userID, orgID, codeID, code string) error {
	return renderURLTemplate(w, tmpl, &PasskeyURLData{userID, orgID, codeID, code})
}

type PasskeyCodeDetails struct {
	*ObjectDetails
	CodeID string
	Code   string
}

type WebAuthNRegistrationDetails struct {
	*ObjectDetails

	ID                                 string
	PublicKeyCredentialCreationOptions []byte
}
