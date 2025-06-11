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
)

func TestGetApplication(t *testing.T) {
	t.Parallel()

	iamOwnerCtx := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)
	p := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.Id, gofakeit.AppName(), false, false)

	t.Cleanup(func() {
		instance.DeleteProject(iamOwnerCtx, t, p.GetId())
	})

	apiAppName := gofakeit.AppName()
	createdApiApp, errAPIAppCreation := instance.Client.AppV2Beta.CreateApplication(iamOwnerCtx, &app.CreateApplicationRequest{
		ProjectId: p.GetId(),
		Name:      apiAppName,
		CreationRequestType: &app.CreateApplicationRequest_ApiRequest{
			ApiRequest: &app.CreateAPIApplicationRequest{
				AuthMethodType: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_BASIC,
			},
		},
	})
	require.Nil(t, errAPIAppCreation)

	samlAppName := gofakeit.AppName()
	createdSAMLApp, errSAMLAppCreation := instance.Client.AppV2Beta.CreateApplication(iamOwnerCtx, &app.CreateApplicationRequest{
		ProjectId: p.GetId(),
		Name:      samlAppName,
		CreationRequestType: &app.CreateApplicationRequest_SamlRequest{
			SamlRequest: &app.CreateSAMLApplicationRequest{
				LoginVersion: &app.LoginVersion{Version: &app.LoginVersion_LoginV1{LoginV1: &app.LoginV1{}}},
				Metadata:     &app.CreateSAMLApplicationRequest_MetadataXml{MetadataXml: samlMetadataGen(gofakeit.URL())},
			},
		},
	})
	require.Nil(t, errSAMLAppCreation)

	oidcAppName := gofakeit.AppName()
	createdOIDCApp, errOIDCAppCreation := instance.Client.AppV2Beta.CreateApplication(iamOwnerCtx, &app.CreateApplicationRequest{
		ProjectId: p.GetId(),
		Name:      oidcAppName,
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
				LoginVersion:           &app.LoginVersion{Version: &app.LoginVersion_LoginV2{LoginV2: &app.LoginV2{BaseUri: &baseURI}}},
			},
		},
	})
	require.Nil(t, errOIDCAppCreation)

	tt := []struct {
		testName     string
		inputRequest *app.GetApplicationRequest

		expectedErrorType       codes.Code
		expectedAppName         string
		expectedAppID           string
		expectedApplicationType string
	}{
		{
			testName: "when unknown app ID should return not found error",
			inputRequest: &app.GetApplicationRequest{
				ApplicationId: gofakeit.Sentence(2),
			},

			expectedErrorType: codes.NotFound,
		},
		{
			testName: "when providing API app ID should return valid API app result",
			inputRequest: &app.GetApplicationRequest{
				ApplicationId: createdApiApp.GetAppId(),
			},

			expectedAppName:         apiAppName,
			expectedAppID:           createdApiApp.GetAppId(),
			expectedApplicationType: fmt.Sprintf("%T", &app.Application_ApiConfig{}),
		},
		{
			testName: "when providing SAML app ID should return valid SAML app result",
			inputRequest: &app.GetApplicationRequest{
				ApplicationId: createdSAMLApp.GetAppId(),
			},

			expectedAppName:         samlAppName,
			expectedAppID:           createdSAMLApp.GetAppId(),
			expectedApplicationType: fmt.Sprintf("%T", &app.Application_SamlConfig{}),
		},
		{
			testName: "when providing OIDC app ID should return valid OIDC app result",
			inputRequest: &app.GetApplicationRequest{
				ApplicationId: createdOIDCApp.GetAppId(),
			},

			expectedAppName:         oidcAppName,
			expectedAppID:           createdOIDCApp.GetAppId(),
			expectedApplicationType: fmt.Sprintf("%T", &app.Application_OidcConfig{}),
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// When
			res, err := instance.Client.AppV2Beta.GetApplication(iamOwnerCtx, tc.inputRequest)

			// Then
			require.Equal(t, tc.expectedErrorType, status.Code(err))
			if tc.expectedErrorType == codes.OK {

				assert.Equal(t, tc.expectedAppID, res.GetApp().GetId())
				assert.Equal(t, tc.expectedAppName, res.GetApp().GetName())
				assert.NotZero(t, res.GetApp().GetCreationDate())
				assert.NotZero(t, res.GetApp().GetChangeDate())

				appType := fmt.Sprintf("%T", res.GetApp().GetConfig())
				assert.Equal(t, tc.expectedApplicationType, appType)
			}
		})
	}
}
