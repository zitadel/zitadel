package domain

import (
	"net/url"
	"strings"
)

// MaxClientIDMetadataDocumentURLLength bounds a client_id URL, both as an allowlist entry and
// as a client_id resolved at runtime, so an unauthenticated caller cannot make ZITADEL fetch or
// cache arbitrarily long URLs under an allowed prefix.
const MaxClientIDMetadataDocumentURLLength = 2048

// IsClientIDMetadataDocumentURL reports whether clientID is a Client Identifier URL as defined
// in section 3 of draft-ietf-oauth-client-id-metadata-document: an https URL with a hostname
// (not only a port) and a path, without userinfo, fragment or dot path segments, and at most
// MaxClientIDMetadataDocumentURLLength long. A regular ZITADEL client_id is a
// numeric snowflake (optionally suffixed) and is never such a URL.
//
// Percent-encoded dots, slashes and backslashes in the path are rejected too. They are not dot
// segments as written, but a server that decodes them before resolving the path would treat
// them as one, which would let a client_id escape the path of an allowed prefix.
func IsClientIDMetadataDocumentURL(clientID string) bool {
	if len(clientID) > MaxClientIDMetadataDocumentURLLength {
		return false
	}
	u, err := url.Parse(clientID)
	if err != nil || !u.IsAbs() || !strings.EqualFold(u.Scheme, "https") || u.Hostname() == "" {
		return false
	}
	if u.User != nil || u.Fragment != "" || strings.Contains(clientID, "#") {
		return false
	}
	path := u.EscapedPath()
	if path == "" {
		return false
	}
	lowerPath := strings.ToLower(path)
	if strings.Contains(lowerPath, "%2f") || strings.Contains(lowerPath, "%5c") {
		return false
	}
	for _, segment := range strings.Split(strings.ReplaceAll(lowerPath, "%2e", "."), "/") {
		if segment == "." || segment == ".." {
			return false
		}
	}
	return true
}

// ClientIDMetadataDocumentURLAllowed reports whether clientID matches an entry of allowed. An
// entry ending with a slash allows every client_id it is a prefix of, any other entry allows
// exactly that client_id. Entries are compared as plain strings, as the specification requires
// for client_id URLs.
func ClientIDMetadataDocumentURLAllowed(allowed []string, clientID string) bool {
	for _, entry := range allowed {
		if strings.HasSuffix(entry, "/") {
			if strings.HasPrefix(clientID, entry) {
				return true
			}
			continue
		}
		if clientID == entry {
			return true
		}
	}
	return false
}
