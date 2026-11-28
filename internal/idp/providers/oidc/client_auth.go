package oidc

import (
	"crypto/rand"
	"time"

	"github.com/zitadel/oidc/v3/pkg/client/rp"
	oidc_crypto "github.com/zitadel/oidc/v3/pkg/crypto"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"golang.org/x/oauth2"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
)

// clientAssertionLifetime is kept short, since a new assertion is created for every token request.
const clientAssertionLifetime = 5 * time.Minute

// ClientAuth decrypts the credentials of the token endpoint auth method.
// It returns the client secret, which is empty for private_key_jwt, and the provider options for the method.
func ClientAuth(method domain.OIDCTokenEndpointAuthMethod, clientSecret, privateKey *crypto.CryptoValue, keyID string, alg crypto.EncryptionAlgorithm) (string, []ProviderOpts, error) {
	switch method {
	case domain.OIDCTokenEndpointAuthMethodPrivateKeyJWT:
		key, err := crypto.Decrypt(privateKey, alg)
		if err != nil {
			return "", nil, err
		}
		return "", []ProviderOpts{WithPrivateKeyJWT(key, keyID)}, nil
	default:
		secret, err := crypto.DecryptString(clientSecret, alg)
		return secret, nil, err
	}
}

// WithPrivateKeyJWT authenticates on the token endpoint with a client assertion (RFC 7523)
// signed by the provided PEM encoded key. The keyID is sent as `kid` header.
func WithPrivateKeyJWT(key []byte, keyID string) ProviderOpts {
	return func(p *Provider) {
		p.options = append(p.options,
			rp.WithJWTProfile(rp.SignerFromKeyAndKeyID(key, keyID)),
			// send the client_id in the body and never a basic auth header with an empty secret
			rp.WithAuthStyle(oauth2.AuthStyleInParams),
		)
	}
}

type clientAssertionClaims struct {
	Issuer    string        `json:"iss"`
	Subject   string        `json:"sub"`
	Audience  oidc.Audience `json:"aud"`
	JWTID     string        `json:"jti"`
	IssuedAt  oidc.Time     `json:"iat"`
	ExpiresAt oidc.Time     `json:"exp"`
}

// clientAssertion creates a signed client assertion if the relying party is configured for private_key_jwt.
// Like zitadel/oidc's CodeExchangeHandler, the audience contains the issuer and the token endpoint.
func clientAssertion(relyingParty rp.RelyingParty) (string, error) {
	signer := relyingParty.Signer()
	if signer == nil {
		return "", nil
	}
	clientID := relyingParty.OAuthConfig().ClientID
	now := time.Now()
	return oidc_crypto.Sign(&clientAssertionClaims{
		Issuer:    clientID,
		Subject:   clientID,
		Audience:  oidc.Audience{relyingParty.Issuer(), relyingParty.OAuthConfig().Endpoint.TokenURL},
		JWTID:     rand.Text(),
		IssuedAt:  oidc.FromTime(now),
		ExpiresAt: oidc.FromTime(now.Add(clientAssertionLifetime)),
	}, signer)
}
