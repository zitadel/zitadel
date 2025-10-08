//go:build integration

package app_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel/internal/integration"
	app "github.com/zitadel/zitadel/pkg/grpc/app/v2beta"
	org "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
)

func TestCreateApplication(t *testing.T) {
	p := instance.CreateProject(IAMOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)

	t.Parallel()

	notExistingProjectID := integration.ID()

	tt := []struct {
		testName        string
		creationRequest *app.CreateApplicationRequest
		inputCtx        context.Context

		expectedResponseType string
		expectedErrorType    codes.Code
	}{
		{
			testName: "when project for API app creation is not found should return failed precondition error",
			inputCtx: IAMOwnerCtx,
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
			inputCtx: IAMOwnerCtx,
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
			inputCtx: IAMOwnerCtx,
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
			inputCtx: IAMOwnerCtx,
			creationRequest: &app.CreateApplicationRequest{
				ProjectId: p.GetId(),
				Name:      integration.ApplicationName(),
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
			inputCtx: IAMOwnerCtx,
			creationRequest: &app.CreateApplicationRequest{
				ProjectId: notExistingProjectID,
				Name:      integration.ApplicationName(),
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
			inputCtx: IAMOwnerCtx,
			creationRequest: &app.CreateApplicationRequest{
				ProjectId: p.GetId(),
				Name:      integration.ApplicationName(),
				CreationRequestType: &app.CreateApplicationRequest_SamlRequest{
					SamlRequest: &app.CreateSAMLApplicationRequest{
						Metadata: &app.CreateSAMLApplicationRequest_MetadataXml{
							MetadataXml: samlMetadataGen(integration.URL()),
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
			res, err := instance.Client.AppV2Beta.CreateApplication(tc.inputCtx, tc.creationRequest)

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

func TestCreateApplication_WithDifferentPermissions(t *testing.T) {
	p, projectOwnerCtx := getProjectAndProjectContext(t, instance, IAMOwnerCtx)

	t.Parallel()

	tt := []struct {
		testName        string
		creationRequest *app.CreateApplicationRequest
		inputCtx        context.Context

		expectedResponseType string
		expectedErrorType    codes.Code
	}{
		// Login User with no project.app.write
		{
			testName: "when user has no project.app.write permission for API request should return permission error",
			inputCtx: LoginUserCtx,
			creationRequest: &app.CreateApplicationRequest{
				ProjectId: p.GetId(),
				Name:      integration.ApplicationName(),
				CreationRequestType: &app.CreateApplicationRequest_ApiRequest{
					ApiRequest: &app.CreateAPIApplicationRequest{
						AuthMethodType: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT,
					},
				},
			},
			expectedErrorType: codes.PermissionDenied,
		},
		{
			testName: "when user has no project.app.write permission for OIDC request should return permission error",
			inputCtx: LoginUserCtx,
			creationRequest: &app.CreateApplicationRequest{
				ProjectId: p.GetId(),
				Name:      integration.ApplicationName(),
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

			expectedErrorType: codes.PermissionDenied,
		},
		{
			testName: "when user has no project.app.write permission for SAML request should return permission error",
			inputCtx: LoginUserCtx,
			creationRequest: &app.CreateApplicationRequest{
				ProjectId: p.GetId(),
				Name:      integration.ApplicationName(),
				CreationRequestType: &app.CreateApplicationRequest_SamlRequest{
					SamlRequest: &app.CreateSAMLApplicationRequest{
						Metadata: &app.CreateSAMLApplicationRequest_MetadataXml{
							MetadataXml: samlMetadataGen(integration.URL()),
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
			expectedErrorType: codes.PermissionDenied,
		},

		// OrgOwner with project.app.write permission
		{
			testName: "when user is OrgOwner API request should succeed",
			inputCtx: OrgOwnerCtx,
			creationRequest: &app.CreateApplicationRequest{
				ProjectId: p.GetId(),
				Name:      integration.ApplicationName(),
				CreationRequestType: &app.CreateApplicationRequest_ApiRequest{
					ApiRequest: &app.CreateAPIApplicationRequest{
						AuthMethodType: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT,
					},
				},
			},
			expectedResponseType: fmt.Sprintf("%T", &app.CreateApplicationResponse_ApiResponse{}),
		},
		{
			testName: "when user is OrgOwner OIDC request should succeed",
			inputCtx: OrgOwnerCtx,
			creationRequest: &app.CreateApplicationRequest{
				ProjectId: p.GetId(),
				Name:      integration.ApplicationName(),
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
			testName: "when user is OrgOwner SAML request should succeed",
			inputCtx: OrgOwnerCtx,
			creationRequest: &app.CreateApplicationRequest{
				ProjectId: p.GetId(),
				Name:      integration.ApplicationName(),
				CreationRequestType: &app.CreateApplicationRequest_SamlRequest{
					SamlRequest: &app.CreateSAMLApplicationRequest{
						Metadata: &app.CreateSAMLApplicationRequest_MetadataXml{
							MetadataXml: samlMetadataGen(integration.URL()),
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

		// Project owner with project.app.write permission
		{
			testName: "when user is ProjectOwner API request should succeed",
			inputCtx: projectOwnerCtx,
			creationRequest: &app.CreateApplicationRequest{
				ProjectId: p.GetId(),
				Name:      integration.ApplicationName(),
				CreationRequestType: &app.CreateApplicationRequest_ApiRequest{
					ApiRequest: &app.CreateAPIApplicationRequest{
						AuthMethodType: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT,
					},
				},
			},
			expectedResponseType: fmt.Sprintf("%T", &app.CreateApplicationResponse_ApiResponse{}),
		},
		{
			testName: "when user is ProjectOwner OIDC request should succeed",
			inputCtx: projectOwnerCtx,
			creationRequest: &app.CreateApplicationRequest{
				ProjectId: p.GetId(),
				Name:      integration.ApplicationName(),
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
			testName: "when user is ProjectOwner SAML request should succeed",
			inputCtx: projectOwnerCtx,
			creationRequest: &app.CreateApplicationRequest{
				ProjectId: p.GetId(),
				Name:      integration.ApplicationName(),
				CreationRequestType: &app.CreateApplicationRequest_SamlRequest{
					SamlRequest: &app.CreateSAMLApplicationRequest{
						Metadata: &app.CreateSAMLApplicationRequest_MetadataXml{
							MetadataXml: samlMetadataGen(integration.URL()),
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
			res, err := instance.Client.AppV2Beta.CreateApplication(tc.inputCtx, tc.creationRequest)

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

func TestUpdateApplication(t *testing.T) {
	orgNotInCtx := instance.CreateOrganization(IAMOwnerCtx, integration.OrganizationName(), integration.Email())
	pNotInCtx := instance.CreateProject(IAMOwnerCtx, t, orgNotInCtx.GetOrganizationId(), integration.ApplicationName(), false, false)

	p := instance.CreateProject(IAMOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)

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

	samlMetas := samlMetadataGen(integration.URL())
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
		ProjectId:           p.GetId(),
		Name:                integration.ApplicationName(),
		CreationRequestType: reqForAppNameCreation,
	})
	require.Nil(t, appNameChangeErr)

	appForAPIConfigChange, appAPIConfigChangeErr := instance.Client.AppV2Beta.CreateApplication(IAMOwnerCtx, &app.CreateApplicationRequest{
		ProjectId:           p.GetId(),
		Name:                integration.ApplicationName(),
		CreationRequestType: reqForAPIAppCreation,
	})
	require.Nil(t, appAPIConfigChangeErr)

	appForOIDCConfigChange, appOIDCConfigChangeErr := instance.Client.AppV2Beta.CreateApplication(IAMOwnerCtx, &app.CreateApplicationRequest{
		ProjectId:           p.GetId(),
		Name:                integration.ApplicationName(),
		CreationRequestType: reqForOIDCAppCreation,
	})
	require.Nil(t, appOIDCConfigChangeErr)

	appForSAMLConfigChange, appSAMLConfigChangeErr := instance.Client.AppV2Beta.CreateApplication(IAMOwnerCtx, &app.CreateApplicationRequest{
		ProjectId:           p.GetId(),
		Name:                integration.ApplicationName(),
		CreationRequestType: reqForSAMLAppCreation,
	})
	require.Nil(t, appSAMLConfigChangeErr)

	t.Parallel()

	tt := []struct {
		testName      string
		inputCtx      context.Context
		updateRequest *app.UpdateApplicationRequest

		expectedErrorType codes.Code
	}{
		{
			testName: "when app for app name change request is not found should return not found error",
			inputCtx: IAMOwnerCtx,
			updateRequest: &app.UpdateApplicationRequest{
				ProjectId: pNotInCtx.GetId(),
				Id:        appForNameChange.GetAppId(),
				Name:      "New name",
			},
			expectedErrorType: codes.NotFound,
		},
		{
			testName: "when request for app name change is valid should return updated timestamp",
			inputCtx: IAMOwnerCtx,
			updateRequest: &app.UpdateApplicationRequest{
				ProjectId: p.GetId(),
				Id:        appForNameChange.GetAppId(),

				Name: "New name",
			},
		},

		{
			testName: "when app for API config change request is not found should return not found error",
			inputCtx: IAMOwnerCtx,
			updateRequest: &app.UpdateApplicationRequest{
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
			inputCtx: IAMOwnerCtx,
			updateRequest: &app.UpdateApplicationRequest{
				ProjectId: p.GetId(),
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
			inputCtx: IAMOwnerCtx,
			updateRequest: &app.UpdateApplicationRequest{
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
			inputCtx: IAMOwnerCtx,
			updateRequest: &app.UpdateApplicationRequest{
				ProjectId: p.GetId(),
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
			inputCtx: IAMOwnerCtx,
			updateRequest: &app.UpdateApplicationRequest{
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
			inputCtx: IAMOwnerCtx,
			updateRequest: &app.UpdateApplicationRequest{
				ProjectId: p.GetId(),
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
			res, err := instance.Client.AppV2Beta.UpdateApplication(tc.inputCtx, tc.updateRequest)

			require.Equal(t, tc.expectedErrorType, status.Code(err))
			if tc.expectedErrorType == codes.OK {
				assert.NotZero(t, res.GetChangeDate())
			}
		})
	}
}

func TestUpdateApplication_WithDifferentPermissions(t *testing.T) {
	baseURI := "http://example.com"

	p, projectOwnerCtx := getProjectAndProjectContext(t, instance, IAMOwnerCtx)

	reqForAppNameCreation := &app.CreateApplicationRequest_ApiRequest{
		ApiRequest: &app.CreateAPIApplicationRequest{AuthMethodType: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT},
	}

	appForNameChange, appNameChangeErr := instance.Client.AppV2Beta.CreateApplication(IAMOwnerCtx, &app.CreateApplicationRequest{
		ProjectId:           p.GetId(),
		Name:                integration.ApplicationName(),
		CreationRequestType: reqForAppNameCreation,
	})
	require.Nil(t, appNameChangeErr)

	appForAPIConfigChangeForProjectOwner := createAPIApp(t, IAMOwnerCtx, instance, p.GetId())
	appForAPIConfigChangeForOrgOwner := createAPIApp(t, IAMOwnerCtx, instance, p.GetId())
	appForAPIConfigChangeForLoginUser := createAPIApp(t, IAMOwnerCtx, instance, p.GetId())

	appForOIDCConfigChangeForProjectOwner := createOIDCApp(t, baseURI, p.GetId())
	appForOIDCConfigChangeForOrgOwner := createOIDCApp(t, baseURI, p.GetId())
	appForOIDCConfigChangeForLoginUser := createOIDCApp(t, baseURI, p.GetId())

	samlMetasForProjectOwner, appForSAMLConfigChangeForProjectOwner := createSAMLApp(t, baseURI, p.GetId())
	samlMetasForOrgOwner, appForSAMLConfigChangeForOrgOwner := createSAMLApp(t, baseURI, p.GetId())
	samlMetasForLoginUser, appForSAMLConfigChangeForLoginUser := createSAMLApp(t, baseURI, p.GetId())

	t.Parallel()

	tt := []struct {
		testName      string
		inputCtx      context.Context
		updateRequest *app.UpdateApplicationRequest

		expectedErrorType codes.Code
	}{
		// ProjectOwner
		{
			testName: "when user is ProjectOwner app name request should succeed",
			inputCtx: projectOwnerCtx,
			updateRequest: &app.UpdateApplicationRequest{
				ProjectId: p.GetId(),
				Id:        appForNameChange.GetAppId(),

				Name: integration.ApplicationName(),
			},
		},
		{
			testName: "when user is ProjectOwner API app request should succeed",
			inputCtx: projectOwnerCtx,
			updateRequest: &app.UpdateApplicationRequest{
				ProjectId: p.GetId(),
				Id:        appForAPIConfigChangeForProjectOwner.GetAppId(),
				UpdateRequestType: &app.UpdateApplicationRequest_ApiConfigurationRequest{
					ApiConfigurationRequest: &app.UpdateAPIApplicationConfigurationRequest{
						AuthMethodType: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_BASIC,
					},
				},
			},
		},
		{
			testName: "when user is ProjectOwner OIDC app request should succeed",
			inputCtx: projectOwnerCtx,
			updateRequest: &app.UpdateApplicationRequest{
				ProjectId: p.GetId(),
				Id:        appForOIDCConfigChangeForProjectOwner.GetAppId(),
				UpdateRequestType: &app.UpdateApplicationRequest_OidcConfigurationRequest{
					OidcConfigurationRequest: &app.UpdateOIDCApplicationConfigurationRequest{
						PostLogoutRedirectUris: []string{"http://example.com/home2"},
					},
				},
			},
		},
		{
			testName: "when user is ProjectOwner SAML request should succeed",
			inputCtx: projectOwnerCtx,
			updateRequest: &app.UpdateApplicationRequest{
				ProjectId: p.GetId(),
				Id:        appForSAMLConfigChangeForProjectOwner.GetAppId(),
				UpdateRequestType: &app.UpdateApplicationRequest_SamlConfigurationRequest{
					SamlConfigurationRequest: &app.UpdateSAMLApplicationConfigurationRequest{
						Metadata: &app.UpdateSAMLApplicationConfigurationRequest_MetadataXml{
							MetadataXml: samlMetasForProjectOwner,
						},
						LoginVersion: &app.LoginVersion{Version: &app.LoginVersion_LoginV1{LoginV1: &app.LoginV1{}}},
					},
				},
			},
		},

		// OrgOwner context
		{
			testName: "when user is OrgOwner app name request should succeed",
			inputCtx: OrgOwnerCtx,
			updateRequest: &app.UpdateApplicationRequest{
				ProjectId: p.GetId(),
				Id:        appForNameChange.GetAppId(),

				Name: integration.ApplicationName(),
			},
		},
		{
			testName: "when user is OrgOwner API app request should succeed",
			inputCtx: OrgOwnerCtx,
			updateRequest: &app.UpdateApplicationRequest{
				ProjectId: p.GetId(),
				Id:        appForAPIConfigChangeForOrgOwner.GetAppId(),
				UpdateRequestType: &app.UpdateApplicationRequest_ApiConfigurationRequest{
					ApiConfigurationRequest: &app.UpdateAPIApplicationConfigurationRequest{
						AuthMethodType: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_BASIC,
					},
				},
			},
		},
		{
			testName: "when user is OrgOwner OIDC app request should succeed",
			inputCtx: OrgOwnerCtx,
			updateRequest: &app.UpdateApplicationRequest{
				ProjectId: p.GetId(),
				Id:        appForOIDCConfigChangeForOrgOwner.GetAppId(),
				UpdateRequestType: &app.UpdateApplicationRequest_OidcConfigurationRequest{
					OidcConfigurationRequest: &app.UpdateOIDCApplicationConfigurationRequest{
						PostLogoutRedirectUris: []string{"http://example.com/home2"},
					},
				},
			},
		},
		{
			testName: "when user is OrgOwner SAML request should succeed",
			inputCtx: OrgOwnerCtx,
			updateRequest: &app.UpdateApplicationRequest{
				ProjectId: p.GetId(),
				Id:        appForSAMLConfigChangeForOrgOwner.GetAppId(),
				UpdateRequestType: &app.UpdateApplicationRequest_SamlConfigurationRequest{
					SamlConfigurationRequest: &app.UpdateSAMLApplicationConfigurationRequest{
						Metadata: &app.UpdateSAMLApplicationConfigurationRequest_MetadataXml{
							MetadataXml: samlMetasForOrgOwner,
						},
						LoginVersion: &app.LoginVersion{Version: &app.LoginVersion_LoginV1{LoginV1: &app.LoginV1{}}},
					},
				},
			},
		},

		// LoginUser
		{
			testName: "when user has no project.app.write permission for app name change request should return permission error",
			inputCtx: LoginUserCtx,
			updateRequest: &app.UpdateApplicationRequest{
				ProjectId: p.GetId(),
				Id:        appForNameChange.GetAppId(),

				Name: integration.ApplicationName(),
			},
			expectedErrorType: codes.PermissionDenied,
		},
		{
			testName: "when user has no project.app.write permission for API request should return permission error",
			inputCtx: LoginUserCtx,
			updateRequest: &app.UpdateApplicationRequest{
				ProjectId: p.GetId(),
				Id:        appForAPIConfigChangeForLoginUser.GetAppId(),
				UpdateRequestType: &app.UpdateApplicationRequest_ApiConfigurationRequest{
					ApiConfigurationRequest: &app.UpdateAPIApplicationConfigurationRequest{
						AuthMethodType: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_BASIC,
					},
				},
			},
			expectedErrorType: codes.PermissionDenied,
		},
		{
			testName: "when user has no project.app.write permission for OIDC request should return permission error",
			inputCtx: LoginUserCtx,
			updateRequest: &app.UpdateApplicationRequest{
				ProjectId: p.GetId(),
				Id:        appForOIDCConfigChangeForLoginUser.GetAppId(),
				UpdateRequestType: &app.UpdateApplicationRequest_OidcConfigurationRequest{
					OidcConfigurationRequest: &app.UpdateOIDCApplicationConfigurationRequest{
						PostLogoutRedirectUris: []string{"http://example.com/home2"},
					},
				},
			},
			expectedErrorType: codes.PermissionDenied,
		},
		{
			testName: "when user has no project.app.write permission for SAML request should return permission error",
			inputCtx: LoginUserCtx,
			updateRequest: &app.UpdateApplicationRequest{
				ProjectId: p.GetId(),
				Id:        appForSAMLConfigChangeForLoginUser.GetAppId(),
				UpdateRequestType: &app.UpdateApplicationRequest_SamlConfigurationRequest{
					SamlConfigurationRequest: &app.UpdateSAMLApplicationConfigurationRequest{
						Metadata: &app.UpdateSAMLApplicationConfigurationRequest_MetadataXml{
							MetadataXml: samlMetasForLoginUser,
						},
						LoginVersion: &app.LoginVersion{Version: &app.LoginVersion_LoginV1{LoginV1: &app.LoginV1{}}},
					},
				},
			},
			expectedErrorType: codes.PermissionDenied,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			res, err := instance.Client.AppV2Beta.UpdateApplication(tc.inputCtx, tc.updateRequest)

			require.Equal(t, tc.expectedErrorType, status.Code(err))
			if tc.expectedErrorType == codes.OK {
				assert.NotZero(t, res.GetChangeDate())
			}
		})
	}
}

func TestDeleteApplication(t *testing.T) {
	p := instance.CreateProject(IAMOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)

	reqForAppNameCreation := &app.CreateApplicationRequest_ApiRequest{
		ApiRequest: &app.CreateAPIApplicationRequest{AuthMethodType: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT},
	}

	appToDelete, appNameChangeErr := instance.Client.AppV2Beta.CreateApplication(IAMOwnerCtx, &app.CreateApplicationRequest{
		ProjectId:           p.GetId(),
		Name:                integration.ApplicationName(),
		CreationRequestType: reqForAppNameCreation,
	})
	require.Nil(t, appNameChangeErr)

	t.Parallel()
	tt := []struct {
		testName      string
		deleteRequest *app.DeleteApplicationRequest
		inputCtx      context.Context

		expectedErrorType codes.Code
	}{
		{
			testName: "when app to delete is not found should return not found error",
			inputCtx: IAMOwnerCtx,
			deleteRequest: &app.DeleteApplicationRequest{
				ProjectId: p.GetId(),
				Id:        integration.ID(),
			},
			expectedErrorType: codes.NotFound,
		},
		{
			testName: "when app to delete is found should return deletion time",
			inputCtx: IAMOwnerCtx,
			deleteRequest: &app.DeleteApplicationRequest{
				ProjectId: p.GetId(),
				Id:        appToDelete.GetAppId(),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// When
			res, err := instance.Client.AppV2Beta.DeleteApplication(tc.inputCtx, tc.deleteRequest)

			// Then
			require.Equal(t, tc.expectedErrorType, status.Code(err))
			if tc.expectedErrorType == codes.OK {
				assert.NotZero(t, res.GetDeletionDate())
			}
		})
	}
}

func TestDeleteApplication_WithDifferentPermissions(t *testing.T) {
	p, projectOwnerCtx := getProjectAndProjectContext(t, instance, IAMOwnerCtx)

	appToDeleteForLoginUser := createAPIApp(t, IAMOwnerCtx, instance, p.GetId())
	appToDeleteForProjectOwner := createAPIApp(t, IAMOwnerCtx, instance, p.GetId())
	appToDeleteForOrgOwner := createAPIApp(t, IAMOwnerCtx, instance, p.GetId())

	t.Parallel()
	tt := []struct {
		testName      string
		deleteRequest *app.DeleteApplicationRequest
		inputCtx      context.Context

		expectedErrorType codes.Code
	}{
		// Login User
		{
			testName: "when user has no project.app.delete permission for app delete request should return permission error",
			inputCtx: LoginUserCtx,
			deleteRequest: &app.DeleteApplicationRequest{
				ProjectId: p.GetId(),
				Id:        appToDeleteForLoginUser.GetAppId(),
			},
			expectedErrorType: codes.PermissionDenied,
		},

		// Project Owner
		{
			testName: "when user is ProjectOwner delete app request should succeed",
			inputCtx: projectOwnerCtx,
			deleteRequest: &app.DeleteApplicationRequest{
				ProjectId: p.GetId(),
				Id:        appToDeleteForProjectOwner.GetAppId(),
			},
		},

		// Org Owner
		{
			testName: "when user is OrgOwner delete app request should succeed",
			inputCtx: projectOwnerCtx,
			deleteRequest: &app.DeleteApplicationRequest{
				ProjectId: p.GetId(),
				Id:        appToDeleteForOrgOwner.GetAppId(),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// When
			res, err := instance.Client.AppV2Beta.DeleteApplication(tc.inputCtx, tc.deleteRequest)

			// Then
			require.Equal(t, tc.expectedErrorType, status.Code(err))
			if tc.expectedErrorType == codes.OK {
				assert.NotZero(t, res.GetDeletionDate())
			}
		})
	}
}

func TestDeactivateApplication(t *testing.T) {
	p := instance.CreateProject(IAMOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)

	reqForAppNameCreation := &app.CreateApplicationRequest_ApiRequest{
		ApiRequest: &app.CreateAPIApplicationRequest{AuthMethodType: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT},
	}

	appToDeactivate, appCreateErr := instance.Client.AppV2Beta.CreateApplication(IAMOwnerCtx, &app.CreateApplicationRequest{
		ProjectId:           p.GetId(),
		Name:                integration.ApplicationName(),
		CreationRequestType: reqForAppNameCreation,
	})
	require.NoError(t, appCreateErr)

	t.Parallel()

	tt := []struct {
		testName      string
		inputCtx      context.Context
		deleteRequest *app.DeactivateApplicationRequest

		expectedErrorType codes.Code
	}{
		{
			testName: "when app to deactivate is not found should return not found error",
			inputCtx: IAMOwnerCtx,
			deleteRequest: &app.DeactivateApplicationRequest{
				ProjectId: p.GetId(),
				Id:        integration.ID(),
			},
			expectedErrorType: codes.NotFound,
		},
		{
			testName: "when app to deactivate is found should return deactivation time",
			inputCtx: IAMOwnerCtx,
			deleteRequest: &app.DeactivateApplicationRequest{
				ProjectId: p.GetId(),
				Id:        appToDeactivate.GetAppId(),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// When
			res, err := instance.Client.AppV2Beta.DeactivateApplication(tc.inputCtx, tc.deleteRequest)

			// Then
			require.Equal(t, tc.expectedErrorType, status.Code(err))
			if tc.expectedErrorType == codes.OK {
				assert.NotZero(t, res.GetDeactivationDate())
			}
		})
	}
}

func TestDeactivateApplication_WithDifferentPermissions(t *testing.T) {
	p, projectOwnerCtx := getProjectAndProjectContext(t, instance, IAMOwnerCtx)

	appToDeactivateForLoginUser := createAPIApp(t, IAMOwnerCtx, instance, p.GetId())
	appToDeactivateForPrjectOwner := createAPIApp(t, IAMOwnerCtx, instance, p.GetId())
	appToDeactivateForOrgOwner := createAPIApp(t, IAMOwnerCtx, instance, p.GetId())

	t.Parallel()

	tt := []struct {
		testName      string
		inputCtx      context.Context
		deleteRequest *app.DeactivateApplicationRequest

		expectedErrorType codes.Code
	}{
		// Login User
		{
			testName: "when user has no project.app.write permission for app deactivate request should return permission error",
			inputCtx: IAMOwnerCtx,
			deleteRequest: &app.DeactivateApplicationRequest{
				ProjectId: p.GetId(),
				Id:        appToDeactivateForLoginUser.GetAppId(),
			},
		},

		// Project Owner
		{
			testName: "when user is ProjectOwner deactivate app request should succeed",
			inputCtx: projectOwnerCtx,
			deleteRequest: &app.DeactivateApplicationRequest{
				ProjectId: p.GetId(),
				Id:        appToDeactivateForPrjectOwner.GetAppId(),
			},
		},

		// Org Owner
		{
			testName: "when user is OrgOwner deactivate app request should succeed",
			inputCtx: OrgOwnerCtx,
			deleteRequest: &app.DeactivateApplicationRequest{
				ProjectId: p.GetId(),
				Id:        appToDeactivateForOrgOwner.GetAppId(),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// When
			res, err := instance.Client.AppV2Beta.DeactivateApplication(tc.inputCtx, tc.deleteRequest)

			// Then
			require.Equal(t, tc.expectedErrorType, status.Code(err))
			if tc.expectedErrorType == codes.OK {
				assert.NotZero(t, res.GetDeactivationDate())
			}
		})
	}
}

func TestReactivateApplication(t *testing.T) {
	p := instance.CreateProject(IAMOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)

	reqForAppNameCreation := &app.CreateApplicationRequest_ApiRequest{
		ApiRequest: &app.CreateAPIApplicationRequest{AuthMethodType: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT},
	}

	appToReactivate, appCreateErr := instance.Client.AppV2Beta.CreateApplication(IAMOwnerCtx, &app.CreateApplicationRequest{
		ProjectId:           p.GetId(),
		Name:                integration.ApplicationName(),
		CreationRequestType: reqForAppNameCreation,
	})
	require.Nil(t, appCreateErr)

	_, appDeactivateErr := instance.Client.AppV2Beta.DeactivateApplication(IAMOwnerCtx, &app.DeactivateApplicationRequest{
		ProjectId: p.GetId(),
		Id:        appToReactivate.GetAppId(),
	})
	require.Nil(t, appDeactivateErr)

	t.Parallel()

	tt := []struct {
		testName          string
		inputCtx          context.Context
		reactivateRequest *app.ReactivateApplicationRequest

		expectedErrorType codes.Code
	}{
		{
			testName: "when app to reactivate is not found should return not found error",
			inputCtx: IAMOwnerCtx,
			reactivateRequest: &app.ReactivateApplicationRequest{
				ProjectId: p.GetId(),
				Id:        integration.ID(),
			},
			expectedErrorType: codes.NotFound,
		},
		{
			testName: "when app to reactivate is found should return deactivation time",
			inputCtx: IAMOwnerCtx,
			reactivateRequest: &app.ReactivateApplicationRequest{
				ProjectId: p.GetId(),
				Id:        appToReactivate.GetAppId(),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// When
			res, err := instance.Client.AppV2Beta.ReactivateApplication(tc.inputCtx, tc.reactivateRequest)

			// Then
			require.Equal(t, tc.expectedErrorType, status.Code(err))
			if tc.expectedErrorType == codes.OK {
				assert.NotZero(t, res.GetReactivationDate())
			}
		})
	}
}

func TestReactivateApplication_WithDifferentPermissions(t *testing.T) {
	p, projectOwnerCtx := getProjectAndProjectContext(t, instance, IAMOwnerCtx)

	appToReactivateForLoginUser := createAPIApp(t, IAMOwnerCtx, instance, p.GetId())
	deactivateApp(t, appToReactivateForLoginUser, p.GetId())

	appToReactivateForProjectOwner := createAPIApp(t, IAMOwnerCtx, instance, p.GetId())
	deactivateApp(t, appToReactivateForProjectOwner, p.GetId())

	appToReactivateForOrgOwner := createAPIApp(t, IAMOwnerCtx, instance, p.GetId())
	deactivateApp(t, appToReactivateForOrgOwner, p.GetId())

	t.Parallel()

	tt := []struct {
		testName          string
		inputCtx          context.Context
		reactivateRequest *app.ReactivateApplicationRequest

		expectedErrorType codes.Code
	}{
		// Login User
		{
			testName: "when user has no project.app.write permission for app reactivate request should return permission error",
			inputCtx: LoginUserCtx,
			reactivateRequest: &app.ReactivateApplicationRequest{
				ProjectId: p.GetId(),
				Id:        appToReactivateForLoginUser.GetAppId(),
			},
			expectedErrorType: codes.PermissionDenied,
		},

		// Project Owner
		{
			testName: "when user is ProjectOwner reactivate app request should succeed",
			inputCtx: projectOwnerCtx,
			reactivateRequest: &app.ReactivateApplicationRequest{
				ProjectId: p.GetId(),
				Id:        appToReactivateForProjectOwner.GetAppId(),
			},
		},

		// Org Owner
		{
			testName: "when user is OrgOwner reactivate app request should succeed",
			inputCtx: OrgOwnerCtx,
			reactivateRequest: &app.ReactivateApplicationRequest{
				ProjectId: p.GetId(),
				Id:        appToReactivateForOrgOwner.GetAppId(),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// When
			res, err := instance.Client.AppV2Beta.ReactivateApplication(tc.inputCtx, tc.reactivateRequest)

			// Then
			require.Equal(t, tc.expectedErrorType, status.Code(err))
			if tc.expectedErrorType == codes.OK {
				assert.NotZero(t, res.GetReactivationDate())
			}
		})
	}
}

func TestRegenerateClientSecret(t *testing.T) {
	p := instance.CreateProject(IAMOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)

	reqForApiAppCreation := &app.CreateApplicationRequest_ApiRequest{
		ApiRequest: &app.CreateAPIApplicationRequest{AuthMethodType: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT},
	}

	apiAppToRegen, apiAppCreateErr := instance.Client.AppV2Beta.CreateApplication(IAMOwnerCtx, &app.CreateApplicationRequest{
		ProjectId:           p.GetId(),
		Name:                integration.ApplicationName(),
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
		ProjectId:           p.GetId(),
		Name:                integration.ApplicationName(),
		CreationRequestType: reqForOIDCAppCreation,
	})
	require.Nil(t, oidcAppCreateErr)

	t.Parallel()

	tt := []struct {
		testName     string
		inputCtx     context.Context
		regenRequest *app.RegenerateClientSecretRequest

		expectedErrorType codes.Code
		oldSecret         string
	}{
		{
			testName: "when app to regen is not expected type should return invalid argument error",
			inputCtx: IAMOwnerCtx,
			regenRequest: &app.RegenerateClientSecretRequest{
				ProjectId:     p.GetId(),
				ApplicationId: integration.ID(),
			},
			expectedErrorType: codes.InvalidArgument,
		},
		{
			testName: "when app to regen is not found should return not found error",
			inputCtx: IAMOwnerCtx,
			regenRequest: &app.RegenerateClientSecretRequest{
				ProjectId:     p.GetId(),
				ApplicationId: integration.ID(),
				AppType:       &app.RegenerateClientSecretRequest_IsApi{},
			},
			expectedErrorType: codes.NotFound,
		},
		{
			testName: "when API app to regen is found should return different secret",
			inputCtx: IAMOwnerCtx,
			regenRequest: &app.RegenerateClientSecretRequest{
				ProjectId:     p.GetId(),
				ApplicationId: apiAppToRegen.GetAppId(),
				AppType:       &app.RegenerateClientSecretRequest_IsApi{},
			},
			oldSecret: apiAppToRegen.GetApiResponse().GetClientSecret(),
		},
		{
			testName: "when OIDC app to regen is found should return different secret",
			inputCtx: IAMOwnerCtx,
			regenRequest: &app.RegenerateClientSecretRequest{
				ProjectId:     p.GetId(),
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
			res, err := instance.Client.AppV2Beta.RegenerateClientSecret(tc.inputCtx, tc.regenRequest)

			// Then
			require.Equal(t, tc.expectedErrorType, status.Code(err))
			if tc.expectedErrorType == codes.OK {
				assert.NotZero(t, res.GetCreationDate())
				assert.NotEqual(t, tc.oldSecret, res.GetClientSecret())
			}
		})
	}

}

func TestRegenerateClientSecret_WithDifferentPermissions(t *testing.T) {
	p, projectOwnerCtx := getProjectAndProjectContext(t, instance, IAMOwnerCtx)

	apiAppToRegenForLoginUser := createAPIApp(t, IAMOwnerCtx, instance, p.GetId())
	apiAppToRegenForProjectOwner := createAPIApp(t, IAMOwnerCtx, instance, p.GetId())
	apiAppToRegenForOrgOwner := createAPIApp(t, IAMOwnerCtx, instance, p.GetId())

	oidcAppToRegenForLoginUser := createOIDCApp(t, baseURI, p.GetId())
	oidcAppToRegenForProjectOwner := createOIDCApp(t, baseURI, p.GetId())
	oidcAppToRegenForOrgOwner := createOIDCApp(t, baseURI, p.GetId())

	t.Parallel()

	tt := []struct {
		testName     string
		inputCtx     context.Context
		regenRequest *app.RegenerateClientSecretRequest

		expectedErrorType codes.Code
		oldSecret         string
	}{
		// Login user
		{
			testName: "when user has no project.app.write permission for API app secret regen request should return permission error",
			inputCtx: LoginUserCtx,
			regenRequest: &app.RegenerateClientSecretRequest{
				ProjectId:     p.GetId(),
				ApplicationId: apiAppToRegenForLoginUser.GetAppId(),
				AppType:       &app.RegenerateClientSecretRequest_IsApi{},
			},
			expectedErrorType: codes.PermissionDenied,
		},
		{
			testName: "when user has no project.app.write permission for OIDC app secret regen request should return permission error",
			inputCtx: LoginUserCtx,
			regenRequest: &app.RegenerateClientSecretRequest{
				ProjectId:     p.GetId(),
				ApplicationId: oidcAppToRegenForLoginUser.GetAppId(),
				AppType:       &app.RegenerateClientSecretRequest_IsOidc{},
			},
			expectedErrorType: codes.PermissionDenied,
		},

		// Project Owner
		{
			testName: "when user is ProjectOwner regen API app secret request should succeed",
			inputCtx: projectOwnerCtx,
			regenRequest: &app.RegenerateClientSecretRequest{
				ProjectId:     p.GetId(),
				ApplicationId: apiAppToRegenForProjectOwner.GetAppId(),
				AppType:       &app.RegenerateClientSecretRequest_IsApi{},
			},
			oldSecret: apiAppToRegenForProjectOwner.GetApiResponse().GetClientSecret(),
		},
		{
			testName: "when user is ProjectOwner regen OIDC app secret request should succeed",
			inputCtx: projectOwnerCtx,
			regenRequest: &app.RegenerateClientSecretRequest{
				ProjectId:     p.GetId(),
				ApplicationId: oidcAppToRegenForProjectOwner.GetAppId(),
				AppType:       &app.RegenerateClientSecretRequest_IsOidc{},
			},
			oldSecret: oidcAppToRegenForProjectOwner.GetOidcResponse().GetClientSecret(),
		},

		// Org Owner
		{
			testName: "when user is OrgOwner regen API app secret request should succeed",
			inputCtx: OrgOwnerCtx,
			regenRequest: &app.RegenerateClientSecretRequest{
				ProjectId:     p.GetId(),
				ApplicationId: apiAppToRegenForOrgOwner.GetAppId(),
				AppType:       &app.RegenerateClientSecretRequest_IsApi{},
			},
			oldSecret: apiAppToRegenForOrgOwner.GetApiResponse().GetClientSecret(),
		},
		{
			testName: "when user is OrgOwner regen OIDC app secret request should succeed",
			inputCtx: OrgOwnerCtx,
			regenRequest: &app.RegenerateClientSecretRequest{
				ProjectId:     p.GetId(),
				ApplicationId: oidcAppToRegenForOrgOwner.GetAppId(),
				AppType:       &app.RegenerateClientSecretRequest_IsOidc{},
			},
			oldSecret: oidcAppToRegenForOrgOwner.GetOidcResponse().GetClientSecret(),
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// When
			res, err := instance.Client.AppV2Beta.RegenerateClientSecret(tc.inputCtx, tc.regenRequest)

			// Then
			require.Equal(t, tc.expectedErrorType, status.Code(err))
			if tc.expectedErrorType == codes.OK {
				assert.NotZero(t, res.GetCreationDate())
				assert.NotEqual(t, tc.oldSecret, res.GetClientSecret())
			}
		})
	}

}
