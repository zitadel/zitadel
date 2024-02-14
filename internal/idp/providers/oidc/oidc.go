package oidc

import (
	"context"

	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"golang.org/x/oauth2"

	"github.com/zitadel/zitadel/internal/idp"
)

var _ idp.Provider = (*Provider)(nil)

// Provider is the [idp.Provider] implementation for a generic OIDC provider
type Provider struct {
	rp.RelyingParty
	options           []rp.Option
	name              string
	isLinkingAllowed  bool
	isCreationAllowed bool
	isAutoCreation    bool
	isAutoUpdate      bool
	useIDToken        bool
	userInfoMapper    func(info *oidc.UserInfo) idp.User
	authOptions       []func(bool) rp.AuthURLOpt
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

// WithIDTokenMapping enables that information to map the user is retrieved from the id_token and not the userinfo endpoint.
func WithIDTokenMapping() ProviderOpts {
	return func(p *Provider) {
		p.useIDToken = true
	}
}

// WithRelyingPartyOption allows to set an additional [rp.Option] like [rp.WithPKCE].
func WithRelyingPartyOption(option rp.Option) ProviderOpts {
	return func(p *Provider) {
		p.options = append(p.options, option)
	}
}

// WithSelectAccount adds the select_account prompt to the auth request (if no login_hint is set)
func WithSelectAccount() ProviderOpts {
	return func(p *Provider) {
		p.authOptions = append(p.authOptions, func(loginHintSet bool) rp.AuthURLOpt {
			if loginHintSet {
				return nil
			}
			return rp.WithPrompt(oidc.PromptSelectAccount)
		})
	}
}

// WithResponseMode sets the `response_mode` params in the auth request
func WithResponseMode(mode oidc.ResponseMode) ProviderOpts {
	return func(p *Provider) {
		paramOpt := rp.WithResponseModeURLParam(mode)
		p.authOptions = append(p.authOptions, func(_ bool) rp.AuthURLOpt {
			return rp.AuthURLOpt(paramOpt)
		})
	}
}

type UserInfoMapper func(info *oidc.UserInfo) idp.User

var DefaultMapper UserInfoMapper = func(info *oidc.UserInfo) idp.User {
	return NewUser(info)
}

// New creates a generic OIDC provider
func New(name, issuer, clientID, clientSecret, redirectURI string, scopes []string, userInfoMapper UserInfoMapper, options ...ProviderOpts) (provider *Provider, err error) {
	provider = &Provider{
		name:           name,
		userInfoMapper: userInfoMapper,
	}
	for _, option := range options {
		option(provider)
	}
	provider.RelyingParty, err = rp.NewRelyingPartyOIDC(context.TODO(), issuer, clientID, clientSecret, redirectURI, setDefaultScope(scopes), provider.options...)
	if err != nil {
		return nil, err
	}
	return provider, nil
}

// setDefaultScope ensures that at least openid ist set
// if none is provided it will request `openid profile email phone`
func setDefaultScope(scopes []string) []string {
	if len(scopes) == 0 {
		return []string{oidc.ScopeOpenID, oidc.ScopeProfile, oidc.ScopeEmail, oidc.ScopePhone}
	}
	for _, scope := range scopes {
		if scope == oidc.ScopeOpenID {
			return scopes
		}
	}
	return append(scopes, oidc.ScopeOpenID)
}

// Name implements the [idp.Provider] interface
func (p *Provider) Name() string {
	return p.name
}

// BeginAuth implements the [idp.Provider] interface.
// It will create a [Session] with an OIDC authorization request as AuthURL.
func (p *Provider) BeginAuth(ctx context.Context, state string, params ...idp.Parameter) (idp.Session, error) {
	opts := make([]rp.AuthURLOpt, 0)
	var loginHintSet bool
	for _, param := range params {
		if username, ok := param.(idp.LoginHintParam); ok {
			loginHintSet = true
			opts = append(opts, loginHint(string(username)))
		}
	}
	for _, option := range p.authOptions {
		if opt := option(loginHintSet); opt != nil {
			opts = append(opts, opt)
		}
	}
	url := rp.AuthURL(state, p.RelyingParty, opts...)
	return &Session{AuthURL: url, Provider: p}, nil
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
