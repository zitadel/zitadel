package azuread

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/idp"
	oidc2 "github.com/zitadel/zitadel/internal/idp/providers/oidc"
)

func TestProvider_BeginAuth(t *testing.T) {
	type fields struct {
		name         string
		clientID     string
		clientSecret string
		redirectURI  string
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
			want: &oidc2.Session{
				AuthURL: "https://login.microsoftonline.com/common/oauth2/v2.0/authorize?client_id=clientID&redirect_uri=redirectURI&response_type=code&scope=openid+profile+email&state=testState",
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
			want: &oidc2.Session{
				AuthURL: "https://login.microsoftonline.com/consumers/oauth2/v2.0/authorize?client_id=clientID&redirect_uri=redirectURI&response_type=code&scope=openid+profile+email&state=testState",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			provider, err := New(tt.fields.name, tt.fields.clientID, tt.fields.clientSecret, tt.fields.redirectURI, tt.fields.options...)
			a.NoError(err)

			session, err := provider.BeginAuth("testState")
			a.NoError(err)

			a.Equal(tt.want.GetAuthURL(), session.GetAuthURL())
		})
	}
}
