package oidc

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/domain"
)

func Test_clientRegistrationRequest_toOIDCApp(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		req     *clientRegistrationRequest
		want    *domain.OIDCApp
		wantErr string // expected registrationError.ErrorType, empty means success
	}{
		{
			name: "defaults applied",
			req: &clientRegistrationRequest{
				RedirectURIs: []string{"https://app.example.com/callback"},
			},
			want: &domain.OIDCApp{
				RedirectUris:    []string{"https://app.example.com/callback"},
				ResponseTypes:   []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
				GrantTypes:      []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
				ApplicationType: gu.Ptr(domain.OIDCApplicationTypeWeb),
				AuthMethodType:  gu.Ptr(domain.OIDCAuthMethodTypeBasic),
				OIDCVersion:     gu.Ptr(domain.OIDCVersionV1),
				AccessTokenType: gu.Ptr(domain.OIDCTokenTypeBearer),
			},
		},
		{
			name: "explicit values mapped and trimmed",
			req: &clientRegistrationRequest{
				RedirectURIs:            []string{" https://app.example.com/callback "},
				ResponseTypes:           []string{"code"},
				GrantTypes:              []string{"authorization_code", "refresh_token"},
				ApplicationType:         "native",
				ClientName:              "My MCP Client",
				TokenEndpointAuthMethod: "none",
				PostLogoutRedirectURIs:  []string{"https://app.example.com/logout"},
			},
			want: &domain.OIDCApp{
				AppName:                "My MCP Client",
				RedirectUris:           []string{"https://app.example.com/callback"},
				ResponseTypes:          []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
				GrantTypes:             []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode, domain.OIDCGrantTypeRefreshToken},
				ApplicationType:        gu.Ptr(domain.OIDCApplicationTypeNative),
				AuthMethodType:         gu.Ptr(domain.OIDCAuthMethodTypeNone),
				PostLogoutRedirectUris: []string{"https://app.example.com/logout"},
				OIDCVersion:            gu.Ptr(domain.OIDCVersionV1),
				AccessTokenType:        gu.Ptr(domain.OIDCTokenTypeBearer),
			},
		},
		{
			name: "client_secret_post auth method mapped",
			req: &clientRegistrationRequest{
				RedirectURIs:            []string{"https://app.example.com/callback"},
				TokenEndpointAuthMethod: "client_secret_post",
			},
			want: &domain.OIDCApp{
				RedirectUris:    []string{"https://app.example.com/callback"},
				ResponseTypes:   []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
				GrantTypes:      []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
				ApplicationType: gu.Ptr(domain.OIDCApplicationTypeWeb),
				AuthMethodType:  gu.Ptr(domain.OIDCAuthMethodTypePost),
				OIDCVersion:     gu.Ptr(domain.OIDCVersionV1),
				AccessTokenType: gu.Ptr(domain.OIDCTokenTypeBearer),
			},
		},
		{
			name: "native loopback http redirect allowed",
			req: &clientRegistrationRequest{
				RedirectURIs:            []string{"http://127.0.0.1:8080/callback"},
				ApplicationType:         "native",
				TokenEndpointAuthMethod: "none",
			},
			want: &domain.OIDCApp{
				RedirectUris:    []string{"http://127.0.0.1:8080/callback"},
				ResponseTypes:   []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
				GrantTypes:      []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
				ApplicationType: gu.Ptr(domain.OIDCApplicationTypeNative),
				AuthMethodType:  gu.Ptr(domain.OIDCAuthMethodTypeNone),
				OIDCVersion:     gu.Ptr(domain.OIDCVersionV1),
				AccessTokenType: gu.Ptr(domain.OIDCTokenTypeBearer),
			},
		},
		{
			name: "jwks_uri rejected",
			req: &clientRegistrationRequest{
				RedirectURIs: []string{"https://app.example.com/callback"},
				JWKsURI:      "https://app.example.com/jwks.json",
			},
			wantErr: registrationErrorInvalidClientMetadata,
		},
		{
			name: "private_key_jwt rejected",
			req: &clientRegistrationRequest{
				RedirectURIs:            []string{"https://app.example.com/callback"},
				TokenEndpointAuthMethod: "private_key_jwt",
			},
			wantErr: registrationErrorInvalidClientMetadata,
		},
		{
			name: "unknown grant_type rejected",
			req: &clientRegistrationRequest{
				RedirectURIs: []string{"https://app.example.com/callback"},
				GrantTypes:   []string{"client_credentials"},
			},
			wantErr: registrationErrorInvalidClientMetadata,
		},
		{
			name: "unknown response_type rejected",
			req: &clientRegistrationRequest{
				RedirectURIs:  []string{"https://app.example.com/callback"},
				ResponseTypes: []string{"token"},
			},
			wantErr: registrationErrorInvalidClientMetadata,
		},
		{
			name: "unknown application_type rejected",
			req: &clientRegistrationRequest{
				RedirectURIs:    []string{"https://app.example.com/callback"},
				ApplicationType: "service",
			},
			wantErr: registrationErrorInvalidClientMetadata,
		},
		{
			name: "missing redirect_uris rejected",
			req: &clientRegistrationRequest{
				ApplicationType: "web",
			},
			wantErr: registrationErrorInvalidRedirectURI,
		},
		{
			name: "custom scheme redirect for web app rejected",
			req: &clientRegistrationRequest{
				RedirectURIs:    []string{"myapp://callback"},
				ApplicationType: "web",
			},
			wantErr: registrationErrorInvalidRedirectURI,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.req.toOIDCApp()
			if tt.wantErr != "" {
				require.NotNil(t, err)
				assert.Equal(t, tt.wantErr, err.ErrorType)
				return
			}
			require.Nil(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_normalizeResponseType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in   string
		want string
	}{
		{"code", "code"},
		{"id_token", "id_token"},
		{"id_token token", "id_token token"},
		{"token id_token", "id_token token"},
		{"token", "token"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, normalizeResponseType(tt.in))
		})
	}
}

