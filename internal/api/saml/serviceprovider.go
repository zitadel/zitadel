package saml

import (
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"github.com/beevik/etree"
	"github.com/caos/zitadel/internal/api/saml/signature"
	"github.com/caos/zitadel/internal/api/saml/xml"
	"github.com/caos/zitadel/internal/api/saml/xml/md"
	"github.com/caos/zitadel/internal/api/saml/xml/samlp"
	"net/url"
)

type ServiceProviderConfig struct {
	Metadata string
	URL      string
}

type ServiceProvider struct {
	ID              string
	metadata        *md.EntityDescriptorType
	url             string
	signerPublicKey interface{}
	defaultLoginURL string
}

func (sp *ServiceProvider) GetEntityID() string {
	return string(sp.metadata.EntityID)
}

func (sp *ServiceProvider) LoginURL(id string) string {
	return sp.defaultLoginURL + id
}

func NewServiceProvider(id string, config *ServiceProviderConfig, defaultLoginURL string) (*ServiceProvider, error) {
	metadataData := make([]byte, 0)
	if config.URL != "" {
		body, err := xml.ReadMetadataFromURL(config.URL)
		if err != nil {
			return nil, err
		}
		metadataData = body
	} else {
		metadataData = []byte(config.Metadata)
	}
	metadata, err := xml.ParseMetadataXmlIntoStruct(metadataData)
	if err != nil {
		return nil, err
	}

	var signerPublicKey interface{}
	certs, err := getSigningCertsFromMetadata(metadata)
	if err != nil {
		return nil, err
	}
	if len(certs) > 1 {
		return nil, fmt.Errorf("currently more than one signing certificate for service providers not supported")
	}
	if len(certs) == 1 {
		signerPublicKey = certs[0].PublicKey
	}

	return &ServiceProvider{
		ID:              id,
		metadata:        metadata,
		url:             config.URL,
		signerPublicKey: signerPublicKey,
		defaultLoginURL: defaultLoginURL,
	}, nil
}

func (sp *ServiceProvider) verifyRequest(request *samlp.AuthnRequestType) error {
	if string(sp.metadata.EntityID) != request.Issuer.Text {
		return fmt.Errorf("request contains unknown issuer")
	}

	return nil
}

func getSigningCertsFromMetadata(metadata *md.EntityDescriptorType) ([]*x509.Certificate, error) {
	return signature.ParseCertificates(xml.GetCertsFromKeyDescriptors(metadata.SPSSODescriptor.KeyDescriptor))
}

func (sp *ServiceProvider) validatePostSignature(authRequest string) error {
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes([]byte(authRequest)); err != nil {
		return err
	}

	if doc.Root() == nil {
		return fmt.Errorf("error while parsing request")
	}

	certs, err := getSigningCertsFromMetadata(sp.metadata)
	if err != nil {
		return err
	}

	return signature.ValidatePost(certs, doc.Root())
}

func (sp *ServiceProvider) validateRedirectSignature(request, relayState, sigAlg, expectedSig string) error {
	if sp.signerPublicKey == nil {
		return fmt.Errorf("error can not validate signature if no certificate is present for this service provider")
	}

	elementToSign := []byte{}
	if url.QueryEscape(relayState) != "" {
		elementToSign = []byte(fmt.Sprintf("SAMLRequest=%s&RelayState=%s&SigAlg=%s", url.QueryEscape(request), url.QueryEscape(relayState), url.QueryEscape(sigAlg)))
	} else {
		elementToSign = []byte(fmt.Sprintf("SAMLRequest=%s&SigAlg=%s", url.QueryEscape(request), url.QueryEscape(sigAlg)))
	}
	signatureValue, err := base64.StdEncoding.DecodeString(expectedSig)
	if err != nil {
		return err
	}

	return signature.ValidateRedirect(sigAlg, elementToSign, signatureValue, sp.signerPublicKey)
}
