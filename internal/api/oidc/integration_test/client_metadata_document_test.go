//go:build integration

package oidc_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/admin"
	"github.com/zitadel/zitadel/pkg/grpc/settings/v2"
	settings_v2beta "github.com/zitadel/zitadel/pkg/grpc/settings/v2beta"
)

// allowedClientIDMetadataAddress is the address the integration API allows client id metadata
// documents to be fetched from, see OIDC.ClientIDMetadataDocument.AllowedURLs in
// apps/api/test-integration-api.yaml.
const allowedClientIDMetadataAddress = "127.0.0.1:8091"

// TestServer_ClientIDMetadataDocument covers Client ID Metadata Documents (CIMD) with the
// setting enabled. The resolver itself (document fetch, origin match, SSRF handling,
// allowlist, public-only enforcement and caching) is covered end to end by the unit tests. A full
// authorize to token flow against a live document is not exercised here for one reason only:
// CIMD requires https and the server's outbound HTTP client validates certificates against the
// system roots, so it does not trust the self-signed certificate of an in-process test server
// (the empty integration denylist and the http issuer are not the blockers). These tests
// therefore assert that support is advertised and that the interception actually runs in the
// live server by observing the outbound fetch.
func TestServer_ClientIDMetadataDocument(t *testing.T) {
	enableClientIDMetadataDocument(t, CTXIAM, Instance)

	issuer := Instance.OIDCIssuer()

	t.Run("discovery advertises client id metadata document support", func(t *testing.T) {
		discovery := fetchDiscoveryRaw(t, issuer)
		assert.Equal(t, true, discovery["client_id_metadata_document_supported"])
	})

	t.Run("an allowed client_id triggers an outbound document fetch", func(t *testing.T) {
		// With the setting on, an allowed client_id is intercepted and the resolver fetches the
		// document, so the server observes a connection (the TLS handshake then fails because
		// the cert is not trusted, which is why the authorization itself errors). A database
		// lookup would never produce an outbound connection, so this distinguishes "the resolver
		// ran" from a plain unknown-client error.
		listener, err := net.Listen("tcp", allowedClientIDMetadataAddress)
		require.NoError(t, err)
		server, connections := newClientIDMetadataDocumentServer(listener)
		defer server.Close()

		_, _, err = Instance.CreateOIDCAuthRequest(CTX, server.URL+"/client", Instance.Users.Get(integration.UserTypeLogin).ID, redirectURI, oidc.ScopeOpenID)
		require.Error(t, err)
		assert.GreaterOrEqual(t, connections.Load(), int32(1), "the resolver must have attempted to fetch the metadata document")
	})

	t.Run("a client_id the system does not allow is not fetched, even if the instance allows any url", func(t *testing.T) {
		// The shared instance allows any URL, so the system allowlist alone keeps this address
		// out.
		listener, err := net.Listen("tcp", "127.0.0.1:0")
		require.NoError(t, err)
		server, connections := newClientIDMetadataDocumentServer(listener)
		defer server.Close()

		_, _, err = Instance.CreateOIDCAuthRequest(CTX, server.URL+"/client", Instance.Users.Get(integration.UserTypeLogin).ID, redirectURI, oidc.ScopeOpenID)
		require.Error(t, err)
		assert.Equal(t, int32(0), connections.Load(), "a client_id the system does not allow must not be fetched")
	})

	// Non-URL client_ids keep going through the regular database lookup unchanged: the rest of
	// the OIDC integration suite runs against this same instance with the setting enabled, so a
	// regression on the normal client path would surface there.
}

// TestServer_ClientIDMetadataDocument_disabled verifies that without the setting CIMD is
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

