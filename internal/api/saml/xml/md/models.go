package md

import (
	"encoding/xml"
	"github.com/caos/zitadel/internal/api/saml/xml/saml"
	"github.com/caos/zitadel/internal/api/saml/xml/xenc"
	"github.com/caos/zitadel/internal/api/saml/xml/xml_dsig"
)

type LocalizedNameType struct {
	XMLName xml.Name
	XmlLang string `xml:"lang,attr"`
	Text    string `xml:",chardata"`
	//InnerXml string `xml:",innerxml"`
}

type LocalizedURIType struct {
	XMLName xml.Name
	XmlLang string `xml:"lang,attr"`
	Text    string `xml:",chardata"`
	//InnerXml string `xml:",innerxml"`
}

type ExtensionsType struct {
	XMLName xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:metadata Extensions"`
	//InnerXml string   `xml:",innerxml"`
}

type EndpointType struct {
	XMLName          xml.Name
	Binding          string `xml:"Binding,attr"`
	Location         string `xml:"Location,attr"`
	ResponseLocation string `xml:"ResponseLocation,attr,omitempty"`
	//InnerXml         string `xml:",innerxml"`
}

type IndexedEndpointType struct {
	XMLName          xml.Name
	Index            string `xml:"index,attr"`
	IsDefault        string `xml:"isDefault,attr,omitempty"`
	Binding          string `xml:"Binding,attr"`
	Location         string `xml:"Location,attr"`
	ResponseLocation string `xml:"ResponseLocation,attr,omitempty"`
	//InnerXml         string `xml:",innerxml"`
}

type EntitiesDescriptorType struct {
	XMLName            xml.Name                 `xml:"urn:oasis:names:tc:SAML:2.0:metadata EntitiesDescriptor"`
	ValidUntil         string                   `xml:"validUntil,attr,omitempty"`
	CacheDuration      string                   `xml:"cacheDuration,attr,omitempty"`
	Id                 string                   `xml:"ID,attr,omitempty"`
	Name               string                   `xml:"Name,attr,omitempty"`
	Signature          *xml_dsig.SignatureType  `xml:"Signature"`
	Extensions         *ExtensionsType          `xml:"Extensions"`
	EntityDescriptor   []EntityDescriptorType   `xml:"EntityDescriptor"`
	EntitiesDescriptor []EntitiesDescriptorType `xml:"EntitiesDescriptor"`
	//InnerXml           string                   `xml:",innerxml"`
}

type EntityDescriptorType struct {
	XMLName                      xml.Name                          `xml:"urn:oasis:names:tc:SAML:2.0:metadata EntityDescriptor"`
	EntityID                     EntityIDType                      `xml:"entityID,attr"`
	ValidUntil                   string                            `xml:"validUntil,attr,omitempty"`
	CacheDuration                string                            `xml:"cacheDuration,attr,omitempty"`
	Id                           string                            `xml:"ID,attr,omitempty"`
	Signature                    *xml_dsig.SignatureType           `xml:"Signature"`
	Extensions                   *ExtensionsType                   `xml:"Extensions"`
	Organization                 *OrganizationType                 `xml:"Organization"`
	ContactPerson                []ContactType                     `xml:"ContactPerson"`
	AdditionalMetadataLocation   []AdditionalMetadataLocationType  `xml:"AdditionalMetadataLocation"`
	RoleDescriptor               *RoleDescriptorType               `xml:"RoleDescriptor,omitempty"`
	IDPSSODescriptor             *IDPSSODescriptorType             `xml:"IDPSSODescriptor,omitempty"`
	SPSSODescriptor              *SPSSODescriptorType              `xml:"SPSSODescriptor,omitempty"`
	AuthnAuthorityDescriptor     *AuthnAuthorityDescriptorType     `xml:"AuthnAuthorityDescriptor,omitempty"`
	AttributeAuthorityDescriptor *AttributeAuthorityDescriptorType `xml:"AttributeAuthorityDescriptor,omitempty"`
	PDPDescriptor                *PDPDescriptorType                `xml:"PDPDescriptor,omitempty"`
	AffiliationDescriptor        *AffiliationDescriptorType        `xml:"AffiliationDescriptor"`
	//InnerXml                     string                            `xml:",innerxml"`
}

