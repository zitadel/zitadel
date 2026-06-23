//go:build integration

package oidc_test

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v3/pkg/oidc"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/feature/v2"
)

// TestServer_ClientIDMetadataDocument covers the Client ID Metadata Document (CIMD) feature
// with the feature enabled. The resolver itself (document fetch, origin match, SSRF handling,
// public-only enforcement and caching) is covered end to end by the unit tests. A full
// authorize to token flow against a live document is not exercised here for one reason only:
// CIMD requires https and the server's outbound HTTP client validates certificates against the
// system roots, so it does not trust the self-signed certificate of an in-process test server
// (the empty integration denylist and the http issuer are not the blockers). These tests
// therefore assert that the feature is advertised and that the interception actually runs in
// the live server by observing the outbound fetch.
func TestServer_ClientIDMetadataDocument(t *testing.T) {
	integration.EnsureInstanceFeature(t, CTX, Instance,
		&feature.SetInstanceFeaturesRequest{
			OidcClientIdMetadataDocument: gu.Ptr(true),
		},
		func(tt *assert.CollectT, got *feature.GetInstanceFeaturesResponse) {
			assert.True(tt, got.GetOidcClientIdMetadataDocument().GetEnabled())
		},
	)

	issuer := Instance.OIDCIssuer()

	t.Run("discovery advertises client id metadata document support", func(t *testing.T) {
		discovery := fetchDiscoveryRaw(t, issuer)
		assert.Equal(t, true, discovery["client_id_metadata_document_supported"])
	})

	t.Run("an https client_id triggers an outbound document fetch", func(t *testing.T) {
		// Serve a metadata document over TLS and count inbound connections. With the feature
		// on, an https client_id is intercepted and the resolver fetches the document, so the
		// server observes a connection (the TLS handshake then fails because the cert is not
		// trusted, which is why the authorization itself errors). A database lookup would never
		// produce an outbound connection, so this distinguishes "the resolver ran" from a plain
		// unknown-client error.
		var connections atomic.Int32
		server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"redirect_uris":              []string{redirectURI},
				"token_endpoint_auth_method": "none",
			})
		}))
		server.Config.ConnState = func(_ net.Conn, state http.ConnState) {
			if state == http.StateNew {
				connections.Add(1)
			}
		}
		server.StartTLS()
		defer server.Close()

		_, _, err := Instance.CreateOIDCAuthRequest(CTX, server.URL, Instance.Users.Get(integration.UserTypeLogin).ID, redirectURI, oidc.ScopeOpenID)
		require.Error(t, err)
		assert.GreaterOrEqual(t, connections.Load(), int32(1), "the resolver must have attempted to fetch the metadata document")
	})

	// Non-URL client_ids keep going through the regular database lookup unchanged: the rest of
	// the OIDC integration suite runs against this same instance with the feature enabled, so a
	// regression on the normal client path would surface there.
}

// TestServer_ClientIDMetadataDocument_disabled verifies that without the feature CIMD is
// neither advertised nor active: an https-url client_id is treated as a regular (unknown)
// client and the discovery document does not carry the support metadata.
func TestServer_ClientIDMetadataDocument_disabled(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	instance := integration.NewInstance(ctx)
	issuer := instance.OIDCIssuer()

	discovery := fetchDiscoveryRaw(t, issuer)
	_, advertised := discovery["client_id_metadata_document_supported"]
	assert.False(t, advertised, "CIMD support must not be advertised when the feature is disabled")

	_, _, err := instance.CreateOIDCAuthRequest(ctx, "https://app.example.com/client", instance.Users.Get(integration.UserTypeLogin).ID, redirectURI, oidc.ScopeOpenID)
	require.Error(t, err, "an https client_id must not resolve when the feature is disabled")
}

func fetchDiscoveryRaw(t testing.TB, issuer string) map[string]any {
	t.Helper()
	resp, err := http.Get(issuer + "/.well-known/openid-configuration")
	require.NoError(t, err)
	defer resp.Body.Close()
	var discovery map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&discovery))
	return discovery
}
