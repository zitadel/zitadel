package messages

import (
	"net/textproto"
	"strings"

	"golang.org/x/net/http/httpguts"
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

// IsValidEmailHeaderName reports whether the name is a valid header field name.
// It is restricted to the token characters of RFC 7230 (a subset of RFC 5322),
// as only these names are canonicalized by [textproto.CanonicalMIMEHeaderKey].
func IsValidEmailHeaderName(name string) bool {
	return httpguts.ValidHeaderFieldName(name)
}

// IsValidEmailHeaderValue reports whether the value can safely be written into a header.
// Line breaks are not allowed, as they would allow to inject additional headers.
func IsValidEmailHeaderValue(value string) bool {
	return !strings.ContainsAny(value, "\r\n")
}
