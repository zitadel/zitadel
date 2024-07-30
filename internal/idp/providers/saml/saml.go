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

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var _ idp.Provider = (*Provider)(nil)

// Provider is the [idp.Provider] implementation for a generic SAML provider
type Provider struct {
	name string

	requestTracker samlsp.RequestTracker
	Certificate    []byte

	spOptions *samlsp.Options

	binding                       string
	nameIDFormat                  saml.NameIDFormat
	transientMappingAttributeName string

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

func WithNameIDFormat(format domain.SAMLNameIDFormat) ProviderOpts {
	return func(p *Provider) {
		p.nameIDFormat = nameIDFormatFromDomain(format)
	}
}

func WithTransientMappingAttributeName(attribute string) ProviderOpts {
	return func(p *Provider) {
		p.transientMappingAttributeName = attribute
	}
}

func WithCustomRequestTracker(tracker samlsp.RequestTracker) ProviderOpts {
	return func(p *Provider) {
		p.requestTracker = tracker
	}
}

func WithEntityID(entityID string) ProviderOpts {
	return func(p *Provider) {
		p.spOptions.EntityID = entityID
	}
}

func New(
	name string,
	rootURLStr string,
	metadata []byte,
	certificate []byte,
	key []byte,
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
		name:        name,
		spOptions:   &opts,
		Certificate: certificate,
		// the library uses transient as default, which does not make sense for federating accounts
		nameIDFormat: saml.PersistentNameIDFormat,
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
		return nil, zerrors.ThrowInternal(err, "SAML-qee09ffuq5", "Errors.Intent.IDPInvalid")
	}
	sp.ServiceProvider.AuthnNameIDFormat = p.nameIDFormat
	if p.requestTracker != nil {
		sp.RequestTracker = p.requestTracker
	}
	if p.binding != "" {
		sp.Binding = p.binding
	}
	return sp, nil
}

func (p *Provider) BeginAuth(ctx context.Context, state string, _ ...idp.Parameter) (idp.Session, error) {
	m, err := p.GetSP()
	if err != nil {
		return nil, err
	}

	return &Session{
		ServiceProvider: m,
		state:           state,
	}, nil
}

func (p *Provider) TransientMappingAttributeName() string {
	return p.transientMappingAttributeName
}

func nameIDFormatFromDomain(format domain.SAMLNameIDFormat) saml.NameIDFormat {
	switch format {
	case domain.SAMLNameIDFormatUnspecified:
		return saml.UnspecifiedNameIDFormat
	case domain.SAMLNameIDFormatEmailAddress:
		return saml.EmailAddressNameIDFormat
	case domain.SAMLNameIDFormatPersistent:
		return saml.PersistentNameIDFormat
	case domain.SAMLNameIDFormatTransient:
		return saml.TransientNameIDFormat
	default:
		return saml.UnspecifiedNameIDFormat
	}
}
