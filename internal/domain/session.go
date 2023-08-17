package domain

import "io"

type SessionState int32

const (
	SessionStateUnspecified SessionState = iota
	SessionStateActive
	SessionStateTerminated
)

type OTPEmailURLData struct {
	UserID string
	Code   string
}

// RenderOTPEmailURLTemplate parses and renders tmpl.
// userID, orgID, codeID and code are passed into the [OTPEmailURLData].
func RenderOTPEmailURLTemplate(w io.Writer, tmpl, userID, code string) error {
	return renderURLTemplate(w, tmpl, &OTPEmailURLData{userID, code})
}
