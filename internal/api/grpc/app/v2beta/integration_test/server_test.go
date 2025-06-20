//go:build integration

package instance_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	app "github.com/zitadel/zitadel/pkg/grpc/app/v2beta"
	project "github.com/zitadel/zitadel/pkg/grpc/project/v2beta"
)

var (
	ctx             context.Context
	NoPermissionCtx context.Context
	LoginUserCtx    context.Context
	ProjectOwnerCtx context.Context
	OrgOwnerCtx     context.Context
	IAMOwnerCtx     context.Context

	instance             *integration.Instance
	instancePermissionV2 *integration.Instance

	baseURI = "http://example.com"
	Project *project.CreateProjectResponse
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()

		instance = integration.NewInstance(ctx)
		IAMOwnerCtx = instance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		Project = instance.CreateProject(IAMOwnerCtx, &testing.T{}, instance.DefaultOrg.GetId(), gofakeit.Name(), false, false)

		LoginUserCtx = instance.WithAuthorization(ctx, integration.UserTypeLogin)
		OrgOwnerCtx = instance.WithAuthorization(ctx, integration.UserTypeOrgOwner)
		NoPermissionCtx = instance.WithAuthorization(ctx, integration.UserTypeNoPermission)

		userResp := instance.CreateMachineUser(IAMOwnerCtx)
		patResp := instance.CreatePersonalAccessToken(IAMOwnerCtx, userResp.GetUserId())
		instance.CreateProjectMembership(&testing.T{}, IAMOwnerCtx, Project.GetId(), userResp.GetUserId())
		ProjectOwnerCtx = integration.WithAuthorizationToken(ctx, patResp.Token)

		return m.Run()
	}())
}

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

func createSAMLApp(t *testing.T, baseURI string) ([]byte, *app.CreateApplicationResponse) {
	samlMetas := samlMetadataGen(gofakeit.URL())
	appForSAMLConfigChange, appSAMLConfigChangeErr := instance.Client.AppV2Beta.CreateApplication(IAMOwnerCtx, &app.CreateApplicationRequest{
		ProjectId: Project.GetId(),
		Name:      gofakeit.AppName(),
		CreationRequestType: &app.CreateApplicationRequest_SamlRequest{
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
		},
	})
	require.Nil(t, appSAMLConfigChangeErr)

	return samlMetas, appForSAMLConfigChange
}

func createOIDCApp(t *testing.T, baseURI string) *app.CreateApplicationResponse {
	appForOIDCConfigChange, appOIDCConfigChangeErr := instance.Client.AppV2Beta.CreateApplication(IAMOwnerCtx, &app.CreateApplicationRequest{
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
	})
	require.Nil(t, appOIDCConfigChangeErr)

	return appForOIDCConfigChange
}

func createAPIAppWithName(t *testing.T) (*app.CreateApplicationResponse, string) {
	appName := gofakeit.AppName()

	reqForAPIAppCreation := &app.CreateApplicationRequest_ApiRequest{
		ApiRequest: &app.CreateAPIApplicationRequest{AuthMethodType: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT},
	}

	appForAPIConfigChange, appAPIConfigChangeErr := instance.Client.AppV2Beta.CreateApplication(IAMOwnerCtx, &app.CreateApplicationRequest{
		ProjectId:           Project.GetId(),
		Name:                appName,
		CreationRequestType: reqForAPIAppCreation,
	})
	require.Nil(t, appAPIConfigChangeErr)

	return appForAPIConfigChange, appName
}

func createAPIApp(t *testing.T) *app.CreateApplicationResponse {
	res, _ := createAPIAppWithName(t)
	return res
}

func deactivateApp(t *testing.T, appToReactivate *app.CreateApplicationResponse) {
	_, appDeactivateErr := instance.Client.AppV2Beta.DeactivateApplication(IAMOwnerCtx, &app.DeactivateApplicationRequest{
		ProjectId: Project.GetId(),
		Id:        appToReactivate.GetAppId(),
	})
	require.Nil(t, appDeactivateErr)
}