// newClientIDMetadataDocumentServer serves a metadata document over TLS on listener and counts
// the inbound connections.
func newClientIDMetadataDocumentServer(listener net.Listener) (*httptest.Server, *atomic.Int32) {
	connections := new(atomic.Int32)
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"client_id":                  "https://" + r.Host + r.URL.Path,
			"redirect_uris":              []string{redirectURI},
			"token_endpoint_auth_method": "none",
		})
	}))
	_ = server.Listener.Close()
	server.Listener = listener
	server.Config.ConnState = func(_ net.Conn, state http.ConnState) {
		if state == http.StateNew {
			connections.Add(1)
		}
	}
	server.StartTLS()
	return server, connections
}

// TestServer_ClientIDMetadataDocument_instanceAllowedURLs manages the client_id URLs an instance
// allows through the API, without any change to the runtime configuration, and asserts on the
// live server which ones are fetched.
func TestServer_ClientIDMetadataDocument_instanceAllowedURLs(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	instance := integration.NewInstance(ctx)
	iamCTX := instance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
	issuer := instance.OIDCIssuer()
	allowedURL := "https://" + allowedClientIDMetadataAddress + "/"

	// fetchCount starts a document server on the address the system allows, authorizes with a
	// client_id on it and returns how many connections the server saw. Every call uses its own
	// client_id: a failed resolution is cached for a short window, and the fetches here always
	// fail on the untrusted test certificate, so a reused client_id would be answered from that
	// cache instead of being fetched.
	var fetches int
	fetchCount := func(t *testing.T) int32 {
		t.Helper()
		listener, err := net.Listen("tcp", allowedClientIDMetadataAddress)
		require.NoError(t, err)
		server, connections := newClientIDMetadataDocumentServer(listener)
		defer server.Close()
		fetches++
		clientID := fmt.Sprintf("%s/client-%d", server.URL, fetches)
		_, _, err = instance.CreateOIDCAuthRequest(ctx, clientID, instance.Users.Get(integration.UserTypeLogin).ID, redirectURI, oidc.ScopeOpenID)
		require.Error(t, err)
		return connections.Load()
	}

	setClientIDMetadataDocumentSettings(t, iamCTX, instance, &settings.ClientIDMetadataDocumentSettings{Enabled: true})

	t.Run("enabled without allowed urls, nothing is fetched or advertised", func(t *testing.T) {
		waitForClientIDMetadataDocumentSupport(t, ctx, issuer, false)
		assert.Equal(t, int32(0), fetchCount(t))
	})

	t.Run("an invalid url cannot be added", func(t *testing.T) {
		_, err := instance.Client.SettingsV2.AddClientIDMetadataDocumentAllowedURL(iamCTX, &settings.AddClientIDMetadataDocumentAllowedURLRequest{
			Url: "https://" + allowedClientIDMetadataAddress + "/oauth/../",
		})
		assert.Equal(t, codes.InvalidArgument, status.Code(err))
	})

	t.Run("setting duplicate or overlong urls is rejected", func(t *testing.T) {
		for name, urls := range map[string][]string{
			"duplicate": {allowedURL, allowedURL},
			"overlong":  {allowedURL + strings.Repeat("a", 2048)},
		} {
			t.Run(name, func(t *testing.T) {
				_, err := instance.Client.SettingsV2.SetSecuritySettings(iamCTX, &settings.SetSecuritySettingsRequest{
					ClientIdMetadataDocument: &settings.ClientIDMetadataDocumentSettings{Enabled: true, AllowedUrls: urls},
				})
				assert.Equal(t, codes.InvalidArgument, status.Code(err))
			})
		}
	})

	t.Run("an added url is fetched", func(t *testing.T) {
		resp, err := instance.Client.SettingsV2.AddClientIDMetadataDocumentAllowedURL(iamCTX, &settings.AddClientIDMetadataDocumentAllowedURLRequest{Url: allowedURL})
		require.NoError(t, err)
		assert.NotEmpty(t, resp.GetCreationDate())

		waitForClientIDMetadataDocumentSupport(t, ctx, issuer, true)
		assert.GreaterOrEqual(t, fetchCount(t), int32(1))

		current, err := instance.Client.SettingsV2.GetSecuritySettings(iamCTX, &settings.GetSecuritySettingsRequest{})
		require.NoError(t, err)
		assert.Equal(t, []string{allowedURL}, current.GetSettings().GetClientIdMetadataDocument().GetAllowedUrls())
	})

	t.Run("an allowed url cannot be added twice", func(t *testing.T) {
		_, err := instance.Client.SettingsV2.AddClientIDMetadataDocumentAllowedURL(iamCTX, &settings.AddClientIDMetadataDocumentAllowedURLRequest{Url: allowedURL})
		assert.Equal(t, codes.AlreadyExists, status.Code(err))
	})

	t.Run("a removed url is no longer fetched", func(t *testing.T) {
		resp, err := instance.Client.SettingsV2.RemoveClientIDMetadataDocumentAllowedURL(iamCTX, &settings.RemoveClientIDMetadataDocumentAllowedURLRequest{Url: allowedURL})
		require.NoError(t, err)
		assert.NotEmpty(t, resp.GetDeletionDate())

		waitForClientIDMetadataDocumentSupport(t, ctx, issuer, false)
		assert.Equal(t, int32(0), fetchCount(t))
	})

	t.Run("removing a url that is not allowed is not an error", func(t *testing.T) {
		resp, err := instance.Client.SettingsV2.RemoveClientIDMetadataDocumentAllowedURL(iamCTX, &settings.RemoveClientIDMetadataDocumentAllowedURLRequest{Url: allowedURL})
		require.NoError(t, err)
		assert.Nil(t, resp.GetDeletionDate())
	})

	t.Run("concurrent adds of the same url allow it once", func(t *testing.T) {
		const concurrentURL = "https://127.0.0.1:8091/concurrent"
		const adds = 5
		codesByAdd := make([]codes.Code, adds)
		var wg sync.WaitGroup
		for i := range adds {
			wg.Go(func() {
				_, err := instance.Client.SettingsV2.AddClientIDMetadataDocumentAllowedURL(iamCTX, &settings.AddClientIDMetadataDocumentAllowedURLRequest{Url: concurrentURL})
				codesByAdd[i] = status.Code(err)
			})
		}
		wg.Wait()
		assert.Equal(t, 1, countCode(codesByAdd, codes.OK), "exactly one add succeeds")
		assert.Equal(t, adds-1, countCode(codesByAdd, codes.AlreadyExists), "every other add reports the url as already allowed")

		_, err := instance.Client.SettingsV2.RemoveClientIDMetadataDocumentAllowedURL(iamCTX, &settings.RemoveClientIDMetadataDocumentAllowedURLRequest{Url: concurrentURL})
		require.NoError(t, err)
	})

	t.Run("concurrent replacements leave only the stored urls taken", func(t *testing.T) {
		const sets = 8
		current, err := instance.Client.SettingsV2.GetSecuritySettings(iamCTX, &settings.GetSecuritySettingsRequest{})
		require.NoError(t, err)
		replacementURLs := make([]string, sets)
		setErrs := make([]error, sets)
		var wg sync.WaitGroup
		for i := range sets {
			replacementURLs[i] = fmt.Sprintf("https://127.0.0.1:8091/replacement-%d", i)
			wg.Go(func() {
				_, setErrs[i] = instance.Client.SettingsV2.SetSecuritySettings(iamCTX, &settings.SetSecuritySettingsRequest{
					EmbeddedIframe:            current.GetSettings().GetEmbeddedIframe(),
					EnableImpersonation:       current.GetSettings().GetEnableImpersonation(),
					DynamicClientRegistration: current.GetSettings().GetDynamicClientRegistration(),
					ClientIdMetadataDocument:  &settings.ClientIDMetadataDocumentSettings{Enabled: true, AllowedUrls: []string{replacementURLs[i]}},
				})
			})
		}
		wg.Wait()
		for _, err := range setErrs {
			require.NoError(t, err)
		}

		// Exactly the url of the replacement stored last is still taken; every other url can be
		// added again, so no replacement left a constraint behind.
		addCodes := make([]codes.Code, sets)
		for i, replacementURL := range replacementURLs {
			_, err := instance.Client.SettingsV2.AddClientIDMetadataDocumentAllowedURL(iamCTX, &settings.AddClientIDMetadataDocumentAllowedURLRequest{Url: replacementURL})
			addCodes[i] = status.Code(err)
		}
		assert.Equal(t, 1, countCode(addCodes, codes.AlreadyExists), "only the stored url is still taken: %v", addCodes)
		assert.Equal(t, sets-1, countCode(addCodes, codes.OK), "every replaced url is released: %v", addCodes)

		setClientIDMetadataDocumentSettings(t, iamCTX, instance, &settings.ClientIDMetadataDocumentSettings{Enabled: true, AllowedUrls: []string{}})
	})

	t.Run("allow any url fetches every url the system allows", func(t *testing.T) {
		setClientIDMetadataDocumentSettings(t, iamCTX, instance, &settings.ClientIDMetadataDocumentSettings{Enabled: true, AllowAnyUrl: true})
		waitForClientIDMetadataDocumentSupport(t, ctx, issuer, true)
		assert.GreaterOrEqual(t, fetchCount(t), int32(1))
	})
}