type OrganizationType struct {
	XMLName                 xml.Name            `xml:"urn:oasis:names:tc:SAML:2.0:metadata Organization"`
	Extensions              *ExtensionsType     `xml:"Extensions"`
	OrganizationName        []LocalizedNameType `xml:"urn:oasis:names:tc:SAML:2.0:metadata OrganizationName,omitempty"`
	OrganizationDisplayName []LocalizedNameType `xml:"urn:oasis:names:tc:SAML:2.0:metadata OrganizationDisplayName,omitempty"`
	OrganizationURL         []LocalizedURIType  `xml:"urn:oasis:names:tc:SAML:2.0:metadata OrganizationURL,omitempty"`
	//InnerXml                string              `xml:",innerxml"`
}

type ContactType struct {
	XMLName         xml.Name        `xml:"urn:oasis:names:tc:SAML:2.0:metadata ContactPerson"`
	ContactType     ContactTypeType `xml:"contactType,attr"`
	Extensions      *ExtensionsType `xml:"Extensions"`
	Company         string          `xml:"Company,omitempty"`
	GivenName       string          `xml:"GivenName,omitempty"`
	SurName         string          `xml:"SurName,omitempty"`
	EmailAddress    []string        `xml:"EmailAddress,omitempty"`
	TelephoneNumber []string        `xml:"TelephoneNumber,omitempty"`
	//InnerXml        string          `xml:",innerxml"`
}

type AdditionalMetadataLocationType struct {
	XMLName   xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:metadata AdditionalMetadataLocation"`
	Namespace string   `xml:"namespace,attr"`
	Text      string   `xml:",chardata"`
	//	InnerXml  string   `xml:",innerxml"`
}

type RoleDescriptorType struct {
	XMLName                    xml.Name                `xml:"urn:oasis:names:tc:SAML:2.0:metadata RoleDescriptor"`
	Id                         string                  `xml:"ID,attr,omitempty"`
	ValidUntil                 string                  `xml:"validUntil,attr,omitempty"`
	CacheDuration              string                  `xml:"cacheDuration,attr,omitempty"`
	ProtocolSupportEnumeration AnyURIListType          `xml:"protocolSupportEnumeration,attr"`
	ErrorURL                   string                  `xml:"errorURL,attr,omitempty"`
	Signature                  *xml_dsig.SignatureType `xml:"Signature"`
	Extensions                 *ExtensionsType         `xml:"Extensions"`
	KeyDescriptor              []KeyDescriptorType     `xml:"KeyDescriptor"`
	Organization               *OrganizationType       `xml:"Organization"`
	ContactPerson              []ContactType           `xml:"ContactPerson"`
	//InnerXml                   string                  `xml:",innerxml"`
}

type KeyDescriptorType struct {
	XMLName          xml.Name                    `xml:"urn:oasis:names:tc:SAML:2.0:metadata KeyDescriptor"`
	Use              KeyTypes                    `xml:"use,attr,omitempty"`
	KeyInfo          xml_dsig.KeyInfoType        `xml:"http://www.w3.org/2000/09/xmldsig# KeyInfo"`
	EncryptionMethod []xenc.EncryptionMethodType `xml:"EncryptionMethod"`
	//InnerXml         string                      `xml:",innerxml"`
}

