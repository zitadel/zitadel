package oauth

import (
	"context"

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
	user              func() idp.User
	isLinkingAllowed  bool
	isCreationAllowed bool
	isAutoCreation    bool
	isAutoUpdate      bool
	generateVerifier  func() string
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
func New(config *oauth2.Config, name, userEndpoint string, user func() idp.User, options ...ProviderOpts) (provider *Provider, err error) {
	provider = &Provider{
		name:             name,
		userEndpoint:     userEndpoint,
		user:             user,
		generateVerifier: oauth2.GenerateVerifier,
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
func (p *Provider) BeginAuth(ctx context.Context, state string, params ...idp.Parameter) (idp.Session, error) {
	opts := make([]rp.AuthURLOpt, 0)
	var loginHintSet bool
	for _, param := range params {
		if username, ok := param.(idp.LoginHintParam); ok {
			loginHintSet = true
			opts = append(opts, loginHint(string(username)))
		}
	}
	if !loginHintSet {
		opts = append(opts, rp.WithPrompt(oidc.PromptSelectAccount))
	}

	var codeVerifier string
	if p.RelyingParty.IsPKCE() {
		codeVerifier = p.generateVerifier()
		opts = append(opts, rp.WithCodeChallenge(oidc.NewSHACodeChallenge(codeVerifier)))
	}

	url := rp.AuthURL(state, p.RelyingParty, opts...)
	return &Session{AuthURL: url, Provider: p, CodeVerifier: codeVerifier}, nil
}

func loginHint(hint string) rp.AuthURLOpt {
	return func() []oauth2.AuthCodeOption {
		return []oauth2.AuthCodeOption{oauth2.SetAuthURLParam("login_hint", hint)}
	}
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

func (p *Provider) User() idp.User {
	return p.user()
}
