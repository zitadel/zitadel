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

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/feature/v2"
)

// TestServer_DynamicClientRegistration_Management covers the OAuth 2.0 Dynamic Client
// Registration Management Protocol (RFC 7592): read, update and delete of a dynamically
// registered client, authenticated with the registration access token issued at registration.
func TestServer_DynamicClientRegistration_Management(t *testing.T) {
	integration.EnsureInstanceFeature(t, CTX, Instance,
		&feature.SetInstanceFeaturesRequest{
			OidcDynamicClientRegistration: gu.Ptr(true),
		},
		func(tt *assert.CollectT, got *feature.GetInstanceFeaturesResponse) {
			assert.True(tt, got.GetOidcDynamicClientRegistration().GetEnabled())
		},
	)

	issuer := Instance.OIDCIssuer()

	status, reg := registerDynamicClient(t, issuer, "", map[string]any{
		"client_name":                "lifecycle client",
		"redirect_uris":              []string{redirectURI},
		"token_endpoint_auth_method": "none",
	})
	require.Equal(t, http.StatusCreated, status)
	clientID, _ := reg["client_id"].(string)
	require.NotEmpty(t, clientID)
	token, _ := reg["registration_access_token"].(string)
	require.NotEmpty(t, token)
	clientURI, _ := reg["registration_client_uri"].(string)
	require.Equal(t, issuer+"/oauth/v2/register/"+clientID, clientURI)

	// The client and its dedicated project were just created, so wait for the projection to
	// catch up before the first read.
	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
	require.EventuallyWithT(t, func(tt *assert.CollectT) {
		getStatus, body := manageDynamicClient(t, http.MethodGet, clientURI, token, nil)
		assert.Equal(tt, http.StatusOK, getStatus)
		assert.Equal(tt, clientID, body["client_id"])
		assert.Equal(tt, []any{redirectURI}, body["redirect_uris"])
		// A read does not rotate the token, so it is not returned again.
		assert.Empty(tt, body["registration_access_token"])
	}, retryDuration, tick, "registered client not readable")

	t.Run("missing or wrong registration access token is unauthorized", func(t *testing.T) {
		noToken, _ := manageDynamicClient(t, http.MethodGet, clientURI, "", nil)
		assert.Equal(t, http.StatusUnauthorized, noToken)

		wrongToken, _ := manageDynamicClient(t, http.MethodGet, clientURI, token+"tampered", nil)
		assert.Equal(t, http.StatusUnauthorized, wrongToken)
	})

	var rotatedToken string
	t.Run("update replaces metadata and rotates the token", func(t *testing.T) {
		putStatus, body := manageDynamicClient(t, http.MethodPut, clientURI, token, map[string]any{
			"redirect_uris":              []string{redirectURI + "/updated"},
			"token_endpoint_auth_method": "none",
		})
		require.Equal(t, http.StatusOK, putStatus)
		assert.Equal(t, clientID, body["client_id"])
		assert.Equal(t, []any{redirectURI + "/updated"}, body["redirect_uris"])
		rotatedToken, _ = body["registration_access_token"].(string)
		require.NotEmpty(t, rotatedToken)
		assert.NotEqual(t, token, rotatedToken)
	})

	t.Run("the previous token is eventually rejected after rotation", func(t *testing.T) {
		// Rotation invalidates the previous token, but verification is served from the
		// eventually consistent projection (the new token is accepted immediately through a
		// strongly consistent fallback). The previous token therefore stops working once the
		// rotation has been projected.
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		require.EventuallyWithT(t, func(tt *assert.CollectT) {
			oldStatus, _ := manageDynamicClient(t, http.MethodGet, clientURI, token, nil)
			assert.Equal(tt, http.StatusUnauthorized, oldStatus)
		}, retryDuration, tick, "previous token still accepted after rotation")
	})

	t.Run("read reflects the update with the rotated token", func(t *testing.T) {
		// The update is read back from the projection, which is eventually consistent, so
		// retry until the changed redirect_uris show up.
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		require.EventuallyWithT(t, func(tt *assert.CollectT) {
			getStatus, body := manageDynamicClient(t, http.MethodGet, clientURI, rotatedToken, nil)
			assert.Equal(tt, http.StatusOK, getStatus)
			assert.Equal(tt, []any{redirectURI + "/updated"}, body["redirect_uris"])
		}, retryDuration, tick, "update not reflected in read")
	})

	t.Run("a token bound to another client is rejected", func(t *testing.T) {
		_, other := registerDynamicClient(t, issuer, "", map[string]any{
			"redirect_uris":              []string{redirectURI},
			"token_endpoint_auth_method": "none",
		})
		otherURI, _ := other["registration_client_uri"].(string)
		require.NotEmpty(t, otherURI)
		// Using the first client's token against another client's endpoint must be rejected
		// before any lookup: the token's client_id does not match the path.
		status, _ := manageDynamicClient(t, http.MethodGet, otherURI, rotatedToken, nil)
		assert.Equal(t, http.StatusUnauthorized, status)
	})

	t.Run("delete removes the client", func(t *testing.T) {
		delStatus, _ := manageDynamicClient(t, http.MethodDelete, clientURI, rotatedToken, nil)
		require.Equal(t, http.StatusNoContent, delStatus)

		// After deletion the client is gone; wait for the projection to catch up.
		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		require.EventuallyWithT(t, func(tt *assert.CollectT) {
			getStatus, _ := manageDynamicClient(t, http.MethodGet, clientURI, rotatedToken, nil)
			assert.Equal(tt, http.StatusNotFound, getStatus)
		}, retryDuration, tick, "client still readable after deletion")
	})
}

// TestServer_DynamicClientRegistration_Management_disabled verifies that without the feature
// the management endpoints are not served.
func TestServer_DynamicClientRegistration_Management_disabled(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	instance := integration.NewInstance(ctx)
	issuer := instance.OIDCIssuer()

	status, _ := manageDynamicClient(t, http.MethodGet, issuer+"/oauth/v2/register/whatever", "sometoken", nil)
	assert.Equal(t, http.StatusNotFound, status)
}

func manageDynamicClient(t testing.TB, method, uri, token string, body map[string]any) (int, map[string]any) {
	t.Helper()
	var reader io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		require.NoError(t, err)
		reader = bytes.NewReader(raw)
	}
	req, err := http.NewRequest(method, uri, reader)
	require.NoError(t, err)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()
	var respBody map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&respBody)
	return resp.StatusCode, respBody
}
