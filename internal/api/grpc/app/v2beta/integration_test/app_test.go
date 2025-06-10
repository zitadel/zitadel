//go:build integration

package instance_test

import (
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel/internal/integration"
	app "github.com/zitadel/zitadel/pkg/grpc/app/v2beta"
	org "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
	project "github.com/zitadel/zitadel/pkg/grpc/project/v2beta"
)

const (
	samlMetadataString = `<?xml version="1.0"?>
<md:EntityDescriptor xmlns:md="urn:oasis:names:tc:SAML:2.0:metadata"
                     validUntil="2022-08-26T14:08:16Z"
                     cacheDuration="PT604800S"
                     entityID="https://test.com/saml/metadata">
    <md:SPSSODescriptor AuthnRequestsSigned="false" WantAssertionsSigned="false" protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol">
        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>
        <md:AssertionConsumerService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST"
                                     Location="https://test.com/saml/acs"
                                     index="1" />
        
    </md:SPSSODescriptor>
</md:EntityDescriptor>
`
)

var (
	baseURI      = "http://example.com"
	samlMetadata = []byte(samlMetadataString)
)

func TestCreateApplication(t *testing.T) {
	t.Parallel()

	iamOwnerCtx := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)

	orgNotInCtx := instance.CreateOrganization(iamOwnerCtx, gofakeit.Name(), gofakeit.Email())
	p := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.Id, gofakeit.AppName(), false, false)
	pNotInCtx := instance.CreateProject(iamOwnerCtx, t, orgNotInCtx.GetOrganizationId(), gofakeit.AppName(), false, false)

	tt := []struct {
		testName        string
		creationRequest *app.CreateApplicationRequest

		expectedResponseType string
		expectedErrorType    codes.Code
	}{
		{
			testName: "when project for API app creation is not found should return failed precondition error",
			creationRequest: &app.CreateApplicationRequest{
				ProjectId: pNotInCtx.GetId(),
				Name:      "App Name",
				CreationRequestType: &app.CreateApplicationRequest_ApiRequest{
					ApiRequest: &app.CreateAPIApplicationRequest{
						AuthMethodType: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT,
					},
				},
			},
			expectedErrorType: codes.FailedPrecondition,
		},
		{
			testName: "when CreateAPIApp request is valid should create app and return no error",
			creationRequest: &app.CreateApplicationRequest{
				ProjectId: p.GetId(),
				Name:      "App Name",
				CreationRequestType: &app.CreateApplicationRequest_ApiRequest{
					ApiRequest: &app.CreateAPIApplicationRequest{
						AuthMethodType: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT,
					},
				},
			},
			expectedResponseType: fmt.Sprintf("%T", &app.CreateApplicationResponse_ApiResponse{}),
		},
		{
			testName: "when project for OIDC app creation is not found should return failed precondition error",
			creationRequest: &app.CreateApplicationRequest{
				ProjectId: pNotInCtx.GetId(),
				Name:      "App Name",
				CreationRequestType: &app.CreateApplicationRequest_OidcRequest{
					OidcRequest: &app.CreateOIDCApplicationRequest{
						RedirectUris:           []string{"http://example.com"},
						ResponseTypes:          []app.OIDCResponseType{app.OIDCResponseType_OIDC_RESPONSE_TYPE_CODE},
						GrantTypes:             []app.OIDCGrantType{app.OIDCGrantType_OIDC_GRANT_TYPE_AUTHORIZATION_CODE},
						AppType:                app.OIDCAppType_OIDC_APP_TYPE_WEB,
						AuthMethodType:         app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_BASIC,
						PostLogoutRedirectUris: []string{"http://example.com/home"},
						Version:                app.OIDCVersion_OIDC_VERSION_1_0,
						AccessTokenType:        app.OIDCTokenType_OIDC_TOKEN_TYPE_JWT,
						BackChannelLogoutUri:   "http://example.com/logout",
						LoginVersion: &app.LoginVersion{
							Version: &app.LoginVersion_LoginV2{
								LoginV2: &app.LoginV2{
									BaseUri: &baseURI,
								},
							},
						},
					},
				},
			},
			expectedErrorType: codes.FailedPrecondition,
		},
		{
			testName: "when CreateOIDCApp request is valid should create app and return no error",
			creationRequest: &app.CreateApplicationRequest{
				ProjectId: p.GetId(),
				Name:      gofakeit.AppName(),
				CreationRequestType: &app.CreateApplicationRequest_OidcRequest{
					OidcRequest: &app.CreateOIDCApplicationRequest{
						RedirectUris:           []string{"http://example.com"},
						ResponseTypes:          []app.OIDCResponseType{app.OIDCResponseType_OIDC_RESPONSE_TYPE_CODE},
						GrantTypes:             []app.OIDCGrantType{app.OIDCGrantType_OIDC_GRANT_TYPE_AUTHORIZATION_CODE},
						AppType:                app.OIDCAppType_OIDC_APP_TYPE_WEB,
						AuthMethodType:         app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_BASIC,
						PostLogoutRedirectUris: []string{"http://example.com/home"},
						Version:                app.OIDCVersion_OIDC_VERSION_1_0,
						AccessTokenType:        app.OIDCTokenType_OIDC_TOKEN_TYPE_JWT,
						BackChannelLogoutUri:   "http://example.com/logout",
						LoginVersion: &app.LoginVersion{
							Version: &app.LoginVersion_LoginV2{
								LoginV2: &app.LoginV2{
									BaseUri: &baseURI,
								},
							},
						},
					},
				},
			},

			expectedResponseType: fmt.Sprintf("%T", &app.CreateApplicationResponse_OidcResponse{}),
		},
		{
			testName: "when project for SAML app creation is not found should return failed precondition error",
			creationRequest: &app.CreateApplicationRequest{
				ProjectId: pNotInCtx.GetId(),
				Name:      gofakeit.AppName(),
				CreationRequestType: &app.CreateApplicationRequest_SamlRequest{
					SamlRequest: &app.CreateSAMLApplicationRequest{
						Metadata: &app.CreateSAMLApplicationRequest_MetadataUrl{
							MetadataUrl: "http://example.com/metas",
						},
						LoginVersion: &app.LoginVersion{
							Version: &app.LoginVersion_LoginV2{
								LoginV2: &app.LoginV2{
									BaseUri: &baseURI,
								},
							},
						},
					},
				},
			},
			expectedErrorType: codes.FailedPrecondition,
		},
		{
			testName: "when CreateSAMLApp request is valid should create app and return no error",
			creationRequest: &app.CreateApplicationRequest{
				ProjectId: p.GetId(),
				Name:      gofakeit.AppName(),
				CreationRequestType: &app.CreateApplicationRequest_SamlRequest{
					SamlRequest: &app.CreateSAMLApplicationRequest{
						Metadata: &app.CreateSAMLApplicationRequest_MetadataXml{
							MetadataXml: samlMetadata,
						},
						LoginVersion: &app.LoginVersion{
							Version: &app.LoginVersion_LoginV2{
								LoginV2: &app.LoginV2{
									BaseUri: &baseURI,
								},
							},
						},
					},
				},
			},
			expectedResponseType: fmt.Sprintf("%T", &app.CreateApplicationResponse_SamlResponse{}),
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			res, err := instance.Client.AppV2Beta.CreateApplication(iamOwnerCtx, tc.creationRequest)

			require.Equal(t, tc.expectedErrorType, status.Code(err))
			if tc.expectedErrorType == codes.OK {
				resType := fmt.Sprintf("%T", res.GetCreationResponseType())
				assert.Equal(t, tc.expectedResponseType, resType)
				assert.NotZero(t, res.GetAppId())
				assert.NotZero(t, res.GetCreationDate())
			}
		})
	}
}

