package saml

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"net/url"

	"github.com/crewjam/saml"
	"github.com/crewjam/saml/samlsp"

	"github.com/zitadel/zitadel/internal/idp"
)

var _ idp.Provider = (*Provider)(nil)

// Provider is the [idp.Provider] implementation for a generic SAML provider
type Provider struct {
	name string

	spOptions *samlsp.Options
	metadata  *saml.EntityDescriptor

	rootURL *url.URL
	keypair *tls.Certificate
	binding string

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

func WithSignedRequest() ProviderOpts {
	return func(p *Provider) {
		p.spOptions.SignRequest = true
	}
}

func WithBinding(binding string) ProviderOpts {
	return func(p *Provider) {
		p.binding = binding
	}
}

func New(
	name string,
	rootURLStr string,
	metadata *saml.EntityDescriptor,
	certificate []byte,
	key []byte,
	options ...ProviderOpts,
) (*Provider, error) {
	keyPair, err := tls.X509KeyPair(certificate, key)
	if err != nil {
		return nil, err
	}
	keyPair.Leaf, err = x509.ParseCertificate(keyPair.Certificate[0])
	if err != nil {
		return nil, err
	}
	rootURL, err := url.Parse(rootURLStr)
	if err != nil {
		return nil, err
	}
	opts := samlsp.Options{
		URL:         *rootURL,
		Key:         keyPair.PrivateKey.(*rsa.PrivateKey),
		Certificate: keyPair.Leaf,
		IDPMetadata: metadata,
		SignRequest: false,
		RelayStateFunc: func(w http.ResponseWriter, r *http.Request) string {
			return r.URL.String()
		},
	}
	provider := &Provider{
		name:      name,
		spOptions: &opts,
	}
	for _, option := range options {
		option(provider)
	}
	return provider, nil
}

func (p *Provider) Name() string {
	return p.name
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

func (p *Provider) BeginAuth(ctx context.Context, state string, params ...any) (idp.Session, error) {
	m, err := samlsp.New(*p.spOptions)
	if err != nil {
		return nil, err
	}

	return &Session{
		serviceProvider: m,
		state:           state,
	}, nil
}