type SSODescriptorType struct {
	XMLName                    xml.Name                `xml:"urn:oasis:names:tc:SAML:2.0:metadata SSODescriptor"`
	Id                         string                  `xml:"ID,attr,omitempty"`
	ValidUntil                 string                  `xml:"validUntil,attr,omitempty"`
	CacheDuration              string                  `xml:"cacheDuration,attr,omitempty"`
	ProtocolSupportEnumeration AnyURIListType          `xml:"protocolSupportEnumeration,attr"`
	ErrorURL                   string                  `xml:"errorURL,attr,omitempty"`
	ArtifactResolutionService  []IndexedEndpointType   `xml:"urn:oasis:names:tc:SAML:2.0:metadata ArtifactResolutionService"`
	SingleLogoutService        []EndpointType          `xml:"urn:oasis:names:tc:SAML:2.0:metadata SingleLogoutService"`
	ManageNameIDService        []EndpointType          `xml:"urn:oasis:names:tc:SAML:2.0:metadata ManageNameIDService"`
	NameIDFormat               []string                `xml:"NameIDFormat"`
	Signature                  *xml_dsig.SignatureType `xml:"Signature"`
	Extensions                 *ExtensionsType         `xml:"Extensions"`
	KeyDescriptor              []KeyDescriptorType     `xml:"KeyDescriptor"`
	Organization               *OrganizationType       `xml:"Organization"`
	ContactPerson              []ContactType           `xml:"ContactPerson"`
	//	InnerXml                   string                  `xml:",innerxml"`
}

type IDPSSODescriptorType struct {
	XMLName                    xml.Name                `xml:"urn:oasis:names:tc:SAML:2.0:metadata IDPSSODescriptor"`
	WantAuthnRequestsSigned    string                  `xml:"WantAuthnRequestsSigned,attr,omitempty"`
	Id                         string                  `xml:"ID,attr,omitempty"`
	ValidUntil                 string                  `xml:"validUntil,attr,omitempty"`
	CacheDuration              string                  `xml:"cacheDuration,attr,omitempty"`
	ProtocolSupportEnumeration AnyURIListType          `xml:"protocolSupportEnumeration,attr"`
	ErrorURL                   string                  `xml:"errorURL,attr,omitempty"`
	SingleSignOnService        []EndpointType          `xml:"urn:oasis:names:tc:SAML:2.0:metadata SingleSignOnService"`
	NameIDMappingService       []EndpointType          `xml:"urn:oasis:names:tc:SAML:2.0:metadata NameIDMappingService"`
	AssertionIDRequestService  []EndpointType          `xml:"urn:oasis:names:tc:SAML:2.0:metadata AssertionIDRequestService"`
	AttributeProfile           []string                `xml:"AttributeProfile"`
	Attribute                  []*saml.AttributeType   `xml:"Attribute"`
	ArtifactResolutionService  []IndexedEndpointType   `xml:"urn:oasis:names:tc:SAML:2.0:metadata ArtifactResolutionService"`
	SingleLogoutService        []EndpointType          `xml:"urn:oasis:names:tc:SAML:2.0:metadata SingleLogoutService"`
	ManageNameIDService        []EndpointType          `xml:"urn:oasis:names:tc:SAML:2.0:metadata ManageNameIDService"`
	NameIDFormat               []string                `xml:"NameIDFormat"`
	Signature                  *xml_dsig.SignatureType `xml:"Signature"`
	Extensions                 *ExtensionsType         `xml:"Extensions"`
	KeyDescriptor              []KeyDescriptorType     `xml:"KeyDescriptor"`
	Organization               *OrganizationType       `xml:"Organization"`
	ContactPerson              []ContactType           `xml:"ContactPerson"`
	//InnerXml                   string                  `xml:",innerxml"`
}

