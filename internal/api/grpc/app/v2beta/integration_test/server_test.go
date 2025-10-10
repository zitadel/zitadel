//go:build integration

package app_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	app "github.com/zitadel/zitadel/pkg/grpc/app/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/feature/v2"
	project_v2beta "github.com/zitadel/zitadel/pkg/grpc/project/v2beta"
)

var (
	NoPermissionCtx context.Context
	LoginUserCtx    context.Context
	OrgOwnerCtx     context.Context
	IAMOwnerCtx     context.Context

	instance             *integration.Instance
	instancePermissionV2 *integration.Instance

	baseURI = "http://example.com"
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()

		instance = integration.NewInstance(ctx)
		instancePermissionV2 = integration.NewInstance(ctx)

		IAMOwnerCtx = instance.WithAuthorization(ctx, integration.UserTypeIAMOwner)

		LoginUserCtx = instance.WithAuthorization(ctx, integration.UserTypeLogin)
		OrgOwnerCtx = instance.WithAuthorization(ctx, integration.UserTypeOrgOwner)
		NoPermissionCtx = instance.WithAuthorization(ctx, integration.UserTypeNoPermission)

		return m.Run()
	}())
}

func getProjectAndProjectContext(t *testing.T, inst *integration.Instance, ctx context.Context) (*project_v2beta.CreateProjectResponse, context.Context) {
	project := inst.CreateProject(ctx, t, inst.DefaultOrg.GetId(), integration.ProjectName(), false, false)
	userResp := inst.CreateMachineUser(ctx)
	patResp := inst.CreatePersonalAccessToken(ctx, userResp.GetUserId())
	inst.CreateProjectMembership(t, ctx, project.GetId(), userResp.GetUserId())
	projectOwnerCtx := integration.WithAuthorizationToken(context.Background(), patResp.Token)

	return project, projectOwnerCtx
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

func createSAMLAppWithName(t *testing.T, baseURI, projectID string) ([]byte, *app.CreateApplicationResponse, string) {
	samlMetas := samlMetadataGen(integration.URL())
	appName := integration.ApplicationName()

	appForSAMLConfigChange, appSAMLConfigChangeErr := instance.Client.AppV2Beta.CreateApplication(IAMOwnerCtx, &app.CreateApplicationRequest{
		ProjectId: projectID,
		Name:      appName,
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

	return samlMetas, appForSAMLConfigChange, appName
}

func createSAMLApp(t *testing.T, baseURI, projectID string) ([]byte, *app.CreateApplicationResponse) {
	metas, app, _ := createSAMLAppWithName(t, baseURI, projectID)
	return metas, app
}

func createOIDCAppWithName(t *testing.T, baseURI, projectID string) (*app.CreateApplicationResponse, string) {
	appName := integration.ApplicationName()

	appForOIDCConfigChange, appOIDCConfigChangeErr := instance.Client.AppV2Beta.CreateApplication(IAMOwnerCtx, &app.CreateApplicationRequest{
		ProjectId: projectID,
		Name:      appName,
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

	return appForOIDCConfigChange, appName
}

func createOIDCApp(t *testing.T, baseURI, projctID string) *app.CreateApplicationResponse {
	app, _ := createOIDCAppWithName(t, baseURI, projctID)

	return app
}

func createAPIAppWithName(t *testing.T, ctx context.Context, inst *integration.Instance, projectID string) (*app.CreateApplicationResponse, string) {
	appName := integration.ApplicationName()

	reqForAPIAppCreation := &app.CreateApplicationRequest_ApiRequest{
		ApiRequest: &app.CreateAPIApplicationRequest{AuthMethodType: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT},
	}

	appForAPIConfigChange, appAPIConfigChangeErr := inst.Client.AppV2Beta.CreateApplication(ctx, &app.CreateApplicationRequest{
		ProjectId:           projectID,
		Name:                appName,
		CreationRequestType: reqForAPIAppCreation,
	})
	require.Nil(t, appAPIConfigChangeErr)

	return appForAPIConfigChange, appName
}

func createAPIApp(t *testing.T, ctx context.Context, inst *integration.Instance, projectID string) *app.CreateApplicationResponse {
	res, _ := createAPIAppWithName(t, ctx, inst, projectID)
	return res
}

func deactivateApp(t *testing.T, appToDeactivate *app.CreateApplicationResponse, projectID string) {
	_, appDeactivateErr := instance.Client.AppV2Beta.DeactivateApplication(IAMOwnerCtx, &app.DeactivateApplicationRequest{
		ProjectId: projectID,
		Id:        appToDeactivate.GetAppId(),
	})
	require.Nil(t, appDeactivateErr)
}

func ensureFeaturePermissionV2Enabled(t *testing.T, instance *integration.Instance) {
	ctx := instance.WithAuthorization(context.Background(), integration.UserTypeIAMOwner)
	f, err := instance.Client.FeatureV2.GetInstanceFeatures(ctx, &feature.GetInstanceFeaturesRequest{
		Inheritance: true,
	})
	require.NoError(t, err)

	if f.PermissionCheckV2.GetEnabled() {
		return
	}

	_, err = instance.Client.FeatureV2.SetInstanceFeatures(ctx, &feature.SetInstanceFeaturesRequest{
		PermissionCheckV2: gu.Ptr(true),
	})
	require.NoError(t, err)

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, 5*time.Minute)
	require.EventuallyWithT(t, func(tt *assert.CollectT) {
		f, err := instance.Client.FeatureV2.GetInstanceFeatures(ctx, &feature.GetInstanceFeaturesRequest{Inheritance: true})
		require.NoError(tt, err)
		assert.True(tt, f.PermissionCheckV2.GetEnabled())
	}, retryDuration, tick, "timed out waiting for ensuring instance feature")
}

func createAppKey(t *testing.T, ctx context.Context, inst *integration.Instance, projectID, appID string, expirationDate time.Time) *app.CreateApplicationKeyResponse {
	res, err := inst.Client.AppV2Beta.CreateApplicationKey(ctx,
		&app.CreateApplicationKeyRequest{
			AppId:          appID,
			ProjectId:      projectID,
			ExpirationDate: timestamppb.New(expirationDate.UTC()),
		},
	)

	require.Nil(t, err)

	return res
}
