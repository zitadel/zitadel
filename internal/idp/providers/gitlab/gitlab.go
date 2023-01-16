package gitlab

import (
	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/oidc"
)

const (
	issuer = "https://gitlab.com"
	name   = "GitLab"
)

var _ idp.Provider = (*Provider)(nil)

// Provider is the idp.Provider implementation for Gitlab
type Provider struct {
	*oidc.Provider
}

func New(clientID, clientSecret, redirectURI string, options ...oidc.ProviderOpts) (*Provider, error) {
	return NewCustomIssuer(name, issuer, clientID, clientSecret, redirectURI, options...)
}

func NewCustomIssuer(name, issuer, clientID, clientSecret, redirectURI string, options ...oidc.ProviderOpts) (*Provider, error) {
	rp, err := oidc.New(name, issuer, clientID, clientSecret, redirectURI, options...)
	if err != nil {
		return nil, err
	}
	return &Provider{
		Provider: rp,
	}, nil
}
