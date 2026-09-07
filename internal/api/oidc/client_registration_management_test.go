package oidc

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
)

// fakeRegistrationTokenAlg is a reversible stand-in for the instance AuthAlgorithm. It only
// has to round-trip a string so the binding logic can be exercised without real keys.
type fakeRegistrationTokenAlg struct{}

func (fakeRegistrationTokenAlg) EncryptToken(data string) (string, error) {
	return "enc(" + data + ")", nil
}

func (fakeRegistrationTokenAlg) DecryptToken(token string) (string, error) {
	if !strings.HasPrefix(token, "enc(") || !strings.HasSuffix(token, ")") {
		return "", errors.New("not a registration access token")
	}
	return strings.TrimSuffix(strings.TrimPrefix(token, "enc("), ")"), nil
}

func (fakeRegistrationTokenAlg) LegacyTokenEnabled() bool { return false }

func Test_registrationAccessToken_roundTrip(t *testing.T) {
	s := &Server{encAlg: fakeRegistrationTokenAlg{}}

	tests := []struct {
		name     string
		clientID string
		orgID    string
		secret   string
	}{
		{"plain", "client1", "org1", "topsecret"},
		{"secret with colon is preserved", "client1", "org1", "top:sec:ret"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := s.registrationAccessToken(tt.clientID, tt.orgID, tt.secret)
			require.NoError(t, err)

			binding, err := s.parseRegistrationAccessToken(token)
			require.NoError(t, err)
			assert.Equal(t, tt.clientID, binding.clientID)
			assert.Equal(t, tt.orgID, binding.orgID)
			assert.Equal(t, tt.secret, binding.secret)
		})
	}
}

func Test_parseRegistrationAccessToken_errors(t *testing.T) {
	s := &Server{encAlg: fakeRegistrationTokenAlg{}}

	t.Run("undecryptable token errors", func(t *testing.T) {
		_, err := s.parseRegistrationAccessToken("not-encrypted")
		assert.Error(t, err)
	})

	t.Run("wrong number of parts is an invalid token format", func(t *testing.T) {
		token, err := fakeRegistrationTokenAlg{}.EncryptToken("client1:org1")
		require.NoError(t, err)
		_, err = s.parseRegistrationAccessToken(token)
		assert.ErrorIs(t, err, ErrInvalidTokenFormat)
	})
}

func Test_clientRegistrationResponseFromClient(t *testing.T) {
	client := &query.OIDCClient{
		ClientID:               "client1",
		RedirectURIs:           []string{"https://client.example.com/callback"},
		ResponseTypes:          []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
		GrantTypes:             []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
		ApplicationType:        domain.OIDCApplicationTypeWeb,
		AuthMethodType:         domain.OIDCAuthMethodTypeNone,
		PostLogoutRedirectURIs: []string{"https://client.example.com/logout"},
		HashedSecret:           "$plain$$secret",
	}

	resp := clientRegistrationResponseFromClient(client)

	assert.Equal(t, "client1", resp.ClientID)
	assert.Equal(t, []string{"https://client.example.com/callback"}, resp.RedirectURIs)
	assert.Equal(t, []string{"code"}, resp.ResponseTypes)
	assert.Equal(t, []string{"authorization_code"}, resp.GrantTypes)
	assert.Equal(t, "web", resp.ApplicationType)
	assert.Equal(t, "none", resp.TokenEndpointAuthMethod)
	assert.Equal(t, []string{"https://client.example.com/logout"}, resp.PostLogoutRedirectURIs)
	// The stored secret hash must never be returned, and a read does not (re)issue a token.
	assert.Empty(t, resp.ClientSecret)
	assert.Empty(t, resp.RegistrationAccessToken)
}
