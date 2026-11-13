package github

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/oauth"
)

func TestProvider_BeginAuth(t *testing.T) {
	type fields struct {
		clientID     string
		clientSecret string
		redirectURI  string
		scopes       []string
		options      []oauth.ProviderOpts
	}
	tests := []struct {
		name   string
		fields fields
		want   idp.Session
	}{
		{
			name: "successful auth",
			fields: fields{
				clientID:     "clientID",
				clientSecret: "clientSecret",
				redirectURI:  "redirectURI",
			},
			want: &oauth.Session{
				AuthURL: "https://github.com/login/oauth/authorize?client_id=clientID&prompt=select_account&redirect_uri=redirectURI&response_type=code&state=testState",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			r := require.New(t)

			provider, err := New(tt.fields.clientID, tt.fields.clientSecret, tt.fields.redirectURI, tt.fields.scopes, http.DefaultClient, tt.fields.options...)
			r.NoError(err)

			ctx := context.Background()
			session, err := provider.BeginAuth(ctx, "testState")
			r.NoError(err)

			wantAuth, wantErr := tt.want.GetAuth(ctx)
			gotAuth, gotErr := session.GetAuth(ctx)
			a.ErrorIs(gotErr, wantErr)
			// PKCE is used by default
			assertRedirectAuthWithPKCE(t, wantAuth, gotAuth)
			a.NotEmpty(session.PersistentParameters()[oauth.CodeVerifier])
		})
	}
}

func TestProvider_BeginAuth_EnterpriseServerWithoutPKCE(t *testing.T) {
	// GitHub Enterprise Server does not support PKCE (yet)
	provider, err := NewCustomURL("GitHub Enterprise", "clientID", "clientSecret", "redirectURI", "https://ghes.com/login/oauth/authorize", "https://ghes.com/login/oauth/access_token", "https://ghes.com/api/v3/user", nil, http.DefaultClient)
	require.NoError(t, err)

	ctx := context.Background()
	session, err := provider.BeginAuth(ctx, "testState")
	require.NoError(t, err)

	auth, err := session.GetAuth(ctx)
	require.NoError(t, err)
	assert.Equal(t, &idp.RedirectAuth{RedirectURL: "https://ghes.com/login/oauth/authorize?client_id=clientID&prompt=select_account&redirect_uri=redirectURI&response_type=code&state=testState"}, auth)
	assert.Empty(t, session.PersistentParameters())
}

// assertRedirectAuthWithPKCE asserts that the actual auth URL contains the PKCE parameters
// and otherwise matches the expected one, since the code challenge is random.
func assertRedirectAuthWithPKCE(t *testing.T, expected, actual idp.Auth) {
	t.Helper()

	expectedRedirect, ok := expected.(*idp.RedirectAuth)
	require.True(t, ok)
	actualRedirect, ok := actual.(*idp.RedirectAuth)
	require.True(t, ok)

	expectedURL, err := url.Parse(expectedRedirect.RedirectURL)
	require.NoError(t, err)
	actualURL, err := url.Parse(actualRedirect.RedirectURL)
	require.NoError(t, err)

	actualQuery := actualURL.Query()
	assert.NotEmpty(t, actualQuery.Get("code_challenge"))
	assert.Equal(t, "S256", actualQuery.Get("code_challenge_method"))

	actualQuery.Del("code_challenge")
	actualQuery.Del("code_challenge_method")

	assert.Equal(t, expectedURL.Scheme, actualURL.Scheme)
	assert.Equal(t, expectedURL.Host, actualURL.Host)
	assert.Equal(t, expectedURL.Path, actualURL.Path)
	assert.Equal(t, expectedURL.Query(), actualQuery)
}