func TestPatchApplication(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)

	orgNotInCtx := instance.CreateOrganization(iamOwnerCtx, gofakeit.Name(), gofakeit.Email())
	p := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.Id, gofakeit.AppName(), false, false)
	pNotInCtx := instance.CreateProject(iamOwnerCtx, t, orgNotInCtx.GetOrganizationId(), gofakeit.AppName(), false, false)

	baseURI := "http://example.com"

	t.Cleanup(func() {
		instance.Client.Projectv2Beta.DeleteProject(iamOwnerCtx, &project.DeleteProjectRequest{
			Id: p.GetId(),
		})
		instance.Client.OrgV2beta.DeleteOrganization(iamOwnerCtx, &org.DeleteOrganizationRequest{
			Id: orgNotInCtx.GetOrganizationId(),
		})
	})

	reqForAppNameCreation := &app.CreateApplicationRequest_ApiRequest{
		ApiRequest: &app.CreateAPIApplicationRequest{AuthMethodType: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT},
	}
	reqForAPIAppCreation := reqForAppNameCreation

	reqForOIDCAppCreation := &app.CreateApplicationRequest_OidcRequest{
		OidcRequest: &app.CreateOIDCApplicationRequest{
			RedirectUris:           []string{"http://example.com"},
			ResponseTypes:          []app.OIDCResponseType{app.OIDCResponseType_OIDC_RESPONSE_TYPE_CODE},
			GrantTypes:             []app.OIDCGrantType{app.OIDCGrantType_OIDC_GRANT_TYPE_AUTHORIZATION_CODE},
			AppType:                app.OIDCAppType_OIDC_APP_TYPE_WEB,
			AuthMethodType:         app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_BASIC,
			PostLogoutRedirectUris: []string{"http://example.com/home"},
			Version:                app.OIDCVersion_OIDC_VERSION_1_0,
			AccessTokenType:        app.OIDCTokenType_OIDC_TOKEN_TYPE_JWT,
			BackChannelLogoutUri:   "http://example.com/logout",
			LoginVersion: &app.LoginVersion{
				Version: &app.LoginVersion_LoginV2{
					LoginV2: &app.LoginV2{
						BaseUri: &baseURI,
					},
				},
			},
		},
	}

	reqForSAMLAppCreation := &app.CreateApplicationRequest_SamlRequest{
		SamlRequest: &app.CreateSAMLApplicationRequest{
			Metadata: &app.CreateSAMLApplicationRequest_MetadataXml{
				MetadataXml: samlMetadata,
			},
			LoginVersion: &app.LoginVersion{
				Version: &app.LoginVersion_LoginV2{
					LoginV2: &app.LoginV2{
						BaseUri: &baseURI,
					},
				},
			},
		},
	}

	appForNameChange, appNameChangeErr := instance.Client.AppV2Beta.CreateApplication(iamOwnerCtx, &app.CreateApplicationRequest{
		ProjectId:           p.GetId(),
		Name:                gofakeit.AppName(),
		CreationRequestType: reqForAppNameCreation,
	})
	require.Nil(t, appNameChangeErr)

	appForAPIConfigChange, appAPIConfigChangeErr := instance.Client.AppV2Beta.CreateApplication(iamOwnerCtx, &app.CreateApplicationRequest{
		ProjectId:           p.GetId(),
		Name:                gofakeit.AppName(),
		CreationRequestType: reqForAPIAppCreation,
	})
	require.Nil(t, appAPIConfigChangeErr)

	appForOIDCConfigChange, appOIDCConfigChangeErr := instance.Client.AppV2Beta.CreateApplication(iamOwnerCtx, &app.CreateApplicationRequest{
		ProjectId:           p.GetId(),
		Name:                gofakeit.AppName(),
		CreationRequestType: reqForOIDCAppCreation,
	})
	require.Nil(t, appOIDCConfigChangeErr)

	appForSAMLConfigChange, appSAMLConfigChangeErr := instance.Client.AppV2Beta.CreateApplication(iamOwnerCtx, &app.CreateApplicationRequest{
		ProjectId:           p.GetId(),
		Name:                gofakeit.AppName(),
		CreationRequestType: reqForSAMLAppCreation,
	})
	require.Nil(t, appSAMLConfigChangeErr)

	tt := []struct {
		testName     string
		patchRequest *app.PatchApplicationRequest

		expectedErrorType codes.Code
	}{
		{
			testName: "when app for app name change request is not found should return not found error",
			patchRequest: &app.PatchApplicationRequest{
				ProjectId:     pNotInCtx.GetId(),
				ApplicationId: appForNameChange.GetAppId(),

				PatchRequestType: &app.PatchApplicationRequest_ApplicationNameRequest{
					ApplicationNameRequest: &app.PatchApplicationNameRequest{
						Name: "New name",
					},
				},
			},
			expectedErrorType: codes.NotFound,
		},
		{
			testName: "when request for app name change is valid should return updated timestamp",
			patchRequest: &app.PatchApplicationRequest{
				ProjectId:     p.GetId(),
				ApplicationId: appForNameChange.GetAppId(),

				PatchRequestType: &app.PatchApplicationRequest_ApplicationNameRequest{
					ApplicationNameRequest: &app.PatchApplicationNameRequest{
						Name: "New name",
					},
				},
			},
		},

		{
			testName: "when app for API config change request is not found should return not found error",
			patchRequest: &app.PatchApplicationRequest{
				ProjectId:     pNotInCtx.GetId(),
				ApplicationId: appForAPIConfigChange.GetAppId(),
				PatchRequestType: &app.PatchApplicationRequest_ApiConfigurationRequest{
					ApiConfigurationRequest: &app.PatchAPIApplicationConfigurationRequest{
						AuthMethodType: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT,
					},
				},
			},
			expectedErrorType: codes.NotFound,
		},
		{
			testName: "when request for API config change is valid should return updated timestamp",
			patchRequest: &app.PatchApplicationRequest{
				ProjectId:     p.GetId(),
				ApplicationId: appForAPIConfigChange.GetAppId(),
				PatchRequestType: &app.PatchApplicationRequest_ApiConfigurationRequest{
					ApiConfigurationRequest: &app.PatchAPIApplicationConfigurationRequest{
						AuthMethodType: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_BASIC,
					},
				},
			},
		},

		{
			testName: "when app for OIDC config change request is not found should return not found error",
			patchRequest: &app.PatchApplicationRequest{
				ProjectId:     pNotInCtx.GetId(),
				ApplicationId: appForOIDCConfigChange.GetAppId(),
				PatchRequestType: &app.PatchApplicationRequest_OidcConfigurationRequest{
					OidcConfigurationRequest: &app.PatchOIDCApplicationConfigurationRequest{
						PostLogoutRedirectUris: []string{"http://example.com/home2"},
					},
				},
			},
			expectedErrorType: codes.NotFound,
		},
		{
			testName: "when request for OIDC config change is valid should return updated timestamp",
			patchRequest: &app.PatchApplicationRequest{
				ProjectId:     p.GetId(),
				ApplicationId: appForOIDCConfigChange.GetAppId(),
				PatchRequestType: &app.PatchApplicationRequest_OidcConfigurationRequest{
					OidcConfigurationRequest: &app.PatchOIDCApplicationConfigurationRequest{
						PostLogoutRedirectUris: []string{"http://example.com/home2"},
					},
				},
			},
		},

		{
			testName: "when app for SAML config change request is not found should return not found error",
			patchRequest: &app.PatchApplicationRequest{
				ProjectId:     pNotInCtx.GetId(),
				ApplicationId: appForSAMLConfigChange.GetAppId(),
				PatchRequestType: &app.PatchApplicationRequest_SamlConfigurationRequest{
					SamlConfigurationRequest: &app.PatchSAMLApplicationConfigurationRequest{
						Metadata: &app.PatchSAMLApplicationConfigurationRequest_MetadataXml{
							MetadataXml: samlMetadata,
						},
						LoginVersion: &app.LoginVersion{Version: &app.LoginVersion_LoginV1{LoginV1: &app.LoginV1{}}},
					},
				},
			},
			expectedErrorType: codes.NotFound,
		},
		{
			testName: "when request for SAML config change is valid should return updated timestamp",
			patchRequest: &app.PatchApplicationRequest{
				ProjectId:     p.GetId(),
				ApplicationId: appForSAMLConfigChange.GetAppId(),
				PatchRequestType: &app.PatchApplicationRequest_SamlConfigurationRequest{
					SamlConfigurationRequest: &app.PatchSAMLApplicationConfigurationRequest{
						Metadata: &app.PatchSAMLApplicationConfigurationRequest_MetadataXml{
							MetadataXml: samlMetadata,
						},
						LoginVersion: &app.LoginVersion{Version: &app.LoginVersion_LoginV1{LoginV1: &app.LoginV1{}}},
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			res, err := instance.Client.AppV2Beta.PatchApplication(iamOwnerCtx, tc.patchRequest)

			require.Equal(t, tc.expectedErrorType, status.Code(err))
			if tc.expectedErrorType == codes.OK {
				assert.NotZero(t, res.GetChangeDate())
			}
		})
	}
}
