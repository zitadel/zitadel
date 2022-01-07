package saml

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/amdonov/xmlsig"
	"github.com/caos/zitadel/internal/api/saml/xml/metadata/md"
	"github.com/caos/zitadel/internal/api/saml/xml/metadata/saml"
	"github.com/caos/zitadel/internal/api/saml/xml/protocol/samlp"
	"io/ioutil"
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

	Certificate         *Certificate
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
	tlsConfig  *tls.Config

	LoginService                 string
	SingleSignOnService          string
	SingleLogoutService          string
	ArtifactResulationService    string
	SLOArtifactResulationService string
	NameIDMappingService         string
	AttributeService             string

	ServiceProviders []*ServiceProvider
}

func NewIdentityProvider(entityID string, conf *IdentityProviderConfig, storage IDPStorage) (*IdentityProvider, error) {
	certData, err := ioutil.ReadFile(conf.Certificate.Path)
	if err != nil {
		return nil, err
	}

	tlsConfig, err := ConfigureTLS(conf.Certificate.Path, conf.Certificate.PrivateKeyPath, conf.Certificate.CaPath)
	if err != nil {
		return nil, err
	}

	signer, err := xmlsig.NewSignerWithOptions(tlsConfig.Certificates[0], xmlsig.SignerOptions{
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

	metadata, aaMetadata := conf.getMetadata(entityID, certData)
	return &IdentityProvider{
		storage:                      storage,
		EntityID:                     entityID,
		Metadata:                     metadata,
		AAMetadata:                   aaMetadata,
		tlsConfig:                    tlsConfig,
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
	return p.LoginService + "?requestId=" + requestID
}

func (p *IdentityProvider) GetServiceProvider(entityID string) *ServiceProvider {
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
		return p.ServiceProviders[index]
	}
	return nil
}

func (p *IdentityProvider) AddServiceProvider(config *ServiceProviderConfig) error {
	sp, err := NewServiceProvider(config)
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
		if request.Destination != sso.Location {
			foundEndpoint = true
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
