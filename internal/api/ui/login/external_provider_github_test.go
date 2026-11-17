package login

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/query"
)

// Test_githubEnterpriseProvider ensures the provider is built from the GitHub Enterprise Server template,
// which is the only one set by the query for such an IdP.
func Test_githubEnterpriseProvider(t *testing.T) {
	l := &Login{
		idpConfigAlg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
		httpClient:   http.DefaultClient,
	}
	identityProvider := &query.IDPTemplate{
		Name: "GitHub Enterprise",
		GitHubEnterpriseIDPTemplate: &query.GitHubEnterpriseIDPTemplate{
			IDPID:    "idp",
			ClientID: "clientID",
			ClientSecret: &crypto.CryptoValue{
				CryptoType: crypto.TypeEncryption,
				Algorithm:  "enc",
				KeyID:      "id",
				Crypted:    []byte("clientSecret"),
			},
			AuthorizationEndpoint: "https://ghes.example.com/login/oauth/authorize",
			TokenEndpoint:         "https://ghes.example.com/login/oauth/access_token",
			UserEndpoint:          "https://ghes.example.com/api/v3/user",
			Scopes:                []string{"read:user"},
		},
	}

	provider, err := l.githubEnterpriseProvider(context.Background(), identityProvider)
	require.NoError(t, err)

	config := provider.OAuthConfig()
	assert.Equal(t, "GitHub Enterprise", provider.Name())
	assert.Equal(t, "clientID", config.ClientID)
	assert.Equal(t, "clientSecret", config.ClientSecret)
	assert.Equal(t, []string{"read:user"}, config.Scopes)
	assert.Equal(t, "https://ghes.example.com/login/oauth/authorize", config.Endpoint.AuthURL)
	assert.Equal(t, "https://ghes.example.com/login/oauth/access_token", config.Endpoint.TokenURL)
}
