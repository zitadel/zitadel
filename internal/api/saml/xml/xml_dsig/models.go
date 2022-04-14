package xml_dsig

import "encoding/xml"

type SignatureType struct {
	XMLName        xml.Name           `xml:"http://www.w3.org/2000/09/xmldsig# Signature"`
	Id             string             `xml:"Id,attr,omitempty"`
	SignedInfo     SignedInfoType     `xml:"SignedInfo"`
	SignatureValue SignatureValueType `xml:"SignatureValue"`
	KeyInfo        *KeyInfoType       `xml:"http://www.w3.org/2000/09/xmldsig# KeyInfo"`
	Object         []ObjectType       `xml:"Object"`
	//InnerXml       string             `xml:",innerxml"`
}

type SignatureValueType struct {
	XMLName xml.Name `xml:"http://www.w3.org/2000/09/xmldsig# SignatureValue"`
	Id      string   `xml:"Id,attr,omitempty"`
	Text    string   `xml:",chardata"`
	//InnerXml string   `xml:",innerxml"`
}

type SignedInfoType struct {
	XMLName                xml.Name                   `xml:"http://www.w3.org/2000/09/xmldsig# SignedInfo"`
	Id                     string                     `xml:"Id,attr,omitempty"`
	CanonicalizationMethod CanonicalizationMethodType `xml:"CanonicalizationMethod"`
	SignatureMethod        SignatureMethodType        `xml:"SignatureMethod"`
	Reference              []ReferenceType            `xml:"Reference"`
	//InnerXml               string                     `xml:",innerxml"`
}

type CanonicalizationMethodType struct {
	XMLName   xml.Name `xml:"http://www.w3.org/2000/09/xmldsig# CanonicalizationMethod"`
	Algorithm string   `xml:"Algorithm,attr"`
	//InnerXml  string   `xml:",innerxml"`
}

type SignatureMethodType struct {
	XMLName          xml.Name              `xml:"http://www.w3.org/2000/09/xmldsig# SignatureMethod"`
	Algorithm        string                `xml:"Algorithm,attr"`
	HMACOutputLength *HMACOutputLengthType `xml:",any"`
	//InnerXml         string                `xml:",innerxml"`
}

type ReferenceType struct {
	XMLName      xml.Name         `xml:"http://www.w3.org/2000/09/xmldsig# Reference"`
	Id           string           `xml:"Id,attr,omitempty"`
	URI          string           `xml:"URI,attr,omitempty"`
	Type         string           `xml:"Type,attr,omitempty"`
	Transforms   *TransformsType  `xml:"Transforms"`
	DigestMethod DigestMethodType `xml:"DigestMethod"`
	DigestValue  DigestValueType  `xml:"http://www.w3.org/2000/09/xmldsig# DigestValue"`
	//InnerXml     string           `xml:",innerxml"`
}

type TransformsType struct {
	XMLName   xml.Name        `xml:"http://www.w3.org/2000/09/xmldsig# Transforms"`
	Transform []TransformType `xml:",any"`
	//InnerXml  string          `xml:",innerxml"`
}

type TransformType struct {
	XMLName   xml.Name `xml:"http://www.w3.org/2000/09/xmldsig# Transform"`
	Algorithm string   `xml:"Algorithm,attr"`
	XPath     []string `xml:"XPath"`
	//InnerXml  string   `xml:",innerxml"`
}

type DigestMethodType struct {
	XMLName   xml.Name `xml:"http://www.w3.org/2000/09/xmldsig# DigestMethod"`
	Algorithm string   `xml:"Algorithm,attr"`
	//InnerXml  string   `xml:",innerxml"`
}

type KeyInfoType struct {
	XMLName         xml.Name
	Id              string                `xml:"Id,attr,omitempty"`
	KeyName         []string              `xml:"KeyName"`
	KeyValue        []KeyValueType        `xml:"KeyValue"`
	RetrievalMethod []RetrievalMethodType `xml:"RetrievalMethod"`
	X509Data        []X509DataType        `xml:"X509Data"`
	PGPData         []PGPDataType         `xml:"PGPData"`
	SPKIData        []SPKIDataType        `xml:"SPKIData"`
	MgmtData        []string              `xml:"MgmtData"`
	//InnerXml        string                `xml:",innerxml"`
}

type KeyValueType struct {
	XMLName     xml.Name         `xml:"http://www.w3.org/2000/09/xmldsig# KeyValue"`
	DSAKeyValue *DSAKeyValueType `xml:"DSAKeyValue"`
	RSAKeyValue *RSAKeyValueType `xml:"RSAKeyValue"`
	//InnerXml    string           `xml:",innerxml"`
}

