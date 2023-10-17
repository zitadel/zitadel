package apple

import (
	"crypto/x509"
	"encoding/pem"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/zitadel/oidc/v3/pkg/crypto"
	openid "github.com/zitadel/oidc/v3/pkg/oidc"

	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/oidc"
)

const (
	name   = "Apple"
	issuer = "https://appleid.apple.com"
)

var _ idp.Provider = (*Provider)(nil)

// Provider is the [idp.Provider] implementation for Apple
type Provider struct {
	*oidc.Provider
}

func New(clientID, teamID, keyID, callbackURL string, key []byte, scopes []string, options ...oidc.ProviderOpts) (*Provider, error) {
	secret, err := clientSecretFromPrivateKey(key, teamID, clientID, keyID)
	if err != nil {
		return nil, err
	}
	options = append(options, oidc.WithResponseMode("form_post"))
	rp, err := oidc.New(name, issuer, clientID, secret, callbackURL, scopes, oidc.DefaultMapper, options...)
	if err != nil {
		return nil, err
	}
	return &Provider{
		Provider: rp,
	}, nil
}

// clientSecretFromPrivateKey uses the private key to create and sign a JWT, which has to be used as client_secret at Apple.
func clientSecretFromPrivateKey(key []byte, teamID, clientID, keyID string) (string, error) {
	block, _ := pem.Decode(key)
	b := block.Bytes
	pk, err := x509.ParsePKCS8PrivateKey(b)
	if err != nil {
		return "", err
	}
	signingKey := jose.SigningKey{
		Algorithm: jose.ES256,
		Key:       &jose.JSONWebKey{Key: pk, KeyID: keyID},
	}
	signer, err := jose.NewSigner(signingKey, &jose.SignerOptions{})
	if err != nil {
		return "", err
	}
	iat := time.Now()
	exp := iat.Add(time.Hour)
	return crypto.Sign(&openid.JWTTokenRequest{
		Issuer:    teamID,
		Subject:   clientID,
		Audience:  []string{issuer},
		ExpiresAt: openid.FromTime(exp),
		IssuedAt:  openid.FromTime(iat),
	}, signer)
}
