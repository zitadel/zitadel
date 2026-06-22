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

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v3/pkg/oidc"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/feature/v2"
	oidc_pb "github.com/zitadel/zitadel/pkg/grpc/oidc/v2"
)

// TestServer_DynamicClientRegistration covers the OAuth 2.0 Dynamic Client Registration
// endpoint (RFC 7591) with the feature enabled and open registration allowed (the
// integration server sets OIDC.DynamicClientRegistration.AllowUnauthenticated).
func TestServer_DynamicClientRegistration(t *testing.T) {
	integration.EnsureInstanceFeature(t, CTX, Instance,
		&feature.SetInstanceFeaturesRequest{
			OidcDynamicClientRegistration: gu.Ptr(true),
		},
		func(tt *assert.CollectT, got *feature.GetInstanceFeaturesResponse) {
			assert.True(tt, got.GetOidcDynamicClientRegistration().GetEnabled())
		},
	)

	issuer := Instance.OIDCIssuer()

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

	t.Run("missing redirect uri is rejected", func(t *testing.T) {
		status, body := registerDynamicClient(t, issuer, "", map[string]any{
			"client_name": "no redirects",
		})
		assert.Equal(t, http.StatusBadRequest, status)
		assert.Equal(t, "invalid_redirect_uri", body["error"])
	})
}

// TestServer_DynamicClientRegistration_disabled verifies that without the feature the
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
