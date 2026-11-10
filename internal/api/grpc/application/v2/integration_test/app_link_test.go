//go:build integration

package app_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/application/v2"
)

func TestOIDCAppLinkConfig(t *testing.T) {
	p := instance.CreateProject(IAMOwnerCtx, t, instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)

	ios := &application.IOSAppLinkConfig{
		TeamId:   "ABCDE12345",
		BundleId: "com.example.app",
	}
	android := &application.AndroidAppLinkConfig{
		PackageName:            "com.example.app",
		Sha256CertFingerprints: []string{"AA:BB:CC:DD:EE:FF:00:11:22:33:44:55:66:77:88:99:AA:BB:CC:DD:EE:FF:00:11:22:33:44:55:66:77:88:99"},
	}

	t.Run("create with incomplete ios config is rejected", func(t *testing.T) {
		_, err := instance.Client.ApplicationV2.CreateApplication(IAMOwnerCtx, &application.CreateApplicationRequest{
			ProjectId: p.GetId(),
			Name:      integration.ApplicationName(),
			ApplicationType: &application.CreateApplicationRequest_OidcConfiguration{
				OidcConfiguration: baseOIDCCreateRequest(iosOnlyTeamID()),
			},
		})
		require.Equal(t, codes.InvalidArgument, status.Code(err))
	})

	t.Run("create get update omit and clear round trip", func(t *testing.T) {
		created, err := instance.Client.ApplicationV2.CreateApplication(IAMOwnerCtx, &application.CreateApplicationRequest{
			ProjectId: p.GetId(),
			Name:      integration.ApplicationName(),
			ApplicationType: &application.CreateApplicationRequest_OidcConfiguration{
				OidcConfiguration: baseOIDCCreateRequestWithAppLinks(ios, android),
			},
		})
		require.NoError(t, err)
		require.NotEmpty(t, created.GetApplicationId())

		assertOIDCAppLinks(t, created.GetApplicationId(), ios, android)

		updatedIos := &application.IOSAppLinkConfig{
			TeamId:   "NEWTEAM123",
			BundleId: "com.example.app.updated",
		}
		updatedAndroid := &application.AndroidAppLinkConfig{
			PackageName:            "com.example.app.updated",
			Sha256CertFingerprints: []string{"11:22:33:44:55:66:77:88:99:AA:BB:CC:DD:EE:FF:00:11:22:33:44:55:66:77:88:99:AA:BB:CC:DD:EE:FF:00"},
		}
		_, err = instance.Client.ApplicationV2.UpdateApplication(IAMOwnerCtx, &application.UpdateApplicationRequest{
			ProjectId:     p.GetId(),
			ApplicationId: created.GetApplicationId(),
			ApplicationType: &application.UpdateApplicationRequest_OidcConfiguration{
				OidcConfiguration: &application.UpdateOIDCApplicationConfigurationRequest{
					Ios:     updatedIos,
					Android: updatedAndroid,
				},
			},
		})
		require.NoError(t, err)
		assertOIDCAppLinks(t, created.GetApplicationId(), updatedIos, updatedAndroid)

		_, err = instance.Client.ApplicationV2.UpdateApplication(IAMOwnerCtx, &application.UpdateApplicationRequest{
			ProjectId:     p.GetId(),
			ApplicationId: created.GetApplicationId(),
			ApplicationType: &application.UpdateApplicationRequest_OidcConfiguration{
				OidcConfiguration: &application.UpdateOIDCApplicationConfigurationRequest{
					PostLogoutRedirectUris: []string{"http://example.com/home-omitted"},
				},
			},
		})
		require.NoError(t, err)
		assertOIDCAppLinks(t, created.GetApplicationId(), updatedIos, updatedAndroid)

		_, err = instance.Client.ApplicationV2.UpdateApplication(IAMOwnerCtx, &application.UpdateApplicationRequest{
			ProjectId:     p.GetId(),
			ApplicationId: created.GetApplicationId(),
			ApplicationType: &application.UpdateApplicationRequest_OidcConfiguration{
				OidcConfiguration: &application.UpdateOIDCApplicationConfigurationRequest{
					Ios:     &application.IOSAppLinkConfig{},
					Android: &application.AndroidAppLinkConfig{},
				},
			},
		})
		require.NoError(t, err)
		assertOIDCAppLinks(t, created.GetApplicationId(), nil, nil)
	})
}

func baseOIDCCreateRequest(opts ...func(*application.CreateOIDCApplicationRequest)) *application.CreateOIDCApplicationRequest {
	req := &application.CreateOIDCApplicationRequest{
		RedirectUris:           []string{"http://example.com"},
		ResponseTypes:          []application.OIDCResponseType{application.OIDCResponseType_OIDC_RESPONSE_TYPE_CODE},
		GrantTypes:             []application.OIDCGrantType{application.OIDCGrantType_OIDC_GRANT_TYPE_AUTHORIZATION_CODE},
		ApplicationType:        application.OIDCApplicationType_OIDC_APP_TYPE_NATIVE,
		AuthMethodType:         application.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_NONE,
		PostLogoutRedirectUris: []string{"http://example.com/home"},
		Version:                application.OIDCVersion_OIDC_VERSION_1_0,
		AccessTokenType:        application.OIDCTokenType_OIDC_TOKEN_TYPE_BEARER,
	}
	for _, opt := range opts {
		opt(req)
	}
	return req
}

func baseOIDCCreateRequestWithAppLinks(ios *application.IOSAppLinkConfig, android *application.AndroidAppLinkConfig) *application.CreateOIDCApplicationRequest {
	return baseOIDCCreateRequest(func(req *application.CreateOIDCApplicationRequest) {
		req.Ios = ios
		req.Android = android
	})
}

func iosOnlyTeamID() func(*application.CreateOIDCApplicationRequest) {
	return func(req *application.CreateOIDCApplicationRequest) {
		req.Ios = &application.IOSAppLinkConfig{TeamId: "ABCDE12345"}
	}
}

func assertOIDCAppLinks(t *testing.T, appID string, wantIOS *application.IOSAppLinkConfig, wantAndroid *application.AndroidAppLinkConfig) {
	t.Helper()
	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(IAMOwnerCtx, time.Minute)
	require.EventuallyWithT(t, func(ct *assert.CollectT) {
		res, err := instance.Client.ApplicationV2.GetApplication(IAMOwnerCtx, &application.GetApplicationRequest{
			ApplicationId: appID,
		})
		require.NoError(ct, err)
		cfg := res.GetApplication().GetOidcConfiguration()
		require.NotNil(ct, cfg)

		if wantIOS == nil {
			assert.Nil(ct, cfg.GetIos())
		} else {
			require.NotNil(ct, cfg.GetIos())
			assert.Equal(ct, wantIOS.GetTeamId(), cfg.GetIos().GetTeamId())
			assert.Equal(ct, wantIOS.GetBundleId(), cfg.GetIos().GetBundleId())
		}

		if wantAndroid == nil {
			assert.Nil(ct, cfg.GetAndroid())
		} else {
			require.NotNil(ct, cfg.GetAndroid())
			assert.Equal(ct, wantAndroid.GetPackageName(), cfg.GetAndroid().GetPackageName())
			assert.Equal(ct, wantAndroid.GetSha256CertFingerprints(), cfg.GetAndroid().GetSha256CertFingerprints())
		}
	}, retryDuration, tick)
}
