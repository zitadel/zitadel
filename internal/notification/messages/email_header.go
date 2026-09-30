package messages

import (
	"net/textproto"
	"strings"
)

// reservedEmailHeaders are set by [Email.GetContent] itself
// and must not be overwritten by additional headers.
var reservedEmailHeaders = map[string]struct{}{
	"From":         {},
	"To":           {},
	"Cc":           {},
	"Bcc":          {},
	"Reply-To":     {},
	"Return-Path":  {},
	"Date":         {},
	"Subject":      {},
	"Mime-Version": {},
	"Content-Type": {},
}

// IsReservedEmailHeader reports whether the header is set by ZITADEL itself.
// The check is case-insensitive.
func IsReservedEmailHeader(name string) bool {
	_, ok := reservedEmailHeaders[textproto.CanonicalMIMEHeaderKey(name)]
	return ok
}

// IsValidEmailHeaderName reports whether the name is a valid header field name
// according to RFC 5322: printable US-ASCII characters except the colon.
func IsValidEmailHeaderName(name string) bool {
	if name == "" {
		return false
	}
	for _, r := range name {
		if r < 33 || r > 126 || r == ':' {
			return false
		}
	}
	return true
}

// IsValidEmailHeaderValue reports whether the value can safely be written into a header.
// Line breaks are not allowed, as they would allow to inject additional headers.
func IsValidEmailHeaderValue(value string) bool {
	return !strings.ContainsAny(value, "\r\n")
}
