package signature

import (
	"crypto"
	"crypto/dsa"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/xml"
	"fmt"
	"github.com/amdonov/xmlsig"
	"github.com/beevik/etree"
	"github.com/caos/zitadel/internal/api/saml/xml/xml_dsig"
	dsig "github.com/russellhaering/goxmldsig"
	"github.com/russellhaering/goxmldsig/etreeutils"
	"math/big"
)

/*
commented as russellhaering/goxmldsig produces invalid singatures for responses currently

func Create(signingContext *dsig.SigningContext, element interface{}) (*xml_dsig.SignatureType, error) {
	data, _, err := canonicalize(element)
	if err != nil {
		return nil, err
	}

	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(data); err != nil {
		return nil, err
	}

	signedEl, err := signingContext.SignEnveloped(doc.Root())
	if err != nil {
		return nil, err
	}

	sigEl := signedEl.Child[len(signedEl.Child)-1]
	sigTyped := sigEl.(*etree.Element)

	sigDoc := etree.NewDocument()
	sigDoc.SetRoot(sigTyped)

	reqBuf, err := sigDoc.WriteToBytes()
	if err != nil {
		return nil, err
	}

	sig, err := xml.DecodeSignature("", string(reqBuf))
	if err != nil {
		return nil, err
	}

	// unfortunately as the unmarshilling is correct but the innerXML attributes still contain the element with namespace they have to be cleaned out
	sig.InnerXml = ""
	sig.SignedInfo.InnerXml = ""
	sig.SignedInfo.CanonicalizationMethod.InnerXml = ""
	sig.SignedInfo.SignatureMethod.InnerXml = ""
	for i := range sig.SignedInfo.Reference {
		ref := sig.SignedInfo.Reference[i]
		for j := range ref.Transforms.Transform {
			ref.Transforms.Transform[j].InnerXml = ""
		}
		ref.Transforms.InnerXml = ""
		ref.InnerXml = ""
		sig.SignedInfo.Reference[i] = ref
	}
	sig.SignatureValue.InnerXml = ""
	sig.KeyInfo.InnerXml = ""
	for i := range sig.KeyInfo.X509Data {
		d := sig.KeyInfo.X509Data[i]
		d.InnerXml = ""
		sig.KeyInfo.X509Data[i] = d
	}

	return sig, nil
}*/

func Create(signer xmlsig.Signer, data interface{}) (*xml_dsig.SignatureType, error) {
	sig, err := signer.CreateSignature(data)
	if err != nil {
		return nil, err
	}
	transforms := []xml_dsig.TransformType{}
	for _, t := range sig.SignedInfo.Reference.Transforms.Transform {
		transforms = append(transforms, xml_dsig.TransformType{
			XMLName:   xml.Name{},
			Algorithm: t.Algorithm,
		})
	}

	return &xml_dsig.SignatureType{
		XMLName: xml.Name{},
		SignedInfo: xml_dsig.SignedInfoType{
			XMLName: xml.Name{},
			CanonicalizationMethod: xml_dsig.CanonicalizationMethodType{
				XMLName:   xml.Name{},
				Algorithm: sig.SignedInfo.CanonicalizationMethod.Algorithm,
			},
			SignatureMethod: xml_dsig.SignatureMethodType{
				XMLName:   xml.Name{},
				Algorithm: sig.SignedInfo.SignatureMethod.Algorithm,
			},
			Reference: []xml_dsig.ReferenceType{{
				DigestValue: xml_dsig.DigestValueType(sig.SignedInfo.Reference.DigestValue),
				DigestMethod: xml_dsig.DigestMethodType{
					XMLName:   xml.Name{},
					Algorithm: sig.SignedInfo.Reference.DigestMethod.Algorithm,
				},
				Transforms: &xml_dsig.TransformsType{
					Transform: transforms,
				},
				URI: sig.SignedInfo.Reference.URI,
			}},
		},
		SignatureValue: xml_dsig.SignatureValueType{
			Text: sig.SignatureValue,
		},
		KeyInfo: &xml_dsig.KeyInfoType{
			XMLName: xml.Name{},
			X509Data: []xml_dsig.X509DataType{{
				X509Certificate: sig.KeyInfo.X509Data.X509Certificate,
			}},
		},
	}, nil
}

func ValidatePost(certs []*x509.Certificate, el *etree.Element) error {
	certificateStore := dsig.MemoryX509CertificateStore{
		Roots: certs,
	}

	validationContext := dsig.NewDefaultValidationContext(&certificateStore)
	validationContext.IdAttribute = "ID"

	if el.FindElement("./Signature/KeyInfo/X509Data/X509Certificate") == nil {
		if sigEl := el.FindElement("./Signature"); sigEl != nil {
			if keyInfo := sigEl.FindElement("KeyInfo"); keyInfo != nil {
				sigEl.RemoveChild(keyInfo)
			}
		}
	}

	ctx, err := etreeutils.NSBuildParentContext(el)
	if err != nil {
		return err
	}
	ctx, err = ctx.SubContext(el)
	if err != nil {
		return err
	}
	el, err = etreeutils.NSDetatch(ctx, el)
	if err != nil {
		return err
	}

	_, err = validationContext.Validate(el)
	return err
}

func ValidateRedirect(sigAlg string, elementToSign []byte, signature []byte, pubKey interface{}) error {
	switch sigAlg {
	case "http://www.w3.org/2009/xmldsig11#dsa-sha256":
		sum := sha256Sum(elementToSign)
		return verifyDSA(signature, sum, pubKey)
	case "http://www.w3.org/2000/09/xmldsig#dsa-sha1":
		sum := sha1Sum(elementToSign)
		return verifyDSA(signature, sum, pubKey)
	case "http://www.w3.org/2000/09/xmldsig#rsa-sha1":
		sum := sha1Sum(elementToSign)
		return rsa.VerifyPKCS1v15(pubKey.(*rsa.PublicKey), crypto.SHA1, sum, signature)
	case "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256":
		sum := sha256Sum(elementToSign)
		return rsa.VerifyPKCS1v15(pubKey.(*rsa.PublicKey), crypto.SHA256, sum, signature)
	default:
		return fmt.Errorf("unsupported signature algorithm, %s", sigAlg)
	}
}

type dsaSignature struct {
	R, S *big.Int
}

func verifyDSA(signature, sum []byte, pubKey interface{}) error {
	dsaSig := new(dsaSignature)
	if rest, err := asn1.Unmarshal(signature, dsaSig); err != nil {
		return err
	} else if len(rest) != 0 {
		return fmt.Errorf("trailing data after DSA signature")
	}
	if dsaSig.R.Sign() <= 0 || dsaSig.S.Sign() <= 0 {
		return fmt.Errorf("DSA signature contained zero or negative values")
	}
	if !dsa.Verify(pubKey.(*dsa.PublicKey), sum, dsaSig.R, dsaSig.S) {
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
