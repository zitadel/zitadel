package http

import (
	"crypto/tls"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/denylist"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const testVerifier = "challenge-token-value"

func TestValidateDomainHTTP_Success(t *testing.T) {
	t.Parallel()

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/.well-known/zitadel-challenge/"+testVerifier+".txt", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(testVerifier))
	}))
	t.Cleanup(server.Close)

	client := newDomainTestClient(t, &ClientConfig{
		MaxBodySize:  1024,
		Timeout:      2 * time.Second,
		MaxRedirects: 3,
		DenyList:     []denylist.AddressChecker{},
	})

	err := ValidateDomainHTTP(hostFromURL(t, server.URL), testVerifier, testVerifier, client)
	assert.NoError(t, err)
}

func TestValidateDomainHTTP_RedirectToDenylistedURL(t *testing.T) {
	t.Parallel()

	var blockedHits atomic.Int32
	blocked := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		blockedHits.Add(1)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(testVerifier))
	}))
	t.Cleanup(blocked.Close)

	// Redirect target uses hostname "localhost" while the challenge host is 127.0.0.1,
	// so a domain-only denylist entry blocks the redirect without blocking the initial dial.
	// Keep https:// so HTTPS-downgrade checks do not fire before the denylist.
	blockedURL, err := url.Parse(blocked.URL)
	require.NoError(t, err)
	blockedURL.Host = "localhost:" + blockedURL.Port()

	public := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, blockedURL.String(), http.StatusFound)
	}))
	t.Cleanup(public.Close)

	client := newDomainTestClient(t, &ClientConfig{
		MaxBodySize:  1024,
		Timeout:      2 * time.Second,
		MaxRedirects: 3,
		DenyList:     []denylist.AddressChecker{denylist.NewHostChecker("localhost")},
	})

	err = ValidateDomainHTTP(hostFromURL(t, public.URL), testVerifier, testVerifier, client)
	require.Error(t, err)
	assert.True(t, zerrors.IsInternal(err))
	assert.ErrorIs(t, err, denylist.NewAddressDeniedError("localhost"))
	assert.Equal(t, int32(0), blockedHits.Load(), "denylisted redirect target must not be dialed")
}

func TestValidateDomainHTTP_HTTPSDowngradeBlocked(t *testing.T) {
	t.Parallel()

	httpTarget := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(testVerifier))
	}))
	t.Cleanup(httpTarget.Close)

	public := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, httpTarget.URL, http.StatusFound)
	}))
	t.Cleanup(public.Close)

	client := newDomainTestClient(t, &ClientConfig{
		MaxBodySize:         1024,
		Timeout:             2 * time.Second,
		MaxRedirects:        3,
		AllowHTTPSDowngrade: false,
		DenyList:            []denylist.AddressChecker{},
	})

	err := ValidateDomainHTTP(hostFromURL(t, public.URL), testVerifier, testVerifier, client)
	require.Error(t, err)
	assert.True(t, zerrors.IsInternal(err))
	assert.ErrorIs(t, err, ErrHTTPsDowngrade)
}

func TestValidateDomainHTTP_OversizedBody(t *testing.T) {
	t.Parallel()

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("0123456789")) // 10 bytes
	}))
	t.Cleanup(server.Close)

	client := newDomainTestClient(t, &ClientConfig{
		MaxBodySize:  5,
		Timeout:      2 * time.Second,
		MaxRedirects: 3,
		DenyList:     []denylist.AddressChecker{},
	})

	err := ValidateDomainHTTP(hostFromURL(t, server.URL), testVerifier, testVerifier, client)
	require.Error(t, err)
	assert.True(t, zerrors.IsInternal(err))
	assert.True(t, errors.Is(err, ErrResponseTooLarge), "got: %v", err)
}

func TestValidateDomainHTTP_NilClient(t *testing.T) {
	t.Parallel()

	err := ValidateDomainHTTP("example.com", testVerifier, testVerifier, nil)
	require.Error(t, err)
	assert.True(t, zerrors.IsInternal(err))
}

func TestValidateDomain_HTTPUsesClient(t *testing.T) {
	t.Parallel()

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(testVerifier))
	}))
	t.Cleanup(server.Close)

	client := newDomainTestClient(t, &ClientConfig{
		MaxBodySize:  1024,
		Timeout:      2 * time.Second,
		MaxRedirects: 3,
		DenyList:     []denylist.AddressChecker{},
	})

	err := ValidateDomain(hostFromURL(t, server.URL), testVerifier, testVerifier, CheckTypeHTTP, client)
	assert.NoError(t, err)
}

func newDomainTestClient(t *testing.T, cfg *ClientConfig) *http.Client {
	t.Helper()
	client := cfg.NewClient()

	// Challenge URLs are always https://…; trust httptest TLS certs.
	transport, ok := client.Transport.(*MaxBytesRoundTripper)
	require.True(t, ok)
	httpTransport, ok := transport.Underlying.(*http.Transport)
	require.True(t, ok)
	httpTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec // test-only

	return client
}

func hostFromURL(t *testing.T, raw string) string {
	t.Helper()
	u, err := url.Parse(raw)
	require.NoError(t, err)
	return u.Host
}
