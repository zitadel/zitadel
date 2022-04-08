package saml

import (
	"crypto"
	"crypto/dsa"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"fmt"
	"github.com/beevik/etree"
	"github.com/caos/zitadel/internal/api/saml/signature"
	"github.com/caos/zitadel/internal/api/saml/xml"
	"github.com/caos/zitadel/internal/api/saml/xml/md"
	"github.com/caos/zitadel/internal/api/saml/xml/samlp"
	"math/big"
	"net/url"
	"regexp"
	"strings"
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

	certStr := ""
	cert := &x509.Certificate{}
	if metadata.SPSSODescriptor.KeyDescriptor != nil && len(metadata.SPSSODescriptor.KeyDescriptor) > 0 {
		for _, keydesc := range metadata.SPSSODescriptor.KeyDescriptor {
			if keydesc.Use == md.KeyTypesSigning {
				certStr = keydesc.KeyInfo.X509Data[0].X509Certificate
			}
		}

		if certStr != "" {
			certStr = strings.ReplaceAll(certStr, "\n", "")
			certStr = strings.ReplaceAll(certStr, " ", "")
			block, err := base64.StdEncoding.DecodeString(certStr)
			if err != nil {
				return nil, fmt.Errorf("failed to parse PEM block containing the public key")
			}
			certT, err := x509.ParseCertificate(block)
			if err != nil {
				return nil, fmt.Errorf("failed to parse certificate: " + err.Error())
			}
			cert = certT
		}
	}

	return &ServiceProvider{
		ID:              id,
		metadata:        metadata,
		url:             config.URL,
		signerPublicKey: cert.PublicKey,
		defaultLoginURL: defaultLoginURL,
	}, nil
}

func (sp *ServiceProvider) verifyRequest(request *samlp.AuthnRequestType) error {
	if string(sp.metadata.EntityID) != request.Issuer.Text {
		return fmt.Errorf("request contains unknown issuer")
	}

	return nil
}

func (sp *ServiceProvider) verifyRedirectSignature(request, relayState, sigAlg, expectedSig string) error {
	// Validate the signature
	sig := []byte{}
	if url.QueryEscape(relayState) != "" {
		sig = []byte(fmt.Sprintf("SAMLRequest=%s&RelayState=%s&SigAlg=%s", url.QueryEscape(request), url.QueryEscape(relayState), url.QueryEscape(sigAlg)))
	} else {
		sig = []byte(fmt.Sprintf("SAMLRequest=%s&SigAlg=%s", url.QueryEscape(request), url.QueryEscape(sigAlg)))
	}
	signature, err := base64.StdEncoding.DecodeString(expectedSig)
	if err != nil {
		return err
	}

	return sp.verifySignature(sigAlg, sig, signature)
}
func (sp *ServiceProvider) getSigningCerts() ([]*x509.Certificate, error) {
	var certStrs []string

	for _, keyDescriptor := range sp.metadata.SPSSODescriptor.KeyDescriptor {
		for _, x509Data := range keyDescriptor.KeyInfo.X509Data {
			if len(x509Data.X509Certificate) != 0 {
				switch keyDescriptor.Use {
				case "", "signing":
					certStrs = append(certStrs, x509Data.X509Certificate)
				}
			}
		}
	}

	if len(certStrs) == 0 {
		return nil, fmt.Errorf("cannot find any signing certificate in the IDP SSO descriptor")
	}

	var certs []*x509.Certificate

	regex := regexp.MustCompile(`\s+`)
	for _, certStr := range certStrs {
		certStr = regex.ReplaceAllString(certStr, "")
		certBytes, err := base64.StdEncoding.DecodeString(certStr)
		if err != nil {
			return nil, fmt.Errorf("cannot parse certificate: %s", err)
		}

		parsedCert, err := x509.ParseCertificate(certBytes)
		if err != nil {
			return nil, err
		}
		certs = append(certs, parsedCert)
	}

	return certs, nil
}

func (sp *ServiceProvider) verifyPostSignature(authRequest string) error {
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes([]byte(authRequest)); err != nil {
		return err
	}

	certs, err := sp.getSigningCerts()
	if err != nil {
		return err
	}

	return signature.Validate(certs, doc.Root())
}

func (sp *ServiceProvider) verifySignature(sigAlg string, elementToSign []byte, signature []byte) error {
	switch sigAlg {
	case "http://www.w3.org/2009/xmldsig11#dsa-sha256":
		sum := sha256Sum(elementToSign)
		return verifyDSA(sp, signature, sum)
	case "http://www.w3.org/2000/09/xmldsig#dsa-sha1":
		sum := sha1Sum(elementToSign)
		return verifyDSA(sp, signature, sum)
	case "http://www.w3.org/2000/09/xmldsig#rsa-sha1":
		sum := sha1Sum(elementToSign)
		return rsa.VerifyPKCS1v15(sp.signerPublicKey.(*rsa.PublicKey), crypto.SHA1, sum, signature)
	case "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256":
		sum := sha256Sum(elementToSign)
		return rsa.VerifyPKCS1v15(sp.signerPublicKey.(*rsa.PublicKey), crypto.SHA256, sum, signature)
	default:
		return fmt.Errorf("unsupported signature algorithm, %s", sigAlg)
	}
}

type dsaSignature struct {
	R, S *big.Int
}

func verifyDSA(sp *ServiceProvider, signature, sum []byte) error {
	dsaSig := new(dsaSignature)
	if rest, err := asn1.Unmarshal(signature, dsaSig); err != nil {
		return err
	} else if len(rest) != 0 {
		return fmt.Errorf("trailing data after DSA signature")
	}
	if dsaSig.R.Sign() <= 0 || dsaSig.S.Sign() <= 0 {
		return fmt.Errorf("DSA signature contained zero or negative values")
	}
	if !dsa.Verify(sp.signerPublicKey.(*dsa.PublicKey), sum, dsaSig.R, dsaSig.S) {
		return fmt.Errorf("DSA verification failure")
	}
	return nil
}

func sha1Sum(sig []byte) []byte {
	h := sha1.New() // nolint: gosec
	_, err := h.Write(sig)
	if err != nil {
		return nil
	}
	return h.Sum(nil)
}

func sha256Sum(sig []byte) []byte {
	h := sha256.New()
	_, err := h.Write(sig)
	if err != nil {
		return nil
	}
	return h.Sum(nil)
}