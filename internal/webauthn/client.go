package webauthn

import (
	"fmt"

	"github.com/descope/virtualwebauthn"
)

type Client struct {
	rp         virtualwebauthn.RelyingParty
	auth       virtualwebauthn.Authenticator
	credential virtualwebauthn.Credential
}

func NewClient(name, domain, origin string) *Client {
	rp := virtualwebauthn.RelyingParty{
		Name:   name,
		ID:     domain,
		Origin: origin,
	}
	return &Client{
		rp:         rp,
		auth:       virtualwebauthn.NewAuthenticator(),
		credential: virtualwebauthn.NewCredential(virtualwebauthn.KeyTypeEC2),
	}
}

func (c *Client) CreateAttestationResponse(options []byte) ([]byte, error) {
	parsedAttestationOptions, err := virtualwebauthn.ParseAttestationOptions(string(options))
	if err != nil {
		return nil, fmt.Errorf("webauthn.Client.CreateAttestationResponse: %w", err)
	}
	return []byte(virtualwebauthn.CreateAttestationResponse(
		c.rp, c.auth, c.credential, *parsedAttestationOptions,
	)), nil
}