type SPSSODescriptorType struct {
	XMLName                    xml.Name                        `xml:"urn:oasis:names:tc:SAML:2.0:metadata SPSSODescriptor"`
	AuthnRequestsSigned        string                          `xml:"AuthnRequestsSigned,attr,omitempty"`
	WantAssertionsSigned       string                          `xml:"WantAssertionsSigned,attr,omitempty"`
	Id                         string                          `xml:"ID,attr,omitempty"`
	ValidUntil                 string                          `xml:"validUntil,attr,omitempty"`
	CacheDuration              string                          `xml:"cacheDuration,attr,omitempty"`
	ProtocolSupportEnumeration AnyURIListType                  `xml:"protocolSupportEnumeration,attr"`
	ErrorURL                   string                          `xml:"errorURL,attr,omitempty"`
	AssertionConsumerService   []IndexedEndpointType           `xml:"urn:oasis:names:tc:SAML:2.0:metadata AssertionConsumerService"`
	AttributeConsumingService  []AttributeConsumingServiceType `xml:"AttributeConsumingService"`
	ArtifactResolutionService  []IndexedEndpointType           `xml:"urn:oasis:names:tc:SAML:2.0:metadata ArtifactResolutionService"`
	SingleLogoutService        []EndpointType                  `xml:"urn:oasis:names:tc:SAML:2.0:metadata SingleLogoutService"`
	ManageNameIDService        []EndpointType                  `xml:"urn:oasis:names:tc:SAML:2.0:metadata ManageNameIDService"`
	NameIDFormat               []string                        `xml:"NameIDFormat"`
	Signature                  *xml_dsig.SignatureType         `xml:"Signature"`
	Extensions                 *ExtensionsType                 `xml:"Extensions"`
	KeyDescriptor              []KeyDescriptorType             `xml:"KeyDescriptor"`
	Organization               *OrganizationType               `xml:"Organization"`
	ContactPerson              []ContactType                   `xml:"ContactPerson"`
	//	InnerXml                   string                          `xml:",innerxml"`
}

type AttributeConsumingServiceType struct {
	XMLName            xml.Name                 `xml:"urn:oasis:names:tc:SAML:2.0:metadata AttributeConsumingService"`
	Index              uint64                   `xml:"index,attr"`
	IsDefault          bool                     `xml:"isDefault,attr,omitempty"`
	ServiceName        []LocalizedNameType      `xml:"urn:oasis:names:tc:SAML:2.0:metadata ServiceName"`
	ServiceDescription []LocalizedNameType      `xml:"urn:oasis:names:tc:SAML:2.0:metadata ServiceDescription"`
	RequestedAttribute []RequestedAttributeType `xml:"RequestedAttribute"`
	//InnerXml           string                   `xml:",innerxml"`
}

type RequestedAttributeType struct {
	XMLName        xml.Name `xml:"urn:oasis:names:tc:SAML:2.0:metadata RequestedAttribute"`
	IsRequired     string   `xml:"isRequired,attr,omitempty"`
	Name           string   `xml:"Name,attr"`
	NameFormat     string   `xml:"NameFormat,attr,omitempty"`
	FriendlyName   string   `xml:"FriendlyName,attr,omitempty"`
	AttributeValue []string `xml:",any"`
	//InnerXml       string   `xml:",innerxml"`
}

type AuthnAuthorityDescriptorType struct {
	XMLName                    xml.Name                `xml:"urn:oasis:names:tc:SAML:2.0:metadata AuthnAuthorityDescriptor"`
	Id                         string                  `xml:"ID,attr,omitempty"`
	ValidUntil                 string                  `xml:"validUntil,attr,omitempty"`
	CacheDuration              string                  `xml:"cacheDuration,attr,omitempty"`
	ProtocolSupportEnumeration AnyURIListType          `xml:"protocolSupportEnumeration,attr"`
	ErrorURL                   string                  `xml:"errorURL,attr,omitempty"`
	AuthnQueryService          []EndpointType          `xml:"urn:oasis:names:tc:SAML:2.0:metadata AuthnQueryService"`
	AssertionIDRequestService  []EndpointType          `xml:"urn:oasis:names:tc:SAML:2.0:metadata AssertionIDRequestService"`
	NameIDFormat               []string                `xml:"NameIDFormat"`
	Signature                  *xml_dsig.SignatureType `xml:"Signature"`
	Extensions                 *ExtensionsType         `xml:"Extensions"`
	KeyDescriptor              []KeyDescriptorType     `xml:"KeyDescriptor"`
	Organization               *OrganizationType       `xml:"Organization"`
	ContactPerson              []ContactType           `xml:"ContactPerson"`
	//InnerXml                   string                  `xml:",innerxml"`
}

