//go:build integration

package oidc_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v3/pkg/oidc"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/integration"
	oidc_pb "github.com/zitadel/zitadel/pkg/grpc/oidc/v2"
	"github.com/zitadel/zitadel/pkg/grpc/settings/v2"
)

// TestServer_DynamicClientRegistration covers the OAuth 2.0 Dynamic Client Registration
// endpoint (RFC 7591) in open mode, where no access token is required.
func TestServer_DynamicClientRegistration(t *testing.T) {
	issuer := Instance.OIDCIssuer()
	enableDynamicClientRegistration(t, CTXIAM, Instance, true)

	t.Run("discovery advertises the registration endpoint", func(t *testing.T) {
		discovery := fetchDiscovery(t, issuer)
		assert.Equal(t, issuer+"/oauth/v2/register", discovery.RegistrationEndpoint)
	})

	t.Run("public client without secret, usable for authorization code with PKCE", func(t *testing.T) {
		status, body := registerDynamicClient(t, issuer, "", map[string]any{
			"client_name":                "integration public client",
			"redirect_uris":              []string{redirectURI},
			"token_endpoint_auth_method": "none",
		})
		require.Equal(t, http.StatusCreated, status)
		clientID, _ := body["client_id"].(string)
		require.NotEmpty(t, clientID)
		assert.Empty(t, body["client_secret"])

		// The returned client_id must be usable for an authorization code + PKCE flow.
		// The client and its dedicated project were just created, so wait for the
		// projection to catch up before starting the authorization flow.
		var authRequestID string
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		require.EventuallyWithT(t, func(tt *assert.CollectT) {
			_, id, err := Instance.CreateOIDCAuthRequest(CTX, clientID, Instance.Users.Get(integration.UserTypeLogin).ID, redirectURI, oidc.ScopeOpenID)
			assert.NoError(tt, err)
			authRequestID = id
		}, retryDuration, tick, "registered client not usable for authorization")

		sessionID, sessionToken, _, _ := Instance.CreateVerifiedWebAuthNSession(t, CTXLOGIN, User.GetUserId())
		linkResp, err := Instance.Client.OIDCv2.CreateCallback(CTXLOGIN, &oidc_pb.CreateCallbackRequest{
			AuthRequestId: authRequestID,
			CallbackKind: &oidc_pb.CreateCallbackRequest_Session{
				Session: &oidc_pb.Session{
					SessionId:    sessionID,
					SessionToken: sessionToken,
				},
			},
		})
		require.NoError(t, err)

		code := assertCodeResponse(t, linkResp.GetCallbackUrl())
		tokens, err := exchangeTokens(t, Instance, clientID, code, redirectURI)
		require.NoError(t, err)
		assert.NotEmpty(t, tokens.AccessToken)

		// An Authorization header is always treated as a registration token, even with open
		// registration enabled: the request must not fall back to anonymous registration.
		// This token belongs to an ordinary user without project.app.register_dynamic, so
		// it is rejected rather than silently registering a client.
		tokenStatus, tokenBody := registerDynamicClient(t, issuer, tokens.AccessToken, map[string]any{
			"client_name":   "integration token-mode client",
			"redirect_uris": []string{redirectURI},
		})
		require.Equal(t, http.StatusForbidden, tokenStatus)
		assert.Equal(t, "insufficient_scope", tokenBody["error"])
		assert.Empty(t, tokenBody["client_id"])
	})

	t.Run("confidential client returns a secret", func(t *testing.T) {
		status, body := registerDynamicClient(t, issuer, "", map[string]any{
			"client_name":                "integration confidential client",
			"redirect_uris":              []string{redirectURI},
			"token_endpoint_auth_method": "client_secret_basic",
		})
		require.Equal(t, http.StatusCreated, status)
		assert.NotEmpty(t, body["client_id"])
		assert.NotEmpty(t, body["client_secret"])
	})

	t.Run("invalid redirect uri is rejected", func(t *testing.T) {
		status, body := registerDynamicClient(t, issuer, "", map[string]any{
			"redirect_uris":    []string{"myapp://callback"},
			"application_type": "web",
		})
		assert.Equal(t, http.StatusBadRequest, status)
		assert.Equal(t, "invalid_redirect_uri", body["error"])
	})

	t.Run("malformed request body is rejected", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, issuer+"/oauth/v2/register", bytes.NewReader([]byte("{not valid json")))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("missing redirect uri is rejected", func(t *testing.T) {
		status, body := registerDynamicClient(t, issuer, "", map[string]any{
			"client_name": "no redirects",
		})
		assert.Equal(t, http.StatusBadRequest, status)
		assert.Equal(t, "invalid_redirect_uri", body["error"])
	})
}

