package saml

import (
	"encoding/xml"
	"github.com/caos/zitadel/internal/api/saml/xml/xenc"
	"github.com/caos/zitadel/internal/api/saml/xml/xml_dsig"
)

type BaseIDAbstractType struct {
	XMLName  xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:assertion BaseID"`
	InnerXml string   `xml:",innerxml"`
}

type NameIDType struct {
	XMLName         xml.Name
	Format          string `xml:"Format,attr,omitempty"`
	SPProvidedID    string `xml:"SPProvidedID,attr,omitempty"`
	NameQualifier   string `xml:"NameQualifier,attr,omitempty"`
	SPNameQualifier string `xml:"SPNameQualifier,attr,omitempty"`
	Text            string `xml:",chardata"`
	//InnerXml        string `xml:",innerxml"`
}

type EncryptedElementType struct {
	XMLName       xml.Name
	EncryptedData xenc.EncryptedDataType  `xml:"EncryptedData"`
	EncryptedKey  []xenc.EncryptedKeyType `xml:"EncryptedKey"`
	//InnerXml      string                  `xml:",innerxml"`
}

type AssertionType struct {
	XMLName                xml.Name                     `xml:"urn:oasis:names:tc:SAML:2.0:assertion Assertion"`
	Version                string                       `xml:"Version,attr"`
	Id                     string                       `xml:"ID,attr"`
	IssueInstant           string                       `xml:"IssueInstant,attr"`
	Issuer                 NameIDType                   `xml:"urn:oasis:names:tc:SAML:2.0:assertion Issuer"`
	Signature              *xml_dsig.SignatureType      `xml:"Signature"`
	Subject                *SubjectType                 `xml:"Subject"`
	Conditions             *ConditionsType              `xml:"Conditions"`
	Advice                 *AdviceType                  `xml:"Advice"`
	Statement              []StatementAbstractType      `xml:"urn:oasis:names:tc:SAML:2.0:assertion Statement"`
	AuthnStatement         []AuthnStatementType         `xml:"AuthnStatement"`
	AuthzDecisionStatement []AuthzDecisionStatementType `xml:"AuthzDecisionStatement"`
	AttributeStatement     []AttributeStatementType     `xml:"AttributeStatement"`
	//InnerXml               string                       `xml:",innerxml"`
}

type SubjectType struct {
	XMLName             xml.Name                  `xml:"urn:oasis:names:tc:SAML:2.0:assertion Subject"`
	BaseID              *BaseIDAbstractType       `xml:"BaseID"`
	NameID              *NameIDType               `xml:"urn:oasis:names:tc:SAML:2.0:assertion NameID"`
	EncryptedID         *EncryptedElementType     `xml:"urn:oasis:names:tc:SAML:2.0:assertion EncryptedID"`
	SubjectConfirmation []SubjectConfirmationType `xml:"SubjectConfirmation"`
	//InnerXml            string                    `xml:",innerxml"`
}

type SubjectConfirmationType struct {
	XMLName                 xml.Name                     `xml:"urn:oasis:names:tc:SAML:2.0:assertion SubjectConfirmation"`
	Method                  string                       `xml:"Method,attr"`
	SubjectConfirmationData *SubjectConfirmationDataType `xml:"SubjectConfirmationData"`
	BaseID                  *BaseIDAbstractType          `xml:"BaseID"`
	NameID                  *NameIDType                  `xml:"urn:oasis:names:tc:SAML:2.0:assertion NameID"`
	EncryptedID             *EncryptedElementType        `xml:"urn:oasis:names:tc:SAML:2.0:assertion EncryptedID"`
	//InnerXml                string                       `xml:",innerxml"`
}

type SubjectConfirmationDataType struct {
	XMLName      xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:assertion SubjectConfirmationData"`
	NotBefore    string   `xml:"NotBefore,attr,omitempty"`
	NotOnOrAfter string   `xml:"NotOnOrAfter,attr,omitempty"`
	Recipient    string   `xml:"Recipient,attr,omitempty"`
	InResponseTo string   `xml:"InResponseTo,attr,omitempty"`
	Address      string   `xml:"Address,attr,omitempty"`
	//InnerXml     string   `xml:",innerxml"`
}

type KeyInfoConfirmationDataType struct {
	XMLName      xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:assertion KeyInfoConfirmationData"`
	NotBefore    string   `xml:"NotBefore,attr,omitempty"`
	NotOnOrAfter string   `xml:"NotOnOrAfter,attr,omitempty"`
	Recipient    string   `xml:"Recipient,attr,omitempty"`
	InResponseTo string   `xml:"InResponseTo,attr,omitempty"`
	Address      string   `xml:"Address,attr,omitempty"`
	//InnerXml     string   `xml:",innerxml"`
}

type ConditionsType struct {
	XMLName             xml.Name                  `xml:"urn:oasis:names:tc:SAML:2.0:assertion Conditions"`
	NotBefore           string                    `xml:"NotBefore,attr,omitempty"`
	NotOnOrAfter        string                    `xml:"NotOnOrAfter,attr,omitempty"`
	Condition           []ConditionAbstractType   `xml:"urn:oasis:names:tc:SAML:2.0:assertion Condition"`
	AudienceRestriction []AudienceRestrictionType `xml:"AudienceRestriction"`
	OneTimeUse          []OneTimeUseType          `xml:"OneTimeUse"`
	ProxyRestriction    []ProxyRestrictionType    `xml:"ProxyRestriction"`
	//InnerXml            string                    `xml:",innerxml"`
}

