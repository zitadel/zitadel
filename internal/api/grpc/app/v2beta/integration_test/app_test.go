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

func TestCreateApplication(t *testing.T) {
	iamOwnerCtx := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)

	orgNotInCtx := instance.CreateOrganization(iamOwnerCtx, gofakeit.Name(), gofakeit.Email())
	p := instance.CreateProject(iamOwnerCtx, t, instance.DefaultOrg.Id, gofakeit.AppName(), false, false)
	pNotInCtx := instance.CreateProject(iamOwnerCtx, t, orgNotInCtx.GetOrganizationId(), gofakeit.AppName(), false, false)

	baseURI := "http://example.com"
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
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			res, err := instance.Client.AppV2Beta.CreateApplication(iamOwnerCtx, tc.creationRequest)

			require.Equal(t, tc.expectedErrorType, status.Code(err))
			if tc.expectedErrorType == codes.OK {
				resType := fmt.Sprintf("%T", res.GetCreationResponseType())
				assert.Equal(t, tc.expectedResponseType, resType)
				assert.NotEmpty(t, res.GetAppId())
				assert.NotEmpty(t, res.GetCreationDate())
			}
		})
	}
}
