package denylist

import (
	"net"
	"net/http"
	"net/url"
	"syscall"
	"time"
)

// NewHTTPTransport returns a cloned default transport that enforces denylist checks
// right before each TCP dial to avoid DNS rebinding TOCTOU gaps.
//
// When a request goes through a proxy, the dial only reaches the proxy, so the dial-time
// check never sees the target. Such requests are checked against the denylist by URL,
// before the proxy is contacted. That check resolves the hostname locally while the proxy
// resolves it again, so it cannot rule out a name that resolves differently at the proxy:
// a proxy in front of ZITADEL has to enforce the same deny rules itself.
func NewHTTPTransport(denyList []AddressChecker) http.RoundTripper {
	return newHTTPTransport(denyList, http.DefaultTransport.(*http.Transport).Clone())
}

// newHTTPTransport applies the denylist checks to base.
func newHTTPTransport(denyList []AddressChecker, base *http.Transport) http.RoundTripper {
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
	base.Proxy = denyProxiedTargets(denyList, base.Proxy)

	return base
}

// denyProxiedTargets wraps proxy so that a request which would be sent through a proxy is
// refused when its target is on the denyList. The target is resolved with the request
// context, so a stalled lookup cannot outlive the request. A request without a proxy is left
// to the dial-time check.
func denyProxiedTargets(denyList []AddressChecker, proxy func(*http.Request) (*url.URL, error)) func(*http.Request) (*url.URL, error) {
	if proxy == nil {
		return nil
	}
	return func(req *http.Request) (*url.URL, error) {
		proxyURL, err := proxy(req)
		if err != nil || proxyURL == nil {
			return proxyURL, err
		}
		lookup := func(host string) ([]net.IP, error) {
			return net.DefaultResolver.LookupIP(req.Context(), "ip", host)
		}
		if err := IsURLBlocked(denyList, req.URL, lookup); err != nil {
			return nil, err
		}
		return proxyURL, nil
	}
}