type PDPDescriptorType struct {
	XMLName                    xml.Name                `xml:"urn:oasis:names:tc:SAML:2.0:metadata PDPDescriptor"`
	Id                         string                  `xml:"ID,attr,omitempty"`
	ValidUntil                 string                  `xml:"validUntil,attr,omitempty"`
	CacheDuration              string                  `xml:"cacheDuration,attr,omitempty"`
	ProtocolSupportEnumeration AnyURIListType          `xml:"protocolSupportEnumeration,attr"`
	ErrorURL                   string                  `xml:"errorURL,attr,omitempty"`
	AuthzService               []EndpointType          `xml:"urn:oasis:names:tc:SAML:2.0:metadata AuthzService"`
	AssertionIDRequestService  []EndpointType          `xml:"urn:oasis:names:tc:SAML:2.0:metadata AssertionIDRequestService"`
	NameIDFormat               []string                `xml:"NameIDFormat"`
	Signature                  *xml_dsig.SignatureType `xml:"Signature"`
	Extensions                 *ExtensionsType         `xml:"Extensions"`
	KeyDescriptor              []KeyDescriptorType     `xml:"KeyDescriptor"`
	Organization               *OrganizationType       `xml:"Organization"`
	ContactPerson              []ContactType           `xml:"ContactPerson"`
	//InnerXml                   string                  `xml:",innerxml"`
}

type AttributeAuthorityDescriptorType struct {
	XMLName                    xml.Name                `xml:"urn:oasis:names:tc:SAML:2.0:metadata AttributeAuthorityDescriptor"`
	Id                         string                  `xml:"ID,attr,omitempty"`
	ValidUntil                 string                  `xml:"validUntil,attr,omitempty"`
	CacheDuration              string                  `xml:"cacheDuration,attr,omitempty"`
	ProtocolSupportEnumeration AnyURIListType          `xml:"protocolSupportEnumeration,attr"`
	ErrorURL                   string                  `xml:"errorURL,attr,omitempty"`
	AttributeService           []EndpointType          `xml:"urn:oasis:names:tc:SAML:2.0:metadata AttributeService"`
	AssertionIDRequestService  []EndpointType          `xml:"urn:oasis:names:tc:SAML:2.0:metadata AssertionIDRequestService"`
	NameIDFormat               []string                `xml:"NameIDFormat"`
	AttributeProfile           []string                `xml:"AttributeProfile"`
	Attribute                  []*saml.AttributeType   `xml:"Attribute"`
	Signature                  *xml_dsig.SignatureType `xml:"Signature"`
	Extensions                 *ExtensionsType         `xml:"Extensions"`
	KeyDescriptor              []KeyDescriptorType     `xml:"KeyDescriptor"`
	Organization               *OrganizationType       `xml:"Organization"`
	ContactPerson              []ContactType           `xml:"ContactPerson"`
	//InnerXml                   string                  `xml:",innerxml"`
}

type AffiliationDescriptorType struct {
	XMLName            xml.Name                `xml:"urn:oasis:names:tc:SAML:2.0:metadata AffiliationDescriptor"`
	AffiliationOwnerID EntityIDType            `xml:"affiliationOwnerID,attr"`
	ValidUntil         string                  `xml:"validUntil,attr,omitempty"`
	CacheDuration      string                  `xml:"cacheDuration,attr,omitempty"`
	Id                 string                  `xml:"ID,attr,omitempty"`
	Signature          *xml_dsig.SignatureType `xml:"Signature"`
	Extensions         *ExtensionsType         `xml:"Extensions"`
	AffiliateMember    []EntityIDType          `xml:"AffiliateMember"`
	KeyDescriptor      []KeyDescriptorType     `xml:"KeyDescriptor"`
	//InnerXml           string                  `xml:",innerxml"`
}

// XSD SimpleType declarations

type EntityIDType string
type ContactTypeType string

const ContactTypeTypeTechnical ContactTypeType = "technical"
const ContactTypeTypeSupport ContactTypeType = "support"
const ContactTypeTypeAdministrative ContactTypeType = "administrative"
const ContactTypeTypeBilling ContactTypeType = "billing"
const ContactTypeTypeOther ContactTypeType = "other"

type AnyURIListType string
type KeyTypes string

const KeyTypesEncryption KeyTypes = "encryption"
const KeyTypesSigning KeyTypes = "signing"
