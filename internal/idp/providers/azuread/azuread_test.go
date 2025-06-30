package azuread

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	openid "github.com/zitadel/oidc/v3/pkg/oidc"

	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/oauth"
	"github.com/zitadel/zitadel/internal/idp/providers/oidc"
)

func TestProvider_BeginAuth(t *testing.T) {
	type fields struct {
		name         string
		clientID     string
		clientSecret string
		redirectURI  string
		scopes       []string
		options      []ProviderOptions
	}
	tests := []struct {
		name   string
		fields fields
		want   idp.Session
	}{
		{
			name: "default common tenant",
			fields: fields{
				clientID:     "clientID",
				clientSecret: "clientSecret",
				redirectURI:  "redirectURI",
			},
			want: &oidc.Session{
				AuthURL: "https://login.microsoftonline.com/common/oauth2/v2.0/authorize?client_id=clientID&prompt=select_account&redirect_uri=redirectURI&response_type=code&scope=openid+profile+email+phone+User.Read&state=testState",
			},
		},
		{
			name: "tenant",
			fields: fields{
				clientID:     "clientID",
				clientSecret: "clientSecret",
				redirectURI:  "redirectURI",
				options: []ProviderOptions{
					WithTenant(ConsumersTenant),
				},
			},
			want: &oidc.Session{
				AuthURL: "https://login.microsoftonline.com/consumers/oauth2/v2.0/authorize?client_id=clientID&prompt=select_account&redirect_uri=redirectURI&response_type=code&scope=openid+profile+email+phone+User.Read&state=testState",
			},
		},
		{
			name: "custom scopes",
			fields: fields{
				clientID:     "clientID",
				clientSecret: "clientSecret",
				redirectURI:  "redirectURI",
				scopes:       []string{openid.ScopeOpenID, openid.ScopeProfile, "custom"},
				options: []ProviderOptions{
					WithTenant(ConsumersTenant),
				},
			},
			want: &oidc.Session{
				AuthURL: "https://login.microsoftonline.com/consumers/oauth2/v2.0/authorize?client_id=clientID&prompt=select_account&redirect_uri=redirectURI&response_type=code&scope=openid+profile+custom+User.Read&state=testState",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			r := require.New(t)

			provider, err := New(tt.fields.name, tt.fields.clientID, tt.fields.clientSecret, tt.fields.redirectURI, tt.fields.scopes, tt.fields.options...)
			r.NoError(err)

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
		name         string
		clientID     string
		clientSecret string
		redirectURI  string
		scopes       []string
		options      []ProviderOptions
	}
	type want struct {
		name            string
		tenant          TenantType
		emailVerified   bool
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
			name: "default common tenant",
			fields: fields{
				name:         "default common tenant",
				clientID:     "clientID",
				clientSecret: "clientSecret",
				redirectURI:  "redirectURI",
				scopes:       nil,
				options:      nil,
			},
			want: want{
				name:            "default common tenant",
				tenant:          CommonTenant,
				emailVerified:   false,
				linkingAllowed:  false,
				creationAllowed: false,
				autoCreation:    false,
				autoUpdate:      false,
				pkce:            false,
			},
		},
		{
			name: "all set",
			fields: fields{
				name:         "custom tenant",
				clientID:     "clientID",
				clientSecret: "clientSecret",
				redirectURI:  "redirectURI",
				options: []ProviderOptions{
					WithTenant("tenant"),
					WithEmailVerified(),
					WithOAuthOptions(
						oauth.WithLinkingAllowed(),
						oauth.WithCreationAllowed(),
						oauth.WithAutoCreation(),
						oauth.WithAutoUpdate(),
						oauth.WithRelyingPartyOption(rp.WithPKCE(nil)),
					),
				},
			},
			want: want{
				name:            "custom tenant",
				tenant:          "tenant",
				emailVerified:   true,
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

			provider, err := New(tt.fields.name, tt.fields.clientID, tt.fields.clientSecret, tt.fields.redirectURI, tt.fields.scopes, tt.fields.options...)
			require.NoError(t, err)

			a.Equal(tt.want.name, provider.Name())
			a.Equal(tt.want.tenant, provider.tenant)
			a.Equal(tt.want.emailVerified, provider.emailVerified)
			a.Equal(tt.want.linkingAllowed, provider.IsLinkingAllowed())
			a.Equal(tt.want.creationAllowed, provider.IsCreationAllowed())
			a.Equal(tt.want.autoCreation, provider.IsAutoCreation())
			a.Equal(tt.want.autoUpdate, provider.IsAutoUpdate())
			a.Equal(tt.want.pkce, provider.RelyingParty.IsPKCE())
		})
	}
}