// TestServer_ClientIDMetadataDocument_legacySecuritySettings verifies that setting the security
// settings through the admin v1 and settings v2beta APIs, which do not know client ID metadata
// documents or dynamic client registration, leaves both untouched.
func TestServer_ClientIDMetadataDocument_legacySecuritySettings(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	instance := integration.NewInstance(ctx)
	iamCTX := instance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
	issuer := instance.OIDCIssuer()

	enableDynamicClientRegistration(t, iamCTX, instance, true)
	setClientIDMetadataDocumentSettings(t, iamCTX, instance, &settings.ClientIDMetadataDocumentSettings{
		Enabled:     true,
		AllowedUrls: []string{"https://" + allowedClientIDMetadataAddress + "/"},
		AllowAnyUrl: true,
	})
	waitForClientIDMetadataDocumentSupport(t, ctx, issuer, true)

	assertUnchanged := func(t *testing.T) {
		t.Helper()
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
		require.EventuallyWithT(t, func(tt *assert.CollectT) {
			current, err := instance.Client.SettingsV2.GetSecuritySettings(iamCTX, &settings.GetSecuritySettingsRequest{})
			if !assert.NoError(tt, err) {
				return
			}
			assert.True(tt, current.GetSettings().GetEmbeddedIframe().GetEnabled(), "the legacy change must be applied")
			assert.Equal(tt, &settings.ClientIDMetadataDocumentSettings{
				Enabled:     true,
				AllowedUrls: []string{"https://" + allowedClientIDMetadataAddress + "/"},
				AllowAnyUrl: true,
			}, current.GetSettings().GetClientIdMetadataDocument())
			assert.Equal(tt, &settings.DynamicClientRegistrationSettings{Enabled: true, AllowUnauthenticated: true}, current.GetSettings().GetDynamicClientRegistration())
		}, retryDuration, tick, "security settings not as expected")
		waitForClientIDMetadataDocumentSupport(t, ctx, issuer, true)
	}

	t.Run("admin v1", func(t *testing.T) {
		_, err := instance.Client.Admin.SetSecurityPolicy(iamCTX, &admin.SetSecurityPolicyRequest{
			EnableIframeEmbedding: true,
			AllowedOrigins:        []string{"https://origin.example.com"},
		})
		require.NoError(t, err)
		assertUnchanged(t)
	})

	t.Run("settings v2beta", func(t *testing.T) {
		_, err := instance.Client.SettingsV2beta.SetSecuritySettings(iamCTX, &settings_v2beta.SetSecuritySettingsRequest{
			EmbeddedIframe: &settings_v2beta.EmbeddedIframeSettings{
				Enabled:        true,
				AllowedOrigins: []string{"https://other-origin.example.com"},
			},
		})
		require.NoError(t, err)
		assertUnchanged(t)
	})
}

