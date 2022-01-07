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
	"encoding/xml"
	"fmt"
	"github.com/caos/zitadel/internal/api/saml/xml/metadata/md"
	"github.com/caos/zitadel/internal/api/saml/xml/protocol/saml"
	"github.com/caos/zitadel/internal/api/saml/xml/protocol/samlp"
	"io"
	"math/big"
	"net/http"
)

type ServiceProviderConfig struct {
	Metadata string
	URL      string
}

type ServiceProvider struct {
	ID              string
	metadata        *md.EntityDescriptor
	url             string
	signerPublicKey interface{}
}

func (sp *ServiceProvider) GetEntityID() string {
	return string(sp.metadata.EntityID)
}

func NewServiceProvider(config *ServiceProviderConfig) (*ServiceProvider, error) {
	var metadata *md.EntityDescriptor
	metadataData := make([]byte, 0)
	if config.URL != "" {
		resp, err := http.Get(config.URL)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		metadataData = body
	} else {
		metadataData = []byte(config.Metadata)
	}

	err := xml.Unmarshal(metadataData, metadata)
	if err != nil {
		return nil, err
	}

	certStr := ""
	for _, keydesc := range metadata.SPSSODescriptor.KeyDescriptor {
		if keydesc.Use == md.KeyTypesSigning {
			certStr = keydesc.KeyInfo.X509Data[0].X509Certificate
		}
	}

	block, err := base64.StdEncoding.DecodeString(certStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the public key")
	}
	cert, err := x509.ParseCertificate(block)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: " + err.Error())
	}

	return &ServiceProvider{
		metadata:        metadata,
		url:             config.URL,
		signerPublicKey: cert.PublicKey,
	}, nil
}

func (sp *ServiceProvider) verifyRequest(request *samlp.AuthnRequest) error {
	if string(sp.metadata.EntityID) != request.Issuer.Text {
		return fmt.Errorf("request contains unknown issuer")
	}

	return nil
}

func (sp *ServiceProvider) getIssuer() *saml.Issuer {
	return &saml.Issuer{
		Format: "urn:oasis:names:tc:SAML:2.0:nameid-format:entity",
		Text:   string(sp.metadata.EntityID),
	}
}

func (sp *ServiceProvider) verifySignature(request, relayState, sigAlg, expectedSig string) error {
	// Validate the signature
	sig := []byte(fmt.Sprintf("SAMLRequest=%s&RelayState=%s&SigAlg=%s", request, relayState, sigAlg))

	signature, err := base64.StdEncoding.DecodeString(expectedSig)

	if err != nil {
		return err
	}
	switch sigAlg {
	case "http://www.w3.org/2009/xmldsig11#dsa-sha256":
		sum := sha256Sum(sig)
		return verifyDSA(sp, signature, sum)
	case "http://www.w3.org/2000/09/xmldsig#dsa-sha1":
		sum := sha1Sum(sig)
		return verifyDSA(sp, signature, sum)
	case "http://www.w3.org/2000/09/xmldsig#rsa-sha1":
		sum := sha1Sum(sig)
		return rsa.VerifyPKCS1v15(sp.signerPublicKey.(*rsa.PublicKey), crypto.SHA1, sum, signature)
	case "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256":
		sum := sha256Sum(sig)
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
