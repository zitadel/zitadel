package xenc

import (
	"encoding/xml"
	"github.com/caos/zitadel/internal/api/saml/xml/xml_dsig"
)

type EncryptedType struct {
	XMLName              xml.Name
	Id                   string                    `xml:"Id,attr,omitempty"`
	Type                 string                    `xml:"Type,attr,omitempty"`
	MimeType             string                    `xml:"MimeType,attr,omitempty"`
	Encoding             string                    `xml:"Encoding,attr,omitempty"`
	EncryptionMethod     *EncryptionMethodType     `xml:"EncryptionMethod"`
	KeyInfo              *xml_dsig.KeyInfoType     `xml:"http://www.w3.org/2000/09/xmldsig# KeyInfo"`
	CipherData           CipherDataType            `xml:"CipherData"`
	EncryptionProperties *EncryptionPropertiesType `xml:"EncryptionProperties"`
	//InnerXml             string                    `xml:",innerxml"`
}

type EncryptionMethodType struct {
	XMLName    xml.Name
	Algorithm  string       `xml:"Algorithm,attr"`
	KeySize    *KeySizeType `xml:"KeySize"`
	OAEPparams string       `xml:"OAEPparams"`
	//InnerXml   string       `xml:",innerxml"`
}

type CipherDataType struct {
	XMLName         xml.Name
	CipherValue     string               `xml:"CipherValue"`
	CipherReference *CipherReferenceType `xml:"CipherReference"`
	//InnerXml        string               `xml:",innerxml"`
}

type CipherReferenceType struct {
	XMLName    xml.Name
	URI        string          `xml:"URI,attr"`
	Transforms *TransformsType `xml:"Transforms"`
	//InnerXml   string          `xml:",innerxml"`
}

type TransformsType struct {
	XMLName   xml.Name
	Transform []xml_dsig.TransformType `xml:",any"`
	//InnerXml  string                   `xml:",innerxml"`
}

type EncryptedDataType struct {
	XMLName              xml.Name
	Id                   string                    `xml:"Id,attr,omitempty"`
	Type                 string                    `xml:"Type,attr,omitempty"`
	MimeType             string                    `xml:"MimeType,attr,omitempty"`
	Encoding             string                    `xml:"Encoding,attr,omitempty"`
	EncryptionMethod     *EncryptionMethodType     `xml:"EncryptionMethod"`
	KeyInfo              *xml_dsig.KeyInfoType     `xml:"http://www.w3.org/2000/09/xmldsig# KeyInfo"`
	CipherData           CipherDataType            `xml:"CipherData"`
	EncryptionProperties *EncryptionPropertiesType `xml:"EncryptionProperties"`
	//InnerXml             string                    `xml:",innerxml"`
}

type EncryptedKeyType struct {
	XMLName              xml.Name
	Recipient            string                    `xml:"Recipient,attr,omitempty"`
	Id                   string                    `xml:"Id,attr,omitempty"`
	Type                 string                    `xml:"Type,attr,omitempty"`
	MimeType             string                    `xml:"MimeType,attr,omitempty"`
	Encoding             string                    `xml:"Encoding,attr,omitempty"`
	ReferenceList        *ReferenceListType        `xml:"ReferenceList"`
	CarriedKeyName       string                    `xml:"CarriedKeyName"`
	EncryptionMethod     *EncryptionMethodType     `xml:"EncryptionMethod"`
	KeyInfo              *xml_dsig.KeyInfoType     `xml:"http://www.w3.org/2000/09/xmldsig# KeyInfo"`
	CipherData           CipherDataType            `xml:"CipherData"`
	EncryptionProperties *EncryptionPropertiesType `xml:"EncryptionProperties"`
	//InnerXml             string                    `xml:",innerxml"`
}

type AgreementMethodType struct {
	XMLName           xml.Name
	Algorithm         string                `xml:"Algorithm,attr"`
	KANonce           string                `xml:"KA-Nonce"`
	OriginatorKeyInfo *xml_dsig.KeyInfoType `xml:"http://www.w3.org/2000/09/xmldsig# OriginatorKeyInfo"`
	RecipientKeyInfo  *xml_dsig.KeyInfoType `xml:"http://www.w3.org/2000/09/xmldsig# RecipientKeyInfo"`
	//InnerXml          string                `xml:",innerxml"`
}

type ReferenceType struct {
	XMLName xml.Name
	URI     string `xml:"URI,attr"`
	//InnerXml string `xml:",innerxml"`
}

type EncryptionPropertiesType struct {
	XMLName            xml.Name
	Id                 string                   `xml:"Id,attr,omitempty"`
	EncryptionProperty []EncryptionPropertyType `xml:",any"`
	//InnerXml           string                   `xml:",innerxml"`
}

type EncryptionPropertyType struct {
	XMLName xml.Name
	Target  string `xml:"Target,attr,omitempty"`
	Id      string `xml:"Id,attr,omitempty"`
	//	InnerXml string `xml:",innerxml"`
}

type ReferenceListType struct {
	XMLName       xml.Name        `xml:"ReferenceList"`
	DataReference []ReferenceType `xml:"DataReference"`
	KeyReference  []ReferenceType `xml:"KeyReference"`
}

// XSD SimpleType declarations

type KeySizeType int64
