package oidc

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	oidc_crypto "github.com/zitadel/oidc/v3/pkg/crypto"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
)

// TestSession_Authorize_ClientAuth checks the token request ZITADEL sends for each token endpoint auth method.
// The token endpoint answers without an id_token, so a successful exchange ends in [rp.ErrMissingIDToken],
// while rejected credentials end in an oauth2 error.
func TestSession_Authorize_ClientAuth(t *testing.T) {
	const clientID = "client-id"
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	pkcs8, err := x509.MarshalPKCS8PrivateKey(key)
	require.NoError(t, err)
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: pkcs8})

	tests := []struct {
		name         string
		method       domain.OIDCTokenEndpointAuthMethod
		clientSecret *crypto.CryptoValue
		privateKey   *crypto.CryptoValue
		keyID        string
		check        func(t *testing.T, r *http.Request, issuer string) bool
	}{
		{
			name:         "client_secret uses basic auth",
			method:       domain.OIDCTokenEndpointAuthMethodClientSecret,
			clientSecret: encrypted([]byte("secret")),
			check: func(t *testing.T, r *http.Request, _ string) bool {
				user, password, ok := r.BasicAuth()
				return assert.True(t, ok) &&
					assert.Equal(t, clientID, user) &&
					assert.Equal(t, "secret", password) &&
					assert.Empty(t, r.PostForm.Get("client_assertion"))
			},
		},
		{
			name:       "private_key_jwt sends a signed client assertion",
			method:     domain.OIDCTokenEndpointAuthMethodPrivateKeyJWT,
			privateKey: encrypted(keyPEM),
			keyID:      "key-id",
			check: func(t *testing.T, r *http.Request, issuer string) bool {
				_, _, basic := r.BasicAuth()
				if !assert.False(t, basic) ||
					!assert.Equal(t, clientID, r.PostForm.Get("client_id")) ||
					!assert.False(t, r.PostForm.Has("client_secret")) ||
					!assert.Equal(t, oidc.ClientAssertionTypeJWTAssertion, r.PostForm.Get("client_assertion_type")) {
					return false
				}
				jws, err := jose.ParseSigned(r.PostForm.Get("client_assertion"), []jose.SignatureAlgorithm{jose.RS256})
				if !assert.NoError(t, err) || !assert.Equal(t, "key-id", jws.Signatures[0].Header.KeyID) {
					return false
				}
				payload, err := jws.Verify(&key.PublicKey)
				if !assert.NoError(t, err) {
					return false
				}
				var claims clientAssertionClaims
				return assert.NoError(t, json.Unmarshal(payload, &claims)) &&
					assert.Equal(t, clientID, claims.Issuer) &&
					assert.Equal(t, clientID, claims.Subject) &&
					assert.Contains(t, claims.Audience, issuer+"/token") &&
					assert.NotEmpty(t, claims.JWTID) &&
					assert.LessOrEqual(t, claims.ExpiresAt.AsTime().Sub(claims.IssuedAt.AsTime()), clientAssertionLifetime) &&
					assert.WithinDuration(t, time.Now(), claims.IssuedAt.AsTime(), time.Minute)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var issuer string
			mux := http.NewServeMux()
			mux.HandleFunc("/.well-known/openid-configuration", func(w http.ResponseWriter, _ *http.Request) {
				_ = json.NewEncoder(w).Encode(map[string]string{
					"issuer":                 issuer,
					"authorization_endpoint": issuer + "/authorize",
					"token_endpoint":         issuer + "/token",
					"jwks_uri":               issuer + "/keys",
				})
			})
			mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
				// assert, not require: FailNow must not be called from the server's handler goroutine
				if !assert.NoError(t, r.ParseForm()) || !tt.check(t, r, issuer) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusUnauthorized)
					_, _ = w.Write([]byte(`{"error":"invalid_client"}`))
					return
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"access_token":"access","token_type":"Bearer","expires_in":3600}`))
			})
			server := httptest.NewServer(mux)
			defer server.Close()
			issuer = server.URL

			secret, opts, err := ClientAuth(tt.method, tt.clientSecret, tt.privateKey, tt.keyID, crypto.CreateMockEncryptionAlg(gomock.NewController(t)))
			require.NoError(t, err)
			provider, err := New("name", issuer, clientID, secret, "https://zitadel.cloud/callback", nil, DefaultMapper, server.Client(), opts...)
			require.NoError(t, err)

			begin, err := provider.BeginAuth(context.Background(), "state")
			require.NoError(t, err)
			session := NewSession(provider, "code", begin.(*Session).PersistentParameters())

			err = session.Authorize(context.Background())
			assert.ErrorIs(t, err, rp.ErrMissingIDToken)
		})
	}
}

func TestClientAuth_invalidCredentials(t *testing.T) {
	alg := crypto.CreateMockEncryptionAlg(gomock.NewController(t))

	_, _, err := ClientAuth(domain.OIDCTokenEndpointAuthMethodClientSecret, nil, nil, "", alg)
	assert.Error(t, err, "client secret missing")

	_, _, err = ClientAuth(domain.OIDCTokenEndpointAuthMethodPrivateKeyJWT, nil, nil, "key-id", alg)
	assert.Error(t, err, "private key missing")

	_, opts, err := ClientAuth(domain.OIDCTokenEndpointAuthMethodPrivateKeyJWT, nil, encrypted([]byte("no key")), "key-id", alg)
	require.NoError(t, err)
	_, err = New("name", "https://issuer.invalid", "client", "", "https://zitadel.cloud/callback", nil, DefaultMapper, http.DefaultClient, opts...)
	assert.ErrorIs(t, err, oidc_crypto.ErrPEMDecode)
}

func encrypted(value []byte) *crypto.CryptoValue {
	return &crypto.CryptoValue{
		CryptoType: crypto.TypeEncryption,
		Algorithm:  "enc",
		KeyID:      "id",
		Crypted:    value,
	}
}
