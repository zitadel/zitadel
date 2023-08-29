package saml

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/xml"
	"net/url"

	"github.com/crewjam/saml"
	"github.com/crewjam/saml/samlsp"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/idp"
)

var _ idp.Provider = (*Provider)(nil)

type GetRequest func(ctx context.Context, intentID string) (*samlsp.TrackedRequest, error)
type AddRequest func(ctx context.Context, intentID, requestID string) error

// Provider is the [idp.Provider] implementation for a generic SAML provider
type Provider struct {
	name string

	getRequest GetRequest
	addRequest AddRequest

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
	metadata []byte,
	certificate []byte,
	key []byte,
	getRequest GetRequest,
	addRequest AddRequest,
	options ...ProviderOpts,
) (*Provider, error) {
	entityDescriptor := new(saml.EntityDescriptor)
	if err := xml.Unmarshal(metadata, entityDescriptor); err != nil {
		return nil, err
	}
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
		IDPMetadata: entityDescriptor,
		SignRequest: false,
	}
	provider := &Provider{
		name:       name,
		addRequest: addRequest,
		getRequest: getRequest,
		spOptions:  &opts,
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

func (p *Provider) GetSP() (*samlsp.Middleware, error) {
	sp, err := samlsp.New(*p.spOptions)
	if err != nil {
		return nil, errors.ThrowInternal(err, "SAML-x1v0hlrcjd", "Errors.Intent.IDPInvalid")
	}
	sp.RequestTracker = &RequestTracker{
		getRequest: p.getRequest,
		addRequest: p.addRequest,
	}
	return sp, nil
}

func (p *Provider) BeginAuth(ctx context.Context, state string, params ...any) (idp.Session, error) {
	m, err := p.GetSP()
	if err != nil {
		return nil, err
	}

	return &Session{
		serviceProvider: m,
		state:           state,
	}, nil
}
