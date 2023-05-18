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

	PasskeyID                          string
	PublicKeyCredentialCreationOptions []byte
}
