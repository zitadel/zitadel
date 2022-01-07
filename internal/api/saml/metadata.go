package saml

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"github.com/caos/zitadel/internal/api/saml/xml/metadata/md"
	"github.com/caos/zitadel/internal/api/saml/xml/metadata/xenc"
	"github.com/caos/zitadel/internal/api/saml/xml/metadata/xml_dsig"
	"net/http"
)

func (p *Provider) metadataHandle(w http.ResponseWriter, r *http.Request) {
	err := writeXML(w, p.Metadata)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to respond with metadata").Error(), http.StatusInternalServerError)
		return
	}
}

func writeXML(w http.ResponseWriter, body interface{}) error {
	_, err := w.Write([]byte(xml.Header))
	if err != nil {
		return err
	}
	encoder := xml.NewEncoder(w)

	err = encoder.Encode(body)
	if err != nil {
		return err
	}
	err = encoder.Flush()
	return err
}

func (p *IdentityProviderConfig) getMetadata(
	entityID string,
	idpCertData []byte,
) (*md.IDPSSODescriptorType, *md.AttributeAuthorityDescriptorType) {
	idpKeyDescriptors := []md.KeyDescriptorType{
		{
			Use: md.KeyTypesSigning,
			KeyInfo: xml_dsig.KeyInfoType{
				KeyName: []string{entityID + " IDP " + string(md.KeyTypesSigning)},
				X509Data: []xml_dsig.X509DataType{{
					X509Certificate: base64.StdEncoding.EncodeToString(idpCertData),
				}},
			},
		},
		{
			Use: md.KeyTypesEncryption,
			KeyInfo: xml_dsig.KeyInfoType{
				KeyName: []string{entityID + " IDP " + string(md.KeyTypesEncryption)},
				X509Data: []xml_dsig.X509DataType{{
					X509Certificate: base64.StdEncoding.EncodeToString(idpCertData),
				}},
			},
			EncryptionMethod: []xenc.EncryptionMethodType{{
				Algorithm: p.EncryptionAlgorithm,
			}},
		},
	}

	return &md.IDPSSODescriptorType{
			XMLName:                    xml.Name{},
			WantAuthnRequestsSigned:    p.WantAuthRequestsSigned,
			Id:                         NewID(),
			ValidUntil:                 p.ValidUntil,
			CacheDuration:              p.CacheDuration,
			ProtocolSupportEnumeration: "urn:oasis:names:tc:SAML:2.0:protocol",
			ErrorURL:                   p.ErrorURL,
			SingleSignOnService: []md.EndpointType{{
				Binding:  "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect",
				Location: p.SingleSignOnService,
			}},
			//TODO definition for more profiles
			AttributeProfile: []string{
				"urn:oasis:names:tc:SAML:2.0:profiles:attribute:basic",
			},
			//TODO definition for all provided attributes
			Attribute: nil,
			ArtifactResolutionService: []md.IndexedEndpointType{{
				Index:     "0",
				IsDefault: "true",
				Binding:   "urn:oasis:names:tc:SAML:2.0:bindings:SOAP",
				Location:  p.ArtifactResulationService,
			}},
			SingleLogoutService: []md.EndpointType{
				{
					Binding:  "urn:oasis:names:tc:SAML:2.0:bindings:SOAP",
					Location: p.SLOArtifactResulationService,
				},
				{
					Binding:  "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect",
					Location: p.SingleLogoutService,
				},
			},
			NameIDFormat:  []string{p.NameIDFormat},
			Signature:     nil,
			KeyDescriptor: idpKeyDescriptors,

			Organization:  nil,
			ContactPerson: nil,
			/*
				NameIDMappingService: nil,
				AssertionIDRequestService: nil,
				ManageNameIDService: nil,
			*/
			InnerXml: "",
		},
		&md.AttributeAuthorityDescriptorType{
			XMLName:                    xml.Name{},
			Id:                         NewID(),
			ValidUntil:                 p.ValidUntil,
			CacheDuration:              p.CacheDuration,
			ProtocolSupportEnumeration: "urn:oasis:names:tc:SAML:2.0:protocol",
			ErrorURL:                   p.ErrorURL,
			AttributeService: []md.EndpointType{{
				Binding:  "urn:oasis:names:tc:SAML:2.0:bindings:SOAP",
				Location: p.AttributeService,
			}},
			NameIDFormat: []string{p.NameIDFormat},
			//TODO definition for more profiles
			AttributeProfile: []string{
				"urn:oasis:names:tc:SAML:2.0:profiles:attribute:basic",
			},
			//TODO definition for all provided attributes
			Attribute:     nil,
			Signature:     nil,
			KeyDescriptor: idpKeyDescriptors,

			Organization:  nil,
			ContactPerson: nil,

			/*
				AssertionIDRequestService: nil,
			*/
			InnerXml: "",
		}
}

func (p *ProviderConfig) getMetadata(
	idp *IdentityProvider,
) *md.EntityDescriptor {

	entity := &md.EntityDescriptor{
		XMLName:       xml.Name{Local: "md"},
		EntityID:      md.EntityIDType(p.EntityID),
		Id:            NewID(),
		Signature:     nil,
		Organization:  nil,
		ContactPerson: nil,
		/*
			AuthnAuthorityDescriptor:     nil,
			PDPDescriptor:         nil,
			AffiliationDescriptor: nil,
		*/
	}

	if p.IDP != nil {
		entity.IDPSSODescriptor = idp.Metadata
		entity.AttributeAuthorityDescriptor = idp.AAMetadata
	}

	if p.Organisation != nil {
		org := &md.OrganizationType{
			XMLName:    xml.Name{},
			Extensions: nil,
			OrganizationName: []md.LocalizedNameType{
				{Text: p.Organisation.Name},
			},
			OrganizationDisplayName: []md.LocalizedNameType{
				{Text: p.Organisation.DisplayName},
			},
			OrganizationURL: []md.LocalizedURIType{
				{Text: p.Organisation.URL},
			},
			InnerXml: "",
		}
		entity.AttributeAuthorityDescriptor.Organization = org
		entity.IDPSSODescriptor.Organization = org
	}

	if p.ContactPerson != nil {
		contactPerson := []md.ContactType{
			{
				XMLName:         xml.Name{},
				ContactType:     p.ContactPerson.ContactType,
				Company:         p.ContactPerson.Company,
				GivenName:       p.ContactPerson.GivenName,
				SurName:         p.ContactPerson.SurName,
				EmailAddress:    []string{p.ContactPerson.EmailAddress},
				TelephoneNumber: []string{p.ContactPerson.TelephoneNumber},
				InnerXml:        "",
			},
		}
		entity.AttributeAuthorityDescriptor.ContactPerson = contactPerson
		entity.IDPSSODescriptor.ContactPerson = contactPerson
	}

	return entity
}

func (p *Provider) GetMetadata() (*md.EntityDescriptor, error) {
	metadata := *p.Metadata
	idpSig, err := createSignatureM(p.Signer, metadata)
	if err != nil {
		return nil, err
	}
	metadata.Signature = idpSig
	return &metadata, nil
}
