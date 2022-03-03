package saml

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/amdonov/xmlsig"
	"github.com/caos/logging"
	"github.com/caos/oidc/pkg/op"
	"github.com/caos/zitadel/internal/api/saml/xml/metadata/md"
	"github.com/caos/zitadel/internal/api/saml/xml/metadata/saml"
	"github.com/caos/zitadel/internal/api/saml/xml/protocol/samlp"
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/eventstore/key"
	"gopkg.in/square/go-jose.v2"
	"net/http"
	"text/template"
)

type IDPStorage interface {
	AuthStorage
	EntityStorage
	UserStorage
	Health(context.Context) error
}

type AuthStorage interface {
	CreateAuthRequest(context.Context, *samlp.AuthnRequest, string, string, string) (AuthRequestInt, error)
	AuthRequestByID(context.Context, string) (AuthRequestInt, error)
	AuthRequestByCode(context.Context, string) (AuthRequestInt, error)
	GetAttributesFromNameID(ctx context.Context, nameID string) (map[string]interface{}, error)
}

type MetadataIDP struct {
	ValidUntil    string
	CacheDuration string
	ErrorURL      string
}

type IdentityProviderConfig struct {
	Metadata *MetadataIDP

	SignatureAlgorithm  string
	DigestAlgorithm     string
	EncryptionAlgorithm string

	NameIDFormat           string
	WantAuthRequestsSigned string

	Endpoints *EndpointConfig `yaml:"Endpoints"`
}

type EndpointConfig struct {
	Callback      Endpoint `yaml:"Callback"`
	SingleSignOn  Endpoint `yaml:"SingleSignOn"`
	SingleLogOut  Endpoint `yaml:"SingleLogOut"`
	Artifact      Endpoint `yaml:"Artifact"`
	SLOArtifact   Endpoint `yaml:"SLOArtifact"`
	NameIDMapping Endpoint `yaml:"NameIDMapping"`
	Attribute     Endpoint `yaml:"Attribute"`
}

type Endpoint struct {
	Path string `yaml:"Path"`
	URL  string `yaml:"URL"`
}

type IdentityProvider struct {
	storage        IDPStorage
	postTemplate   *template.Template
	logoutTemplate *template.Template

	EntityID   string
	Metadata   *md.IDPSSODescriptorType
	AAMetadata *md.AttributeAuthorityDescriptorType
	signer     xmlsig.Signer

	CallbackEndpoint              op.Endpoint
	SingleSignOnEndpoint          op.Endpoint
	SingleLogoutEndpoint          op.Endpoint
	ArtifactResulationEndpoint    op.Endpoint
	SLOArtifactResulationEndpoint op.Endpoint
	NameIDMappingEndpoint         op.Endpoint
	AttributeEndpoint             op.Endpoint

	serviceProviders []*ServiceProvider
}

func NewIdentityProvider(metadataEndpoint *op.Endpoint, conf *IdentityProviderConfig, storage IDPStorage) (*IdentityProvider, error) {
	cert, key := getResponseCert(storage)

	certPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert,
		},
	)

	keyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)
	tlsCert, err := tls.X509KeyPair(certPem, keyPem)
	if err != nil {
		return nil, err
	}

	signer, err := xmlsig.NewSignerWithOptions(tlsCert, xmlsig.SignerOptions{
		SignatureAlgorithm: conf.SignatureAlgorithm,
		DigestAlgorithm:    conf.DigestAlgorithm,
	})
	if err != nil {
		return nil, err
	}

	postTemplate, err := template.New("post").Parse(postTemplate)
	if err != nil {
		return nil, err
	}

	logoutTemplate, err := template.New("post").Parse(logoutTemplate)
	if err != nil {
		return nil, err
	}

	metadata, aaMetadata := conf.getMetadata(metadataEndpoint, tlsCert.Certificate[0])
	return &IdentityProvider{
		storage:                       storage,
		EntityID:                      metadataEndpoint.Absolute(""),
		Metadata:                      metadata,
		AAMetadata:                    aaMetadata,
		signer:                        signer,
		CallbackEndpoint:              op.NewEndpointWithURL(conf.Endpoints.Callback.Path, conf.Endpoints.Callback.URL),
		SingleSignOnEndpoint:          op.NewEndpointWithURL(conf.Endpoints.SingleSignOn.Path, conf.Endpoints.SingleSignOn.URL),
		SingleLogoutEndpoint:          op.NewEndpointWithURL(conf.Endpoints.SingleLogOut.Path, conf.Endpoints.SingleLogOut.URL),
		ArtifactResulationEndpoint:    op.NewEndpointWithURL(conf.Endpoints.Artifact.Path, conf.Endpoints.Artifact.URL),
		SLOArtifactResulationEndpoint: op.NewEndpointWithURL(conf.Endpoints.SLOArtifact.Path, conf.Endpoints.SLOArtifact.URL),
		NameIDMappingEndpoint:         op.NewEndpointWithURL(conf.Endpoints.NameIDMapping.Path, conf.Endpoints.NameIDMapping.URL),
		AttributeEndpoint:             op.NewEndpointWithURL(conf.Endpoints.Attribute.Path, conf.Endpoints.Attribute.URL),
		postTemplate:                  postTemplate,
		logoutTemplate:                logoutTemplate,
	}, nil
}

