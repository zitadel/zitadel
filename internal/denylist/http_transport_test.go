package denylist

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
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
