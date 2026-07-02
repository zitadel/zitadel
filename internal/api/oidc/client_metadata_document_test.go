package oidc

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v3/pkg/oidc"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/cache"
	"github.com/zitadel/zitadel/internal/cache/connector/gomap"
	"github.com/zitadel/zitadel/internal/cache/connector/noop"
	"github.com/zitadel/zitadel/internal/denylist"
	"github.com/zitadel/zitadel/internal/domain"
)

func TestLooksLikeClientIDMetadataURL(t *testing.T) {
	tests := []struct {
		clientID string
		want     bool
	}{
		{"https://app.example.com/oauth/client", true},
		{"https://app.example.com", true},
		{"https://app.example.com:8443/client", true},
		{"HTTPS://app.example.com/client", true},
		{"http://app.example.com/client", false},
		{"320948502934092851", false},
		{"320948502934092851@project", false},
		{"app.example.com/client", false},
		{"/oauth/client", false},
		{"", false},
		{"ftp://app.example.com", false},
	}
	for _, tt := range tests {
		t.Run(tt.clientID, func(t *testing.T) {
			assert.Equal(t, tt.want, looksLikeClientIDMetadataURL(tt.clientID))
		})
	}
}

func TestSameOrigin(t *testing.T) {
	clientID := mustParseURL(t, "https://app.example.com/oauth/client")
	tests := []struct {
		name   string
		rawURI string
		want   bool
	}{
		{"same origin", "https://app.example.com/callback", true},
		{"same origin explicit default port", "https://app.example.com:443/callback", true},
		{"same origin upper case host", "https://APP.EXAMPLE.COM/callback", true},
		{"same origin trailing dot host", "https://app.example.com./callback", true},
		{"different host", "https://other.example.com/callback", false},
		{"different scheme", "http://app.example.com/callback", false},
		{"different port", "https://app.example.com:8443/callback", false},
		{"subdomain", "https://sub.app.example.com/callback", false},
		{"relative uri", "/callback", false},
		{"custom scheme", "com.example.app://callback", false},
		{"loopback", "http://127.0.0.1/callback", false},
		{"userinfo is rejected", "https://user@app.example.com/callback", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, sameOrigin(clientID, tt.rawURI))
		})
	}
}

func TestSameOrigin_PortNormalization(t *testing.T) {
	// An explicit default port on the client_id must match an implicit one on the redirect.
	explicitDefault := mustParseURL(t, "https://app.example.com:443/client")
	assert.True(t, sameOrigin(explicitDefault, "https://app.example.com/callback"))

	nonDefault := mustParseURL(t, "https://app.example.com:8443/client")
	assert.True(t, sameOrigin(nonDefault, "https://app.example.com:8443/callback"))
	assert.False(t, sameOrigin(nonDefault, "https://app.example.com/callback"))
}

func TestSameOrigin_IPv6Host(t *testing.T) {
	clientID := mustParseURL(t, "https://[2001:db8::1]/client")
	assert.True(t, sameOrigin(clientID, "https://[2001:db8::1]/callback"))
	assert.False(t, sameOrigin(clientID, "https://[2001:db8::2]/callback"))
}

func TestSameOrigin_PunycodeHost(t *testing.T) {
	// The client_id uses the unicode host while the redirect_uri uses its punycode form.
	clientID := mustParseURL(t, "https://bücher.example/client")
	assert.True(t, sameOrigin(clientID, "https://xn--bcher-kva.example/callback"))
}

func TestOriginMatchedURIs(t *testing.T) {
	clientID := mustParseURL(t, "https://app.example.com/client")
	got := originMatchedURIs(clientID, []string{
		"https://app.example.com/callback",
		" https://app.example.com/second ",
		"https://evil.example.com/callback",
		"http://app.example.com/callback",
		"",
	})
	assert.Equal(t, []string{"https://app.example.com/callback", "https://app.example.com/second"}, got)

	assert.Nil(t, originMatchedURIs(clientID, []string{"https://evil.example.com/callback"}))
	assert.Nil(t, originMatchedURIs(clientID, nil))
}

