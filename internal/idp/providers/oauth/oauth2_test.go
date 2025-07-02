package oauth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"golang.org/x/oauth2"

	"github.com/zitadel/zitadel/internal/idp"
)

func TestProvider_BeginAuth(t *testing.T) {
	type fields struct {
		config       *oauth2.Config
		name         string
		userEndpoint string
		userMapper   func() idp.User
		options      []ProviderOpts
	}
	tests := []struct {
		name   string
		fields fields
		want   idp.Session
	}{
		{
			name: "successful auth without PKCE",
			fields: fields{
				config: &oauth2.Config{
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
					Endpoint: oauth2.Endpoint{
						AuthURL:  "https://oauth2.com/authorize",
						TokenURL: "https://oauth2.com/token",
					},
					RedirectURL: "redirectURI",
					Scopes:      []string{"user"},
				},
			},
			want: &Session{AuthURL: "https://oauth2.com/authorize?client_id=clientID&prompt=select_account&redirect_uri=redirectURI&response_type=code&scope=user&state=testState"},
		},
		{
			name: "successful auth with PKCE",
			fields: fields{
				config: &oauth2.Config{
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
					Endpoint: oauth2.Endpoint{
						AuthURL:  "https://oauth2.com/authorize",
						TokenURL: "https://oauth2.com/token",
					},
					RedirectURL: "redirectURI",
					Scopes:      []string{"user"},
				},
				options: []ProviderOpts{
					WithLinkingAllowed(),
					WithCreationAllowed(),
					WithAutoCreation(),
					WithAutoUpdate(),
					WithRelyingPartyOption(rp.WithPKCE(nil)),
				},
			},
			want: &Session{AuthURL: "https://oauth2.com/authorize?client_id=clientID&code_challenge=2ZoH_a01aprzLkwVbjlPsBo4m8mJ_zOKkaDqYM7Oh5w&code_challenge_method=S256&prompt=select_account&redirect_uri=redirectURI&response_type=code&scope=user&state=testState"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			r := require.New(t)

			provider, err := New(tt.fields.config, tt.fields.name, tt.fields.userEndpoint, tt.fields.userMapper, tt.fields.options...)
			r.NoError(err)
			provider.generateVerifier = func() string {
				return "pkceOAuthVerifier"
			}

			ctx := context.Background()
			session, err := provider.BeginAuth(ctx, "testState")
			r.NoError(err)

			wantAuth, wantErr := tt.want.GetAuth(ctx)
			gotAuth, gotErr := session.GetAuth(ctx)
			a.Equal(wantAuth, gotAuth)
			a.ErrorIs(gotErr, wantErr)
		})
	}
}

func TestProvider_Options(t *testing.T) {
	type fields struct {
		config       *oauth2.Config
		name         string
		userEndpoint string
		userMapper   func() idp.User
		options      []ProviderOpts
	}
	type want struct {
		name            string
		linkingAllowed  bool
		creationAllowed bool
		autoCreation    bool
		autoUpdate      bool
		pkce            bool
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "default",
			fields: fields{
				name: "oauth",
				config: &oauth2.Config{
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
					Endpoint: oauth2.Endpoint{
						AuthURL:  "https://oauth2.com/authorize",
						TokenURL: "https://oauth2.com/token",
					},
					RedirectURL: "redirectURI",
					Scopes:      []string{"user"},
				},
				options: nil,
			},
			want: want{
				name:            "oauth",
				linkingAllowed:  false,
				creationAllowed: false,
				autoCreation:    false,
				autoUpdate:      false,
				pkce:            false,
			},
		},
		{
			name: "all true",
			fields: fields{
				name: "oauth",
				config: &oauth2.Config{
					ClientID:     "clientID",
					ClientSecret: "clientSecret",
					Endpoint: oauth2.Endpoint{
						AuthURL:  "https://oauth2.com/authorize",
						TokenURL: "https://oauth2.com/token",
					},
					RedirectURL: "redirectURI",
					Scopes:      []string{"user"},
				},
				options: []ProviderOpts{
					WithLinkingAllowed(),
					WithCreationAllowed(),
					WithAutoCreation(),
					WithAutoUpdate(),
					WithRelyingPartyOption(rp.WithPKCE(nil)),
				},
			},
			want: want{
				name:            "oauth",
				linkingAllowed:  true,
				creationAllowed: true,
				autoCreation:    true,
				autoUpdate:      true,
				pkce:            true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			provider, err := New(tt.fields.config, tt.fields.name, tt.fields.userEndpoint, tt.fields.userMapper, tt.fields.options...)
			require.NoError(t, err)

			a.Equal(tt.want.name, provider.Name())
			a.Equal(tt.want.linkingAllowed, provider.IsLinkingAllowed())
			a.Equal(tt.want.creationAllowed, provider.IsCreationAllowed())
			a.Equal(tt.want.autoCreation, provider.IsAutoCreation())
			a.Equal(tt.want.autoUpdate, provider.IsAutoUpdate())
			a.Equal(tt.want.pkce, provider.RelyingParty.IsPKCE())
		})
	}
}