type ConditionAbstractType struct {
	XMLName xml.Name
	//InnerXml string `xml:",innerxml"`
}

type AudienceRestrictionType struct {
	XMLName  xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:assertion AudienceRestriction"`
	Audience []string `xml:",any"`
	//InnerXml string   `xml:",innerxml"`
}

type OneTimeUseType struct {
	XMLName xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:assertion OneTimeUse"`
	//InnerXml string   `xml:",innerxml"`
}

type ProxyRestrictionType struct {
	XMLName  xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:assertion ProxyRestriction"`
	Count    int      `xml:"Count,attr,omitempty"`
	Audience []string `xml:",any"`
	//InnerXml string   `xml:",innerxml"`
}

type AdviceType struct {
	XMLName            xml.Name               `xml:"urn:oasis:names:tc:SAML:2.0:assertion Advice"`
	AssertionIDRef     []string               `xml:"AssertionIDRef"`
	AssertionURIRef    []string               `xml:"AssertionURIRef"`
	Assertion          []AssertionType        `xml:"Assertion"`
	EncryptedAssertion []EncryptedElementType `xml:"urn:oasis:names:tc:SAML:2.0:assertion EncryptedAssertion"`
	//InnerXml           string                 `xml:",innerxml"`
}

type StatementAbstractType struct {
	XMLName xml.Name
	//InnerXml string `xml:",innerxml"`
}

type AuthnStatementType struct {
	XMLName             xml.Name             `xml:"urn:oasis:names:tc:SAML:2.0:assertion AuthnStatement"`
	AuthnInstant        string               `xml:"AuthnInstant,attr"`
	SessionIndex        string               `xml:"SessionIndex,attr,omitempty"`
	SessionNotOnOrAfter string               `xml:"SessionNotOnOrAfter,attr,omitempty"`
	SubjectLocality     *SubjectLocalityType `xml:"SubjectLocality"`
	AuthnContext        AuthnContextType     `xml:"AuthnContext"`
	//InnerXml            string               `xml:",innerxml"`
}

type SubjectLocalityType struct {
	XMLName xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:assertion SubjectLocality"`
	Address string   `xml:"Address,attr,omitempty"`
	DNSName string   `xml:"DNSName,attr,omitempty"`
	//InnerXml string   `xml:",innerxml"`
}

type AuthnContextType struct {
	XMLName                 xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:assertion AuthnContext"`
	AuthenticatingAuthority []string `xml:"urn:oasis:names:tc:SAML:2.0:assertion AuthenticatingAuthority"`
	AuthnContextClassRef    string   `xml:"urn:oasis:names:tc:SAML:2.0:assertion AuthnContextClassRef,omitempty"`
	AuthnContextDecl        string   `xml:"urn:oasis:names:tc:SAML:2.0:assertion AuthnContextDecl,omitempty"`
	AuthnContextDeclRef     string   `xml:"urn:oasis:names:tc:SAML:2.0:assertion AuthnContextDeclRef,omitempty"`
	//InnerXml                string   `xml:",innerxml"`
}

type AuthzDecisionStatementType struct {
	XMLName  xml.Name      `xml:"urn:oasis:names:tc:SAML:2.0:assertion AuthzDecisionStatement"`
	Resource string        `xml:"Resource,attr"`
	Decision DecisionType  `xml:"Decision,attr"`
	Action   []ActionType  `xml:"Action"`
	Evidence *EvidenceType `xml:"Evidence"`
	//	InnerXml string        `xml:",innerxml"`
}

type ActionType struct {
	XMLName   xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:assertion Action"`
	Namespace string   `xml:"Namespace,attr"`
	Text      string   `xml:",chardata"`
	//	InnerXml  string   `xml:",innerxml"`
}

type EvidenceType struct {
	XMLName            xml.Name               `xml:"urn:oasis:names:tc:SAML:2.0:assertion Evidence"`
	AssertionIDRef     []string               `xml:"AssertionIDRef"`
	AssertionURIRef    []string               `xml:"AssertionURIRef"`
	Assertion          []AssertionType        `xml:"Assertion"`
	EncryptedAssertion []EncryptedElementType `xml:"urn:oasis:names:tc:SAML:2.0:assertion EncryptedAssertion"`
	//	InnerXml           string                 `xml:",innerxml"`
}

type AttributeStatementType struct {
	XMLName   xml.Name         `xml:"urn:oasis:names:tc:SAML:2.0:assertion AttributeStatement"`
	Attribute []*AttributeType `xml:"Attribute"`
	//InnerXml  string           `xml:",innerxml"`
}

type AttributeType struct {
	XMLName        xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:assertion Attribute"`
	Name           string   `xml:"Name,attr"`
	NameFormat     string   `xml:"NameFormat,attr,omitempty"`
	FriendlyName   string   `xml:"FriendlyName,attr,omitempty"`
	AttributeValue []string `xml:",any"`
	//InnerXml       string   `xml:",innerxml"`
}

// XSD SimpleType declarations

type DecisionType string

const DecisionTypePermit DecisionType = "Permit"

const DecisionTypeDeny DecisionType = "Deny"

const DecisionTypeIndeterminate DecisionType = "Indeterminate"
