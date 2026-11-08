package denylist

import (
	"net"
	"net/http"
	"syscall"
	"time"
)

// NewHTTPTransport returns a cloned default transport that enforces denylist checks
// right before each TCP dial to avoid DNS rebinding TOCTOU gaps.
func NewHTTPTransport(denyList []AddressChecker) http.RoundTripper {
	base := http.DefaultTransport.(*http.Transport).Clone()
	if len(denyList) == 0 {
		return base
	}
	base.DialContext = NewDialer(denyList).DialContext

	return base
}

// NewDialer returns a [net.Dialer] that checks the resolved IP against the denylist
// right before each TCP dial to avoid DNS rebinding TOCTOU gaps.
// A blocked dial fails with an error wrapping an [AddressDeniedError].
// For connections that are not HTTP, e.g. SMTP.
func NewDialer(denyList []AddressChecker) *net.Dialer {
	return &net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 30 * time.Second,
		Control: func(network, address string, c syscall.RawConn) error {
			host, _, err := net.SplitHostPort(address)
			if err != nil {
				return err
			}

			parsedIP := net.ParseIP(host)
			if parsedIP == nil { // at this point, it must be an IP so it should never happen
				return &net.DNSError{Err: "invalid IP address", Name: host}
			}

			return IsHostBlocked(denyList, host, parsedIP)
		},
	}
}
