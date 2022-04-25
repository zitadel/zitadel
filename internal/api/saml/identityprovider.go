package saml

import (
	"bytes"
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/amdonov/xmlsig"
	"github.com/caos/logging"
	"github.com/caos/oidc/pkg/op"
	"github.com/caos/zitadel/internal/api/saml/key"
	"github.com/caos/zitadel/internal/api/saml/serviceprovider"
	"github.com/caos/zitadel/internal/api/saml/signature"
	"github.com/caos/zitadel/internal/api/saml/xml/md"
	"github.com/caos/zitadel/internal/api/saml/xml/samlp"
	dsig "github.com/russellhaering/goxmldsig"
	"gopkg.in/square/go-jose.v2"
	"io"
	"net/http"
	"text/template"
)

type IDPStorage interface {
	AuthStorage
	IdentityProviderStorage
	UserStorage
	Health(context.Context) error
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

	WantAuthRequestsSigned string

	Endpoints *EndpointConfig `yaml:"Endpoints"`
}

type EndpointConfig struct {
	Certificate   Endpoint `yaml:"Certificate"`
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

	EntityID       string
	Metadata       *md.IDPSSODescriptorType
	AAMetadata     *md.AttributeAuthorityDescriptorType
	signer         xmlsig.Signer
	signingContext *dsig.SigningContext

	CertificateEndpoint           op.Endpoint
	CallbackEndpoint              op.Endpoint
	SingleSignOnEndpoint          op.Endpoint
	SingleLogoutEndpoint          op.Endpoint
	ArtifactResulationEndpoint    op.Endpoint
	SLOArtifactResulationEndpoint op.Endpoint
	NameIDMappingEndpoint         op.Endpoint
	AttributeEndpoint             op.Endpoint

	serviceProviders []*serviceprovider.ServiceProvider
}

func NewIdentityProvider(metadataEndpoint *op.Endpoint, conf *IdentityProviderConfig, storage IDPStorage) (*IdentityProvider, error) {
	cert, privateKey := getResponseCert(storage)

	signingContext, signer, err := signature.GetSigningContextAndSigner(cert, privateKey, conf.SignatureAlgorithm)
	if err != nil {
		return nil, err
	}

	postTemplate, err := template.New("post").Parse(postTemplate)
	if err != nil {
		return nil, err
	}

	logoutTemplate, err := template.New("logout").Parse(logoutTemplate)
	if err != nil {
		return nil, err
	}

	tlsCert, err := signature.ParseTlsKeyPair(cert, privateKey)
	if err != nil {
		return nil, err
	}

	metadata, aaMetadata := conf.getMetadata(metadataEndpoint, tlsCert.Certificate[0])
	idp := &IdentityProvider{
		storage:        storage,
		EntityID:       metadataEndpoint.Absolute(""),
		Metadata:       metadata,
		AAMetadata:     aaMetadata,
		signingContext: signingContext,
		signer:         signer,
		postTemplate:   postTemplate,
		logoutTemplate: logoutTemplate,
	}
	if conf.Endpoints != nil {
		idp.CertificateEndpoint = op.NewEndpointWithURL(conf.Endpoints.Certificate.Path, conf.Endpoints.Certificate.URL)
		idp.CallbackEndpoint = op.NewEndpointWithURL(conf.Endpoints.Callback.Path, conf.Endpoints.Callback.URL)
		idp.SingleSignOnEndpoint = op.NewEndpointWithURL(conf.Endpoints.SingleSignOn.Path, conf.Endpoints.SingleSignOn.URL)
		idp.SingleLogoutEndpoint = op.NewEndpointWithURL(conf.Endpoints.SingleLogOut.Path, conf.Endpoints.SingleLogOut.URL)
		idp.ArtifactResulationEndpoint = op.NewEndpointWithURL(conf.Endpoints.Artifact.Path, conf.Endpoints.Artifact.URL)
		idp.SLOArtifactResulationEndpoint = op.NewEndpointWithURL(conf.Endpoints.SLOArtifact.Path, conf.Endpoints.SLOArtifact.URL)
		idp.NameIDMappingEndpoint = op.NewEndpointWithURL(conf.Endpoints.NameIDMapping.Path, conf.Endpoints.NameIDMapping.URL)
		idp.AttributeEndpoint = op.NewEndpointWithURL(conf.Endpoints.Attribute.Path, conf.Endpoints.Attribute.URL)
	}
	return idp, nil
}

type Route struct {
	Endpoint   string
	HandleFunc http.HandlerFunc
}

func (p *IdentityProvider) GetRoutes() []*Route {
	return []*Route{
		{p.CertificateEndpoint.Relative(), p.certificateHandleFunc},
		{p.CallbackEndpoint.Relative(), p.callbackHandleFunc},
		{p.SingleSignOnEndpoint.Relative(), p.ssoHandleFunc},
		{p.SingleLogoutEndpoint.Relative(), p.logoutHandleFunc},
		{p.ArtifactResulationEndpoint.Relative(), notImplementedHandleFunc},
		{p.SLOArtifactResulationEndpoint.Relative(), notImplementedHandleFunc},
		{p.NameIDMappingEndpoint.Relative(), notImplementedHandleFunc},
		{p.AttributeEndpoint.Relative(), p.attributeQueryHandleFunc},
	}
}

func (p *IdentityProvider) GetServiceProvider(ctx context.Context, entityID string) (*serviceprovider.ServiceProvider, error) {
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

func (p *IdentityProvider) verifyRequestDestinationOfAuthRequest(request *samlp.AuthnRequestType) error {
	// google provides no destination in their requests
	if request.Destination != "" {
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
	}
	return nil
}

func (p *IdentityProvider) verifyRequestDestinationOfAttrQuery(request *samlp.AttributeQueryType) error {
	// google provides no destination in their requests
	if request.Destination != "" {
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
	}
	return nil
}

func notImplementedHandleFunc(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, fmt.Sprintf("not implemented yet"), http.StatusNotImplemented)
}

func getResponseCert(storage IdentityProviderStorage) ([]byte, *rsa.PrivateKey) {
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

func (i *IdentityProvider) certificateHandleFunc(w http.ResponseWriter, r *http.Request) {
	cert := i.Metadata.KeyDescriptor[0].KeyInfo.X509Data[0].X509Certificate

	data, err := base64.StdEncoding.DecodeString(cert)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to read certificate: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	certPem := new(bytes.Buffer)
	if err := pem.Encode(certPem, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: []byte(data),
	}); err != nil {
		http.Error(w, fmt.Errorf("failed to pem encode certificate: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=idp.crt")
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	_, err = io.Copy(w, certPem)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to response with certificate: %w", err).Error(), http.StatusInternalServerError)
		return
	}
}
