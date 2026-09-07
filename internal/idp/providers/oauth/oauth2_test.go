package oauth

import (
	"context"
	"net/http"
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
		params       []idp.Parameter
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
		{
			name: "authorization parameters override the default prompt",
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
				options: []ProviderOpts{WithAuthorizationParameters(map[string]string{"prompt": "login", "acr_values": "AAL2_OR_AAL3_ANY"})},
			},
			want: &Session{AuthURL: "https://oauth2.com/authorize?acr_values=AAL2_OR_AAL3_ANY&client_id=clientID&prompt=login&redirect_uri=redirectURI&response_type=code&scope=user&state=testState"},
		},
		{
			name: "an empty authorization parameter removes the default prompt",
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
				options: []ProviderOpts{WithAuthorizationParameters(map[string]string{"prompt": ""})},
			},
			want: &Session{AuthURL: "https://oauth2.com/authorize?client_id=clientID&redirect_uri=redirectURI&response_type=code&scope=user&state=testState"},
		},
		{
			name: "reserved authorization parameters are dropped",
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
				options: []ProviderOpts{WithAuthorizationParameters(map[string]string{"client_id": "evil", "scope": "admin", "acr_values": "loa2"})},
			},
			want: &Session{AuthURL: "https://oauth2.com/authorize?acr_values=loa2&client_id=clientID&prompt=select_account&redirect_uri=redirectURI&response_type=code&scope=user&state=testState"},
		},
		{
			name: "only allow-listed parameters are forwarded",
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
				options: []ProviderOpts{WithForwardedParameters([]string{"acr_values", "max_age"})},
				params:  []idp.Parameter{idp.AuthorizationParameters{"acr_values": "loa3", "max_age": "300", "ui_locales": "de"}},
			},
			want: &Session{AuthURL: "https://oauth2.com/authorize?acr_values=loa3&client_id=clientID&max_age=300&prompt=select_account&redirect_uri=redirectURI&response_type=code&scope=user&state=testState"},
		},
		{
			name: "nothing is forwarded without an allow-list",
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
				options: []ProviderOpts{},
				params:  []idp.Parameter{idp.AuthorizationParameters{"acr_values": "loa3"}},
			},
			want: &Session{AuthURL: "https://oauth2.com/authorize?client_id=clientID&prompt=select_account&redirect_uri=redirectURI&response_type=code&scope=user&state=testState"},
		},
		{
			name: "configured parameters take precedence over forwarded ones",
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
				options: []ProviderOpts{WithForwardedParameters([]string{"acr_values"}), WithAuthorizationParameters(map[string]string{"acr_values": "static"})},
				params:  []idp.Parameter{idp.AuthorizationParameters{"acr_values": "forwarded"}},
			},
			want: &Session{AuthURL: "https://oauth2.com/authorize?acr_values=static&client_id=clientID&prompt=select_account&redirect_uri=redirectURI&response_type=code&scope=user&state=testState"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			r := require.New(t)

			provider, err := New(tt.fields.config, tt.fields.name, tt.fields.userEndpoint, tt.fields.userMapper, http.DefaultClient, tt.fields.options...)
			r.NoError(err)
			provider.generateVerifier = func() string {
				return "pkceOAuthVerifier"
			}

			ctx := context.Background()
			session, err := provider.BeginAuth(ctx, "testState", tt.fields.params...)
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

			provider, err := New(tt.fields.config, tt.fields.name, tt.fields.userEndpoint, tt.fields.userMapper, http.DefaultClient, tt.fields.options...)
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
