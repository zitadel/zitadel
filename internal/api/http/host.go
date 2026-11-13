package http

import (
	"net"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/http/httpguts"

	"github.com/zitadel/zitadel/internal/zerrors"
)

// ValidateHost returns the host in its canonical form if it is a valid host[:port]:
// a host as allowed in a Host header, with an IPv6 address in brackets and an optional port between 1 and 65535.
//
// The host must only contain characters allowed in a Host header (the check net/http applies to the Host header)
// and must parse as a pure URL authority: no userinfo, path, query or fragment.
// This ensures the host checked against the trusted domains is the same host a URL built from it points to,
// e.g. "trusted.example.com:443@evil.example" is rejected instead of being checked as "trusted.example.com".
func ValidateHost(host string) (string, error) {
	if host == "" || !httpguts.ValidHostHeader(host) {
		return "", zerrors.ThrowInvalidArgument(nil, "HTTP-Hq7tN", "invalid host")
	}
	u, err := url.Parse("//" + host)
	// httpguts already rejects the userinfo delimiter "@", the userinfo check is kept as defense in depth.
	if err != nil || u.User != nil || u.Host != host {
		return "", zerrors.ThrowInvalidArgument(err, "HTTP-Pw4zR", "invalid host")
	}
	hostname, port := u.Hostname(), u.Port()
	// url.Parse only accepts a valid IPv6 address in brackets, but outside of brackets it accepts
	// an empty hostname (e.g. ":8080"), which would skip the trusted domain check,
	// and colons or brackets (e.g. "::1" as hostname ":" with port "1", or "]").
	if !strings.HasPrefix(host, "[") && (hostname == "" || strings.ContainsAny(hostname, "[]:")) {
		return "", zerrors.ThrowInvalidArgument(nil, "HTTP-Mv2xK", "invalid host")
	}
	// an empty port (e.g. "example.com:") is allowed by RFC 3986 and means the default port
	if port == "" {
		return strings.TrimSuffix(host, ":"), nil
	}
	portNumber, err := strconv.ParseUint(port, 10, 16)
	if err != nil || portNumber == 0 {
		return "", zerrors.ThrowInvalidArgument(err, "HTTP-Fd9sL", "invalid port")
	}
	// return the port without leading zeros (e.g. "example.com:0443"), JoinHostPort adds the brackets of an IPv6 address
	return net.JoinHostPort(hostname, strconv.FormatUint(portNumber, 10)), nil
}
