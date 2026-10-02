package oidc

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"sync"
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
	"github.com/zitadel/zitadel/internal/feature"
	"github.com/zitadel/zitadel/internal/query"
)

func TestNewClientIDMetadataAllowlist(t *testing.T) {
	t.Run("valid entries are kept, blank entries are skipped", func(t *testing.T) {
		allowlist, err := newClientIDMetadataAllowlist(ClientIDMetadataDocumentConfig{
			AllowedURLs: []string{
				"https://app.example.com/client",
				" https://clients.example.com/ ",
				"",
			},
		})
		require.NoError(t, err)
		assert.Equal(t, clientIDMetadataAllowlist{urls: []string{"https://app.example.com/client", "https://clients.example.com/"}}, allowlist)
		assert.True(t, allowlist.overlaps(nil, true))
	})
	t.Run("nothing allowed", func(t *testing.T) {
		allowlist, err := newClientIDMetadataAllowlist(ClientIDMetadataDocumentConfig{})
		require.NoError(t, err)
		assert.False(t, allowlist.overlaps(nil, true))
		assert.False(t, allowlist.allows("https://app.example.com/client"))
	})
	t.Run("any url allowed", func(t *testing.T) {
		allowlist, err := newClientIDMetadataAllowlist(ClientIDMetadataDocumentConfig{AllowAnyURL: true})
		require.NoError(t, err)
		assert.True(t, allowlist.overlaps(nil, true))
		assert.True(t, allowlist.allows("https://app.example.com/client"))
	})
	for _, entry := range []string{
		"https://app.example.com",
		"http://app.example.com/client",
		"https://app.example.com/oauth/../",
		"https://app.example.com/client#fragment",
		"app.example.com/client",
	} {
		t.Run("invalid entry "+entry, func(t *testing.T) {
			_, err := newClientIDMetadataAllowlist(ClientIDMetadataDocumentConfig{AllowedURLs: []string{"https://app.example.com/client", entry}})
			require.Error(t, err)
		})
	}
}

// TestClientIDMetadataResolver_Handles pins the control plane: resolution needs the instance
// security setting, a valid client_id URL, and both the system and the instance must allow it.
// Anything else must fall through to the regular database lookup.
func TestClientIDMetadataResolver_Handles(t *testing.T) {
	const clientID = "https://app.example.com/client"
	tests := []struct {
		name     string
		system   clientIDMetadataAllowlist
		instance []authz.MockContextInstanceOpts
		clientID string
		want     bool
	}{
		{
			name:     "system and instance allow the url",
			system:   clientIDMetadataAllowlist{urls: []string{"https://app.example.com/"}},
			instance: []authz.MockContextInstanceOpts{authz.WithMockClientIDMetadataDocumentAllowedURLs(clientID)},
			clientID: clientID,
			want:     true,
		},
		{
			name:     "system allows any url, instance allows the url",
			system:   clientIDMetadataAllowlist{allowAny: true},
			instance: []authz.MockContextInstanceOpts{authz.WithMockClientIDMetadataDocumentAllowedURLs("https://app.example.com/")},
			clientID: clientID,
			want:     true,
		},
		{
			name:     "system allows the url, instance allows any url",
			system:   clientIDMetadataAllowlist{urls: []string{clientID}},
			instance: []authz.MockContextInstanceOpts{authz.WithMockClientIDMetadataDocumentAllowAnyURL(true)},
			clientID: clientID,
			want:     true,
		},
		{
			name:     "both allow any url",
			system:   clientIDMetadataAllowlist{allowAny: true},
			instance: []authz.MockContextInstanceOpts{authz.WithMockClientIDMetadataDocumentAllowAnyURL(true)},
			clientID: clientID,
			want:     true,
		},
		{
			name:     "instance allows any url, system does not allow the url",
			system:   clientIDMetadataAllowlist{urls: []string{"https://other.example.com/"}},
			instance: []authz.MockContextInstanceOpts{authz.WithMockClientIDMetadataDocumentAllowAnyURL(true)},
			clientID: clientID,
		},
		{
			name:     "instance allows the url, system allows nothing",
			instance: []authz.MockContextInstanceOpts{authz.WithMockClientIDMetadataDocumentAllowedURLs(clientID)},
			clientID: clientID,
		},
		{
			name:     "system allows any url, instance allows nothing",
			system:   clientIDMetadataAllowlist{allowAny: true},
			clientID: clientID,
		},
		{
			name:     "system allows any url, instance allows another url",
			system:   clientIDMetadataAllowlist{allowAny: true},
			instance: []authz.MockContextInstanceOpts{authz.WithMockClientIDMetadataDocumentAllowedURLs("https://other.example.com/")},
			clientID: clientID,
		},
		{
			name:     "invalid url under allowed prefixes",
			system:   clientIDMetadataAllowlist{allowAny: true},
			instance: []authz.MockContextInstanceOpts{authz.WithMockClientIDMetadataDocumentAllowAnyURL(true)},
			clientID: "https://app.example.com/a/../client",
		},
		{
			name:     "regular client_id",
			system:   clientIDMetadataAllowlist{allowAny: true},
			instance: []authz.MockContextInstanceOpts{authz.WithMockClientIDMetadataDocumentAllowAnyURL(true)},
			clientID: "320948502934092851",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := newClientIDMetadataResolver(http.DefaultClient, tt.system, nil, time.Hour, time.Hour, nil)
			enabled := append([]authz.MockContextInstanceOpts{authz.WithMockClientIDMetadataDocument(true)}, tt.instance...)
			assert.Equal(t, tt.want, resolver.Handles(authz.NewMockContext("instance", "org", "", enabled...), tt.clientID), "enabled")
			assert.False(t, resolver.Handles(authz.NewMockContext("instance", "org", "", tt.instance...), tt.clientID), "disabled")
		})
	}
}