type RetrievalMethodType struct {
	XMLName    xml.Name        `xml:"http://www.w3.org/2000/09/xmldsig# RetrievalMethod"`
	URI        string          `xml:"URI,attr"`
	Type       string          `xml:"Type,attr,omitempty"`
	Transforms *TransformsType `xml:",any"`
	//InnerXml   string          `xml:",innerxml"`
}

type X509DataType struct {
	XMLName          xml.Name              `xml:"http://www.w3.org/2000/09/xmldsig# X509Data"`
	X509IssuerSerial *X509IssuerSerialType `xml:"X509IssuerSerial"`
	X509SKI          string                `xml:"X509SKI,omitempty"`
	X509SubjectName  string                `xml:"X509SubjectName,omitempty"`
	X509Certificate  string                `xml:"http://www.w3.org/2000/09/xmldsig# X509Certificate"`
	X509CRL          string                `xml:"X509CRL,omitempty"`
	//InnerXml         string                `xml:",innerxml"`
}

type X509IssuerSerialType struct {
	XMLName          xml.Name `xml:"http://www.w3.org/2000/09/xmldsig# X509IssuerSerial"`
	X509IssuerName   string   `xml:"X509IssuerName"`
	X509SerialNumber int64    `xml:"X509SerialNumber"`
	//InnerXml         string   `xml:",innerxml"`
}

type PGPDataType struct {
	XMLName      xml.Name `xml:"http://www.w3.org/2000/09/xmldsig# PGPData"`
	PGPKeyID     string   `xml:"PGPKeyID"`
	PGPKeyPacket string   `xml:"PGPKeyPacket"`
	//InnerXml     string   `xml:",innerxml"`
}

type SPKIDataType struct {
	XMLName  xml.Name `xml:"http://www.w3.org/2000/09/xmldsig# SPKIData"`
	SPKISexp string   `xml:",any"`
	//InnerXml string   `xml:",innerxml"`
}

type ObjectType struct {
	XMLName  xml.Name `xml:"http://www.w3.org/2000/09/xmldsig# Object"`
	Id       string   `xml:"Id,attr,omitempty"`
	MimeType string   `xml:"MimeType,attr,omitempty"`
	Encoding string   `xml:"Encoding,attr,omitempty"`
	//InnerXml string   `xml:",innerxml"`
}

type ManifestType struct {
	XMLName   xml.Name        `xml:"http://www.w3.org/2000/09/xmldsig# Manifest"`
	Id        string          `xml:"Id,attr,omitempty"`
	Reference []ReferenceType `xml:",any"`
	//InnerXml  string          `xml:",innerxml"`
}

type SignaturePropertiesType struct {
	XMLName           xml.Name                `xml:"http://www.w3.org/2000/09/xmldsig# SignatureProperties"`
	Id                string                  `xml:"Id,attr,omitempty"`
	SignatureProperty []SignaturePropertyType `xml:",any"`
	//InnerXml          string                  `xml:",innerxml"`
}

type SignaturePropertyType struct {
	XMLName xml.Name `xml:"http://www.w3.org/2000/09/xmldsig# SignatureProperty"`
	Target  string   `xml:"Target,attr"`
	Id      string   `xml:"Id,attr,omitempty"`
	//InnerXml string   `xml:",innerxml"`
}

type DSAKeyValueType struct {
	XMLName xml.Name      `xml:"http://www.w3.org/2000/09/xmldsig# DSAKeyValue"`
	G       *CryptoBinary `xml:"G"`
	Y       CryptoBinary  `xml:"Y"`
	J       *CryptoBinary `xml:"J"`
	//InnerXml string        `xml:",innerxml"`
}

type RSAKeyValueType struct {
	XMLName  xml.Name     `xml:"http://www.w3.org/2000/09/xmldsig# RSAKeyValue"`
	Modulus  CryptoBinary `xml:"Modulus"`
	Exponent CryptoBinary `xml:"Exponent"`
	//InnerXml string       `xml:",innerxml"`
}

type CryptoBinary string

type DigestValueType string

type HMACOutputLengthType int64

const (
	DigestMethodSHA256 = "http://www.w3.org/2001/04/xmlenc#sha256"
	DigestMethodSHA1   = "http://www.w3.org/2000/09/xmldsig#sha1"
	DigestMethodSHA512 = "http://www.w3.org/2001/04/xmlenc#sha512"
)
