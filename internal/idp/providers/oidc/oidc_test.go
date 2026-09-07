package oidc

import (
	"context"
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/oidc"

	"github.com/zitadel/zitadel/internal/idp"
)

func TestProvider_BeginAuth(t *testing.T) {
	type fields struct {
		name         string
		issuer       string
		clientID     string
		clientSecret string
		redirectURI  string
		scopes       []string
		userMapper   func(info *oidc.UserInfo) idp.User
		httpMock     func(issuer string)
		opts         []ProviderOpts
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
				name:         "oidc",
				issuer:       "https://issuer.com",
				clientID:     "clientID",
				clientSecret: "clientSecret",
				redirectURI:  "redirectURI",
				scopes:       []string{"openid"},
				userMapper:   DefaultMapper,
				httpMock: func(issuer string) {
					gock.New(issuer).
						Get(oidc.DiscoveryEndpoint).
						Reply(200).
						JSON(&oidc.DiscoveryConfiguration{
							Issuer:                issuer,
							AuthorizationEndpoint: issuer + "/authorize",
							TokenEndpoint:         issuer + "/token",
							UserinfoEndpoint:      issuer + "/userinfo",
						})
				},
				opts: []ProviderOpts{WithSelectAccount()},
			},
			want: &Session{AuthURL: "https://issuer.com/authorize?client_id=clientID&prompt=select_account&redirect_uri=redirectURI&response_type=code&scope=openid&state=testState"},
		},
		{
			name: "successful auth with PKCE",
			fields: fields{
				name:         "oidc",
				issuer:       "https://issuer.com",
				clientID:     "clientID",
				clientSecret: "clientSecret",
				redirectURI:  "redirectURI",
				scopes:       []string{"openid"},
				userMapper:   DefaultMapper,
				httpMock: func(issuer string) {
					gock.New(issuer).
						Get(oidc.DiscoveryEndpoint).
						Reply(200).
						JSON(&oidc.DiscoveryConfiguration{
							Issuer:                issuer,
							AuthorizationEndpoint: issuer + "/authorize",
							TokenEndpoint:         issuer + "/token",
							UserinfoEndpoint:      issuer + "/userinfo",
						})
				},
				opts: []ProviderOpts{WithSelectAccount(), WithRelyingPartyOption(rp.WithPKCE(nil))},
			},
			want: &Session{AuthURL: "https://issuer.com/authorize?client_id=clientID&code_challenge=2ZoH_a01aprzLkwVbjlPsBo4m8mJ_zOKkaDqYM7Oh5w&code_challenge_method=S256&prompt=select_account&redirect_uri=redirectURI&response_type=code&scope=openid&state=testState"},
		},
		{
			name: "authorization parameters override the default prompt",
			fields: fields{
				name:         "oidc",
				issuer:       "https://issuer.com",
				clientID:     "clientID",
				clientSecret: "clientSecret",
				redirectURI:  "redirectURI",
				scopes:       []string{"openid"},
				userMapper:   DefaultMapper,
				httpMock: func(issuer string) {
					gock.New(issuer).
						Get(oidc.DiscoveryEndpoint).
						Reply(200).
						JSON(&oidc.DiscoveryConfiguration{
							Issuer:                issuer,
							AuthorizationEndpoint: issuer + "/authorize",
							TokenEndpoint:         issuer + "/token",
							UserinfoEndpoint:      issuer + "/userinfo",
						})
				},
				opts: []ProviderOpts{WithSelectAccount(), WithAuthorizationParameters(map[string]string{"prompt": "login", "acr_values": "AAL2_OR_AAL3_ANY"})},
			},
			want: &Session{AuthURL: "https://issuer.com/authorize?acr_values=AAL2_OR_AAL3_ANY&client_id=clientID&prompt=login&redirect_uri=redirectURI&response_type=code&scope=openid&state=testState"},
		},
		{
			name: "an empty authorization parameter removes the default prompt",
			fields: fields{
				name:         "oidc",
				issuer:       "https://issuer.com",
				clientID:     "clientID",
				clientSecret: "clientSecret",
				redirectURI:  "redirectURI",
				scopes:       []string{"openid"},
				userMapper:   DefaultMapper,
				httpMock: func(issuer string) {
					gock.New(issuer).
						Get(oidc.DiscoveryEndpoint).
						Reply(200).
						JSON(&oidc.DiscoveryConfiguration{
							Issuer:                issuer,
							AuthorizationEndpoint: issuer + "/authorize",
							TokenEndpoint:         issuer + "/token",
							UserinfoEndpoint:      issuer + "/userinfo",
						})
				},
				opts: []ProviderOpts{WithSelectAccount(), WithAuthorizationParameters(map[string]string{"prompt": ""})},
			},
			want: &Session{AuthURL: "https://issuer.com/authorize?client_id=clientID&redirect_uri=redirectURI&response_type=code&scope=openid&state=testState"},
		},
		{
			name: "authorization parameters are escaped",
			fields: fields{
				name:         "oidc",
				issuer:       "https://issuer.com",
				clientID:     "clientID",
				clientSecret: "clientSecret",
				redirectURI:  "redirectURI",
				scopes:       []string{"openid"},
				userMapper:   DefaultMapper,
				httpMock: func(issuer string) {
					gock.New(issuer).
						Get(oidc.DiscoveryEndpoint).
						Reply(200).
						JSON(&oidc.DiscoveryConfiguration{
							Issuer:                issuer,
							AuthorizationEndpoint: issuer + "/authorize",
							TokenEndpoint:         issuer + "/token",
							UserinfoEndpoint:      issuer + "/userinfo",
						})
				},
				opts: []ProviderOpts{WithSelectAccount(), WithAuthorizationParameters(map[string]string{"acr_values": "level 2&3"})},
			},
			want: &Session{AuthURL: "https://issuer.com/authorize?acr_values=level+2%263&client_id=clientID&prompt=select_account&redirect_uri=redirectURI&response_type=code&scope=openid&state=testState"},
		},
		{
			name: "reserved authorization parameters are dropped",
			fields: fields{
				name:         "oidc",
				issuer:       "https://issuer.com",
				clientID:     "clientID",
				clientSecret: "clientSecret",
				redirectURI:  "redirectURI",
				scopes:       []string{"openid"},
				userMapper:   DefaultMapper,
				httpMock: func(issuer string) {
					gock.New(issuer).
						Get(oidc.DiscoveryEndpoint).
						Reply(200).
						JSON(&oidc.DiscoveryConfiguration{
							Issuer:                issuer,
							AuthorizationEndpoint: issuer + "/authorize",
							TokenEndpoint:         issuer + "/token",
							UserinfoEndpoint:      issuer + "/userinfo",
						})
				},
				opts: []ProviderOpts{WithSelectAccount(), WithAuthorizationParameters(map[string]string{"client_id": "evil", "redirect_uri": "https://evil.com", "acr_values": "loa2"})},
			},
			want: &Session{AuthURL: "https://issuer.com/authorize?acr_values=loa2&client_id=clientID&prompt=select_account&redirect_uri=redirectURI&response_type=code&scope=openid&state=testState"},
		},
		{
			name: "only allow-listed parameters are forwarded",
			fields: fields{
				name:         "oidc",
				issuer:       "https://issuer.com",
				clientID:     "clientID",
				clientSecret: "clientSecret",
				redirectURI:  "redirectURI",
				scopes:       []string{"openid"},
				userMapper:   DefaultMapper,
				httpMock: func(issuer string) {
					gock.New(issuer).
						Get(oidc.DiscoveryEndpoint).
						Reply(200).
						JSON(&oidc.DiscoveryConfiguration{
							Issuer:                issuer,
							AuthorizationEndpoint: issuer + "/authorize",
							TokenEndpoint:         issuer + "/token",
							UserinfoEndpoint:      issuer + "/userinfo",
						})
				},
				opts:   []ProviderOpts{WithSelectAccount(), WithForwardedParameters([]string{"acr_values", "max_age"})},
				params: []idp.Parameter{idp.AuthorizationParameters{"acr_values": "loa3", "max_age": "300", "ui_locales": "de"}},
			},
			want: &Session{AuthURL: "https://issuer.com/authorize?acr_values=loa3&client_id=clientID&max_age=300&prompt=select_account&redirect_uri=redirectURI&response_type=code&scope=openid&state=testState"},
		},
		{
			name: "nothing is forwarded without an allow-list",
			fields: fields{
				name:         "oidc",
				issuer:       "https://issuer.com",
				clientID:     "clientID",
				clientSecret: "clientSecret",
				redirectURI:  "redirectURI",
				scopes:       []string{"openid"},
				userMapper:   DefaultMapper,
				httpMock: func(issuer string) {
					gock.New(issuer).
						Get(oidc.DiscoveryEndpoint).
						Reply(200).
						JSON(&oidc.DiscoveryConfiguration{
							Issuer:                issuer,
							AuthorizationEndpoint: issuer + "/authorize",
							TokenEndpoint:         issuer + "/token",
							UserinfoEndpoint:      issuer + "/userinfo",
						})
				},
				opts:   []ProviderOpts{WithSelectAccount()},
				params: []idp.Parameter{idp.AuthorizationParameters{"acr_values": "loa3"}},
			},
			want: &Session{AuthURL: "https://issuer.com/authorize?client_id=clientID&prompt=select_account&redirect_uri=redirectURI&response_type=code&scope=openid&state=testState"},
		},
		{
			name: "forwarded prompt replaces the default",
			fields: fields{
				name:         "oidc",
				issuer:       "https://issuer.com",
				clientID:     "clientID",
				clientSecret: "clientSecret",
				redirectURI:  "redirectURI",
				scopes:       []string{"openid"},
				userMapper:   DefaultMapper,
				httpMock: func(issuer string) {
					gock.New(issuer).
						Get(oidc.DiscoveryEndpoint).
						Reply(200).
						JSON(&oidc.DiscoveryConfiguration{
							Issuer:                issuer,
							AuthorizationEndpoint: issuer + "/authorize",
							TokenEndpoint:         issuer + "/token",
							UserinfoEndpoint:      issuer + "/userinfo",
						})
				},
				opts:   []ProviderOpts{WithSelectAccount(), WithForwardedParameters([]string{"prompt"})},
				params: []idp.Parameter{idp.AuthorizationParameters{"prompt": "login"}},
			},
			want: &Session{AuthURL: "https://issuer.com/authorize?client_id=clientID&prompt=login&redirect_uri=redirectURI&response_type=code&scope=openid&state=testState"},
		},
		{
			name: "configured parameters take precedence over forwarded ones",
			fields: fields{
				name:         "oidc",
				issuer:       "https://issuer.com",
				clientID:     "clientID",
				clientSecret: "clientSecret",
				redirectURI:  "redirectURI",
				scopes:       []string{"openid"},
				userMapper:   DefaultMapper,
				httpMock: func(issuer string) {
					gock.New(issuer).
						Get(oidc.DiscoveryEndpoint).
						Reply(200).
						JSON(&oidc.DiscoveryConfiguration{
							Issuer:                issuer,
							AuthorizationEndpoint: issuer + "/authorize",
							TokenEndpoint:         issuer + "/token",
							UserinfoEndpoint:      issuer + "/userinfo",
						})
				},
				opts:   []ProviderOpts{WithSelectAccount(), WithForwardedParameters([]string{"acr_values"}), WithAuthorizationParameters(map[string]string{"acr_values": "static"})},
				params: []idp.Parameter{idp.AuthorizationParameters{"acr_values": "forwarded"}},
			},
			want: &Session{AuthURL: "https://issuer.com/authorize?acr_values=static&client_id=clientID&prompt=select_account&redirect_uri=redirectURI&response_type=code&scope=openid&state=testState"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Off()
			tt.fields.httpMock(tt.fields.issuer)
			a := assert.New(t)
			r := require.New(t)

			provider, err := New(tt.fields.name, tt.fields.issuer, tt.fields.clientID, tt.fields.clientSecret, tt.fields.redirectURI, tt.fields.scopes, tt.fields.userMapper, http.DefaultClient, tt.fields.opts...)
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
		name         string
		issuer       string
		clientID     string
		clientSecret string
		redirectURI  string
		scopes       []string
		userMapper   func(info *oidc.UserInfo) idp.User
		opts         []ProviderOpts
		httpMock     func(issuer string)
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
				name:         "oidc",
				issuer:       "https://issuer.com",
				clientID:     "clientID",
				clientSecret: "clientSecret",
				redirectURI:  "redirectURI",
				scopes:       []string{"openid"},
				userMapper:   DefaultMapper,
				opts:         nil,
				httpMock: func(issuer string) {
					gock.New(issuer).
						Get(oidc.DiscoveryEndpoint).
						Reply(200).
						JSON(&oidc.DiscoveryConfiguration{
							Issuer:                issuer,
							AuthorizationEndpoint: issuer + "/authorize",
							TokenEndpoint:         issuer + "/token",
							UserinfoEndpoint:      issuer + "/userinfo",
						})
				},
			},
			want: want{
				name:            "oidc",
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
				name:         "oidc",
				issuer:       "https://issuer.com",
				clientID:     "clientID",
				clientSecret: "clientSecret",
				redirectURI:  "redirectURI",
				scopes:       []string{"openid"},
				userMapper:   DefaultMapper,
				opts: []ProviderOpts{
					WithLinkingAllowed(),
					WithCreationAllowed(),
					WithAutoCreation(),
					WithAutoUpdate(),
					WithRelyingPartyOption(rp.WithPKCE(nil)),
				},
				httpMock: func(issuer string) {
					gock.New(issuer).
						Get(oidc.DiscoveryEndpoint).
						Reply(200).
						JSON(&oidc.DiscoveryConfiguration{
							Issuer:                issuer,
							AuthorizationEndpoint: issuer + "/authorize",
							TokenEndpoint:         issuer + "/token",
							UserinfoEndpoint:      issuer + "/userinfo",
						})
				},
			},
			want: want{
				name:            "oidc",
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
			defer gock.Off()
			tt.fields.httpMock(tt.fields.issuer)
			a := assert.New(t)

			provider, err := New(tt.fields.name, tt.fields.issuer, tt.fields.clientID, tt.fields.clientSecret, tt.fields.redirectURI, tt.fields.scopes, tt.fields.userMapper, http.DefaultClient, tt.fields.opts...)
			require.NoError(t, err)

			a.Equal(tt.want.name, provider.Name())
			a.Equal(tt.want.linkingAllowed, provider.IsLinkingAllowed())
			a.Equal(tt.want.creationAllowed, provider.IsCreationAllowed())
			a.Equal(tt.want.autoCreation, provider.IsAutoCreation())
			a.Equal(tt.want.autoUpdate, provider.IsAutoUpdate())
		})
	}
}