// TestServer_DynamicClientRegistration_tokenMode covers the default mode, where a caller
// must present an access token whose user holds project.app.register_dynamic in the token's
// organization.
func TestServer_DynamicClientRegistration_tokenMode(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	instance := integration.NewInstance(ctx)
	iamCTX := instance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
	issuer := instance.OIDCIssuer()

	enableDynamicClientRegistration(t, iamCTX, instance, false)

	metadata := map[string]any{
		"client_name":                "token mode client",
		"redirect_uris":              []string{redirectURI},
		"token_endpoint_auth_method": "none",
	}

	t.Run("without a token registration is unauthorized", func(t *testing.T) {
		status, body := registerDynamicClient(t, issuer, "", metadata)
		assert.Equal(t, http.StatusUnauthorized, status)
		assert.Equal(t, "invalid_token", body["error"])
	})

	t.Run("a token without the permission is forbidden", func(t *testing.T) {
		_, pat, err := instance.CreateMachineUserPATWithMembership(iamCTX)
		require.NoError(t, err)

		status, body := registerDynamicClient(t, issuer, pat, metadata)
		assert.Equal(t, http.StatusForbidden, status)
		assert.Equal(t, "insufficient_scope", body["error"])
	})

	// The dedicated least-privilege role carries project.app.register_dynamic and nothing
	// else, so it is enough on its own to register a client.
	t.Run("a token with ORG_DYNAMIC_CLIENT_REGISTRAR may register", func(t *testing.T) {
		_, pat, err := instance.CreateMachineUserPATWithMembership(iamCTX, domain.RoleOrgDynamicClientRegistrar)
		require.NoError(t, err)

		var body map[string]any
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
		require.EventuallyWithT(t, func(tt *assert.CollectT) {
			var status int
			status, body = registerDynamicClient(t, issuer, pat, metadata)
			assert.Equal(tt, http.StatusCreated, status)
		}, retryDuration, tick, "membership not effective for registration")
		assert.NotEmpty(t, body["client_id"])
	})
}

// TestServer_DynamicClientRegistration_disabled verifies that with the setting off the
// endpoint is neither advertised nor served.
func TestServer_DynamicClientRegistration_disabled(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	instance := integration.NewInstance(ctx)
	issuer := instance.OIDCIssuer()

	discovery := fetchDiscovery(t, issuer)
	assert.Empty(t, discovery.RegistrationEndpoint)

	status, _ := registerDynamicClient(t, issuer, "", map[string]any{
		"redirect_uris": []string{redirectURI},
	})
	assert.Equal(t, http.StatusNotFound, status)
}

// enableDynamicClientRegistration turns the endpoint on through the instance's security
// settings and waits until the change is observable, i.e. until the projection behind the
// cached instance has caught up and discovery advertises the endpoint.
func enableDynamicClientRegistration(t *testing.T, ctx context.Context, instance *integration.Instance, allowUnauthenticated bool) {
	t.Helper()
	_, err := instance.Client.SettingsV2.SetSecuritySettings(ctx, &settings.SetSecuritySettingsRequest{
		DynamicClientRegistration: &settings.DynamicClientRegistrationSettings{
			Enabled:              true,
			AllowUnauthenticated: allowUnauthenticated,
		},
	})
	require.NoError(t, err)

	issuer := instance.OIDCIssuer()
	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
	require.EventuallyWithT(t, func(tt *assert.CollectT) {
		resp, err := http.Get(issuer + "/.well-known/openid-configuration")
		if !assert.NoError(tt, err) {
			return
		}
		defer resp.Body.Close()
		var discovery oidc.DiscoveryConfiguration
		if !assert.NoError(tt, json.NewDecoder(resp.Body).Decode(&discovery)) {
			return
		}
		assert.Equal(tt, issuer+"/oauth/v2/register", discovery.RegistrationEndpoint)
	}, retryDuration, tick, "registration endpoint not advertised")
}

func registerDynamicClient(t testing.TB, issuer, accessToken string, metadata map[string]any) (int, map[string]any) {
	t.Helper()
	raw, err := json.Marshal(metadata)
	require.NoError(t, err)
	req, err := http.NewRequest(http.MethodPost, issuer+"/oauth/v2/register", bytes.NewReader(raw))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()
	var body map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&body)
	return resp.StatusCode, body
}

func fetchDiscovery(t testing.TB, issuer string) oidc.DiscoveryConfiguration {
	t.Helper()
	resp, err := http.Get(issuer + "/.well-known/openid-configuration")
	require.NoError(t, err)
	defer resp.Body.Close()
	var discovery oidc.DiscoveryConfiguration
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&discovery))
	return discovery
}
