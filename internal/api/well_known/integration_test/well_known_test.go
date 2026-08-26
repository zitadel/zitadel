//go:build integration

package well_known_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/well_known"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/application/v2"
)

var (
	CTX         context.Context
	IAMOwnerCtx context.Context
	Instance    *integration.Instance
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()

		Instance = integration.NewInstance(ctx)
		CTX = ctx
		IAMOwnerCtx = Instance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		return m.Run()
	}())
}

func TestWellKnownAppLinks(t *testing.T) {
	issuer := Instance.OIDCIssuer()

	t.Run("empty when no apps configured", func(t *testing.T) {
		assertEventuallyJSON(t, issuer+well_known.AppleAppSiteAssociationPath, map[string]any{
			"webcredentials": map[string]any{
				"apps": []any{},
			},
		})
		assertEventuallyJSON(t, issuer+well_known.AssetLinksPath, []any{})
	})

	project := Instance.CreateProject(IAMOwnerCtx, t, Instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)

	ios1 := &application.IOSAppLinkConfig{TeamId: "TEAMAAAAAA", BundleId: "com.example.one"}
	android1 := &application.AndroidAppLinkConfig{
		PackageName:            "com.example.one",
		Sha256CertFingerprints: []string{"11:11:11:11:11:11:11:11:11:11:11:11:11:11:11:11:11:11:11:11:11:11:11:11:11:11:11:11:11:11:11:11"},
	}
	createOIDCAppWithLinks(t, project.GetId(), ios1, android1)

	ios2 := &application.IOSAppLinkConfig{TeamId: "TEAMBBBBBB", BundleId: "com.example.two"}
	createOIDCAppWithLinks(t, project.GetId(), ios2, nil)

	android2 := &application.AndroidAppLinkConfig{
		PackageName:            "com.example.two",
		Sha256CertFingerprints: []string{"22:22:22:22:22:22:22:22:22:22:22:22:22:22:22:22:22:22:22:22:22:22:22:22:22:22:22:22:22:22:22:22"},
	}
	createOIDCAppWithLinks(t, project.GetId(), nil, android2)

	// duplicate iOS app id should be deduped in AASA
	createOIDCAppWithLinks(t, project.GetId(), ios1, nil)

	assertEventuallyJSONUnorderedApps(t, issuer,
		[]string{
			"TEAMAAAAAA.com.example.one",
			"TEAMBBBBBB.com.example.two",
		},
		[]any{
			map[string]any{
				"relation": []any{"delegate_permission/common.handle_all_urls", "delegate_permission/common.get_login_creds"},
				"target": map[string]any{
					"namespace":                "android_app",
					"package_name":             "com.example.one",
					"sha256_cert_fingerprints": []any{android1.Sha256CertFingerprints[0]},
				},
			},
			map[string]any{
				"relation": []any{"delegate_permission/common.handle_all_urls", "delegate_permission/common.get_login_creds"},
				"target": map[string]any{
					"namespace":                "android_app",
					"package_name":             "com.example.two",
					"sha256_cert_fingerprints": []any{android2.Sha256CertFingerprints[0]},
				},
			},
		},
	)
}

func createOIDCAppWithLinks(t *testing.T, projectID string, ios *application.IOSAppLinkConfig, android *application.AndroidAppLinkConfig) {
	t.Helper()
	_, err := Instance.Client.ApplicationV2.CreateApplication(IAMOwnerCtx, &application.CreateApplicationRequest{
		ProjectId: projectID,
		Name:      integration.ApplicationName(),
		ApplicationType: &application.CreateApplicationRequest_OidcConfiguration{
			OidcConfiguration: &application.CreateOIDCApplicationRequest{
				RedirectUris:    []string{"http://example.com"},
				ResponseTypes:   []application.OIDCResponseType{application.OIDCResponseType_OIDC_RESPONSE_TYPE_CODE},
				GrantTypes:      []application.OIDCGrantType{application.OIDCGrantType_OIDC_GRANT_TYPE_AUTHORIZATION_CODE},
				ApplicationType: application.OIDCApplicationType_OIDC_APP_TYPE_NATIVE,
				AuthMethodType:  application.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_NONE,
				Version:         application.OIDCVersion_OIDC_VERSION_1_0,
				AccessTokenType: application.OIDCTokenType_OIDC_TOKEN_TYPE_BEARER,
				Ios:             ios,
				Android:         android,
			},
		},
	})
	require.NoError(t, err)
}

func assertEventuallyJSON(t *testing.T, url string, want any) {
	t.Helper()
	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
	require.EventuallyWithT(t, func(ct *assert.CollectT) {
		resp, err := http.Get(url)
		if !assert.NoError(ct, err) {
			return
		}
		defer func() {
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
		}()
		if !assert.Equal(ct, http.StatusOK, resp.StatusCode) {
			return
		}
		assert.Contains(ct, resp.Header.Get("Content-Type"), "application/json")
		assert.Equal(ct, "max-age=300, must-revalidate", resp.Header.Get("Cache-Control"))

		body, err := io.ReadAll(resp.Body)
		if !assert.NoError(ct, err) {
			return
		}
		var got any
		if !assert.NoError(ct, json.Unmarshal(body, &got)) {
			return
		}
		assert.Equal(ct, want, got)
	}, retryDuration, tick)
}

func assertEventuallyJSONUnorderedApps(t *testing.T, url string, wantApps []string, wantAssetLinks []any) {
	t.Helper()
	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
	require.EventuallyWithT(t, func(ct *assert.CollectT) {
		aasaResp, err := http.Get(url + well_known.AppleAppSiteAssociationPath)
		if !assert.NoError(ct, err) {
			return
		}
		defer aasaResp.Body.Close()
		if !assert.Equal(ct, http.StatusOK, aasaResp.StatusCode) {
			return
		}
		var aasa struct {
			WebCredentials struct {
				Apps []string `json:"apps"`
			} `json:"webcredentials"`
		}
		if !assert.NoError(ct, json.NewDecoder(aasaResp.Body).Decode(&aasa)) {
			return
		}
		assert.ElementsMatch(ct, wantApps, aasa.WebCredentials.Apps)

		assetResp, err := http.Get(url + well_known.AssetLinksPath)
		if !assert.NoError(ct, err) {
			return
		}
		defer assetResp.Body.Close()
		if !assert.Equal(ct, http.StatusOK, assetResp.StatusCode) {
			return
		}
		var got []any
		if !assert.NoError(ct, json.NewDecoder(assetResp.Body).Decode(&got)) {
			return
		}
		assert.ElementsMatch(ct, wantAssetLinks, got)
	}, retryDuration, tick)
}
