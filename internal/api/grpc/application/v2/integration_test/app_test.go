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
	"github.com/zitadel/zitadel/pkg/grpc/application/v2"
	org "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
)

func TestCreateApplication(t *testing.T) {
	p := instance.CreateProject(IAMOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)

	t.Parallel()

	notExistingProjectID := integration.ID()

	tt := []struct {
		testName        string
		creationRequest *application.CreateApplicationRequest
		inputCtx        context.Context

		expectedResponseType string
		expectedErrorType    codes.Code
	}{
		{
			testName: "when project for API application creation is not found should return failed precondition error",
			inputCtx: IAMOwnerCtx,
			creationRequest: &application.CreateApplicationRequest{
				ProjectId: notExistingProjectID,
				Name:      "App Name",
				ApplicationType: &application.CreateApplicationRequest_ApiConfiguration{
					ApiConfiguration: &application.CreateAPIApplicationRequest{
						AuthMethodType: application.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT,
					},
				},
			},
			expectedErrorType: codes.FailedPrecondition,
		},
		{
			testName: "when CreateAPIApp request is valid should create application and return no error",
			inputCtx: IAMOwnerCtx,
			creationRequest: &application.CreateApplicationRequest{
				ProjectId: p.GetId(),
				Name:      "App Name",
				ApplicationType: &application.CreateApplicationRequest_ApiConfiguration{
					ApiConfiguration: &application.CreateAPIApplicationRequest{
						AuthMethodType: application.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT,
					},
				},
			},
			expectedResponseType: fmt.Sprintf("%T", &application.CreateApplicationResponse_ApiConfiguration{}),
		},
		{
			testName: "when project for OIDC application creation is not found should return failed precondition error",
			inputCtx: IAMOwnerCtx,
			creationRequest: &application.CreateApplicationRequest{
				ProjectId: notExistingProjectID,
				Name:      "App Name",
				ApplicationType: &application.CreateApplicationRequest_OidcConfiguration{
					OidcConfiguration: &application.CreateOIDCApplicationRequest{
						RedirectUris:           []string{"http://example.com"},
						ResponseTypes:          []application.OIDCResponseType{application.OIDCResponseType_OIDC_RESPONSE_TYPE_CODE},
						GrantTypes:             []application.OIDCGrantType{application.OIDCGrantType_OIDC_GRANT_TYPE_AUTHORIZATION_CODE},
						ApplicationType:        application.OIDCApplicationType_OIDC_APP_TYPE_WEB,
						AuthMethodType:         application.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_BASIC,
						PostLogoutRedirectUris: []string{"http://example.com/home"},
						Version:                application.OIDCVersion_OIDC_VERSION_1_0,
						AccessTokenType:        application.OIDCTokenType_OIDC_TOKEN_TYPE_JWT,
						BackChannelLogoutUri:   "http://example.com/logout",
						LoginVersion: &application.LoginVersion{
							Version: &application.LoginVersion_LoginV2{
								LoginV2: &application.LoginV2{
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
			testName: "when CreateOIDCApp request is valid should create application and return no error",
			inputCtx: IAMOwnerCtx,
			creationRequest: &application.CreateApplicationRequest{
				ProjectId: p.GetId(),
				Name:      integration.ApplicationName(),
				ApplicationType: &application.CreateApplicationRequest_OidcConfiguration{
					OidcConfiguration: &application.CreateOIDCApplicationRequest{
						RedirectUris:           []string{"http://example.com"},
						ResponseTypes:          []application.OIDCResponseType{application.OIDCResponseType_OIDC_RESPONSE_TYPE_CODE},
						GrantTypes:             []application.OIDCGrantType{application.OIDCGrantType_OIDC_GRANT_TYPE_AUTHORIZATION_CODE},
						ApplicationType:        application.OIDCApplicationType_OIDC_APP_TYPE_WEB,
						AuthMethodType:         application.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_BASIC,
						PostLogoutRedirectUris: []string{"http://example.com/home"},
						Version:                application.OIDCVersion_OIDC_VERSION_1_0,
						AccessTokenType:        application.OIDCTokenType_OIDC_TOKEN_TYPE_JWT,
						BackChannelLogoutUri:   "http://example.com/logout",
						LoginVersion: &application.LoginVersion{
							Version: &application.LoginVersion_LoginV2{
								LoginV2: &application.LoginV2{
									BaseUri: &baseURI,
								},
							},
						},
					},
				},
			},

			expectedResponseType: fmt.Sprintf("%T", &application.CreateApplicationResponse_OidcConfiguration{}),
		},
		{
			testName: "when project for SAML application creation is not found should return failed precondition error",
			inputCtx: IAMOwnerCtx,
			creationRequest: &application.CreateApplicationRequest{
				ProjectId: notExistingProjectID,
				Name:      integration.ApplicationName(),
				ApplicationType: &application.CreateApplicationRequest_SamlConfiguration{
					SamlConfiguration: &application.CreateSAMLApplicationRequest{
						Metadata: &application.CreateSAMLApplicationRequest_MetadataUrl{
							MetadataUrl: "http://example.com/metas",
						},
						LoginVersion: &application.LoginVersion{
							Version: &application.LoginVersion_LoginV2{
								LoginV2: &application.LoginV2{
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
			testName: "when CreateSAMLApp request is valid should create application and return no error",
			inputCtx: IAMOwnerCtx,
			creationRequest: &application.CreateApplicationRequest{
				ProjectId: p.GetId(),
				Name:      integration.ApplicationName(),
				ApplicationType: &application.CreateApplicationRequest_SamlConfiguration{
					SamlConfiguration: &application.CreateSAMLApplicationRequest{
						Metadata: &application.CreateSAMLApplicationRequest_MetadataXml{
							MetadataXml: samlMetadataGen(integration.URL()),
						},
						LoginVersion: &application.LoginVersion{
							Version: &application.LoginVersion_LoginV2{
								LoginV2: &application.LoginV2{
									BaseUri: &baseURI,
								},
							},
						},
					},
				},
			},
			expectedResponseType: fmt.Sprintf("%T", &application.CreateApplicationResponse_SamlConfiguration{}),
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			res, err := instance.Client.ApplicationV2.CreateApplication(tc.inputCtx, tc.creationRequest)

			require.Equal(t, tc.expectedErrorType, status.Code(err))
			if tc.expectedErrorType == codes.OK {
				resType := fmt.Sprintf("%T", res.GetApplicationType())
				assert.Equal(t, tc.expectedResponseType, resType)
				assert.NotZero(t, res.GetApplicationId())
				assert.NotZero(t, res.GetCreationDate())
			}
		})
	}
}

func TestCreateOIDCApplication_WithProvidedID(t *testing.T) {
	p := instance.CreateProject(IAMOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)

	t.Parallel()

	customID := integration.ID()
	res, err := instance.Client.ApplicationV2.CreateApplication(IAMOwnerCtx, &application.CreateApplicationRequest{
		ProjectId:       p.GetId(),
		Name:            integration.ApplicationName(),
		ApplicationId:   customID,
		ApplicationType: &application.CreateApplicationRequest_OidcConfiguration{
			OidcConfiguration: &application.CreateOIDCApplicationRequest{
				RedirectUris:           []string{"http://example.com"},
				ResponseTypes:          []application.OIDCResponseType{application.OIDCResponseType_OIDC_RESPONSE_TYPE_CODE},
				GrantTypes:             []application.OIDCGrantType{application.OIDCGrantType_OIDC_GRANT_TYPE_AUTHORIZATION_CODE},
				ApplicationType:        application.OIDCApplicationType_OIDC_APP_TYPE_WEB,
				AuthMethodType:         application.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_BASIC,
				PostLogoutRedirectUris: []string{"http://example.com/home"},
				Version:                application.OIDCVersion_OIDC_VERSION_1_0,
				AccessTokenType:        application.OIDCTokenType_OIDC_TOKEN_TYPE_JWT,
				BackChannelLogoutUri:   "http://example.com/logout",
				LoginVersion: &application.LoginVersion{
					Version: &application.LoginVersion_LoginV2{
						LoginV2: &application.LoginV2{BaseUri: &baseURI},
					},
				},
			},
		},
	})
	require.NoError(t, err)

	assert.Equal(t, customID, res.GetApplicationId())
	assert.NotEmpty(t, res.GetOidcConfiguration().GetClientId())
	assert.NotZero(t, res.GetCreationDate())
}

func TestCreateApplication_WithDifferentPermissions(t *testing.T) {
	p, projectOwnerCtx := getProjectAndProjectContext(t, instance, IAMOwnerCtx)

	t.Parallel()

	tt := []struct {
		testName        string
		creationRequest *application.CreateApplicationRequest
		inputCtx        context.Context

		expectedResponseType string
		expectedErrorType    codes.Code
	}{
		// Login User with no project.application.write
		{
			testName: "when user has no project.application.write permission for API request should return permission error",
			inputCtx: LoginUserCtx,
			creationRequest: &application.CreateApplicationRequest{
				ProjectId: p.GetId(),
				Name:      integration.ApplicationName(),
				ApplicationType: &application.CreateApplicationRequest_ApiConfiguration{
					ApiConfiguration: &application.CreateAPIApplicationRequest{
						AuthMethodType: application.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT,
					},
				},
			},
			expectedErrorType: codes.PermissionDenied,
		},
		{
			testName: "when user has no project.application.write permission for OIDC request should return permission error",
			inputCtx: LoginUserCtx,
			creationRequest: &application.CreateApplicationRequest{
				ProjectId: p.GetId(),
				Name:      integration.ApplicationName(),
				ApplicationType: &application.CreateApplicationRequest_OidcConfiguration{
					OidcConfiguration: &application.CreateOIDCApplicationRequest{
						RedirectUris:           []string{"http://example.com"},
						ResponseTypes:          []application.OIDCResponseType{application.OIDCResponseType_OIDC_RESPONSE_TYPE_CODE},
						GrantTypes:             []application.OIDCGrantType{application.OIDCGrantType_OIDC_GRANT_TYPE_AUTHORIZATION_CODE},
						ApplicationType:        application.OIDCApplicationType_OIDC_APP_TYPE_WEB,
						AuthMethodType:         application.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_BASIC,
						PostLogoutRedirectUris: []string{"http://example.com/home"},
						Version:                application.OIDCVersion_OIDC_VERSION_1_0,
						AccessTokenType:        application.OIDCTokenType_OIDC_TOKEN_TYPE_JWT,
						BackChannelLogoutUri:   "http://example.com/logout",
						LoginVersion: &application.LoginVersion{
							Version: &application.LoginVersion_LoginV2{
								LoginV2: &application.LoginV2{
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
			testName: "when user has no project.application.write permission for SAML request should return permission error",
			inputCtx: LoginUserCtx,
			creationRequest: &application.CreateApplicationRequest{
				ProjectId: p.GetId(),
				Name:      integration.ApplicationName(),
				ApplicationType: &application.CreateApplicationRequest_SamlConfiguration{
					SamlConfiguration: &application.CreateSAMLApplicationRequest{
						Metadata: &application.CreateSAMLApplicationRequest_MetadataXml{
							MetadataXml: samlMetadataGen(integration.URL()),
						},
						LoginVersion: &application.LoginVersion{
							Version: &application.LoginVersion_LoginV2{
								LoginV2: &application.LoginV2{
									BaseUri: &baseURI,
								},
							},
						},
					},
				},
			},
			expectedErrorType: codes.PermissionDenied,
		},

		// OrgOwner with project.application.write permission
		{
			testName: "when user is OrgOwner API request should succeed",
			inputCtx: OrgOwnerCtx,
			creationRequest: &application.CreateApplicationRequest{
				ProjectId: p.GetId(),
				Name:      integration.ApplicationName(),
				ApplicationType: &application.CreateApplicationRequest_ApiConfiguration{
					ApiConfiguration: &application.CreateAPIApplicationRequest{
						AuthMethodType: application.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT,
					},
				},
			},
			expectedResponseType: fmt.Sprintf("%T", &application.CreateApplicationResponse_ApiConfiguration{}),
		},
		{
			testName: "when user is OrgOwner OIDC request should succeed",
			inputCtx: OrgOwnerCtx,
			creationRequest: &application.CreateApplicationRequest{
				ProjectId: p.GetId(),
				Name:      integration.ApplicationName(),
				ApplicationType: &application.CreateApplicationRequest_OidcConfiguration{
					OidcConfiguration: &application.CreateOIDCApplicationRequest{
						RedirectUris:           []string{"http://example.com"},
						ResponseTypes:          []application.OIDCResponseType{application.OIDCResponseType_OIDC_RESPONSE_TYPE_CODE},
						GrantTypes:             []application.OIDCGrantType{application.OIDCGrantType_OIDC_GRANT_TYPE_AUTHORIZATION_CODE},
						ApplicationType:        application.OIDCApplicationType_OIDC_APP_TYPE_WEB,
						AuthMethodType:         application.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_BASIC,
						PostLogoutRedirectUris: []string{"http://example.com/home"},
						Version:                application.OIDCVersion_OIDC_VERSION_1_0,
						AccessTokenType:        application.OIDCTokenType_OIDC_TOKEN_TYPE_JWT,
						BackChannelLogoutUri:   "http://example.com/logout",
						LoginVersion: &application.LoginVersion{
							Version: &application.LoginVersion_LoginV2{
								LoginV2: &application.LoginV2{
									BaseUri: &baseURI,
								},
							},
						},
					},
				},
			},

			expectedResponseType: fmt.Sprintf("%T", &application.CreateApplicationResponse_OidcConfiguration{}),
		},
		{
			testName: "when user is OrgOwner SAML request should succeed",
			inputCtx: OrgOwnerCtx,
			creationRequest: &application.CreateApplicationRequest{
				ProjectId: p.GetId(),
				Name:      integration.ApplicationName(),
				ApplicationType: &application.CreateApplicationRequest_SamlConfiguration{
					SamlConfiguration: &application.CreateSAMLApplicationRequest{
						Metadata: &application.CreateSAMLApplicationRequest_MetadataXml{
							MetadataXml: samlMetadataGen(integration.URL()),
						},
						LoginVersion: &application.LoginVersion{
							Version: &application.LoginVersion_LoginV2{
								LoginV2: &application.LoginV2{
									BaseUri: &baseURI,
								},
							},
						},
					},
				},
			},
			expectedResponseType: fmt.Sprintf("%T", &application.CreateApplicationResponse_SamlConfiguration{}),
		},

		// Project owner with project.application.write permission
		{
			testName: "when user is ProjectOwner API request should succeed",
			inputCtx: projectOwnerCtx,
			creationRequest: &application.CreateApplicationRequest{
				ProjectId: p.GetId(),
				Name:      integration.ApplicationName(),
				ApplicationType: &application.CreateApplicationRequest_ApiConfiguration{
					ApiConfiguration: &application.CreateAPIApplicationRequest{
						AuthMethodType: application.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT,
					},
				},
			},
			expectedResponseType: fmt.Sprintf("%T", &application.CreateApplicationResponse_ApiConfiguration{}),
		},
		{
			testName: "when user is ProjectOwner OIDC request should succeed",
			inputCtx: projectOwnerCtx,
			creationRequest: &application.CreateApplicationRequest{
				ProjectId: p.GetId(),
				Name:      integration.ApplicationName(),
				ApplicationType: &application.CreateApplicationRequest_OidcConfiguration{
					OidcConfiguration: &application.CreateOIDCApplicationRequest{
						RedirectUris:           []string{"http://example.com"},
						ResponseTypes:          []application.OIDCResponseType{application.OIDCResponseType_OIDC_RESPONSE_TYPE_CODE},
						GrantTypes:             []application.OIDCGrantType{application.OIDCGrantType_OIDC_GRANT_TYPE_AUTHORIZATION_CODE},
						ApplicationType:        application.OIDCApplicationType_OIDC_APP_TYPE_WEB,
						AuthMethodType:         application.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_BASIC,
						PostLogoutRedirectUris: []string{"http://example.com/home"},
						Version:                application.OIDCVersion_OIDC_VERSION_1_0,
						AccessTokenType:        application.OIDCTokenType_OIDC_TOKEN_TYPE_JWT,
						BackChannelLogoutUri:   "http://example.com/logout",
						LoginVersion: &application.LoginVersion{
							Version: &application.LoginVersion_LoginV2{
								LoginV2: &application.LoginV2{
									BaseUri: &baseURI,
								},
							},
						},
					},
				},
			},

			expectedResponseType: fmt.Sprintf("%T", &application.CreateApplicationResponse_OidcConfiguration{}),
		},
		{
			testName: "when user is ProjectOwner SAML request should succeed",
			inputCtx: projectOwnerCtx,
			creationRequest: &application.CreateApplicationRequest{
				ProjectId: p.GetId(),
				Name:      integration.ApplicationName(),
				ApplicationType: &application.CreateApplicationRequest_SamlConfiguration{
					SamlConfiguration: &application.CreateSAMLApplicationRequest{
						Metadata: &application.CreateSAMLApplicationRequest_MetadataXml{
							MetadataXml: samlMetadataGen(integration.URL()),
						},
						LoginVersion: &application.LoginVersion{
							Version: &application.LoginVersion_LoginV2{
								LoginV2: &application.LoginV2{
									BaseUri: &baseURI,
								},
							},
						},
					},
				},
			},
			expectedResponseType: fmt.Sprintf("%T", &application.CreateApplicationResponse_SamlConfiguration{}),
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			res, err := instance.Client.ApplicationV2.CreateApplication(tc.inputCtx, tc.creationRequest)

			require.Equal(t, tc.expectedErrorType, status.Code(err))
			if tc.expectedErrorType == codes.OK {
				resType := fmt.Sprintf("%T", res.GetApplicationType())
				assert.Equal(t, tc.expectedResponseType, resType)
				assert.NotZero(t, res.GetApplicationId())
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

	reqForAppNameCreation := &application.CreateApplicationRequest_ApiConfiguration{
		ApiConfiguration: &application.CreateAPIApplicationRequest{AuthMethodType: application.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT},
	}
	reqForAPIAppCreation := reqForAppNameCreation

	reqForOIDCAppCreation := &application.CreateApplicationRequest_OidcConfiguration{
		OidcConfiguration: &application.CreateOIDCApplicationRequest{
			RedirectUris:           []string{"http://example.com"},
			ResponseTypes:          []application.OIDCResponseType{application.OIDCResponseType_OIDC_RESPONSE_TYPE_CODE},
			GrantTypes:             []application.OIDCGrantType{application.OIDCGrantType_OIDC_GRANT_TYPE_AUTHORIZATION_CODE},
			ApplicationType:        application.OIDCApplicationType_OIDC_APP_TYPE_WEB,
			AuthMethodType:         application.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_BASIC,
			PostLogoutRedirectUris: []string{"http://example.com/home"},
			Version:                application.OIDCVersion_OIDC_VERSION_1_0,
			AccessTokenType:        application.OIDCTokenType_OIDC_TOKEN_TYPE_JWT,
			BackChannelLogoutUri:   "http://example.com/logout",
			LoginVersion: &application.LoginVersion{
				Version: &application.LoginVersion_LoginV2{
					LoginV2: &application.LoginV2{
						BaseUri: &baseURI,
					},
				},
			},
		},
	}

	samlMetas := samlMetadataGen(integration.URL())
	reqForSAMLAppCreation := &application.CreateApplicationRequest_SamlConfiguration{
		SamlConfiguration: &application.CreateSAMLApplicationRequest{
			Metadata: &application.CreateSAMLApplicationRequest_MetadataXml{
				MetadataXml: samlMetas,
			},
			LoginVersion: &application.LoginVersion{
				Version: &application.LoginVersion_LoginV2{
					LoginV2: &application.LoginV2{
						BaseUri: &baseURI,
					},
				},
			},
		},
	}

	appForNameChange, appNameChangeErr := instance.Client.ApplicationV2.CreateApplication(IAMOwnerCtx, &application.CreateApplicationRequest{
		ProjectId:       p.GetId(),
		Name:            integration.ApplicationName(),
		ApplicationType: reqForAppNameCreation,
	})
	require.Nil(t, appNameChangeErr)

	appForAPIConfigChange, appAPIConfigChangeErr := instance.Client.ApplicationV2.CreateApplication(IAMOwnerCtx, &application.CreateApplicationRequest{
		ProjectId:       p.GetId(),
		Name:            integration.ApplicationName(),
		ApplicationType: reqForAPIAppCreation,
	})
	require.Nil(t, appAPIConfigChangeErr)

	appForOIDCConfigChange, appOIDCConfigChangeErr := instance.Client.ApplicationV2.CreateApplication(IAMOwnerCtx, &application.CreateApplicationRequest{
		ProjectId:       p.GetId(),
		Name:            integration.ApplicationName(),
		ApplicationType: reqForOIDCAppCreation,
	})
	require.Nil(t, appOIDCConfigChangeErr)

	appForSAMLConfigChange, appSAMLConfigChangeErr := instance.Client.ApplicationV2.CreateApplication(IAMOwnerCtx, &application.CreateApplicationRequest{
		ProjectId:       p.GetId(),
		Name:            integration.ApplicationName(),
		ApplicationType: reqForSAMLAppCreation,
	})
	require.Nil(t, appSAMLConfigChangeErr)

	t.Parallel()

	tt := []struct {
		testName      string
		inputCtx      context.Context
		updateRequest *application.UpdateApplicationRequest

		expectedErrorType codes.Code
	}{
		{
			testName: "when application for application name change request is not found should return not found error",
			inputCtx: IAMOwnerCtx,
			updateRequest: &application.UpdateApplicationRequest{
				ProjectId:     pNotInCtx.GetId(),
				ApplicationId: appForNameChange.GetApplicationId(),
				Name:          "New name",
			},
			expectedErrorType: codes.NotFound,
		},
		{
			testName: "when request for application name change is valid should return updated timestamp",
			inputCtx: IAMOwnerCtx,
			updateRequest: &application.UpdateApplicationRequest{
				ProjectId:     p.GetId(),
				ApplicationId: appForNameChange.GetApplicationId(),

				Name: "New name",
			},
		},

		{
			testName: "when application for API config change request is not found should return not found error",
			inputCtx: IAMOwnerCtx,
			updateRequest: &application.UpdateApplicationRequest{
				ProjectId:     pNotInCtx.GetId(),
				ApplicationId: appForAPIConfigChange.GetApplicationId(),
				ApplicationType: &application.UpdateApplicationRequest_ApiConfiguration{
					ApiConfiguration: &application.UpdateAPIApplicationConfigurationRequest{
						AuthMethodType: application.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT,
					},
				},
			},
			expectedErrorType: codes.NotFound,
		},
		{
			testName: "when request for API config change is valid should return updated timestamp",
			inputCtx: IAMOwnerCtx,
			updateRequest: &application.UpdateApplicationRequest{
				ProjectId:     p.GetId(),
				ApplicationId: appForAPIConfigChange.GetApplicationId(),
				ApplicationType: &application.UpdateApplicationRequest_ApiConfiguration{
					ApiConfiguration: &application.UpdateAPIApplicationConfigurationRequest{
						AuthMethodType: application.APIAuthMethodType_API_AUTH_METHOD_TYPE_BASIC,
					},
				},
			},
		},
		{
			testName: "when application for OIDC config change request is not found should return not found error",
			inputCtx: IAMOwnerCtx,
			updateRequest: &application.UpdateApplicationRequest{
				ProjectId:     pNotInCtx.GetId(),
				ApplicationId: appForOIDCConfigChange.GetApplicationId(),
				ApplicationType: &application.UpdateApplicationRequest_OidcConfiguration{
					OidcConfiguration: &application.UpdateOIDCApplicationConfigurationRequest{
						PostLogoutRedirectUris: []string{"http://example.com/home2"},
					},
				},
			},
			expectedErrorType: codes.NotFound,
		},
		{
			testName: "when request for OIDC config change is valid should return updated timestamp",
			inputCtx: IAMOwnerCtx,
			updateRequest: &application.UpdateApplicationRequest{
				ProjectId:     p.GetId(),
				ApplicationId: appForOIDCConfigChange.GetApplicationId(),
				ApplicationType: &application.UpdateApplicationRequest_OidcConfiguration{
					OidcConfiguration: &application.UpdateOIDCApplicationConfigurationRequest{
						PostLogoutRedirectUris: []string{"http://example.com/home2"},
					},
				},
			},
		},

		{
			testName: "when application for SAML config change request is not found should return not found error",
			inputCtx: IAMOwnerCtx,
			updateRequest: &application.UpdateApplicationRequest{
				ProjectId:     pNotInCtx.GetId(),
				ApplicationId: appForSAMLConfigChange.GetApplicationId(),
				ApplicationType: &application.UpdateApplicationRequest_SamlConfiguration{
					SamlConfiguration: &application.UpdateSAMLApplicationConfigurationRequest{
						Metadata: &application.UpdateSAMLApplicationConfigurationRequest_MetadataXml{
							MetadataXml: samlMetas,
						},
						LoginVersion: &application.LoginVersion{Version: &application.LoginVersion_LoginV1{LoginV1: &application.LoginV1{}}},
					},
				},
			},
			expectedErrorType: codes.NotFound,
		},
		{
			testName: "when request for SAML config change is valid should return updated timestamp",
			inputCtx: IAMOwnerCtx,
			updateRequest: &application.UpdateApplicationRequest{
				ProjectId:     p.GetId(),
				ApplicationId: appForSAMLConfigChange.GetApplicationId(),
				ApplicationType: &application.UpdateApplicationRequest_SamlConfiguration{
					SamlConfiguration: &application.UpdateSAMLApplicationConfigurationRequest{
						Metadata: &application.UpdateSAMLApplicationConfigurationRequest_MetadataXml{
							MetadataXml: samlMetas,
						},
						LoginVersion: &application.LoginVersion{Version: &application.LoginVersion_LoginV1{LoginV1: &application.LoginV1{}}},
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			res, err := instance.Client.ApplicationV2.UpdateApplication(tc.inputCtx, tc.updateRequest)

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

	reqForAppNameCreation := &application.CreateApplicationRequest_ApiConfiguration{
		ApiConfiguration: &application.CreateAPIApplicationRequest{AuthMethodType: application.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT},
	}

	appForNameChange, appNameChangeErr := instance.Client.ApplicationV2.CreateApplication(IAMOwnerCtx, &application.CreateApplicationRequest{
		ProjectId:       p.GetId(),
		Name:            integration.ApplicationName(),
		ApplicationType: reqForAppNameCreation,
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
		updateRequest *application.UpdateApplicationRequest

		expectedErrorType codes.Code
	}{
		// ProjectOwner
		{
			testName: "when user is ProjectOwner application name request should succeed",
			inputCtx: projectOwnerCtx,
			updateRequest: &application.UpdateApplicationRequest{
				ProjectId:     p.GetId(),
				ApplicationId: appForNameChange.GetApplicationId(),

				Name: integration.ApplicationName(),
			},
		},
		{
			testName: "when user is ProjectOwner API application request should succeed",
			inputCtx: projectOwnerCtx,
			updateRequest: &application.UpdateApplicationRequest{
				ProjectId:     p.GetId(),
				ApplicationId: appForAPIConfigChangeForProjectOwner.GetApplicationId(),
				ApplicationType: &application.UpdateApplicationRequest_ApiConfiguration{
					ApiConfiguration: &application.UpdateAPIApplicationConfigurationRequest{
						AuthMethodType: application.APIAuthMethodType_API_AUTH_METHOD_TYPE_BASIC,
					},
				},
			},
		},
		{
			testName: "when user is ProjectOwner OIDC application request should succeed",
			inputCtx: projectOwnerCtx,
			updateRequest: &application.UpdateApplicationRequest{
				ProjectId:     p.GetId(),
				ApplicationId: appForOIDCConfigChangeForProjectOwner.GetApplicationId(),
				ApplicationType: &application.UpdateApplicationRequest_OidcConfiguration{
					OidcConfiguration: &application.UpdateOIDCApplicationConfigurationRequest{
						PostLogoutRedirectUris: []string{"http://example.com/home2"},
					},
				},
			},
		},
		{
			testName: "when user is ProjectOwner SAML request should succeed",
			inputCtx: projectOwnerCtx,
			updateRequest: &application.UpdateApplicationRequest{
				ProjectId:     p.GetId(),
				ApplicationId: appForSAMLConfigChangeForProjectOwner.GetApplicationId(),
				ApplicationType: &application.UpdateApplicationRequest_SamlConfiguration{
					SamlConfiguration: &application.UpdateSAMLApplicationConfigurationRequest{
						Metadata: &application.UpdateSAMLApplicationConfigurationRequest_MetadataXml{
							MetadataXml: samlMetasForProjectOwner,
						},
						LoginVersion: &application.LoginVersion{Version: &application.LoginVersion_LoginV1{LoginV1: &application.LoginV1{}}},
					},
				},
			},
		},

		// OrgOwner context
		{
			testName: "when user is OrgOwner application name request should succeed",
			inputCtx: OrgOwnerCtx,
			updateRequest: &application.UpdateApplicationRequest{
				ProjectId:     p.GetId(),
				ApplicationId: appForNameChange.GetApplicationId(),

				Name: integration.ApplicationName(),
			},
		},
		{
			testName: "when user is OrgOwner API application request should succeed",
			inputCtx: OrgOwnerCtx,
			updateRequest: &application.UpdateApplicationRequest{
				ProjectId:     p.GetId(),
				ApplicationId: appForAPIConfigChangeForOrgOwner.GetApplicationId(),
				ApplicationType: &application.UpdateApplicationRequest_ApiConfiguration{
					ApiConfiguration: &application.UpdateAPIApplicationConfigurationRequest{
						AuthMethodType: application.APIAuthMethodType_API_AUTH_METHOD_TYPE_BASIC,
					},
				},
			},
		},
		{
			testName: "when user is OrgOwner OIDC application request should succeed",
			inputCtx: OrgOwnerCtx,
			updateRequest: &application.UpdateApplicationRequest{
				ProjectId:     p.GetId(),
				ApplicationId: appForOIDCConfigChangeForOrgOwner.GetApplicationId(),
				ApplicationType: &application.UpdateApplicationRequest_OidcConfiguration{
					OidcConfiguration: &application.UpdateOIDCApplicationConfigurationRequest{
						PostLogoutRedirectUris: []string{"http://example.com/home2"},
					},
				},
			},
		},
		{
			testName: "when user is OrgOwner SAML request should succeed",
			inputCtx: OrgOwnerCtx,
			updateRequest: &application.UpdateApplicationRequest{
				ProjectId:     p.GetId(),
				ApplicationId: appForSAMLConfigChangeForOrgOwner.GetApplicationId(),
				ApplicationType: &application.UpdateApplicationRequest_SamlConfiguration{
					SamlConfiguration: &application.UpdateSAMLApplicationConfigurationRequest{
						Metadata: &application.UpdateSAMLApplicationConfigurationRequest_MetadataXml{
							MetadataXml: samlMetasForOrgOwner,
						},
						LoginVersion: &application.LoginVersion{Version: &application.LoginVersion_LoginV1{LoginV1: &application.LoginV1{}}},
					},
				},
			},
		},

		// LoginUser
		{
			testName: "when user has no project.application.write permission for application name change request should return permission error",
			inputCtx: LoginUserCtx,
			updateRequest: &application.UpdateApplicationRequest{
				ProjectId:     p.GetId(),
				ApplicationId: appForNameChange.GetApplicationId(),

				Name: integration.ApplicationName(),
			},
			expectedErrorType: codes.PermissionDenied,
		},
		{
			testName: "when user has no project.application.write permission for API request should return permission error",
			inputCtx: LoginUserCtx,
			updateRequest: &application.UpdateApplicationRequest{
				ProjectId:     p.GetId(),
				ApplicationId: appForAPIConfigChangeForLoginUser.GetApplicationId(),
				ApplicationType: &application.UpdateApplicationRequest_ApiConfiguration{
					ApiConfiguration: &application.UpdateAPIApplicationConfigurationRequest{
						AuthMethodType: application.APIAuthMethodType_API_AUTH_METHOD_TYPE_BASIC,
					},
				},
			},
			expectedErrorType: codes.PermissionDenied,
		},
		{
			testName: "when user has no project.application.write permission for OIDC request should return permission error",
			inputCtx: LoginUserCtx,
			updateRequest: &application.UpdateApplicationRequest{
				ProjectId:     p.GetId(),
				ApplicationId: appForOIDCConfigChangeForLoginUser.GetApplicationId(),
				ApplicationType: &application.UpdateApplicationRequest_OidcConfiguration{
					OidcConfiguration: &application.UpdateOIDCApplicationConfigurationRequest{
						PostLogoutRedirectUris: []string{"http://example.com/home2"},
					},
				},
			},
			expectedErrorType: codes.PermissionDenied,
		},
		{
			testName: "when user has no project.application.write permission for SAML request should return permission error",
			inputCtx: LoginUserCtx,
			updateRequest: &application.UpdateApplicationRequest{
				ProjectId:     p.GetId(),
				ApplicationId: appForSAMLConfigChangeForLoginUser.GetApplicationId(),
				ApplicationType: &application.UpdateApplicationRequest_SamlConfiguration{
					SamlConfiguration: &application.UpdateSAMLApplicationConfigurationRequest{
						Metadata: &application.UpdateSAMLApplicationConfigurationRequest_MetadataXml{
							MetadataXml: samlMetasForLoginUser,
						},
						LoginVersion: &application.LoginVersion{Version: &application.LoginVersion_LoginV1{LoginV1: &application.LoginV1{}}},
					},
				},
			},
			expectedErrorType: codes.PermissionDenied,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			res, err := instance.Client.ApplicationV2.UpdateApplication(tc.inputCtx, tc.updateRequest)

			require.Equal(t, tc.expectedErrorType, status.Code(err))
			if tc.expectedErrorType == codes.OK {
				assert.NotZero(t, res.GetChangeDate())
			}
		})
	}
}

func TestDeleteApplication(t *testing.T) {
	p := instance.CreateProject(IAMOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)

	reqForAppNameCreation := &application.CreateApplicationRequest_ApiConfiguration{
		ApiConfiguration: &application.CreateAPIApplicationRequest{AuthMethodType: application.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT},
	}

	appToDelete, appNameChangeErr := instance.Client.ApplicationV2.CreateApplication(IAMOwnerCtx, &application.CreateApplicationRequest{
		ProjectId:       p.GetId(),
		Name:            integration.ApplicationName(),
		ApplicationType: reqForAppNameCreation,
	})
	require.Nil(t, appNameChangeErr)

	t.Parallel()
	tt := []struct {
		testName      string
		deleteRequest *application.DeleteApplicationRequest
		inputCtx      context.Context

		expectedErrorType codes.Code
	}{
		{
			testName: "when application to delete is not found should return not found error",
			inputCtx: IAMOwnerCtx,
			deleteRequest: &application.DeleteApplicationRequest{
				ProjectId:     p.GetId(),
				ApplicationId: integration.ID(),
			},
			expectedErrorType: codes.NotFound,
		},
		{
			testName: "when application to delete is found should return deletion time",
			inputCtx: IAMOwnerCtx,
			deleteRequest: &application.DeleteApplicationRequest{
				ProjectId:     p.GetId(),
				ApplicationId: appToDelete.GetApplicationId(),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// When
			res, err := instance.Client.ApplicationV2.DeleteApplication(tc.inputCtx, tc.deleteRequest)

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
		deleteRequest *application.DeleteApplicationRequest
		inputCtx      context.Context

		expectedErrorType codes.Code
	}{
		// Login User
		{
			testName: "when user has no project.application.delete permission for application delete request should return permission error",
			inputCtx: LoginUserCtx,
			deleteRequest: &application.DeleteApplicationRequest{
				ProjectId:     p.GetId(),
				ApplicationId: appToDeleteForLoginUser.GetApplicationId(),
			},
			expectedErrorType: codes.PermissionDenied,
		},

		// Project Owner
		{
			testName: "when user is ProjectOwner delete application request should succeed",
			inputCtx: projectOwnerCtx,
			deleteRequest: &application.DeleteApplicationRequest{
				ProjectId:     p.GetId(),
				ApplicationId: appToDeleteForProjectOwner.GetApplicationId(),
			},
		},

		// Org Owner
		{
			testName: "when user is OrgOwner delete application request should succeed",
			inputCtx: projectOwnerCtx,
			deleteRequest: &application.DeleteApplicationRequest{
				ProjectId:     p.GetId(),
				ApplicationId: appToDeleteForOrgOwner.GetApplicationId(),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// When
			res, err := instance.Client.ApplicationV2.DeleteApplication(tc.inputCtx, tc.deleteRequest)

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

	reqForAppNameCreation := &application.CreateApplicationRequest_ApiConfiguration{
		ApiConfiguration: &application.CreateAPIApplicationRequest{AuthMethodType: application.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT},
	}

	appToDeactivate, appCreateErr := instance.Client.ApplicationV2.CreateApplication(IAMOwnerCtx, &application.CreateApplicationRequest{
		ProjectId:       p.GetId(),
		Name:            integration.ApplicationName(),
		ApplicationType: reqForAppNameCreation,
	})
	require.NoError(t, appCreateErr)

	t.Parallel()

	tt := []struct {
		testName      string
		inputCtx      context.Context
		deleteRequest *application.DeactivateApplicationRequest

		expectedErrorType codes.Code
	}{
		{
			testName: "when application to deactivate is not found should return not found error",
			inputCtx: IAMOwnerCtx,
			deleteRequest: &application.DeactivateApplicationRequest{
				ProjectId:     p.GetId(),
				ApplicationId: integration.ID(),
			},
			expectedErrorType: codes.NotFound,
		},
		{
			testName: "when application to deactivate is found should return deactivation time",
			inputCtx: IAMOwnerCtx,
			deleteRequest: &application.DeactivateApplicationRequest{
				ProjectId:     p.GetId(),
				ApplicationId: appToDeactivate.GetApplicationId(),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// When
			res, err := instance.Client.ApplicationV2.DeactivateApplication(tc.inputCtx, tc.deleteRequest)

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
		deleteRequest *application.DeactivateApplicationRequest

		expectedErrorType codes.Code
	}{
		// Login User
		{
			testName: "when user has no project.application.write permission for application deactivate request should return permission error",
			inputCtx: IAMOwnerCtx,
			deleteRequest: &application.DeactivateApplicationRequest{
				ProjectId:     p.GetId(),
				ApplicationId: appToDeactivateForLoginUser.GetApplicationId(),
			},
		},

		// Project Owner
		{
			testName: "when user is ProjectOwner deactivate application request should succeed",
			inputCtx: projectOwnerCtx,
			deleteRequest: &application.DeactivateApplicationRequest{
				ProjectId:     p.GetId(),
				ApplicationId: appToDeactivateForPrjectOwner.GetApplicationId(),
			},
		},

		// Org Owner
		{
			testName: "when user is OrgOwner deactivate application request should succeed",
			inputCtx: OrgOwnerCtx,
			deleteRequest: &application.DeactivateApplicationRequest{
				ProjectId:     p.GetId(),
				ApplicationId: appToDeactivateForOrgOwner.GetApplicationId(),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// When
			res, err := instance.Client.ApplicationV2.DeactivateApplication(tc.inputCtx, tc.deleteRequest)

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

	reqForAppNameCreation := &application.CreateApplicationRequest_ApiConfiguration{
		ApiConfiguration: &application.CreateAPIApplicationRequest{AuthMethodType: application.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT},
	}

	appToReactivate, appCreateErr := instance.Client.ApplicationV2.CreateApplication(IAMOwnerCtx, &application.CreateApplicationRequest{
		ProjectId:       p.GetId(),
		Name:            integration.ApplicationName(),
		ApplicationType: reqForAppNameCreation,
	})
	require.Nil(t, appCreateErr)

	_, appDeactivateErr := instance.Client.ApplicationV2.DeactivateApplication(IAMOwnerCtx, &application.DeactivateApplicationRequest{
		ProjectId:     p.GetId(),
		ApplicationId: appToReactivate.GetApplicationId(),
	})
	require.Nil(t, appDeactivateErr)

	t.Parallel()

	tt := []struct {
		testName          string
		inputCtx          context.Context
		reactivateRequest *application.ReactivateApplicationRequest

		expectedErrorType codes.Code
	}{
		{
			testName: "when application to reactivate is not found should return not found error",
			inputCtx: IAMOwnerCtx,
			reactivateRequest: &application.ReactivateApplicationRequest{
				ProjectId:     p.GetId(),
				ApplicationId: integration.ID(),
			},
			expectedErrorType: codes.NotFound,
		},
		{
			testName: "when application to reactivate is found should return deactivation time",
			inputCtx: IAMOwnerCtx,
			reactivateRequest: &application.ReactivateApplicationRequest{
				ProjectId:     p.GetId(),
				ApplicationId: appToReactivate.GetApplicationId(),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// When
			res, err := instance.Client.ApplicationV2.ReactivateApplication(tc.inputCtx, tc.reactivateRequest)

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
		reactivateRequest *application.ReactivateApplicationRequest

		expectedErrorType codes.Code
	}{
		// Login User
		{
			testName: "when user has no project.application.write permission for application reactivate request should return permission error",
			inputCtx: LoginUserCtx,
			reactivateRequest: &application.ReactivateApplicationRequest{
				ProjectId:     p.GetId(),
				ApplicationId: appToReactivateForLoginUser.GetApplicationId(),
			},
			expectedErrorType: codes.PermissionDenied,
		},

		// Project Owner
		{
			testName: "when user is ProjectOwner reactivate application request should succeed",
			inputCtx: projectOwnerCtx,
			reactivateRequest: &application.ReactivateApplicationRequest{
				ProjectId:     p.GetId(),
				ApplicationId: appToReactivateForProjectOwner.GetApplicationId(),
			},
		},

		// Org Owner
		{
			testName: "when user is OrgOwner reactivate application request should succeed",
			inputCtx: OrgOwnerCtx,
			reactivateRequest: &application.ReactivateApplicationRequest{
				ProjectId:     p.GetId(),
				ApplicationId: appToReactivateForOrgOwner.GetApplicationId(),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// When
			res, err := instance.Client.ApplicationV2.ReactivateApplication(tc.inputCtx, tc.reactivateRequest)

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

	reqForApiAppCreation := &application.CreateApplicationRequest_ApiConfiguration{
		ApiConfiguration: &application.CreateAPIApplicationRequest{AuthMethodType: application.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT},
	}

	apiAppToRegen, apiAppCreateErr := instance.Client.ApplicationV2.CreateApplication(IAMOwnerCtx, &application.CreateApplicationRequest{
		ProjectId:       p.GetId(),
		Name:            integration.ApplicationName(),
		ApplicationType: reqForApiAppCreation,
	})
	require.Nil(t, apiAppCreateErr)

	reqForOIDCAppCreation := &application.CreateApplicationRequest_OidcConfiguration{
		OidcConfiguration: &application.CreateOIDCApplicationRequest{
			RedirectUris:           []string{"http://example.com"},
			ResponseTypes:          []application.OIDCResponseType{application.OIDCResponseType_OIDC_RESPONSE_TYPE_CODE},
			GrantTypes:             []application.OIDCGrantType{application.OIDCGrantType_OIDC_GRANT_TYPE_AUTHORIZATION_CODE},
			ApplicationType:        application.OIDCApplicationType_OIDC_APP_TYPE_WEB,
			AuthMethodType:         application.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_BASIC,
			PostLogoutRedirectUris: []string{"http://example.com/home"},
			Version:                application.OIDCVersion_OIDC_VERSION_1_0,
			AccessTokenType:        application.OIDCTokenType_OIDC_TOKEN_TYPE_JWT,
			BackChannelLogoutUri:   "http://example.com/logout",
			LoginVersion: &application.LoginVersion{
				Version: &application.LoginVersion_LoginV2{
					LoginV2: &application.LoginV2{
						BaseUri: &baseURI,
					},
				},
			},
		},
	}

	oidcAppToRegen, oidcAppCreateErr := instance.Client.ApplicationV2.CreateApplication(IAMOwnerCtx, &application.CreateApplicationRequest{
		ProjectId:       p.GetId(),
		Name:            integration.ApplicationName(),
		ApplicationType: reqForOIDCAppCreation,
	})
	require.Nil(t, oidcAppCreateErr)

	t.Parallel()

	tt := []struct {
		testName        string
		inputCtx        context.Context
		generateRequest *application.GenerateClientSecretRequest

		expectedErrorType codes.Code
		oldSecret         string
	}{
		{
			testName: "when application to generate is not found should return not found error",
			inputCtx: IAMOwnerCtx,
			generateRequest: &application.GenerateClientSecretRequest{
				ProjectId:     p.GetId(),
				ApplicationId: integration.ID(),
			},
			expectedErrorType: codes.NotFound,
		},
		{
			testName: "when API application to generate is found should return different secret",
			inputCtx: IAMOwnerCtx,
			generateRequest: &application.GenerateClientSecretRequest{
				ProjectId:     p.GetId(),
				ApplicationId: apiAppToRegen.GetApplicationId(),
			},
			oldSecret: apiAppToRegen.GetApiConfiguration().GetClientSecret(),
		},
		{
			testName: "when OIDC application to generate is found should return different secret",
			inputCtx: IAMOwnerCtx,
			generateRequest: &application.GenerateClientSecretRequest{
				ProjectId:     p.GetId(),
				ApplicationId: oidcAppToRegen.GetApplicationId(),
			},
			oldSecret: oidcAppToRegen.GetOidcConfiguration().GetClientSecret(),
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// When
			res, err := instance.Client.ApplicationV2.GenerateClientSecret(tc.inputCtx, tc.generateRequest)

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
		regenRequest *application.GenerateClientSecretRequest

		expectedErrorType codes.Code
		oldSecret         string
	}{
		// Login user
		{
			testName: "when user has no project.application.write permission for API application secret regen request should return permission error",
			inputCtx: LoginUserCtx,
			regenRequest: &application.GenerateClientSecretRequest{
				ProjectId:     p.GetId(),
				ApplicationId: apiAppToRegenForLoginUser.GetApplicationId(),
			},
			expectedErrorType: codes.PermissionDenied,
		},
		{
			testName: "when user has no project.application.write permission for OIDC application secret regen request should return permission error",
			inputCtx: LoginUserCtx,
			regenRequest: &application.GenerateClientSecretRequest{
				ProjectId:     p.GetId(),
				ApplicationId: oidcAppToRegenForLoginUser.GetApplicationId(),
			},
			expectedErrorType: codes.PermissionDenied,
		},

		// Project Owner
		{
			testName: "when user is ProjectOwner regen API application secret request should succeed",
			inputCtx: projectOwnerCtx,
			regenRequest: &application.GenerateClientSecretRequest{
				ProjectId:     p.GetId(),
				ApplicationId: apiAppToRegenForProjectOwner.GetApplicationId(),
			},
			oldSecret: apiAppToRegenForProjectOwner.GetApiConfiguration().GetClientSecret(),
		},
		{
			testName: "when user is ProjectOwner regen OIDC application secret request should succeed",
			inputCtx: projectOwnerCtx,
			regenRequest: &application.GenerateClientSecretRequest{
				ProjectId:     p.GetId(),
				ApplicationId: oidcAppToRegenForProjectOwner.GetApplicationId(),
			},
			oldSecret: oidcAppToRegenForProjectOwner.GetOidcConfiguration().GetClientSecret(),
		},

		// Org Owner
		{
			testName: "when user is OrgOwner regen API application secret request should succeed",
			inputCtx: OrgOwnerCtx,
			regenRequest: &application.GenerateClientSecretRequest{
				ProjectId:     p.GetId(),
				ApplicationId: apiAppToRegenForOrgOwner.GetApplicationId(),
			},
			oldSecret: apiAppToRegenForOrgOwner.GetApiConfiguration().GetClientSecret(),
		},
		{
			testName: "when user is OrgOwner regen OIDC application secret request should succeed",
			inputCtx: OrgOwnerCtx,
			regenRequest: &application.GenerateClientSecretRequest{
				ProjectId:     p.GetId(),
				ApplicationId: oidcAppToRegenForOrgOwner.GetApplicationId(),
			},
			oldSecret: oidcAppToRegenForOrgOwner.GetOidcConfiguration().GetClientSecret(),
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// When
			res, err := instance.Client.ApplicationV2.GenerateClientSecret(tc.inputCtx, tc.regenRequest)

			// Then
			require.Equal(t, tc.expectedErrorType, status.Code(err))
			if tc.expectedErrorType == codes.OK {
				assert.NotZero(t, res.GetCreationDate())
				assert.NotEqual(t, tc.oldSecret, res.GetClientSecret())
			}
		})
	}

}
