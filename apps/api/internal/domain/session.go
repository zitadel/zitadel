package domain

import (
	"io"

	"golang.org/x/text/language"
)

type SessionState int32

const (
	SessionStateUnspecified SessionState = iota
	SessionStateActive
	SessionStateTerminated
)

type OTPEmailURLData struct {
	Code              string
	UserID            string
	LoginName         string
	DisplayName       string
	PreferredLanguage language.Tag
	SessionID         string
}

// RenderOTPEmailURLTemplate parses and renders tmpl.
// code, userID, (preferred) loginName, displayName and preferredLanguage are passed into the [OTPEmailURLData].
func RenderOTPEmailURLTemplate(w io.Writer, tmpl, code, userID, loginName, displayName, sessionID string, preferredLanguage language.Tag) error {
	return renderURLTemplate(w, tmpl, &OTPEmailURLData{
		Code:              code,
		UserID:            userID,
		LoginName:         loginName,
		DisplayName:       displayName,
		PreferredLanguage: preferredLanguage,
		SessionID:         sessionID,
	})
}
