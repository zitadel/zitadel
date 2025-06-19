package convert

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/query"
	app "github.com/zitadel/zitadel/pkg/grpc/app/v2beta"
)

func samlMetadataGen(entityID string) []byte {
	str := fmt.Sprintf(`<?xml version="1.0"?>
<md:EntityDescriptor xmlns:md="urn:oasis:names:tc:SAML:2.0:metadata"
                     validUntil="2022-08-26T14:08:16Z"
                     cacheDuration="PT604800S"
                     entityID="%s">
    <md:SPSSODescriptor AuthnRequestsSigned="false" WantAssertionsSigned="false" protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol">
        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>
        <md:AssertionConsumerService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST"
                                     Location="https://test.com/saml/acs"
                                     index="1" />
        
    </md:SPSSODescriptor>
</md:EntityDescriptor>
`,
		entityID)

	return []byte(str)
}

func TestCreateSAMLAppRequestToDomain(t *testing.T) {
	t.Parallel()

	genMetaForValidRequest := samlMetadataGen(gofakeit.URL())

	tt := []struct {
		testName  string
		appName   string
		projectID string
		req       *app.CreateSAMLApplicationRequest

		expectedResponse *domain.SAMLApp
		expectedError    error
	}{
		{
			testName:  "login version error",
			appName:   "test-app",
			projectID: "proj-1",
			req: &app.CreateSAMLApplicationRequest{
				Metadata: &app.CreateSAMLApplicationRequest_MetadataXml{
					MetadataXml: samlMetadataGen(gofakeit.URL()),
				},
				LoginVersion: &app.LoginVersion{
					Version: &app.LoginVersion_LoginV2{
						LoginV2: &app.LoginV2{BaseUri: gu.Ptr("%+o")},
					},
				},
			},
			expectedError: &url.Error{
				URL: "%+o",
				Op:  "parse",
				Err: url.EscapeError("%+o"),
			},
		},
		{
			testName:  "valid request",
			appName:   "test-app",
			projectID: "proj-1",
			req: &app.CreateSAMLApplicationRequest{
				Metadata: &app.CreateSAMLApplicationRequest_MetadataXml{
					MetadataXml: genMetaForValidRequest,
				},
				LoginVersion: nil,
			},

			expectedResponse: &domain.SAMLApp{
				ObjectRoot:   models.ObjectRoot{AggregateID: "proj-1"},
				AppName:      "test-app",
				Metadata:     genMetaForValidRequest,
				LoginVersion: 0,
				LoginBaseURI: "",
				State:        0,
			},
		},
		{
			testName:  "nil request",
			appName:   "test-app",
			projectID: "proj-1",
			req:       nil,

			expectedResponse: &domain.SAMLApp{
				AppName:    "test-app",
				ObjectRoot: models.ObjectRoot{AggregateID: "proj-1"},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// When
			res, err := CreateSAMLAppRequestToDomain(tc.appName, tc.projectID, tc.req)

			// Then
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedResponse, res)
		})
	}
}
func TestPatchSAMLAppConfigRequestToDomain(t *testing.T) {
	t.Parallel()

	genMetaForValidRequest := samlMetadataGen(gofakeit.URL())

	tt := []struct {
		testName  string
		appID     string
		projectID string
		req       *app.PatchSAMLApplicationConfigurationRequest

		expectedResponse *domain.SAMLApp
		expectedError    error
	}{
		{
			testName:  "login version error",
			appID:     "app-1",
			projectID: "proj-1",
			req: &app.PatchSAMLApplicationConfigurationRequest{
				Metadata: &app.PatchSAMLApplicationConfigurationRequest_MetadataXml{
					MetadataXml: samlMetadataGen(gofakeit.URL()),
				},
				LoginVersion: &app.LoginVersion{
					Version: &app.LoginVersion_LoginV2{
						LoginV2: &app.LoginV2{BaseUri: gu.Ptr("%+o")},
					},
				},
			},
			expectedError: &url.Error{
				URL: "%+o",
				Op:  "parse",
				Err: url.EscapeError("%+o"),
			},
		},
		{
			testName:  "valid request",
			appID:     "app-1",
			projectID: "proj-1",
			req: &app.PatchSAMLApplicationConfigurationRequest{
				Metadata: &app.PatchSAMLApplicationConfigurationRequest_MetadataXml{
					MetadataXml: genMetaForValidRequest,
				},
				LoginVersion: nil,
			},
			expectedResponse: &domain.SAMLApp{
				ObjectRoot:   models.ObjectRoot{AggregateID: "proj-1"},
				AppID:        "app-1",
				Metadata:     genMetaForValidRequest,
				LoginVersion: 0,
				LoginBaseURI: "",
			},
		},
		{
			testName:  "nil request",
			appID:     "app-1",
			projectID: "proj-1",
			req:       nil,
			expectedResponse: &domain.SAMLApp{
				ObjectRoot: models.ObjectRoot{AggregateID: "proj-1"},
				AppID:      "app-1",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// When
			res, err := PatchSAMLAppConfigRequestToDomain(tc.appID, tc.projectID, tc.req)

			// Then
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedResponse, res)
		})
	}
}

func TestAppSAMLConfigToPb(t *testing.T) {
	t.Parallel()

	metadata := samlMetadataGen(gofakeit.URL())

	tt := []struct {
		name         string
		inputSAMLApp *query.SAMLApp

		expectedPbApp app.ApplicationConfig
	}{
		{
			name: "valid conversion",
			inputSAMLApp: &query.SAMLApp{
				Metadata:     metadata,
				LoginVersion: domain.LoginVersion2,
				LoginBaseURI: gu.Ptr("https://example.com"),
			},
			expectedPbApp: &app.Application_SamlConfig{
				SamlConfig: &app.SAMLConfig{
					Metadata: &app.SAMLConfig_MetadataXml{
						MetadataXml: metadata,
					},
					LoginVersion: &app.LoginVersion{
						Version: &app.LoginVersion_LoginV2{
							LoginV2: &app.LoginV2{BaseUri: gu.Ptr("https://example.com")},
						},
					},
				},
			},
		},
		{
			name:         "nil saml app",
			inputSAMLApp: nil,
			expectedPbApp: &app.Application_SamlConfig{
				SamlConfig: &app.SAMLConfig{
					Metadata:     &app.SAMLConfig_MetadataXml{},
					LoginVersion: &app.LoginVersion{},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// When
			got := appSAMLConfigToPb(tc.inputSAMLApp)

			// Then
			assert.Equal(t, tc.expectedPbApp, got)
		})
	}
}