// enableClientIDMetadataDocument turns CIMD on through the instance's security settings, lets
// the instance resolve every client_id URL the system allows and waits until the change is
// observable, i.e. until the projection behind the cached instance has caught up and discovery
// advertises the support.
//
// SetSecuritySettings replaces the whole policy, so the current settings are read first and
// carried over: several tests share an instance and the dynamic client registration tests
// enable their own switch on it. It is safe to call more than once: a write that changes
// nothing is answered with a failed precondition ("Errors.NoChangesFound"), which here means
// the desired state is already reached rather than a failure.
func enableClientIDMetadataDocument(t *testing.T, ctx context.Context, instance *integration.Instance) {
	t.Helper()
	current, err := instance.Client.SettingsV2.GetSecuritySettings(ctx, &settings.GetSecuritySettingsRequest{})
	require.NoError(t, err)

	_, err = instance.Client.SettingsV2.SetSecuritySettings(ctx, &settings.SetSecuritySettingsRequest{
		EmbeddedIframe:            current.GetSettings().GetEmbeddedIframe(),
		EnableImpersonation:       current.GetSettings().GetEnableImpersonation(),
		DynamicClientRegistration: current.GetSettings().GetDynamicClientRegistration(),
		ClientIdMetadataDocument: &settings.ClientIDMetadataDocumentSettings{
			Enabled:     true,
			AllowAnyUrl: true,
		},
	})
	if status.Code(err) != codes.FailedPrecondition {
		require.NoError(t, err)
	}
	waitForClientIDMetadataDocumentSupport(t, ctx, instance.OIDCIssuer(), true)
}