// TestClientIDMetadataResolver_Supported pins what discovery advertises: support needs the
// instance setting and at least one allowed client_id URL on both the system and the instance,
// so a client is never told to use a mechanism that cannot resolve it.
func TestClientIDMetadataResolver_Supported(t *testing.T) {
	tests := []struct {
		name     string
		system   clientIDMetadataAllowlist
		instance []authz.MockContextInstanceOpts
		want     bool
	}{
		{"both allow urls", clientIDMetadataAllowlist{urls: []string{"https://app.example.com/"}}, []authz.MockContextInstanceOpts{authz.WithMockClientIDMetadataDocumentAllowedURLs("https://app.example.com/client")}, true},
		{"both allow any url", clientIDMetadataAllowlist{allowAny: true}, []authz.MockContextInstanceOpts{authz.WithMockClientIDMetadataDocumentAllowAnyURL(true)}, true},
		{"system allows nothing", clientIDMetadataAllowlist{}, []authz.MockContextInstanceOpts{authz.WithMockClientIDMetadataDocumentAllowAnyURL(true)}, false},
		{"instance allows nothing", clientIDMetadataAllowlist{allowAny: true}, nil, false},
		{"instance url under system prefix", clientIDMetadataAllowlist{urls: []string{"https://app.example.com/"}}, []authz.MockContextInstanceOpts{authz.WithMockClientIDMetadataDocumentAllowedURLs("https://app.example.com/oauth/")}, true},
		{"system url under instance prefix", clientIDMetadataAllowlist{urls: []string{"https://app.example.com/client"}}, []authz.MockContextInstanceOpts{authz.WithMockClientIDMetadataDocumentAllowedURLs("https://app.example.com/")}, true},
		{"system allows any url, instance allows urls", clientIDMetadataAllowlist{allowAny: true}, []authz.MockContextInstanceOpts{authz.WithMockClientIDMetadataDocumentAllowedURLs("https://app.example.com/client")}, true},
		{"instance allows any url, system allows urls", clientIDMetadataAllowlist{urls: []string{"https://app.example.com/client"}}, []authz.MockContextInstanceOpts{authz.WithMockClientIDMetadataDocumentAllowAnyURL(true)}, true},
		{"lists do not overlap", clientIDMetadataAllowlist{urls: []string{"https://a.example.com/"}}, []authz.MockContextInstanceOpts{authz.WithMockClientIDMetadataDocumentAllowedURLs("https://b.example.com/")}, false},
		{"different exact urls", clientIDMetadataAllowlist{urls: []string{"https://app.example.com/a"}}, []authz.MockContextInstanceOpts{authz.WithMockClientIDMetadataDocumentAllowedURLs("https://app.example.com/b")}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := newClientIDMetadataResolver(http.DefaultClient, tt.system, nil, time.Hour, time.Hour, nil)
			enabled := append([]authz.MockContextInstanceOpts{authz.WithMockClientIDMetadataDocument(true)}, tt.instance...)
			assert.Equal(t, tt.want, resolver.Supported(authz.NewMockContext("instance", "org", "", enabled...)), "enabled")
			assert.False(t, resolver.Supported(authz.NewMockContext("instance", "org", "", tt.instance...)), "disabled")
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
	ctx := testResolverContext()

	t.Run("valid document", func(t *testing.T) {
		server := newMetadataServer(t, http.StatusOK, "", func(clientID string) clientRegistrationRequest {
			return clientRegistrationRequest{
				RedirectURIs:            []string{clientID + "/callback"},
				GrantTypes:              []string{"authorization_code", "refresh_token"},
				ResponseTypes:           []string{"code"},
				TokenEndpointAuthMethod: "none",
			}
		})
		resolver := newTestResolver(server.URL, server.Client(), noop.NewCache[clientIDMetadataCacheIndex, string, *clientIDMetadataCacheEntry]())

		client, err := resolver.ResolveClient(ctx, "instance", server.URL+testClientIDPath)
		require.NoError(t, err)
		assert.Equal(t, server.URL+testClientIDPath, client.ClientID)
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
		resolver := newTestResolver(server.URL, server.Client(), noop.NewCache[clientIDMetadataCacheIndex, string, *clientIDMetadataCacheEntry]())

		client, err := resolver.ResolveClient(ctx, "instance", server.URL+testClientIDPath)
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
		resolver := newTestResolver(server.URL, server.Client(), noop.NewCache[clientIDMetadataCacheIndex, string, *clientIDMetadataCacheEntry]())

		_, err := resolver.ResolveClient(ctx, "instance", server.URL+testClientIDPath)
		assertInvalidClient(t, err)
	})

	t.Run("confidential auth method is rejected", func(t *testing.T) {
		server := newMetadataServer(t, http.StatusOK, "", func(clientID string) clientRegistrationRequest {
			return clientRegistrationRequest{
				RedirectURIs:            []string{clientID + "/callback"},
				TokenEndpointAuthMethod: "client_secret_basic",
			}
		})
		resolver := newTestResolver(server.URL, server.Client(), noop.NewCache[clientIDMetadataCacheIndex, string, *clientIDMetadataCacheEntry]())

		_, err := resolver.ResolveClient(ctx, "instance", server.URL+testClientIDPath)
		assertInvalidClient(t, err)
	})

	t.Run("jwks is rejected", func(t *testing.T) {
		server := newMetadataServer(t, http.StatusOK, "", func(clientID string) clientRegistrationRequest {
			return clientRegistrationRequest{
				RedirectURIs: []string{clientID + "/callback"},
				JWKsURI:      "https://app.example.com/jwks",
			}
		})
		resolver := newTestResolver(server.URL, server.Client(), noop.NewCache[clientIDMetadataCacheIndex, string, *clientIDMetadataCacheEntry]())

		_, err := resolver.ResolveClient(ctx, "instance", server.URL+testClientIDPath)
		assertInvalidClient(t, err)
	})

	t.Run("non-200 response is rejected", func(t *testing.T) {
		server := newMetadataServer(t, http.StatusNotFound, "", nil)
		resolver := newTestResolver(server.URL, server.Client(), noop.NewCache[clientIDMetadataCacheIndex, string, *clientIDMetadataCacheEntry]())

		_, err := resolver.ResolveClient(ctx, "instance", server.URL+testClientIDPath)
		assertInvalidClient(t, err)
	})

	t.Run("oversized body is rejected", func(t *testing.T) {
		server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"redirect_uris":["` + strings.Repeat("a", clientIDMetadataMaxBodyBytes) + `"]}`))
		}))
		t.Cleanup(server.Close)
		resolver := newTestResolver(server.URL, server.Client(), noop.NewCache[clientIDMetadataCacheIndex, string, *clientIDMetadataCacheEntry]())

		_, err := resolver.ResolveClient(ctx, "instance", server.URL+testClientIDPath)
		assertInvalidClient(t, err)
	})

	t.Run("invalid json is rejected", func(t *testing.T) {
		server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("not json"))
		}))
		t.Cleanup(server.Close)
		resolver := newTestResolver(server.URL, server.Client(), noop.NewCache[clientIDMetadataCacheIndex, string, *clientIDMetadataCacheEntry]())

		_, err := resolver.ResolveClient(ctx, "instance", server.URL+testClientIDPath)
		assertInvalidClient(t, err)
	})
}

func TestClientIDMetadataResolver_Cache(t *testing.T) {
	ctx := testResolverContext()
	documentCache := gomap.NewCache[clientIDMetadataCacheIndex, string, *clientIDMetadataCacheEntry](
		context.Background(),
		[]clientIDMetadataCacheIndex{clientIDMetadataCacheIndexURL},
		cache.Config{MaxAge: time.Hour, LastUseAge: time.Hour},
	)

	t.Run("cacheable document is fetched once", func(t *testing.T) {
		server := newMetadataServer(t, http.StatusOK, "max-age=600", func(clientID string) clientRegistrationRequest {
			return clientRegistrationRequest{RedirectURIs: []string{clientID + "/callback"}, TokenEndpointAuthMethod: "none"}
		})
		resolver := newTestResolver(server.URL, server.Client(), documentCache)

		_, err := resolver.ResolveClient(ctx, "instance", server.URL+testClientIDPath)
		require.NoError(t, err)
		_, err = resolver.ResolveClient(ctx, "instance", server.URL+testClientIDPath)
		require.NoError(t, err)
		assert.Equal(t, int32(1), server.hits.Load())
	})

	t.Run("no-store document is fetched every time", func(t *testing.T) {
		server := newMetadataServer(t, http.StatusOK, "no-store", func(clientID string) clientRegistrationRequest {
			return clientRegistrationRequest{RedirectURIs: []string{clientID + "/callback"}, TokenEndpointAuthMethod: "none"}
		})
		resolver := newTestResolver(server.URL, server.Client(), documentCache)

		_, err := resolver.ResolveClient(ctx, "instance", server.URL+testClientIDPath)
		require.NoError(t, err)
		_, err = resolver.ResolveClient(ctx, "instance", server.URL+testClientIDPath)
		require.NoError(t, err)
		assert.Equal(t, int32(2), server.hits.Load())
	})

	t.Run("failed resolution is negatively cached", func(t *testing.T) {
		server := newMetadataServer(t, http.StatusNotFound, "", nil)
		resolver := newTestResolver(server.URL, server.Client(), documentCache)

		_, err := resolver.ResolveClient(ctx, "instance", server.URL+testClientIDPath)
		assertInvalidClient(t, err)
		_, err = resolver.ResolveClient(ctx, "instance", server.URL+testClientIDPath)
		assertInvalidClient(t, err)
		assert.Equal(t, int32(1), server.hits.Load(), "an unresolvable client_id must not be fetched again within the negative cache window")
	})

	t.Run("expired entry is re-fetched", func(t *testing.T) {
		server := newMetadataServer(t, http.StatusOK, "max-age=600", func(clientID string) clientRegistrationRequest {
			return clientRegistrationRequest{RedirectURIs: []string{clientID + "/callback"}, TokenEndpointAuthMethod: "none"}
		})
		resolver := newTestResolver(server.URL, server.Client(), documentCache)

		client, err := resolver.ResolveClient(ctx, "instance", server.URL+testClientIDPath)
		require.NoError(t, err)
		require.Equal(t, int32(1), server.hits.Load())

		// Force the cached entry to be expired; the next resolution must re-fetch.
		documentCache.Set(ctx, &clientIDMetadataCacheEntry{
			Key:    clientIDMetadataCacheKey("instance", server.URL+testClientIDPath),
			Client: client,
			Expiry: time.Now().Add(-time.Hour),
		})
		_, err = resolver.ResolveClient(ctx, "instance", server.URL+testClientIDPath)
		require.NoError(t, err)
		assert.Equal(t, int32(2), server.hits.Load())
	})
}

func TestClientIDMetadataResolver_SSRFDenylist(t *testing.T) {
	ctx := testResolverContext()
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
	resolver := newTestResolver(server.URL, config.NewClient(), noop.NewCache[clientIDMetadataCacheIndex, string, *clientIDMetadataCacheEntry]())

	_, err = resolver.ResolveClient(ctx, "instance", server.URL+testClientIDPath)
	assertInvalidClient(t, err)
	assert.Equal(t, int32(0), server.hits.Load(), "the denylist should block the dial before the request reaches the server")
}

// TestClientIDMetadataResolver_Allowlist pins that a client_id the system does not allow is
// rejected before anything else happens: no outbound request and no cache entry, so an
// unauthenticated caller can neither make the server fetch a URL nor fill the cache.
func TestClientIDMetadataResolver_Allowlist(t *testing.T) {
	ctx := testResolverContext()
	server := newMetadataServer(t, http.StatusOK, "", func(clientID string) clientRegistrationRequest {
		return clientRegistrationRequest{RedirectURIs: []string{clientID + "/callback"}, TokenEndpointAuthMethod: "none"}
	})
	documentCache := gomap.NewCache[clientIDMetadataCacheIndex, string, *clientIDMetadataCacheEntry](
		context.Background(),
		[]clientIDMetadataCacheIndex{clientIDMetadataCacheIndexURL},
		cache.Config{MaxAge: time.Hour, LastUseAge: time.Hour},
	)

	for _, allowlist := range []clientIDMetadataAllowlist{
		{},
		{urls: []string{"https://app.example.com/"}},
		{urls: []string{server.URL + "/other/"}},
		{urls: []string{server.URL + testClientIDPath + "/"}},
	} {
		t.Run(strings.Join(allowlist.urls, ","), func(t *testing.T) {
			resolver := newClientIDMetadataResolver(server.Client(), allowlist, documentCache, time.Hour, time.Hour, nil)
			_, err := resolver.ResolveClient(ctx, "instance", server.URL+testClientIDPath)
			assertInvalidClient(t, err)
			assert.Equal(t, int32(0), server.hits.Load(), "a client_id that is not allowed must not be fetched")
			_, found := documentCache.Get(ctx, clientIDMetadataCacheIndexURL, clientIDMetadataCacheKey("instance", server.URL+testClientIDPath))
			assert.False(t, found, "a client_id that is not allowed must not be cached")
		})
	}

	t.Run("exact entry", func(t *testing.T) {
		resolver := newClientIDMetadataResolver(server.Client(), clientIDMetadataAllowlist{urls: []string{server.URL + testClientIDPath}}, documentCache, time.Hour, time.Hour, nil)
		client, err := resolver.ResolveClient(ctx, "instance", server.URL+testClientIDPath)
		require.NoError(t, err)
		assert.Equal(t, server.URL+testClientIDPath, client.ClientID)
	})
}

// TestClientIDMetadataResolver_InstanceAllowlist pins that the instance narrows what the system
// allows: a client_id the system allows but the instance does not is neither fetched nor cached.
func TestClientIDMetadataResolver_InstanceAllowlist(t *testing.T) {
	server := newMetadataServer(t, http.StatusOK, "", func(clientID string) clientRegistrationRequest {
		return clientRegistrationRequest{RedirectURIs: []string{clientID + "/callback"}, TokenEndpointAuthMethod: "none"}
	})
	documentCache := gomap.NewCache[clientIDMetadataCacheIndex, string, *clientIDMetadataCacheEntry](
		context.Background(),
		[]clientIDMetadataCacheIndex{clientIDMetadataCacheIndexURL},
		cache.Config{MaxAge: time.Hour, LastUseAge: time.Hour},
	)
	resolver := newClientIDMetadataResolver(server.Client(), clientIDMetadataAllowlist{allowAny: true}, documentCache, time.Hour, time.Hour, nil)
	clientID := server.URL + testClientIDPath

	for name, instance := range map[string][]authz.MockContextInstanceOpts{
		"nothing allowed":     nil,
		"other url allowed":   {authz.WithMockClientIDMetadataDocumentAllowedURLs(server.URL + "/other/")},
		"longer url allowed":  {authz.WithMockClientIDMetadataDocumentAllowedURLs(clientID + "/")},
		"allow any url false": {authz.WithMockClientIDMetadataDocumentAllowAnyURL(false)},
	} {
		t.Run(name, func(t *testing.T) {
			ctx := authz.NewMockContext("instance", "org", "", instance...)
			_, err := resolver.ResolveClient(ctx, "instance", clientID)
			assertInvalidClient(t, err)
			assert.Equal(t, int32(0), server.hits.Load(), "a client_id the instance does not allow must not be fetched")
			_, found := documentCache.Get(ctx, clientIDMetadataCacheIndexURL, clientIDMetadataCacheKey("instance", clientID))
			assert.False(t, found, "a client_id the instance does not allow must not be cached")
		})
	}

	t.Run("instance allows the url", func(t *testing.T) {
		ctx := authz.NewMockContext("instance", "org", "", authz.WithMockClientIDMetadataDocumentAllowedURLs(server.URL+"/"))
		client, err := resolver.ResolveClient(ctx, "instance", clientID)
		require.NoError(t, err)
		assert.Equal(t, clientID, client.ClientID)
	})
}

// TestClientIDMetadataResolver_LoginV2BaseURI pins that a CIMD client uses the v2 login and,
// like a stored client, the login base URI the instance requires.
func TestClientIDMetadataResolver_LoginV2BaseURI(t *testing.T) {
	server := newMetadataServer(t, http.StatusOK, "", func(clientID string) clientRegistrationRequest {
		return clientRegistrationRequest{RedirectURIs: []string{clientID + "/callback"}, TokenEndpointAuthMethod: "none"}
	})
	resolver := newTestResolver(server.URL, server.Client(), noop.NewCache[clientIDMetadataCacheIndex, string, *clientIDMetadataCacheEntry]())
	baseURI := mustParseURL(t, "https://login.example.com/ui/v2/login")

	t.Run("login v2 required", func(t *testing.T) {
		ctx := authz.NewMockContext("instance", "org", "",
			authz.WithMockClientIDMetadataDocumentAllowAnyURL(true),
			authz.WithMockFeatures(feature.Features{LoginV2: feature.LoginV2{Required: true, BaseURI: baseURI}}),
		)
		client, err := resolver.ResolveClient(ctx, "instance", server.URL+testClientIDPath)
		require.NoError(t, err)
		assert.Equal(t, domain.LoginVersion2, client.LoginVersion)
		require.NotNil(t, client.LoginBaseURI)
		assert.Equal(t, baseURI.String(), (*url.URL)(client.LoginBaseURI).String())
	})
	t.Run("login v2 not required", func(t *testing.T) {
		client, err := resolver.ResolveClient(testResolverContext(), "instance", server.URL+testClientIDPath)
		require.NoError(t, err)
		assert.Equal(t, domain.LoginVersion2, client.LoginVersion)
		assert.Nil(t, client.LoginBaseURI)
	})
}

// TestClientIDMetadataCacheEntry_JSONRoundTrip pins that a resolved client survives the JSON
// round trip the PostgreSQL and Redis cache connectors apply, including the login base URI a
// required Login V2 adds. Without that, every cache read fails and the document is refetched.
func TestClientIDMetadataCacheEntry_JSONRoundTrip(t *testing.T) {
	baseURI := mustParseURL(t, "https://login.example.com/ui/v2/login")
	entry := &clientIDMetadataCacheEntry{
		Key: clientIDMetadataCacheKey("instance", "https://app.example.com/client"),
		Client: &query.OIDCClient{
			ClientID:     "https://app.example.com/client",
			RedirectURIs: []string{"https://app.example.com/callback"},
			LoginVersion: domain.LoginVersion2,
			LoginBaseURI: (*query.URL)(baseURI),
		},
		Expiry: time.Now().Add(time.Minute).UTC().Truncate(time.Second),
	}
	data, err := json.Marshal(entry)
	require.NoError(t, err)

	var got clientIDMetadataCacheEntry
	require.NoError(t, json.Unmarshal(data, &got))
	require.NotNil(t, got.Client)
	require.NotNil(t, got.Client.LoginBaseURI)
	assert.Equal(t, baseURI.String(), got.Client.LoginBaseURI.URL().String())
	assert.Equal(t, entry.Client.RedirectURIs, got.Client.RedirectURIs)
	assert.True(t, entry.Expiry.Equal(got.Expiry))
}

// TestClientIDMetadataResolver_Nil pins that a server built without a resolver treats client
// ID metadata documents as unsupported instead of dereferencing it.
func TestClientIDMetadataResolver_Nil(t *testing.T) {
	var resolver *clientIDMetadataResolver
	ctx := authz.NewMockContext("instance", "org", "",
		authz.WithMockClientIDMetadataDocument(true),
		authz.WithMockClientIDMetadataDocumentAllowAnyURL(true),
	)
	assert.False(t, resolver.Supported(ctx))
	assert.False(t, resolver.Handles(ctx, "https://app.example.com/client"))
}

// blockingCache is a document cache whose writes block until their context is done, like a
// connector whose backend does not answer.
type blockingCache struct {
	cache.Cache[clientIDMetadataCacheIndex, string, *clientIDMetadataCacheEntry]
	setReturned chan struct{}
}

func (c *blockingCache) Set(ctx context.Context, _ *clientIDMetadataCacheEntry) {
	<-ctx.Done()
	close(c.setReturned)
}

// TestClientIDMetadataResolver_ResolveTimeout pins that the shared resolution, which does not
// inherit the caller's cancellation, is still bounded as a whole: a cache write that never
// returns on its own is cancelled by the resolve timeout instead of holding the call open.
func TestClientIDMetadataResolver_ResolveTimeout(t *testing.T) {
	server := newMetadataServer(t, http.StatusOK, "", func(clientID string) clientRegistrationRequest {
		return clientRegistrationRequest{RedirectURIs: []string{clientID + "/callback"}, TokenEndpointAuthMethod: "none"}
	})
	documentCache := &blockingCache{
		Cache:       noop.NewCache[clientIDMetadataCacheIndex, string, *clientIDMetadataCacheEntry](),
		setReturned: make(chan struct{}),
	}
	resolver := newTestResolver(server.URL, server.Client(), documentCache)
	resolver.resolveTimeout = 200 * time.Millisecond

	resolved := make(chan error, 1)
	go func() {
		_, err := resolver.ResolveClient(testResolverContext(), "instance", server.URL+testClientIDPath)
		resolved <- err
	}()
	select {
	case err := <-resolved:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("the resolution was not bounded by the resolve timeout")
	}
	select {
	case <-documentCache.setReturned:
	default:
		t.Fatal("the cache write was not cancelled by the resolve timeout")
	}
}

// TestClientIDMetadataResolver_ClientIDMismatch pins that the document must name the URL it is
// served from as its client_id, compared as plain strings.
func TestClientIDMetadataResolver_ClientIDMismatch(t *testing.T) {
	ctx := testResolverContext()
	for name, documentClientID := range map[string]func(origin string) string{
		"missing":      func(string) string { return "" },
		"other client": func(origin string) string { return origin + "/other" },
		"other host":   func(string) string { return "https://app.example.com" + testClientIDPath },
		"scheme case": func(origin string) string {
			return strings.Replace(origin, "https://", "HTTPS://", 1) + testClientIDPath
		},
		"trailing slash": func(origin string) string { return origin + testClientIDPath + "/" },
	} {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(clientIDMetadataDocument{
					ClientID: documentClientID("https://" + r.Host),
					clientRegistrationRequest: clientRegistrationRequest{
						RedirectURIs:            []string{"https://" + r.Host + "/callback"},
						TokenEndpointAuthMethod: "none",
					},
				})
			}))
			defer server.Close()
			resolver := newTestResolver(server.URL, server.Client(), noop.NewCache[clientIDMetadataCacheIndex, string, *clientIDMetadataCacheEntry]())

			_, err := resolver.ResolveClient(ctx, "instance", server.URL+testClientIDPath)
			assertInvalidClient(t, err)
		})
	}
}

// TestClientIDMetadataResolver_NoRedirects pins that the document is only accepted from the
// client_id URL itself. A redirect is answered with invalid_client and not followed, so an
// allowed client_id cannot lead the fetch to a URL the allowlist does not allow.
func TestClientIDMetadataResolver_NoRedirects(t *testing.T) {
	ctx := testResolverContext()
	target := newMetadataServer(t, http.StatusOK, "", func(clientID string) clientRegistrationRequest {
		return clientRegistrationRequest{RedirectURIs: []string{clientID + "/callback"}, TokenEndpointAuthMethod: "none"}
	})
	for _, status := range []int{http.StatusMovedPermanently, http.StatusFound, http.StatusTemporaryRedirect, http.StatusPermanentRedirect} {
		t.Run(strconv.Itoa(status), func(t *testing.T) {
			redirector := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, target.URL+testClientIDPath, status)
			}))
			defer redirector.Close()
			resolver := newTestResolver(redirector.URL, redirector.Client(), noop.NewCache[clientIDMetadataCacheIndex, string, *clientIDMetadataCacheEntry]())

			_, err := resolver.ResolveClient(ctx, "instance", redirector.URL+testClientIDPath)
			assertInvalidClient(t, err)
			assert.Equal(t, int32(0), target.hits.Load(), "the redirect must not be followed")
		})
	}
}

// helpers

// TestClientIDMetadataResolver_CallerCancellation pins that the shared resolution does not
// inherit the cancellation of whichever caller happened to start it. Binding it to that one
// request would let a single client disconnect abort the fetch for every concurrent waiter and
// store the failure as a negative cache entry, blocking the client_id for everyone until it
// expired.
func TestClientIDMetadataResolver_CallerCancellation(t *testing.T) {
	documentCache := gomap.NewCache[clientIDMetadataCacheIndex, string, *clientIDMetadataCacheEntry](
		context.Background(),
		[]clientIDMetadataCacheIndex{clientIDMetadataCacheIndexURL},
		cache.Config{MaxAge: time.Hour, LastUseAge: time.Hour},
	)

	// The server blocks until released, so the caller is guaranteed to cancel while the shared
	// resolution is still in flight.
	release := make(chan struct{})
	reached := make(chan struct{})
	var once sync.Once
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		once.Do(func() { close(reached) })
		<-release
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(clientIDMetadataDocument{
			ClientID: "https://" + r.Host + r.URL.Path,
			clientRegistrationRequest: clientRegistrationRequest{
				RedirectURIs:            []string{"https://" + r.Host + "/callback"},
				TokenEndpointAuthMethod: "none",
			},
		})
	}))
	defer server.Close()
	resolver := newTestResolver(server.URL, server.Client(), documentCache)

	cancelCtx, cancel := context.WithCancel(testResolverContext())
	cancelled := make(chan error, 1)
	go func() {
		_, err := resolver.ResolveClient(cancelCtx, "instance", server.URL+testClientIDPath)
		cancelled <- err
	}()

	<-reached
	cancel()
	// The abandoning caller returns promptly with its own cancellation rather than blocking on
	// the shared fetch.
	select {
	case err := <-cancelled:
		assert.ErrorIs(t, err, context.Canceled)
	case <-time.After(10 * time.Second):
		t.Fatal("cancelled caller did not return")
	}

	close(release)

	// The resolution survived the cancellation, so a later caller resolves normally instead of
	// hitting a negatively cached failure.
	client, err := resolver.ResolveClient(testResolverContext(), "instance", server.URL+testClientIDPath)
	require.NoError(t, err)
	assert.Equal(t, server.URL+testClientIDPath, client.ClientID)
}

// testResolverContext returns the context of an instance that allows any client_id URL, so the
// system allowlist of the resolver under test decides.
func testResolverContext() context.Context {
	return authz.NewMockContext("instance", "org", "", authz.WithMockClientIDMetadataDocumentAllowAnyURL(true))
}

// testClientIDPath is the path of the client_id URLs the tests resolve against a test server: a
// client_id URL must contain a path.
const testClientIDPath = "/client"

// newTestResolver returns a resolver that allows every client_id URL under serverURL.
func newTestResolver(serverURL string, httpClient *http.Client, documentCache cache.Cache[clientIDMetadataCacheIndex, string, *clientIDMetadataCacheEntry]) *clientIDMetadataResolver {
	return newClientIDMetadataResolver(httpClient, clientIDMetadataAllowlist{urls: []string{serverURL + "/"}}, documentCache, time.Hour, time.Hour, nil)
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
		_ = json.NewEncoder(w).Encode(clientIDMetadataDocument{
			ClientID:                  "https://" + r.Host + r.URL.Path,
			clientRegistrationRequest: build("https://" + r.Host),
		})
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
