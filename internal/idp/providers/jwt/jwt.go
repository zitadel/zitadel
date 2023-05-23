package jwt

import (
	"context"
	"encoding/base64"
	"errors"
	"net/url"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/idp"
)

const (
	queryAuthRequestID = "authRequestID"
	queryUserAgentID   = "userAgentID"
)

var _ idp.Provider = (*Provider)(nil)

var (
	ErrMissingUserAgentID = errors.New("userAgentID missing")
)

// Provider is the [idp.Provider] implementation for a JWT provider
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
	encryptionAlg     crypto.EncryptionAlgorithm
}

type ProviderOpts func(provider *Provider)

// WithLinkingAllowed allows end users to link the federated user to an existing one
func WithLinkingAllowed() ProviderOpts {
	return func(p *Provider) {
		p.isLinkingAllowed = true
	}
}

// WithCreationAllowed allows end users to create a new user using the federated information
func WithCreationAllowed() ProviderOpts {
	return func(p *Provider) {
		p.isCreationAllowed = true
	}
}

// WithAutoCreation enables that federated users are automatically created if not already existing
func WithAutoCreation() ProviderOpts {
	return func(p *Provider) {
		p.isAutoCreation = true
	}
}

// WithAutoUpdate enables that information retrieved from the provider is automatically used to update
// the existing user on each authentication
func WithAutoUpdate() ProviderOpts {
	return func(p *Provider) {
		p.isAutoUpdate = true
	}
}

// New creates a JWT provider
func New(name, issuer, jwtEndpoint, keysEndpoint, headerName string, encryptionAlg crypto.EncryptionAlgorithm, options ...ProviderOpts) (*Provider, error) {
	provider := &Provider{
		name:          name,
		issuer:        issuer,
		jwtEndpoint:   jwtEndpoint,
		keysEndpoint:  keysEndpoint,
		headerName:    headerName,
		encryptionAlg: encryptionAlg,
	}
	for _, option := range options {
		option(provider)
	}

	return provider, nil
}

// Name implements the [idp.Provider] interface
func (p *Provider) Name() string {
	return p.name
}

// BeginAuth implements the [idp.Provider] interface.
// It will create a [Session] with an AuthURL, pointing to the jwtEndpoint
// with the authRequest and encrypted userAgent ids.
func (p *Provider) BeginAuth(ctx context.Context, state string, params ...any) (idp.Session, error) {
	if len(params) < 1 {
		return nil, ErrMissingUserAgentID
	}
	userAgentID, ok := params[0].(string)
	if !ok {
		return nil, ErrMissingUserAgentID
	}
	redirect, err := url.Parse(p.jwtEndpoint)
	if err != nil {
		return nil, err
	}
	q := redirect.Query()
	q.Set(queryAuthRequestID, state)
	nonce, err := p.encryptionAlg.Encrypt([]byte(userAgentID))
	if err != nil {
		return nil, err
	}
	q.Set(queryUserAgentID, base64.RawURLEncoding.EncodeToString(nonce))
	redirect.RawQuery = q.Encode()
	return &Session{AuthURL: redirect.String()}, nil
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