type Route struct {
	Endpoint   string
	HandleFunc http.HandlerFunc
}

func (p *IdentityProvider) GetRoutes() []*Route {
	return []*Route{
		{p.CallbackEndpoint.Relative(), p.callbackHandleFunc},
		{p.SingleSignOnEndpoint.Relative(), p.ssoHandleFunc},
		{p.SingleLogoutEndpoint.Relative(), p.logoutHandleFunc},
		{p.ArtifactResulationEndpoint.Relative(), notImplementedHandleFunc},
		{p.SLOArtifactResulationEndpoint.Relative(), notImplementedHandleFunc},
		{p.NameIDMappingEndpoint.Relative(), notImplementedHandleFunc},
		{p.AttributeEndpoint.Relative(), notImplementedHandleFunc},
	}
}

func (p *IdentityProvider) GetServiceProvider(ctx context.Context, entityID string) (*ServiceProvider, error) {
	index := 0
	found := false
	for i, sp := range p.serviceProviders {
		if sp.GetEntityID() == entityID {
			found = true
			index = i
			break
		}
	}
	if found == true {
		return p.serviceProviders[index], nil
	}

	sp, err := p.storage.GetEntityByID(ctx, entityID)
	if err != nil {
		return nil, err
	}
	if sp != nil {
		p.serviceProviders = append(p.serviceProviders, sp)
	}
	return sp, nil
}

func (p *IdentityProvider) DeleteServiceProvider(entityID string) error {
	index := 0
	found := false
	for i, sp := range p.serviceProviders {
		if sp.GetEntityID() == entityID {
			found = true
			index = i
			break
		}
	}
	if found == true {
		p.serviceProviders = append(p.serviceProviders[:index], p.serviceProviders[index+1:]...)
	}
	return nil
}

func (p *IdentityProvider) getIssuer() *saml.Issuer {
	return &saml.Issuer{
		Format: "urn:oasis:names:tc:SAML:2.0:nameid-format:entity",
		Text:   string(p.EntityID),
	}
}

func (p *IdentityProvider) verifyRequestDestination(request *samlp.AuthnRequest) error {
	foundEndpoint := false
	for _, sso := range p.Metadata.SingleSignOnService {
		if request.Destination == sso.Location {
			foundEndpoint = true
			break
		}
	}
	if !foundEndpoint {
		return fmt.Errorf("destination of request is unknown")
	}

	return nil
}

func notImplementedHandleFunc(w http.ResponseWriter, r *http.Request) {
	http.Error(w, fmt.Sprintf("not implemented yet"), http.StatusNotImplemented)
}

func getResponseCert(storage Storage) ([]byte, *rsa.PrivateKey) {
	ctx := context.Background()
	certAndKeyCh := make(chan key.CertificateAndKey)
	go storage.GetResponseSigningKey(ctx, certAndKeyCh)

	for {
		select {
		case <-ctx.Done():
			//TODO
		case certAndKey := <-certAndKeyCh:
			if certAndKey.Key.Key == nil || certAndKey.Certificate.Key == nil {
				logging.Log("OP-DAvt4").Warn("signer has no key")
				continue
			}
			certWebKey := certAndKey.Certificate.Key.(jose.JSONWebKey)
			keyWebKey := certAndKey.Key.Key.(jose.JSONWebKey)

			return certWebKey.Key.([]byte), keyWebKey.Key.(*rsa.PrivateKey)
		}
	}
}
