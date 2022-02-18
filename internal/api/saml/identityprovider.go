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
	Health(context.Context) error
}

type AuthStorage interface {
	CreateAuthRequest(context.Context, *samlp.AuthnRequest, string, string) (AuthRequestInt, error)
	AuthRequestByID(context.Context, string) (AuthRequestInt, error)
	AuthRequestByCode(context.Context, string) (AuthRequestInt, error)
	GetAttributesFromNameID(ctx context.Context, nameID string) (map[string]interface{}, error)
}

type IdentityProviderConfig struct {
	ValidUntil    string
	CacheDuration string
	ErrorURL      string

	SignatureAlgorithm  string
	DigestAlgorithm     string
	EncryptionAlgorithm string

	NameIDFormat           string
	WantAuthRequestsSigned string

	LoginService                 string
	SingleSignOnService          string
	SingleLogoutService          string
	ArtifactResulationService    string
	SLOArtifactResulationService string
	NameIDMappingService         string
	AttributeService             string
}

type IdentityProvider struct {
	storage      IDPStorage
	postTemplate *template.Template

	EntityID   string
	Metadata   *md.IDPSSODescriptorType
	AAMetadata *md.AttributeAuthorityDescriptorType
	signer     xmlsig.Signer

	baseURL                      string
	LoginService                 string
	SingleSignOnService          string
	SingleLogoutService          string
	ArtifactResulationService    string
	SLOArtifactResulationService string
	NameIDMappingService         string
	AttributeService             string

	ServiceProviders []*ServiceProvider
}

func NewIdentityProvider(baseURL string, conf *IdentityProviderConfig, storage IDPStorage) (*IdentityProvider, error) {
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

	temp, err := template.New("post").Parse(postTemplate)
	if err != nil {
		return nil, err
	}

	metadata, aaMetadata := conf.getMetadata(baseURL, tlsCert.Certificate[0])
	return &IdentityProvider{
		storage:                      storage,
		EntityID:                     baseURL + metadataEndpoint,
		baseURL:                      baseURL,
		Metadata:                     metadata,
		AAMetadata:                   aaMetadata,
		signer:                       signer,
		LoginService:                 conf.LoginService,
		SingleSignOnService:          conf.SingleSignOnService,
		SingleLogoutService:          conf.SingleLogoutService,
		ArtifactResulationService:    conf.ArtifactResulationService,
		SLOArtifactResulationService: conf.SLOArtifactResulationService,
		NameIDMappingService:         conf.NameIDMappingService,
		AttributeService:             conf.AttributeService,
		postTemplate:                 temp,
	}, nil
}

type Route struct {
	Endpoint   string
	HandleFunc http.HandlerFunc
}

func (p *IdentityProvider) GetRoutes() []*Route {
	return []*Route{
		{p.LoginService, p.loginHandleFunc},
		{p.SingleSignOnService, p.ssoHandleFunc},
		{p.SingleLogoutService, p.logoutHandleFunc},
		{p.ArtifactResulationService, notImplementedHandleFunc},
		{p.SLOArtifactResulationService, notImplementedHandleFunc},
		{p.NameIDMappingService, notImplementedHandleFunc},
		{p.AttributeService, notImplementedHandleFunc},
	}
}

func (p *IdentityProvider) GetRedirectURL(requestID string) string {
	//TODO
	return p.baseURL + p.LoginService + "?requestId=" + requestID
}

func (p *IdentityProvider) GetServiceProvider(ctx context.Context, entityID string) (*ServiceProvider, error) {
	index := 0
	found := false
	for i, sp := range p.ServiceProviders {
		if sp.GetEntityID() == entityID {
			found = true
			index = i
			break
		}
	}
	if found == true {
		return p.ServiceProviders[index], nil
	}

	sp, err := p.storage.GetEntityByID(ctx, entityID)
	if err != nil {
		return nil, err
	}
	if sp != nil {
		p.ServiceProviders = append(p.ServiceProviders, sp)
	}
	return sp, nil
}

func (p *IdentityProvider) AddServiceProvider(id string, config *ServiceProviderConfig) error {
	sp, err := NewServiceProvider(id, config)
	if err != nil {
		return err
	}

	p.ServiceProviders = append(p.ServiceProviders, sp)
	return nil
}

func (p *IdentityProvider) DeleteServiceProvider(entityID string) error {
	index := 0
	found := false
	for i, sp := range p.ServiceProviders {
		if sp.GetEntityID() == entityID {
			found = true
			index = i
			break
		}
	}
	if found == true {
		p.ServiceProviders = append(p.ServiceProviders[:index], p.ServiceProviders[index+1:]...)
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
