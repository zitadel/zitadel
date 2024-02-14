package apple

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/oidc"
)

const (
	privateKey = `-----BEGIN PRIVATE KEY-----
MIGTAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBHkwdwIBAQQgXn/LDURaetCoymSj
fRslBiBwzBSa8ifiyfYGIWNStYGgCgYIKoZIzj0DAQehRANCAATymZXIsGrXnl6b
+80miSiVOCcLnyaYa2uQBQvQwgB7GibXhrzF+D/MRTV4P7P8+Lg1K9Khkjc59eNK
4RrQP4g7
-----END PRIVATE KEY-----
`
)

func TestProvider_BeginAuth(t *testing.T) {
	type fields struct {
		clientID    string
		teamID      string
		keyID       string
		privateKey  []byte
		redirectURI string
		scopes      []string
	}
	tests := []struct {
		name   string
		fields fields
		want   idp.Session
	}{
		{
			name: "successful auth",
			fields: fields{
				clientID:    "clientID",
				teamID:      "teamID",
				keyID:       "keyID",
				privateKey:  []byte(privateKey),
				redirectURI: "redirectURI",
				scopes:      []string{"openid"},
			},
			want: &Session{
				Session: &oidc.Session{
					AuthURL: "https://appleid.apple.com/auth/authorize?client_id=clientID&redirect_uri=redirectURI&response_mode=form_post&response_type=code&scope=openid&state=testState",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			r := require.New(t)

			provider, err := New(tt.fields.clientID, tt.fields.teamID, tt.fields.keyID, tt.fields.redirectURI, tt.fields.privateKey, tt.fields.scopes)
			r.NoError(err)
			ctx := context.Background()
			session, err := provider.BeginAuth(ctx, "testState")
			r.NoError(err)
			content, redirect := session.GetAuth(ctx)
			contentExpected, redirectExpected := tt.want.GetAuth(ctx)
			a.Equal(redirectExpected, redirect)
			a.Equal(contentExpected, content)
		})
	}
}