// setClientIDMetadataDocumentSettings sets the client id metadata document settings of the
// instance and carries the other security settings over, as SetSecuritySettings replaces the
// whole policy.
func setClientIDMetadataDocumentSettings(t *testing.T, ctx context.Context, instance *integration.Instance, cimd *settings.ClientIDMetadataDocumentSettings) {
	t.Helper()
	current, err := instance.Client.SettingsV2.GetSecuritySettings(ctx, &settings.GetSecuritySettingsRequest{})
	require.NoError(t, err)
	if cimd.AllowedUrls == nil {
		cimd.AllowedUrls = current.GetSettings().GetClientIdMetadataDocument().GetAllowedUrls()
	}
	_, err = instance.Client.SettingsV2.SetSecuritySettings(ctx, &settings.SetSecuritySettingsRequest{
		EmbeddedIframe:            current.GetSettings().GetEmbeddedIframe(),
		EnableImpersonation:       current.GetSettings().GetEnableImpersonation(),
		DynamicClientRegistration: current.GetSettings().GetDynamicClientRegistration(),
		ClientIdMetadataDocument:  cimd,
	})
	if status.Code(err) != codes.FailedPrecondition {
		require.NoError(t, err)
	}
}

// waitForClientIDMetadataDocumentSupport waits until discovery advertises client id metadata
// document support as wanted, i.e. until the projection behind the cached instance has caught
// up with the last settings change.
func waitForClientIDMetadataDocumentSupport(t *testing.T, ctx context.Context, issuer string, want bool) {
	t.Helper()
	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
	require.EventuallyWithT(t, func(tt *assert.CollectT) {
		resp, err := http.Get(issuer + "/.well-known/openid-configuration")
		if !assert.NoError(tt, err) {
			return
		}
		defer resp.Body.Close()
		var discovery map[string]any
		if !assert.NoError(tt, json.NewDecoder(resp.Body).Decode(&discovery)) {
			return
		}
		supported, _ := discovery["client_id_metadata_document_supported"].(bool)
		assert.Equal(tt, want, supported)
	}, retryDuration, tick, "client id metadata document support not as expected")
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

// countCode returns how many of got equal want.
func countCode(got []codes.Code, want codes.Code) int {
	n := 0
	for _, code := range got {
		if code == want {
			n++
		}
	}
	return n
}
