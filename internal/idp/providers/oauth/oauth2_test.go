package oauth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v2/pkg/client/rp"
	"golang.org/x/oauth2"

	"github.com/zitadel/zitadel/internal/idp"
)

func TestProvider_BeginAuth(t *testing.T) {
	type fields struct {
		config       *oauth2.Config
		name         string
		userEndpoint string
		userMapper   func() idp.User
	}
	tests := []struct {
		name   string
		fields fields
		want   idp.Session
	}{
		{
			name: "successful auth",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			r := require.New(t)

			provider, err := New(tt.fields.config, tt.fields.name, tt.fields.userEndpoint, tt.fields.userMapper)
			r.NoError(err)

			session, err := provider.BeginAuth(context.Background(), "testState")
			r.NoError(err)

			a.Equal(tt.want.GetAuthURL(), session.GetAuthURL())
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