func TestCacheTTLFromResponse(t *testing.T) {
	const max = 15 * time.Minute
	tests := []struct {
		name   string
		header http.Header
		want   time.Duration
	}{
		{"no headers default to cap", http.Header{}, max},
		{"no-store", http.Header{"Cache-Control": {"no-store"}}, 0},
		{"no-cache", http.Header{"Cache-Control": {"no-cache"}}, 0},
		{"max-age below cap", http.Header{"Cache-Control": {"max-age=60"}}, time.Minute},
		{"max-age above cap", http.Header{"Cache-Control": {"max-age=3600"}}, max},
		{"max-age zero", http.Header{"Cache-Control": {"max-age=0"}}, 0},
		{"max-age with other directives", http.Header{"Cache-Control": {"public, max-age=120"}}, 2 * time.Minute},
		{"expires past", http.Header{"Expires": {"Mon, 02 Jan 2006 15:04:05 GMT"}}, 0},
		{"no-store wins over max-age", http.Header{"Cache-Control": {"no-store, max-age=600"}}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, cacheTTLFromResponse(tt.header, max))
		})
	}
}

func TestClientIDMetadataResolver_ResolveClient(t *testing.T) {
	ctx := authz.NewMockContext("instance", "org", "")

	t.Run("valid document", func(t *testing.T) {
		server := newMetadataServer(t, http.StatusOK, "", func(clientID string) clientRegistrationRequest {
			return clientRegistrationRequest{
				RedirectURIs:            []string{clientID + "/callback"},
				GrantTypes:              []string{"authorization_code", "refresh_token"},
				ResponseTypes:           []string{"code"},
				TokenEndpointAuthMethod: "none",
			}
		})
		resolver := newTestResolver(server.Client(), noop.NewCache[clientIDMetadataCacheIndex, string, *clientIDMetadataCacheEntry]())

		client, err := resolver.ResolveClient(ctx, "instance", server.URL)
		require.NoError(t, err)
		assert.Equal(t, server.URL, client.ClientID)
		assert.Equal(t, "instance", client.InstanceID)
		assert.Equal(t, domain.OIDCAuthMethodTypeNone, client.AuthMethodType)
		assert.Equal(t, domain.OIDCApplicationTypeWeb, client.ApplicationType)
		assert.Equal(t, []string{server.URL + "/callback"}, client.RedirectURIs)
		assert.Empty(t, client.ProjectID)
		assert.Empty(t, client.HashedSecret)
		assert.False(t, client.IsDevMode)
		require.NotNil(t, client.Settings)
		assert.Equal(t, time.Hour, client.Settings.AccessTokenLifetime)
	})

	t.Run("only origin matched redirect_uris are kept", func(t *testing.T) {
		server := newMetadataServer(t, http.StatusOK, "", func(clientID string) clientRegistrationRequest {
			return clientRegistrationRequest{
				RedirectURIs:            []string{clientID + "/callback", "https://evil.example.com/callback"},
				TokenEndpointAuthMethod: "none",
			}
		})
		resolver := newTestResolver(server.Client(), noop.NewCache[clientIDMetadataCacheIndex, string, *clientIDMetadataCacheEntry]())

		client, err := resolver.ResolveClient(ctx, "instance", server.URL)
		require.NoError(t, err)
		assert.Equal(t, []string{server.URL + "/callback"}, client.RedirectURIs)
	})

	t.Run("no origin matched redirect_uri is rejected", func(t *testing.T) {
		server := newMetadataServer(t, http.StatusOK, "", func(clientID string) clientRegistrationRequest {
			return clientRegistrationRequest{
				RedirectURIs:            []string{"https://evil.example.com/callback"},
				TokenEndpointAuthMethod: "none",
			}
		})
		resolver := newTestResolver(server.Client(), noop.NewCache[clientIDMetadataCacheIndex, string, *clientIDMetadataCacheEntry]())

		_, err := resolver.ResolveClient(ctx, "instance", server.URL)
		assertInvalidClient(t, err)
	})

	t.Run("confidential auth method is rejected", func(t *testing.T) {
		server := newMetadataServer(t, http.StatusOK, "", func(clientID string) clientRegistrationRequest {
			return clientRegistrationRequest{
				RedirectURIs:            []string{clientID + "/callback"},
				TokenEndpointAuthMethod: "client_secret_basic",
			}
		})
		resolver := newTestResolver(server.Client(), noop.NewCache[clientIDMetadataCacheIndex, string, *clientIDMetadataCacheEntry]())

		_, err := resolver.ResolveClient(ctx, "instance", server.URL)
		assertInvalidClient(t, err)
	})

	t.Run("jwks is rejected", func(t *testing.T) {
		server := newMetadataServer(t, http.StatusOK, "", func(clientID string) clientRegistrationRequest {
			return clientRegistrationRequest{
				RedirectURIs: []string{clientID + "/callback"},
				JWKsURI:      "https://app.example.com/jwks",
			}
		})
		resolver := newTestResolver(server.Client(), noop.NewCache[clientIDMetadataCacheIndex, string, *clientIDMetadataCacheEntry]())

		_, err := resolver.ResolveClient(ctx, "instance", server.URL)
		assertInvalidClient(t, err)
	})

	t.Run("non-200 response is rejected", func(t *testing.T) {
		server := newMetadataServer(t, http.StatusNotFound, "", nil)
		resolver := newTestResolver(server.Client(), noop.NewCache[clientIDMetadataCacheIndex, string, *clientIDMetadataCacheEntry]())

		_, err := resolver.ResolveClient(ctx, "instance", server.URL)
		assertInvalidClient(t, err)
	})

	t.Run("oversized body is rejected", func(t *testing.T) {
		server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"redirect_uris":["` + strings.Repeat("a", clientIDMetadataMaxBodyBytes) + `"]}`))
		}))
		t.Cleanup(server.Close)
		resolver := newTestResolver(server.Client(), noop.NewCache[clientIDMetadataCacheIndex, string, *clientIDMetadataCacheEntry]())

		_, err := resolver.ResolveClient(ctx, "instance", server.URL)
		assertInvalidClient(t, err)
	})

	t.Run("invalid json is rejected", func(t *testing.T) {
		server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("not json"))
		}))
		t.Cleanup(server.Close)
		resolver := newTestResolver(server.Client(), noop.NewCache[clientIDMetadataCacheIndex, string, *clientIDMetadataCacheEntry]())

		_, err := resolver.ResolveClient(ctx, "instance", server.URL)
		assertInvalidClient(t, err)
	})
}

func TestClientIDMetadataResolver_Cache(t *testing.T) {
	ctx := authz.NewMockContext("instance", "org", "")
	documentCache := gomap.NewCache[clientIDMetadataCacheIndex, string, *clientIDMetadataCacheEntry](
		context.Background(),
		[]clientIDMetadataCacheIndex{clientIDMetadataCacheIndexURL},
		cache.Config{MaxAge: time.Hour, LastUseAge: time.Hour},
	)

	t.Run("cacheable document is fetched once", func(t *testing.T) {
		server := newMetadataServer(t, http.StatusOK, "max-age=600", func(clientID string) clientRegistrationRequest {
			return clientRegistrationRequest{RedirectURIs: []string{clientID + "/callback"}, TokenEndpointAuthMethod: "none"}
		})
		resolver := newTestResolver(server.Client(), documentCache)

		_, err := resolver.ResolveClient(ctx, "instance", server.URL)
		require.NoError(t, err)
		_, err = resolver.ResolveClient(ctx, "instance", server.URL)
		require.NoError(t, err)
		assert.Equal(t, int32(1), server.hits.Load())
	})

	t.Run("no-store document is fetched every time", func(t *testing.T) {
		server := newMetadataServer(t, http.StatusOK, "no-store", func(clientID string) clientRegistrationRequest {
			return clientRegistrationRequest{RedirectURIs: []string{clientID + "/callback"}, TokenEndpointAuthMethod: "none"}
		})
		resolver := newTestResolver(server.Client(), documentCache)

		_, err := resolver.ResolveClient(ctx, "instance", server.URL)
		require.NoError(t, err)
		_, err = resolver.ResolveClient(ctx, "instance", server.URL)
		require.NoError(t, err)
		assert.Equal(t, int32(2), server.hits.Load())
	})

	t.Run("failed resolution is negatively cached", func(t *testing.T) {
		server := newMetadataServer(t, http.StatusNotFound, "", nil)
		resolver := newTestResolver(server.Client(), documentCache)

		_, err := resolver.ResolveClient(ctx, "instance", server.URL)
		assertInvalidClient(t, err)
		_, err = resolver.ResolveClient(ctx, "instance", server.URL)
		assertInvalidClient(t, err)
		assert.Equal(t, int32(1), server.hits.Load(), "an unresolvable client_id must not be fetched again within the negative cache window")
	})

	t.Run("expired entry is re-fetched", func(t *testing.T) {
		server := newMetadataServer(t, http.StatusOK, "max-age=600", func(clientID string) clientRegistrationRequest {
			return clientRegistrationRequest{RedirectURIs: []string{clientID + "/callback"}, TokenEndpointAuthMethod: "none"}
		})
		resolver := newTestResolver(server.Client(), documentCache)

		client, err := resolver.ResolveClient(ctx, "instance", server.URL)
		require.NoError(t, err)
		require.Equal(t, int32(1), server.hits.Load())

		// Force the cached entry to be expired; the next resolution must re-fetch.
		documentCache.Set(ctx, &clientIDMetadataCacheEntry{
			Key:    clientIDMetadataCacheKey("instance", server.URL),
			Client: client,
			Expiry: time.Now().Add(-time.Hour),
		})
		_, err = resolver.ResolveClient(ctx, "instance", server.URL)
		require.NoError(t, err)
		assert.Equal(t, int32(2), server.hits.Load())
	})
}

func TestClientIDMetadataResolver_SSRFDenylist(t *testing.T) {
	ctx := authz.NewMockContext("instance", "org", "")
	// The test server listens on a loopback address, which the operator denylist blocks at
	// dial time, so the request must never reach the handler.
	server := newMetadataServer(t, http.StatusOK, "", func(clientID string) clientRegistrationRequest {
		return clientRegistrationRequest{RedirectURIs: []string{clientID + "/callback"}, TokenEndpointAuthMethod: "none"}
	})

	checkers, err := denylist.ParseDenyList([]string{"127.0.0.0/8", "::1/128"})
	require.NoError(t, err)
	config := &http_util.ClientConfig{
		MaxBodySize:  clientIDMetadataMaxBodyBytes,
		Timeout:      5 * time.Second,
		MaxRedirects: 2,
		DenyList:     checkers,
	}
	resolver := newTestResolver(config.NewClient(), noop.NewCache[clientIDMetadataCacheIndex, string, *clientIDMetadataCacheEntry]())

	_, err = resolver.ResolveClient(ctx, "instance", server.URL)
	assertInvalidClient(t, err)
	assert.Equal(t, int32(0), server.hits.Load(), "the denylist should block the dial before the request reaches the server")
}

// helpers

func newTestResolver(httpClient *http.Client, documentCache cache.Cache[clientIDMetadataCacheIndex, string, *clientIDMetadataCacheEntry]) *clientIDMetadataResolver {
	return newClientIDMetadataResolver(httpClient, documentCache, time.Hour, time.Hour)
}

type metadataServer struct {
	*httptest.Server
	hits atomic.Int32
}

// newMetadataServer starts a TLS test server that serves the document returned by build for
// its own URL. When build is nil it only writes the given status code.
func newMetadataServer(t *testing.T, status int, cacheControl string, build func(clientID string) clientRegistrationRequest) *metadataServer {
	t.Helper()
	server := &metadataServer{}
	server.Server = httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server.hits.Add(1)
		if cacheControl != "" {
			w.Header().Set("Cache-Control", cacheControl)
		}
		w.Header().Set("Content-Type", "application/json")
		if status != http.StatusOK {
			w.WriteHeader(status)
			return
		}
		_ = json.NewEncoder(w).Encode(build("https://" + r.Host))
	}))
	t.Cleanup(server.Close)
	return server
}

func mustParseURL(t *testing.T, raw string) *url.URL {
	t.Helper()
	u, err := url.Parse(raw)
	require.NoError(t, err)
	return u
}

func assertInvalidClient(t *testing.T, err error) {
	t.Helper()
	require.Error(t, err)
	var oidcErr *oidc.Error
	require.ErrorAs(t, err, &oidcErr)
	assert.Equal(t, oidc.InvalidClient, oidcErr.ErrorType)
}
