package denylist

import (
	"bufio"
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
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

// connectProxy is a minimal HTTP CONNECT proxy that records the authority of every CONNECT and
// tunnels to it.
type connectProxy struct {
	listener net.Listener
	mu       sync.Mutex
	targets  []string
}

func newConnectProxy(t *testing.T) *connectProxy {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	p := &connectProxy{listener: listener}
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go p.serve(conn)
		}
	}()
	t.Cleanup(func() { _ = listener.Close() })
	return p
}

func (p *connectProxy) serve(conn net.Conn) {
	defer conn.Close()
	req, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil || req.Method != http.MethodConnect {
		return
	}
	p.mu.Lock()
	p.targets = append(p.targets, req.Host)
	p.mu.Unlock()
	upstream, err := net.Dial("tcp", req.Host)
	if err != nil {
		_, _ = conn.Write([]byte("HTTP/1.1 502 Bad Gateway\r\n\r\n"))
		return
	}
	defer upstream.Close()
	_, _ = conn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
	go func() { _, _ = io.Copy(upstream, conn) }()
	_, _ = io.Copy(conn, upstream)
}

func (p *connectProxy) connectTargets() []string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return append([]string(nil), p.targets...)
}

// TestNewHTTPTransport_PinsProxiedHTTPSTargets pins that an HTTPS request through a proxy is
// tunnelled to the vetted address rather than to the hostname, so the proxy cannot resolve
// the name to a denied address on its own. Not parallel: it replaces the package resolver.
func TestNewHTTPTransport_PinsProxiedHTTPSTargets(t *testing.T) {
	target := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))
	defer target.Close()
	_, targetPort, err := net.SplitHostPort(target.Listener.Addr().String())
	require.NoError(t, err)

	resolved := map[string][]net.IP{
		"example.com":      {net.ParseIP("127.0.0.1")},
		"internal.example": {net.ParseIP("10.1.2.3")},
		"mixed.example":    {net.ParseIP("127.0.0.1"), net.ParseIP("10.1.2.3")},
	}
	previous := lookupIP
	lookupIP = func(host string) ([]net.IP, error) {
		if ips, ok := resolved[host]; ok {
			return ips, nil
		}
		return nil, &net.DNSError{Err: "no such host", Name: host}
	}
	t.Cleanup(func() { lookupIP = previous })

	proxy := newConnectProxy(t)
	proxyURL := &url.URL{Scheme: "http", Host: proxy.listener.Addr().String()}
	base := target.Client().Transport.(*http.Transport).Clone()
	base.Proxy = http.ProxyURL(proxyURL)
	client := &http.Client{Transport: newHTTPTransport([]AddressChecker{NewHostChecker("10.0.0.0/8")}, base)}

	t.Run("allowed target is tunnelled to its vetted address", func(t *testing.T) {
		// The test certificate is valid for example.com, so TLS still verifies the hostname.
		resp, err := client.Get("https://example.com:" + targetPort + "/")
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, []string{"127.0.0.1:" + targetPort}, proxy.connectTargets(), "the proxy must receive the vetted address, never the hostname")
	})

	for _, host := range []string{"internal.example", "mixed.example"} {
		t.Run("denied address "+host, func(t *testing.T) {
			before := len(proxy.connectTargets())
			resp, err := client.Get("https://" + host + ":" + targetPort + "/")
			if resp != nil {
				_ = resp.Body.Close()
			}
			var denied *AddressDeniedError
			require.True(t, errors.As(err, &denied), "expected an AddressDeniedError, got %v", err)
			assert.Len(t, proxy.connectTargets(), before, "a denied target must not reach the proxy")
		})
	}
}
