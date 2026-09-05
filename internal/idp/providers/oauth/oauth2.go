package oauth

import (
	"context"
	"maps"
	"net/http"
	"slices"

	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"golang.org/x/oauth2"

	"github.com/zitadel/zitadel/internal/idp"
)

var _ idp.Provider = (*Provider)(nil)

// Provider is the [idp.Provider] implementation for a generic OAuth 2.0 provider
type Provider struct {
	rp.RelyingParty
	options                 []rp.Option
	name                    string
	userEndpoint            string
	user                    func() idp.User
	isLinkingAllowed        bool
	isCreationAllowed       bool
	isAutoCreation          bool
	isAutoUpdate            bool
	authorizationParameters map[string]string
	forwardedParameters     []string
	generateVerifier        func() string
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

// WithAuthorizationParameters sets additional parameters on the authorization request sent to the
// provider. They take precedence over any parameter ZITADEL sets itself, e.g. `prompt=select_account`.
// A parameter with an empty value is omitted from the request altogether.
func WithAuthorizationParameters(parameters map[string]string) ProviderOpts {
	return func(p *Provider) {
		p.authorizationParameters = parameters
	}
}

// WithForwardedParameters allows the listed parameters of the original authorization request
// to be forwarded to the provider. Parameters set by [WithAuthorizationParameters] take precedence.
func WithForwardedParameters(parameters []string) ProviderOpts {
	return func(p *Provider) {
		p.forwardedParameters = parameters
	}
}

// New creates a generic OAuth 2.0 provider
func New(config *oauth2.Config, name, userEndpoint string, user func() idp.User, httpClient *http.Client, options ...ProviderOpts) (provider *Provider, err error) {
	provider = &Provider{
		name:             name,
		userEndpoint:     userEndpoint,
		user:             user,
		generateVerifier: oauth2.GenerateVerifier,
		options:          []rp.Option{rp.WithHTTPClient(httpClient)},
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
			opts = append(opts, urlParam("login_hint", string(username)))
		}
	}
	authParams := idp.ResolveAuthorizationParameters(p.authorizationParameters, p.forwardedParameters, params)
	if _, promptSet := authParams["prompt"]; !promptSet && !loginHintSet {
		opts = append(opts, rp.WithPrompt(oidc.PromptSelectAccount))
	}

	var codeVerifier string
	if p.RelyingParty.IsPKCE() {
		codeVerifier = p.generateVerifier()
		opts = append(opts, rp.WithCodeChallenge(oidc.NewSHACodeChallenge(codeVerifier)))
	}
	// applied last so that the configured parameters take precedence over the ones set by ZITADEL
	opts = append(opts, urlParams(authParams)...)

	url := rp.AuthURL(state, p.RelyingParty, opts...)
	return &Session{AuthURL: url, Provider: p, CodeVerifier: codeVerifier}, nil
}

// urlParams returns the options to set the parameters on the authorization request.
// Parameters without a value are skipped, which removes ZITADEL's default if there is any.
// The keys are sorted to keep the resulting authorization URL stable.
func urlParams(parameters map[string]string) []rp.AuthURLOpt {
	opts := make([]rp.AuthURLOpt, 0, len(parameters))
	for _, key := range slices.Sorted(maps.Keys(parameters)) {
		if parameters[key] == "" {
			continue
		}
		opts = append(opts, urlParam(key, parameters[key]))
	}
	return opts
}

func urlParam(key, value string) rp.AuthURLOpt {
	return func() []oauth2.AuthCodeOption {
		return []oauth2.AuthCodeOption{oauth2.SetAuthURLParam(key, value)}
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
