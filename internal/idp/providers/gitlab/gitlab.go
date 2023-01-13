package google

import (
	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/oidc"
)

const (
	issuer = "https://gitlab.com"
)

var _ idp.Provider = (*Provider)(nil)

type Provider struct {
	oidcProvider *oidc.Provider
}

func New(clientID, clientSecret, redirectURI string) (*Provider, error) {
	rp, err := oidc.New(issuer, clientID, clientSecret, redirectURI)
	if err != nil {
		return nil, err
	}
	provider := &Provider{
		oidcProvider: rp,
	}

	return provider, nil
}

func (p *Provider) Name() string {
	return "gitlab"
}

func (p *Provider) BeginAuth(state string) (idp.Session, error) {
	return p.oidcProvider.BeginAuth(state)
}

func (p *Provider) FetchUser(session idp.Session) (idp.User, error) {
	return p.oidcProvider.FetchUser(session)
}
