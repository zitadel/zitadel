package zitadel

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	openid "github.com/zitadel/oidc/v3/pkg/oidc"

	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/oauth"
	"github.com/zitadel/zitadel/internal/idp/providers/oidc"
)

const issuer = "https://issuer.com"

func discoveryMock(issuer string) {
	gock.New(issuer).
		Get(openid.DiscoveryEndpoint).
		Reply(200).
		JSON(&openid.DiscoveryConfiguration{
			Issuer:                issuer,
			AuthorizationEndpoint: issuer + "/authorize",
			TokenEndpoint:         issuer + "/token",
			UserinfoEndpoint:      issuer + "/userinfo",
		})
}

func TestProvider_BeginAuth(t *testing.T) {
	defer gock.Off()
	discoveryMock(issuer)

	provider, err := New(issuer, "clientID", "clientSecret", "redirectURI", []string{"openid"}, http.DefaultClient, oidc.WithSelectAccount())
	require.NoError(t, err)

	session, err := provider.BeginAuth(context.Background(), "testState")
	require.NoError(t, err)

	auth, err := session.GetAuth(context.Background())
	require.NoError(t, err)
	redirect, ok := auth.(*idp.RedirectAuth)
	require.True(t, ok)

	require.True(t, strings.HasPrefix(redirect.RedirectURL, issuer+"/authorize"), "unexpected auth url %q", redirect.RedirectURL)
	authURL, err := url.Parse(redirect.RedirectURL)
	require.NoError(t, err)
	query := authURL.Query()

	assert.Equal(t, "clientID", query.Get("client_id"))
	assert.Equal(t, "redirectURI", query.Get("redirect_uri"))
	assert.Equal(t, "code", query.Get("response_type"))
	assert.Equal(t, "openid", query.Get("scope"))
	assert.Equal(t, "testState", query.Get("state"))
	assert.Equal(t, "select_account", query.Get("prompt"))
	// the ZITADEL provider always enforces PKCE
	assert.NotEmpty(t, query.Get("code_challenge"))
	assert.Equal(t, "S256", query.Get("code_challenge_method"))
	assert.NotEmpty(t, session.PersistentParameters()[oauth.CodeVerifier])
}

func TestProvider_Options(t *testing.T) {
	type fields struct {
		opts []oidc.ProviderOpts
	}
	type want struct {
		name            string
		linkingAllowed  bool
		creationAllowed bool
		autoCreation    bool
		autoUpdate      bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "default",
			fields: fields{
				opts: nil,
			},
			want: want{
				name:            "ZITADEL",
				linkingAllowed:  false,
				creationAllowed: false,
				autoCreation:    false,
				autoUpdate:      false,
			},
		},
		{
			name: "all true",
			fields: fields{
				opts: []oidc.ProviderOpts{
					oidc.WithLinkingAllowed(),
					oidc.WithCreationAllowed(),
					oidc.WithAutoCreation(),
					oidc.WithAutoUpdate(),
				},
			},
			want: want{
				name:            "ZITADEL",
				linkingAllowed:  true,
				creationAllowed: true,
				autoCreation:    true,
				autoUpdate:      true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Off()
			discoveryMock(issuer)
			a := assert.New(t)

			provider, err := New(issuer, "clientID", "clientSecret", "redirectURI", []string{"openid"}, http.DefaultClient, tt.fields.opts...)
			require.NoError(t, err)

			a.Equal(tt.want.name, provider.Name())
			a.Equal(tt.want.linkingAllowed, provider.IsLinkingAllowed())
			a.Equal(tt.want.creationAllowed, provider.IsCreationAllowed())
			a.Equal(tt.want.autoCreation, provider.IsAutoCreation())
			a.Equal(tt.want.autoUpdate, provider.IsAutoUpdate())
		})
	}
}
