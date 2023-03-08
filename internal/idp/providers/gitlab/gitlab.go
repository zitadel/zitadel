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

// Provider is the [idp.Provider] implementation for Gitlab
type Provider struct {
	*oidc.Provider
}

// New creates a GitLab.com provider using the [oidc.Provider] (OIDC generic provider)
func New(clientID, clientSecret, redirectURI string, scopes []string, options ...oidc.ProviderOpts) (*Provider, error) {
	return NewCustomIssuer(name, issuer, clientID, clientSecret, redirectURI, scopes, options...)
}

// NewCustomIssuer creates a GitLab provider using the [oidc.Provider] (OIDC generic provider)
// with a custom issuer for self-managed instances
func NewCustomIssuer(name, issuer, clientID, clientSecret, redirectURI string, scopes []string, options ...oidc.ProviderOpts) (*Provider, error) {
	rp, err := oidc.New(name, issuer, clientID, clientSecret, redirectURI, scopes, oidc.DefaultMapper, options...)
	if err != nil {
		return nil, err
	}
	return &Provider{
		Provider: rp,
	}, nil
}
