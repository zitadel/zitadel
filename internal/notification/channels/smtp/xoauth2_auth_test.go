package smtp

import (
	"net/smtp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/oauth2"

	"github.com/zitadel/zitadel/internal/notification/channels/mock"
)

func TestXoauth2Auth_Start(t *testing.T) {
	tests := []struct {
		Description string
		Config      XOAuth2AuthConfig
		Host        string
		Server      smtp.ServerInfo
		Token       string
		Expected    string
	}{
		{
			Description: "example from google (https://developers.google.com/workspace/gmail/imap/xoauth2-protocol)",
			Config: XOAuth2AuthConfig{
				User:          "someuser@example.com",
				ClientId:      "my-client",
				ClientSecret:  "secret-value",
				TokenEndpoint: "http://auth.example.com/token",
				Scopes:        []string{"special:token:scope"},
			},
			Host: "example.com",
			Server: smtp.ServerInfo{
				Name: "example.com",
			},
			Token:    "ya29.vF9dft4qmTc2Nvb3RlckBhdHRhdmlzdGEuY29tCg",
			Expected: strings.ReplaceAll("user=someuser@example.com^Aauth=Bearer ya29.vF9dft4qmTc2Nvb3RlckBhdHRhdmlzdGEuY29tCg^A^A", "^A", "\x01"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Description, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ts := mock.NewMockTokenSource(ctrl)

			ts.EXPECT().Token().Return(&oauth2.Token{AccessToken: tt.Token}, nil)

			xoauth := XOAuth2Auth(tt.Config, tt.Host)
			(xoauth.(*xoauth2Auth)).tokenSource = ts

			m, v, err := xoauth.Start(&tt.Server)

			assert.NoError(t, err)
			assert.Equal(t, "XOAUTH2", m)
			assert.Equal(t, tt.Expected, string(v))
		})
	}
}
