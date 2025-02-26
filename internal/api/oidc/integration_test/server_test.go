//go:build integration

package oidc_test

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v3/pkg/client"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/schema"
)

func TestServer_RefreshToken_Status(t *testing.T) {
	clientID, _ := createClient(t, Instance)
	provider, err := Instance.CreateRelyingParty(CTX, clientID, redirectURI)
	require.NoError(t, err)

	tests := []struct {
		name         string
		refreshToken string
	}{
		{
			name:         "invalid base64",
			refreshToken: "~!~@#$%",
		},
		{
			name:         "invalid after decrypt",
			refreshToken: "DEADBEEFDEADBEEF",
		},
		{
			name:         "short input",
			refreshToken: "DEAD",
		},
		{
			name:         "empty input",
			refreshToken: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := rp.RefreshTokenRequest{
				RefreshToken: tt.refreshToken,
				ClientID:     clientID,
				GrantType:    oidc.GrantTypeRefreshToken,
			}
			client.CallTokenEndpoint(CTX, request, tokenEndpointCaller{RelyingParty: provider})

			values := make(url.Values)
			err := schema.NewEncoder().Encode(request, values)
			require.NoError(t, err)

			resp, err := http.Post(provider.OAuthConfig().Endpoint.TokenURL, "application/x-www-form-urlencoded", strings.NewReader(values.Encode()))
			require.NoError(t, err)
			defer resp.Body.Close()
			assert.Less(t, resp.StatusCode, 500)
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			t.Log(string(body))
		})
	}
}

type tokenEndpointCaller struct {
	rp.RelyingParty
}

func (t tokenEndpointCaller) TokenEndpoint() string {
	return t.OAuthConfig().Endpoint.TokenURL
}
