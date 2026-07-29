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
	dialer := &net.Dialer{
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
	base.DialContext = dialer.DialContext

	return base
}
