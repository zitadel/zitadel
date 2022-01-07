package saml

import (
	"encoding/xml"
	"github.com/amdonov/xmlsig"
	"github.com/caos/zitadel/internal/api/saml/xml/metadata/xml_dsig"
	xml_dsigp "github.com/caos/zitadel/internal/api/saml/xml/protocol/xml_dsig"
)

func createSignatureM(signer xmlsig.Signer, data interface{}) (*xml_dsig.SignatureType, error) {
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
			InnerXml: "",
		},
		SignatureValue: xml_dsig.SignatureValueType{
			Text: sig.SignatureValue,
		},
		KeyInfo: &xml_dsig.KeyInfoType{
			XMLName: xml.Name{},
			X509Data: []xml_dsig.X509DataType{{
				X509Certificate: sig.KeyInfo.X509Data.X509Certificate,
			}},
			InnerXml: "",
		},
		InnerXml: "",
	}, nil
}

func createSignatureP(signer xmlsig.Signer, data interface{}) (*xml_dsigp.SignatureType, error) {
	sig, err := signer.CreateSignature(data)
	if err != nil {
		return nil, err
	}
	transforms := []xml_dsigp.TransformType{}
	for _, t := range sig.SignedInfo.Reference.Transforms.Transform {
		transforms = append(transforms, xml_dsigp.TransformType{
			XMLName:   xml.Name{},
			Algorithm: t.Algorithm,
		})
	}

	return &xml_dsigp.SignatureType{
		XMLName: xml.Name{},
		SignedInfo: xml_dsigp.SignedInfoType{
			XMLName: xml.Name{},
			CanonicalizationMethod: xml_dsigp.CanonicalizationMethodType{
				XMLName:   xml.Name{},
				Algorithm: sig.SignedInfo.CanonicalizationMethod.Algorithm,
			},
			SignatureMethod: xml_dsigp.SignatureMethodType{
				XMLName:   xml.Name{},
				Algorithm: sig.SignedInfo.SignatureMethod.Algorithm,
			},
			Reference: []xml_dsigp.ReferenceType{{
				DigestValue: xml_dsigp.DigestValueType(sig.SignedInfo.Reference.DigestValue),
				DigestMethod: xml_dsigp.DigestMethodType{
					XMLName:   xml.Name{},
					Algorithm: sig.SignedInfo.Reference.DigestMethod.Algorithm,
				},
				Transforms: &xml_dsigp.TransformsType{
					Transform: transforms,
				},
				URI: sig.SignedInfo.Reference.URI,
			}},
			InnerXml: "",
		},
		SignatureValue: xml_dsigp.SignatureValueType{
			Text: sig.SignatureValue,
		},
		KeyInfo: &xml_dsigp.KeyInfoType{
			XMLName: xml.Name{},
			X509Data: []xml_dsigp.X509DataType{{
				X509Certificate: sig.KeyInfo.X509Data.X509Certificate,
			}},
			InnerXml: "",
		},
		InnerXml: "",
	}, nil
}
