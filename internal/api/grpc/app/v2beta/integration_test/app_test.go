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

	app "github.com/zitadel/zitadel/pkg/grpc/app/v2beta"
	org "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
)

func TestCreateApplication(t *testing.T) {
	t.Parallel()

	notExistingProjectID := gofakeit.UUID()

	tt := []struct {
		testName        string
		creationRequest *app.CreateApplicationRequest

		expectedResponseType string
		expectedErrorType    codes.Code
	}{
		{
			testName: "when project for API app creation is not found should return failed precondition error",
			creationRequest: &app.CreateApplicationRequest{
				ProjectId: notExistingProjectID,
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
				ProjectId: Project.GetId(),
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
				ProjectId: notExistingProjectID,
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
				ProjectId: Project.GetId(),
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
				ProjectId: notExistingProjectID,
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
				ProjectId: Project.GetId(),
				Name:      gofakeit.AppName(),
				CreationRequestType: &app.CreateApplicationRequest_SamlRequest{
					SamlRequest: &app.CreateSAMLApplicationRequest{
						Metadata: &app.CreateSAMLApplicationRequest_MetadataXml{
							MetadataXml: samlMetadataGen(gofakeit.URL()),
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
			t.Parallel()
			res, err := instance.Client.AppV2Beta.CreateApplication(IAMOwnerCtx, tc.creationRequest)

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
	t.Parallel()

	orgNotInCtx := instance.CreateOrganization(IAMOwnerCtx, gofakeit.Name(), gofakeit.Email())
	pNotInCtx := instance.CreateProject(IAMOwnerCtx, t, orgNotInCtx.GetOrganizationId(), gofakeit.AppName(), false, false)

	baseURI := "http://example.com"

	t.Cleanup(func() {
		instance.Client.OrgV2beta.DeleteOrganization(IAMOwnerCtx, &org.DeleteOrganizationRequest{
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

	samlMetas := samlMetadataGen(gofakeit.URL())
	reqForSAMLAppCreation := &app.CreateApplicationRequest_SamlRequest{
		SamlRequest: &app.CreateSAMLApplicationRequest{
			Metadata: &app.CreateSAMLApplicationRequest_MetadataXml{
				MetadataXml: samlMetas,
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

	appForNameChange, appNameChangeErr := instance.Client.AppV2Beta.CreateApplication(IAMOwnerCtx, &app.CreateApplicationRequest{
		ProjectId:           Project.GetId(),
		Name:                gofakeit.AppName(),
		CreationRequestType: reqForAppNameCreation,
	})
	require.Nil(t, appNameChangeErr)

	appForAPIConfigChange, appAPIConfigChangeErr := instance.Client.AppV2Beta.CreateApplication(IAMOwnerCtx, &app.CreateApplicationRequest{
		ProjectId:           Project.GetId(),
		Name:                gofakeit.AppName(),
		CreationRequestType: reqForAPIAppCreation,
	})
	require.Nil(t, appAPIConfigChangeErr)

	appForOIDCConfigChange, appOIDCConfigChangeErr := instance.Client.AppV2Beta.CreateApplication(IAMOwnerCtx, &app.CreateApplicationRequest{
		ProjectId:           Project.GetId(),
		Name:                gofakeit.AppName(),
		CreationRequestType: reqForOIDCAppCreation,
	})
	require.Nil(t, appOIDCConfigChangeErr)

	appForSAMLConfigChange, appSAMLConfigChangeErr := instance.Client.AppV2Beta.CreateApplication(IAMOwnerCtx, &app.CreateApplicationRequest{
		ProjectId:           Project.GetId(),
		Name:                gofakeit.AppName(),
		CreationRequestType: reqForSAMLAppCreation,
	})
	require.Nil(t, appSAMLConfigChangeErr)

	tt := []struct {
		testName     string
		patchRequest *app.UpdateApplicationRequest

		expectedErrorType codes.Code
	}{
		{
			testName: "when app for app name change request is not found should return not found error",
			patchRequest: &app.UpdateApplicationRequest{
				ProjectId: pNotInCtx.GetId(),
				Id:        appForNameChange.GetAppId(),
				Name:      "New name",
			},
			expectedErrorType: codes.NotFound,
		},
		{
			testName: "when request for app name change is valid should return updated timestamp",
			patchRequest: &app.UpdateApplicationRequest{
				ProjectId: Project.GetId(),
				Id:        appForNameChange.GetAppId(),

				Name: "New name",
			},
		},

		{
			testName: "when app for API config change request is not found should return not found error",
			patchRequest: &app.UpdateApplicationRequest{
				ProjectId: pNotInCtx.GetId(),
				Id:        appForAPIConfigChange.GetAppId(),
				UpdateRequestType: &app.UpdateApplicationRequest_ApiConfigurationRequest{
					ApiConfigurationRequest: &app.UpdateAPIApplicationConfigurationRequest{
						AuthMethodType: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT,
					},
				},
			},
			expectedErrorType: codes.NotFound,
		},
		{
			testName: "when request for API config change is valid should return updated timestamp",
			patchRequest: &app.UpdateApplicationRequest{
				ProjectId: Project.GetId(),
				Id:        appForAPIConfigChange.GetAppId(),
				UpdateRequestType: &app.UpdateApplicationRequest_ApiConfigurationRequest{
					ApiConfigurationRequest: &app.UpdateAPIApplicationConfigurationRequest{
						AuthMethodType: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_BASIC,
					},
				},
			},
		},

		{
			testName: "when app for OIDC config change request is not found should return not found error",
			patchRequest: &app.UpdateApplicationRequest{
				ProjectId: pNotInCtx.GetId(),
				Id:        appForOIDCConfigChange.GetAppId(),
				UpdateRequestType: &app.UpdateApplicationRequest_OidcConfigurationRequest{
					OidcConfigurationRequest: &app.UpdateOIDCApplicationConfigurationRequest{
						PostLogoutRedirectUris: []string{"http://example.com/home2"},
					},
				},
			},
			expectedErrorType: codes.NotFound,
		},
		{
			testName: "when request for OIDC config change is valid should return updated timestamp",
			patchRequest: &app.UpdateApplicationRequest{
				ProjectId: Project.GetId(),
				Id:        appForOIDCConfigChange.GetAppId(),
				UpdateRequestType: &app.UpdateApplicationRequest_OidcConfigurationRequest{
					OidcConfigurationRequest: &app.UpdateOIDCApplicationConfigurationRequest{
						PostLogoutRedirectUris: []string{"http://example.com/home2"},
					},
				},
			},
		},

		{
			testName: "when app for SAML config change request is not found should return not found error",
			patchRequest: &app.UpdateApplicationRequest{
				ProjectId: pNotInCtx.GetId(),
				Id:        appForSAMLConfigChange.GetAppId(),
				UpdateRequestType: &app.UpdateApplicationRequest_SamlConfigurationRequest{
					SamlConfigurationRequest: &app.UpdateSAMLApplicationConfigurationRequest{
						Metadata: &app.UpdateSAMLApplicationConfigurationRequest_MetadataXml{
							MetadataXml: samlMetas,
						},
						LoginVersion: &app.LoginVersion{Version: &app.LoginVersion_LoginV1{LoginV1: &app.LoginV1{}}},
					},
				},
			},
			expectedErrorType: codes.NotFound,
		},
		{
			testName: "when request for SAML config change is valid should return updated timestamp",
			patchRequest: &app.UpdateApplicationRequest{
				ProjectId: Project.GetId(),
				Id:        appForSAMLConfigChange.GetAppId(),
				UpdateRequestType: &app.UpdateApplicationRequest_SamlConfigurationRequest{
					SamlConfigurationRequest: &app.UpdateSAMLApplicationConfigurationRequest{
						Metadata: &app.UpdateSAMLApplicationConfigurationRequest_MetadataXml{
							MetadataXml: samlMetas,
						},
						LoginVersion: &app.LoginVersion{Version: &app.LoginVersion_LoginV1{LoginV1: &app.LoginV1{}}},
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			res, err := instance.Client.AppV2Beta.UpdateApplication(IAMOwnerCtx, tc.patchRequest)

			require.Equal(t, tc.expectedErrorType, status.Code(err))
			if tc.expectedErrorType == codes.OK {
				assert.NotZero(t, res.GetChangeDate())
			}
		})
	}
}

func TestDeleteApplication(t *testing.T) {
	t.Parallel()

	reqForAppNameCreation := &app.CreateApplicationRequest_ApiRequest{
		ApiRequest: &app.CreateAPIApplicationRequest{AuthMethodType: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT},
	}

	appToDelete, appNameChangeErr := instance.Client.AppV2Beta.CreateApplication(IAMOwnerCtx, &app.CreateApplicationRequest{
		ProjectId:           Project.GetId(),
		Name:                gofakeit.AppName(),
		CreationRequestType: reqForAppNameCreation,
	})
	require.Nil(t, appNameChangeErr)

	tt := []struct {
		testName      string
		deleteRequest *app.DeleteApplicationRequest

		expectedErrorType codes.Code
	}{
		{
			testName: "when app to delete is not found should return not found error",
			deleteRequest: &app.DeleteApplicationRequest{
				ProjectId: Project.GetId(),
				Id:        gofakeit.Sentence(2),
			},
			expectedErrorType: codes.NotFound,
		},
		{
			testName: "when app to delete is found should return deletion time",
			deleteRequest: &app.DeleteApplicationRequest{
				ProjectId: Project.GetId(),
				Id:        appToDelete.GetAppId(),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// When
			res, err := instance.Client.AppV2Beta.DeleteApplication(IAMOwnerCtx, tc.deleteRequest)

			// Then
			require.Equal(t, tc.expectedErrorType, status.Code(err))
			if tc.expectedErrorType == codes.OK {
				assert.NotZero(t, res.GetDeletionDate())
			}
		})
	}
}

func TestDeactivateApplication(t *testing.T) {
	t.Parallel()

	reqForAppNameCreation := &app.CreateApplicationRequest_ApiRequest{
		ApiRequest: &app.CreateAPIApplicationRequest{AuthMethodType: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT},
	}

	appToDeactivate, appCreateErr := instance.Client.AppV2Beta.CreateApplication(IAMOwnerCtx, &app.CreateApplicationRequest{
		ProjectId:           Project.GetId(),
		Name:                gofakeit.AppName(),
		CreationRequestType: reqForAppNameCreation,
	})
	require.NoError(t, appCreateErr)

	tt := []struct {
		testName      string
		deleteRequest *app.DeactivateApplicationRequest

		expectedErrorType codes.Code
	}{
		{
			testName: "when app to deactivate is not found should return not found error",
			deleteRequest: &app.DeactivateApplicationRequest{
				ProjectId: Project.GetId(),
				Id:        gofakeit.Sentence(2),
			},
			expectedErrorType: codes.NotFound,
		},
		{
			testName: "when app to deactivate is found should return deactivation time",
			deleteRequest: &app.DeactivateApplicationRequest{
				ProjectId: Project.GetId(),
				Id:        appToDeactivate.GetAppId(),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// When
			res, err := instance.Client.AppV2Beta.DeactivateApplication(IAMOwnerCtx, tc.deleteRequest)

			// Then
			require.Equal(t, tc.expectedErrorType, status.Code(err))
			if tc.expectedErrorType == codes.OK {
				assert.NotZero(t, res.GetDeactivationDate())
			}
		})
	}
}

func TestReactivateApplication(t *testing.T) {
	t.Parallel()

	reqForAppNameCreation := &app.CreateApplicationRequest_ApiRequest{
		ApiRequest: &app.CreateAPIApplicationRequest{AuthMethodType: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT},
	}

	appToReactivate, appCreateErr := instance.Client.AppV2Beta.CreateApplication(IAMOwnerCtx, &app.CreateApplicationRequest{
		ProjectId:           Project.GetId(),
		Name:                gofakeit.AppName(),
		CreationRequestType: reqForAppNameCreation,
	})
	require.Nil(t, appCreateErr)

	_, appDeactivateErr := instance.Client.AppV2Beta.DeactivateApplication(IAMOwnerCtx, &app.DeactivateApplicationRequest{
		ProjectId: Project.GetId(),
		Id:        appToReactivate.GetAppId(),
	})
	require.Nil(t, appDeactivateErr)

	tt := []struct {
		testName          string
		reactivateRequest *app.ReactivateApplicationRequest

		expectedErrorType codes.Code
	}{
		{
			testName: "when app to reactivate is not found should return not found error",
			reactivateRequest: &app.ReactivateApplicationRequest{
				ProjectId: Project.GetId(),
				Id:        gofakeit.Sentence(2),
			},
			expectedErrorType: codes.NotFound,
		},
		{
			testName: "when app to reactivate is found should return deactivation time",
			reactivateRequest: &app.ReactivateApplicationRequest{
				ProjectId: Project.GetId(),
				Id:        appToReactivate.GetAppId(),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// When
			res, err := instance.Client.AppV2Beta.ReactivateApplication(IAMOwnerCtx, tc.reactivateRequest)

			// Then
			require.Equal(t, tc.expectedErrorType, status.Code(err))
			if tc.expectedErrorType == codes.OK {
				assert.NotZero(t, res.GetReactivationDate())
			}
		})
	}
}

func TestRegenerateClientSecret(t *testing.T) {
	t.Parallel()

	reqForApiAppCreation := &app.CreateApplicationRequest_ApiRequest{
		ApiRequest: &app.CreateAPIApplicationRequest{AuthMethodType: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT},
	}

	apiAppToRegen, apiAppCreateErr := instance.Client.AppV2Beta.CreateApplication(IAMOwnerCtx, &app.CreateApplicationRequest{
		ProjectId:           Project.GetId(),
		Name:                gofakeit.AppName(),
		CreationRequestType: reqForApiAppCreation,
	})
	require.Nil(t, apiAppCreateErr)

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

	oidcAppToRegen, oidcAppCreateErr := instance.Client.AppV2Beta.CreateApplication(IAMOwnerCtx, &app.CreateApplicationRequest{
		ProjectId:           Project.GetId(),
		Name:                gofakeit.AppName(),
		CreationRequestType: reqForOIDCAppCreation,
	})
	require.Nil(t, oidcAppCreateErr)

	tt := []struct {
		testName     string
		regenRequest *app.RegenerateClientSecretRequest

		expectedErrorType codes.Code
		oldSecret         string
	}{
		{
			testName: "when app to regen is not expected type should return invalid argument error",
			regenRequest: &app.RegenerateClientSecretRequest{
				ProjectId:     Project.GetId(),
				ApplicationId: gofakeit.Sentence(2),
			},
			expectedErrorType: codes.InvalidArgument,
		},
		{
			testName: "when app to regen is not found should return not found error",
			regenRequest: &app.RegenerateClientSecretRequest{
				ProjectId:     Project.GetId(),
				ApplicationId: gofakeit.Sentence(2),
				AppType:       &app.RegenerateClientSecretRequest_IsApi{},
			},
			expectedErrorType: codes.NotFound,
		},
		{
			testName: "when API app to regen is found should return different secret",
			regenRequest: &app.RegenerateClientSecretRequest{
				ProjectId:     Project.GetId(),
				ApplicationId: apiAppToRegen.GetAppId(),
				AppType:       &app.RegenerateClientSecretRequest_IsApi{},
			},
			oldSecret: apiAppToRegen.GetApiResponse().GetClientSecret(),
		},
		{
			testName: "when OIDC app to regen is found should return different secret",
			regenRequest: &app.RegenerateClientSecretRequest{
				ProjectId:     Project.GetId(),
				ApplicationId: oidcAppToRegen.GetAppId(),
				AppType:       &app.RegenerateClientSecretRequest_IsOidc{},
			},
			oldSecret: oidcAppToRegen.GetOidcResponse().GetClientSecret(),
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// When
			res, err := instance.Client.AppV2Beta.RegenerateClientSecret(IAMOwnerCtx, tc.regenRequest)

			// Then
			require.Equal(t, tc.expectedErrorType, status.Code(err))
			if tc.expectedErrorType == codes.OK {
				assert.NotZero(t, res.GetCreationDate())
				assert.NotEqual(t, tc.oldSecret, res.GetClientSecret())
			}
		})
	}

}
