package denylist

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHTTPTransport_Allowed(t *testing.T) {
	t.Parallel()

	// Create a denylist that does NOT block our local test server
	denyList := []AddressChecker{
		NewHostChecker("192.168.1.0/24"),
		NewHostChecker("10.0.0.1"),
	}

	// Spin up a dummy local test server. It will listen on a loopback address like 127.0.0.1 or ::1
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("success"))
	}))
	defer server.Close()

	// Instantiate the transport with the denylist rules
	transport := NewHTTPTransport(denyList)
	client := &http.Client{Transport: transport}

	// Execute the request to the allowed local server target
	resp, err := client.Get(server.URL)
	assert.NoError(t, err)
	if resp != nil {
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
}

func TestNewHTTPTransport_BlockedByControl(t *testing.T) {
	t.Parallel()

	// Spin up a test server so we can dynamically extract its specific local listening loopback IP
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Parse out the target host/ip text from the server URL
	u, err := net.Dial("tcp", server.Listener.Addr().String())
	assert.NoError(t, err)
	localIP, _, err := net.SplitHostPort(u.RemoteAddr().String())
	assert.NoError(t, err)
	_ = u.Close()

	// Add the explicitly discovered test server IP to the denylist
	denyList := []AddressChecker{
		NewHostChecker(localIP),
	}

	transport := NewHTTPTransport(denyList)
	client := &http.Client{Transport: transport}

	// This must fail because the Control hook interceptor throws an address denied error
	// right before the low-level TCP handshake begins.
	_, err = client.Get(server.URL) //nolint:bodyclose
	assert.Error(t, err)

	// Because Go wraps dialer context execution errors inside a generic url.Error,
	// we verify that the true root cause matches our address denial definition.
	var addressDeniedErr *AddressDeniedError
	assert.ErrorAs(t, err, &addressDeniedErr)
}

func TestNewHTTPTransport_EmptyDenylistReturnsDefaultClone(t *testing.T) {
	t.Parallel()

	// An empty list should bypass adding a dialer wrapper entirely for performance
	transport := NewHTTPTransport([]AddressChecker{})

	assert.NotNil(t, transport)
	httpTransport, ok := transport.(*http.Transport)
	assert.True(t, ok)

	// Ensure that DialContext is still populated (inherited cleanly from http.DefaultTransport Clone)
	assert.NotNil(t, httpTransport.DialContext)
}

func TestNewHTTPTransport_MalformedAddressHandling(t *testing.T) {
	t.Parallel()

	denyList := []AddressChecker{
		NewHostChecker("127.0.0.1"),
	}

	transport := NewHTTPTransport(denyList)
	httpTransport, ok := transport.(*http.Transport)
	assert.True(t, ok)

	// Explicitly invoke DialContext with an unparseable address format (missing a valid port boundary)
	// to test how the underlying Control code deals with structural string parsing failure exceptions.
	_, err := httpTransport.DialContext(context.Background(), "tcp", "malformed-address-without-port")
	assert.Error(t, err)

	// Ensure it didn't crash and returns a clean, standard connection framework error context
	var netErr net.Error
	assert.True(t, errors.As(err, &netErr))
}

// TestNewHTTPTransport_BlockedThroughProxy pins that a proxy cannot be used to reach a
// denied target: the dial only sees the proxy address, so the target is checked by URL.
func TestNewHTTPTransport_BlockedThroughProxy(t *testing.T) {
	t.Parallel()

	var proxyHits atomic.Int32
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxyHits.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	defer proxy.Close()
	proxyURL, err := url.Parse(proxy.URL)
	assert.NoError(t, err)

	// The proxy itself is reachable; only the target is denied. The base transport stands in
	// for http.DefaultTransport with HTTP_PROXY set.
	base := http.DefaultTransport.(*http.Transport).Clone()
	base.Proxy = http.ProxyURL(proxyURL)
	client := &http.Client{Transport: newHTTPTransport([]AddressChecker{NewHostChecker("10.0.0.0/8")}, base)}

	resp, err := client.Get("http://10.1.2.3/metadata")
	if resp != nil {
		_ = resp.Body.Close()
	}
	assert.Error(t, err)
	var denied *AddressDeniedError
	assert.True(t, errors.As(err, &denied), "expected an AddressDeniedError, got %v", err)
	assert.Equal(t, int32(0), proxyHits.Load(), "a denied target must not be sent to the proxy")

	// An allowed target still goes through the proxy.
	resp, err = client.Get("http://203.0.113.10/metadata")
	assert.NoError(t, err)
	if resp != nil {
		_ = resp.Body.Close()
	}
	assert.Equal(t, int32(1), proxyHits.Load())
}

// TestNewHTTPTransport_KeepsEnvironmentProxy pins that the denylist does not disable
// proxies: deployments that need one for egress keep using it.
func TestNewHTTPTransport_KeepsEnvironmentProxy(t *testing.T) {
	t.Parallel()

	transport := NewHTTPTransport([]AddressChecker{NewHostChecker("10.0.0.0/8")}).(*http.Transport)
	assert.NotNil(t, transport.Proxy)
}

// TestNewHTTPTransport_TLSRejectsProxyRedirectedTarget pins why a proxy that resolves an HTTPS
// target differently cannot reach another host on ZITADEL's behalf: the TLS handshake verifies
// the certificate against the requested hostname, so the request is never sent to a host that
// cannot prove it is that name, even one whose certificate the client trusts.
func TestNewHTTPTransport_TLSRejectsProxyRedirectedTarget(t *testing.T) {
	t.Parallel()

	var internalHits atomic.Int32
	internal := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		internalHits.Add(1)
	}))
	defer internal.Close()

	// The proxy ignores the requested authority and tunnels to the internal server, as a proxy
	// resolving the hostname to a private address would.
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodConnect {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		upstream, err := net.Dial("tcp", internal.Listener.Addr().String())
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		w.WriteHeader(http.StatusOK)
		conn, _, err := http.NewResponseController(w).Hijack()
		if err != nil {
			_ = upstream.Close()
			return
		}
		go func() { _, _ = io.Copy(upstream, conn); _ = upstream.Close() }()
		_, _ = io.Copy(conn, upstream)
		_ = conn.Close()
	}))
	defer proxy.Close()
	proxyURL, err := url.Parse(proxy.URL)
	require.NoError(t, err)

	base := http.DefaultTransport.(*http.Transport).Clone()
	base.Proxy = http.ProxyURL(proxyURL)
	roots := x509.NewCertPool()
	roots.AddCert(internal.Certificate())
	base.TLSClientConfig = &tls.Config{RootCAs: roots, MinVersion: tls.VersionTLS12}
	client := &http.Client{Transport: newHTTPTransport([]AddressChecker{NewHostChecker("10.0.0.0/8")}, base)}

	resp, err := client.Get("https://203.0.113.10/metadata")
	if resp != nil {
		_ = resp.Body.Close()
	}
	var hostnameErr x509.HostnameError
	assert.True(t, errors.As(err, &hostnameErr), "expected a certificate hostname error, got %v", err)
	assert.Equal(t, int32(0), internalHits.Load(), "the redirected host must not receive the request")
}
