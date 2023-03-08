package gitlab

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/oidc"
)

func TestProvider_BeginAuth(t *testing.T) {
	type fields struct {
		clientID     string
		clientSecret string
		redirectURI  string
		scopes       []string
		opts         []oidc.ProviderOpts
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
				scopes:       []string{"openid"},
			},
			want: &oidc.Session{
				AuthURL: "https://gitlab.com/oauth/authorize?client_id=clientID&prompt=select_account&redirect_uri=redirectURI&response_type=code&scope=openid&state=testState",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			r := require.New(t)

			provider, err := New(tt.fields.clientID, tt.fields.clientSecret, tt.fields.redirectURI, tt.fields.scopes, tt.fields.opts...)
			r.NoError(err)

			session, err := provider.BeginAuth(context.Background(), "testState")
			r.NoError(err)

			a.Equal(tt.want.GetAuthURL(), session.GetAuthURL())
		})
	}
}
