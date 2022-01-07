package main

import (
	"encoding/xml"
	"fmt"
	"github.com/caos/zitadel/internal/api/saml"
	"github.com/caos/zitadel/internal/api/saml/xml/metadata/md"
	"io/ioutil"
	"os"
)

func main() {
	entityID := "https://saml.caos.ch"

	conf := &saml.ProviderConfig{
		EntityID: entityID,
		MetadataCertificate: &saml.Certificate{
			Path:           "idp.pem",
			PrivateKeyPath: "idp-key.pem",
			CaPath:         "ca.pem",
		},
		SignatureAlgorithm: "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256",
		DigestAlgorithm:    "http://www.w3.org/2000/09/xmldsig#sha1",
		Organisation: &saml.Organisation{
			Name:        "caos AG",
			DisplayName: "caos AG",
			URL:         "https://caos.ch",
		},
		ContactPerson: &saml.ContactPerson{
			ContactType:     md.ContactTypeTypeTechnical,
			Company:         "caos AG",
			GivenName:       "Stefan",
			SurName:         "Benz",
			EmailAddress:    "stefan@caos.ch",
			TelephoneNumber: "+41",
		},
		ValidUntil:    "2021-01-01T00:00:00",
		CacheDuration: "PT30S",
		ErrorURL:      "https://caos.ch",
		IDP: &saml.IdentityProviderConfig{
			ValidUntil:    "2021-01-01T00:00:00",
			CacheDuration: "PT30S",
			ErrorURL:      "https://caos.ch",
			Certificate: &saml.Certificate{
				Path:           "idp.pem",
				PrivateKeyPath: "idp-key.pem",
				CaPath:         "ca.pem",
			},
			SignatureAlgorithm:           "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256",
			DigestAlgorithm:              "http://www.w3.org/2000/09/xmldsig#sha1",
			EncryptionAlgorithm:          "",
			NameIDFormat:                 "urn:oasis:names:tc:SAML:2.0:nameid-format:persistent",
			WantAuthRequestsSigned:       "true",
			LoginService:                 "/saml/login",
			SingleSignOnService:          "/saml/SSO",
			SingleLogoutService:          "/saml/SLO",
			ArtifactResulationService:    "/saml/artifact",
			SLOArtifactResulationService: "/saml/SLOartifact",
			NameIDMappingService:         "/saml/namid",
			AttributeService:             "/saml/attribute",
		},
	}

	provider, err := saml.NewProvider(conf)
	if err != nil {
		fmt.Println(err.Error())
	}

	meta, err := provider.GetMetadata()
	if err != nil {
		fmt.Println(err.Error())
	}

	data, err := xml.Marshal(meta)
	if err != nil {
		fmt.Println(err.Error())
	}

	ioutil.WriteFile("test.xml", data, os.ModePerm)
}
