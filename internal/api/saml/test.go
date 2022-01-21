package saml

import (
	"github.com/caos/zitadel/internal/api/saml/xml/metadata/md"
)

func GetTestConfig(basePath string) *ProviderConfig {
	baseURL := "http://localhost:50002/saml"

	conf := &ProviderConfig{
		BaseURL:            baseURL,
		SignatureAlgorithm: "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256",
		DigestAlgorithm:    "http://www.w3.org/2000/09/xmldsig#sha1",
		Organisation: &Organisation{
			Name:        "caos AG",
			DisplayName: "caos AG",
			URL:         "https://caos.ch",
		},
		ContactPerson: &ContactPerson{
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
		IDP: &IdentityProviderConfig{
			ValidUntil:                   "2021-01-01T00:00:00",
			CacheDuration:                "PT30S",
			ErrorURL:                     "https://caos.ch",
			SignatureAlgorithm:           "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256",
			DigestAlgorithm:              "http://www.w3.org/2000/09/xmldsig#sha1",
			EncryptionAlgorithm:          "",
			NameIDFormat:                 "urn:oasis:names:tc:SAML:2.0:nameid-format:persistent",
			WantAuthRequestsSigned:       "false",
			LoginService:                 "/login",
			SingleSignOnService:          "/SSO",
			SingleLogoutService:          "/SLO",
			ArtifactResulationService:    "/artifact",
			SLOArtifactResulationService: "/SLOartifact",
			NameIDMappingService:         "/namid",
			AttributeService:             "/attribute",
		},
	}
	return conf
}

func AddTestSP(p *Provider) error {
	conf := &ServiceProviderConfig{
		URL: "http://service.example.org/simplesaml/module.php/saml/sp/metadata.php/default-sp",
	}

	return p.IdentityProvider.AddServiceProvider(conf)
}
