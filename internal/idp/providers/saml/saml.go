package saml

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
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
	metadata *saml.EntityDescriptor,
	certificate []byte,
	key []byte,
	options ...ProviderOpts,
) *Provider {
	keyPair, err := tls.X509KeyPair(certificate, key)
	if err != nil {
		panic(err)
	}
	keyPair.Leaf, err = x509.ParseCertificate(keyPair.Certificate[0])
	if err != nil {
		panic(err)
	}
	rootURL, err := url.Parse("http://localhost:8000")
	if err != nil {
		panic(err)
	}
	opts := samlsp.Options{
		URL:         *rootURL,
		Key:         keyPair.PrivateKey.(*rsa.PrivateKey),
		Certificate: keyPair.Leaf,
		IDPMetadata: metadata,
		SignRequest: false,
	}
	provider := &Provider{
		name:      name,
		spOptions: &opts,
	}
	for _, option := range options {
		option(provider)
	}
	return provider
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
	sp, err := samlsp.New(*p.spOptions)
	if err != nil {
		return nil, err
	}

	authReq, err := sp.ServiceProvider.MakeAuthenticationRequest(sp.ServiceProvider.GetSSOBindingLocation(p.binding), p.binding, sp.ResponseBinding)
	if err != nil {
		return nil, err
	}

	var authURL string
	if p.binding == saml.HTTPRedirectBinding {
		redirectURL, err := authReq.Redirect(state, &sp.ServiceProvider)
		if err != nil {
			return nil, err
		}
		authURL = redirectURL.String()
	} else if p.binding == saml.HTTPPostBinding {
		authURL = `<!DOCTYPE html><html><body>` +
			string(authReq.Post(state)) +
			`</body></html>`
	}

	return &Session{
		Provider: p,
		AuthURL:  authURL,
	}, nil
}
