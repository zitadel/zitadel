package oauth

import (
	"context"
	"os"

	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"golang.org/x/oauth2"

	"github.com/zitadel/zitadel/internal/idp"
)

var _ idp.Provider = (*Provider)(nil)

// Provider is the [idp.Provider] implementation for a generic OAuth 2.0 provider
type Provider struct {
	rp.RelyingParty
	options           []rp.Option
	name              string
	userEndpoint      string
	userMapper        func() idp.User
	isLinkingAllowed  bool
	isCreationAllowed bool
	isAutoCreation    bool
	isAutoUpdate      bool
}

type ProviderOpts func(provider *Provider)

// WithLinkingAllowed allows end users to link the federated user to an existing one.
func WithLinkingAllowed() ProviderOpts {
	return func(p *Provider) {
		p.isLinkingAllowed = true
	}
}

// WithCreationAllowed allows end users to create a new user using the federated information.
func WithCreationAllowed() ProviderOpts {
	return func(p *Provider) {
		p.isCreationAllowed = true
	}
}

// WithAutoCreation enables that federated users are automatically created if not already existing.
func WithAutoCreation() ProviderOpts {
	return func(p *Provider) {
		p.isAutoCreation = true
	}
}

// WithAutoUpdate enables that information retrieved from the provider is automatically used to update
// the existing user on each authentication.
func WithAutoUpdate() ProviderOpts {
	return func(p *Provider) {
		p.isAutoUpdate = true
	}
}

// WithRelyingPartyOption allows to set an additional [rp.Option] like [rp.WithPKCE].
func WithRelyingPartyOption(option rp.Option) ProviderOpts {
	return func(p *Provider) {
		p.options = append(p.options, option)
	}
}

// New creates a generic OAuth 2.0 provider
func New(config *oauth2.Config, name, userEndpoint string, userMapper func() idp.User, options ...ProviderOpts) (provider *Provider, err error) {
	provider = &Provider{
		name:         name,
		userEndpoint: userEndpoint,
		userMapper:   userMapper,
	}
	for _, option := range options {
		option(provider)
	}
	provider.RelyingParty, err = rp.NewRelyingPartyOAuth(config, provider.options...)
	if err != nil {
		return nil, err
	}
	return provider, nil
}

// Name implements the [idp.Provider] interface
func (p *Provider) Name() string {
	return p.name
}

// BeginAuth implements the [idp.Provider] interface.
// It will create a [Session] with an OAuth2.0 authorization request as AuthURL.
func (p *Provider) BeginAuth(ctx context.Context, state string, params ...any) (idp.Session, error) {
	opts := []rp.AuthURLOpt{}
	// https://github.com/zitadel/zitadel/blob/main/proto/zitadel/oidc/v2beta/authorization.proto
	switch prompt := os.Getenv("ZITADEL_GENERIC_OAUTH_AUTH_PROMPT"); prompt {
	case "PROMPT_UNSPECIFIED":
		// When set to "UNSPECIFIED" send no "prompt" param
	case "PROMPT_NONE":
		opts = append(opts, rp.WithPrompt(oidc.PromptNone))
	case "PROMPT_LOGIN":
		opts = append(opts, rp.WithPrompt(oidc.PromptLogin))
	case "PROMPT_CONSENT":
		opts = append(opts, rp.WithPrompt(oidc.PromptConsent))
	case "PROMPT_SELECT_ACCOUNT":
		opts = append(opts, rp.WithPrompt(oidc.PromptSelectAccount))
	default:
		// previous behavior must be default to prevent a breaking change
		opts = append(opts, rp.WithPrompt(oidc.PromptSelectAccount))
	}

	for _, param := range params {
		if option, ok := param.(rp.AuthURLOpt); ok {
			opts = append(opts, option)
		}
	}
	url := rp.AuthURL(state, p.RelyingParty, opts...)
	return &Session{AuthURL: url, Provider: p}, nil
}

// IsLinkingAllowed implements the [idp.Provider] interface.
func (p *Provider) IsLinkingAllowed() bool {
	return p.isLinkingAllowed
}

// IsCreationAllowed implements the [idp.Provider] interface.
func (p *Provider) IsCreationAllowed() bool {
	return p.isCreationAllowed
}

// IsAutoCreation implements the [idp.Provider] interface.
func (p *Provider) IsAutoCreation() bool {
	return p.isAutoCreation
}

// IsAutoUpdate implements the [idp.Provider] interface.
func (p *Provider) IsAutoUpdate() bool {
	return p.isAutoUpdate
}
