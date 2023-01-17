package google

import (
	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/oidc"
)

const (
	issuer = "https://accounts.google.com"
	name   = "Google"
)

var _ idp.Provider = (*Provider)(nil)

// Provider is the idp.Provider implementation for Google
type Provider struct {
	*oidc.Provider
}

// New creates a Google provider using the oidc.Provider (OIDC generic provider)
func New(clientID, clientSecret, redirectURI string, opts ...oidc.ProviderOpts) (*Provider, error) {
	rp, err := oidc.New(name, issuer, clientID, clientSecret, redirectURI, opts...)
	if err != nil {
		return nil, err
	}
	provider := &Provider{
		Provider: rp,
	}

	return provider, nil
}
