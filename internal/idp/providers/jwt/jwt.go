package jwt

import (
	"errors"
	"net/url"

	"github.com/zitadel/zitadel/internal/idp"
)

const queryAuthRequestID = "authRequestID"

var _ idp.Provider = (*Provider)(nil)

var ErrNoTokens = errors.New("no tokens")

// Provider is the idp.Provider implementation for a JWT provider
type Provider struct {
	name              string
	headerName        string
	issuer            string
	jwtEndpoint       string
	keysEndpoint      string
	isLinkingAllowed  bool
	isCreationAllowed bool
	isAutoCreation    bool
	isAutoUpdate      bool
}

type ProviderOpts func(provider *Provider)

func WithLinkingAllowed() ProviderOpts {
	return func(p *Provider) {
		p.isLinkingAllowed = true
	}
}
func WithCreationAllowed() ProviderOpts {
	return func(p *Provider) {
		p.isCreationAllowed = true
	}
}
func WithAutoCreation() ProviderOpts {
	return func(p *Provider) {
		p.isAutoCreation = true
	}
}
func WithAutoUpdate() ProviderOpts {
	return func(p *Provider) {
		p.isAutoUpdate = true
	}
}

func New(name, issuer, jwtEndpoint, keysEndpoint, headerName string, options ...ProviderOpts) (*Provider, error) {
	provider := &Provider{
		name:         name,
		issuer:       issuer,
		jwtEndpoint:  jwtEndpoint,
		keysEndpoint: keysEndpoint,
		headerName:   headerName,
	}
	for _, option := range options {
		option(provider)
	}

	return provider, nil
}

func (p *Provider) Name() string {
	return p.name
}

func (p *Provider) BeginAuth(state string) (idp.Session, error) {
	redirect, err := url.Parse(p.jwtEndpoint)
	if err != nil {
		return nil, err
	}
	q := redirect.Query()
	q.Set(queryAuthRequestID, state)
	//TODO: userAgentID
	redirect.RawQuery = q.Encode()
	return &Session{AuthURL: redirect.String()}, nil
}

func (p *Provider) IsLinkingAllowed() bool {
	return p.isLinkingAllowed
}

func (p *Provider) IsCreationAllowed() bool {
	return p.isCreationAllowed
}

func (p *Provider) IsAutoCreation() bool {
	return p.isAutoCreation
}

func (p *Provider) IsAutoUpdate() bool {
	return p.isAutoUpdate
}
