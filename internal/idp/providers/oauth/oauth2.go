package oauth

import (
	"context"
	"errors"

	"github.com/zitadel/oidc/v2/pkg/client/rp"
	"golang.org/x/oauth2"

	"github.com/zitadel/zitadel/internal/idp"
)

var _ idp.Provider = (*Provider)(nil)

var ErrCodeMissing = errors.New("no auth code provided")

// Provider is the idp.Provider implementation for a generic OAuth 2.0 provider
type Provider struct {
	rp.RelyingParty
	options           []rp.Option
	name              string
	userEndpoint      string
	userMapper        func() UserInfoMapper
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

func WithRelyingPartyOption(option rp.Option) ProviderOpts {
	return func(p *Provider) {
		p.options = append(p.options, option)
	}
}

func New(config *oauth2.Config, name, userEndpoint string, userMapper func() UserInfoMapper, options ...ProviderOpts) (*Provider, error) {
	provider := &Provider{
		name:         name,
		userEndpoint: userEndpoint,
		userMapper:   userMapper,
	}
	for _, option := range options {
		option(provider)
	}
	relyingParty, err := rp.NewRelyingPartyOAuth(config, provider.options...)
	if err != nil {
		return nil, err
	}
	provider.RelyingParty = relyingParty
	return provider, nil
}

func (p *Provider) Name() string {
	return p.name
}

func (p *Provider) BeginAuth(ctx context.Context, state string, _ ...any) (idp.Session, error) {
	url := rp.AuthURL(state, p.RelyingParty)
	return &Session{AuthURL: url, Provider: p}, nil
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
