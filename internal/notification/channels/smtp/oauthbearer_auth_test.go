package smtp

import (
	"net/smtp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOauthBearerAuth_Start(t *testing.T) {
	tests := []struct {
		Description string
		Config      OAuthBearerAuthConfig
		Host        string
		Server      smtp.ServerInfo
		Expected    string
	}{
		{
			Description: "example from rfc (https://datatracker.ietf.org/doc/html/rfc7628#section-4.1)",
			Config: OAuthBearerAuthConfig{
				User:        "user@example.com",
				BearerToken: "vF9dft4qmTc2Nvb3RlckBhbHRhdmlzdGEuY29tCg==",
			},
			Host: "server.example.com",
			Server: smtp.ServerInfo{
				Name: "server.example.com",
			},
			Expected: strings.ReplaceAll("n,a=user@example.com,^Ahost=server.example.com^Aport=587^Aauth=Bearer vF9dft4qmTc2Nvb3RlckBhbHRhdmlzdGEuY29tCg==^A^A", "^A", "\x01"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Description, func(t *testing.T) {
			xoauth := OAuthBearerAuth(tt.Config, tt.Host)

			m, v, err := xoauth.Start(&tt.Server)

			assert.NoError(t, err)
			assert.Equal(t, "OAUTHBEARER", m)
			assert.Equal(t, tt.Expected, string(v))
		})
	}
}