func Test_bearerToken(t *testing.T) {
	t.Parallel()
	newRequest := func(authorization string) *http.Request {
		r := httptest.NewRequest(http.MethodPost, "/oauth/v2/register", nil)
		if authorization != "" {
			r.Header.Set("Authorization", authorization)
		}
		return r
	}
	tests := []struct {
		name          string
		authorization string
		want          string
	}{
		{"no header", "", ""},
		{"bearer token", "Bearer abc123", "abc123"},
		{"case insensitive scheme", "bearer abc123", "abc123"},
		{"basic auth ignored", "Basic dXNlcjpwYXNz", ""},
		{"trims spaces", "Bearer   abc123  ", "abc123"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, bearerToken(newRequest(tt.authorization)))
		})
	}
}

func Test_newClientRegistrationResponse(t *testing.T) {
	t.Parallel()

	t.Run("confidential client returns secret and expiry", func(t *testing.T) {
		t.Parallel()
		app := &domain.OIDCApp{
			ClientID:           "client1",
			ClientSecretString: "secret",
			RedirectUris:       []string{"https://app.example.com/callback"},
			ResponseTypes:      []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
			GrantTypes:         []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
			ApplicationType:    gu.Ptr(domain.OIDCApplicationTypeWeb),
			AuthMethodType:     gu.Ptr(domain.OIDCAuthMethodTypeBasic),
		}
		got := newClientRegistrationResponse(app, "My Client", 1700000000)
		assert.Equal(t, &clientRegistrationResponse{
			ClientID:                "client1",
			ClientSecret:            "secret",
			ClientIDIssuedAt:        1700000000,
			ClientSecretExpiresAt:   gu.Ptr(int64(0)),
			RedirectURIs:            []string{"https://app.example.com/callback"},
			ResponseTypes:           []string{"code"},
			GrantTypes:              []string{"authorization_code"},
			ApplicationType:         "web",
			ClientName:              "My Client",
			TokenEndpointAuthMethod: "client_secret_basic",
		}, got)
	})

	t.Run("public client returns no secret nor expiry", func(t *testing.T) {
		t.Parallel()
		app := &domain.OIDCApp{
			ClientID:        "client2",
			RedirectUris:    []string{"https://app.example.com/callback"},
			ResponseTypes:   []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
			GrantTypes:      []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
			ApplicationType: gu.Ptr(domain.OIDCApplicationTypeNative),
			AuthMethodType:  gu.Ptr(domain.OIDCAuthMethodTypeNone),
		}
		got := newClientRegistrationResponse(app, "", 1700000000)
		assert.Empty(t, got.ClientSecret)
		assert.Nil(t, got.ClientSecretExpiresAt)
		assert.Equal(t, "none", got.TokenEndpointAuthMethod)
		assert.Equal(t, "native", got.ApplicationType)
	})
}
