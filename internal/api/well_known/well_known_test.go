package well_known

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/query"
)

func TestBuildAppleAppSiteAssociation(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		configs  []*query.OIDCAppLinkConfig
		wantApps []string
	}{
		{
			name:     "empty",
			configs:  nil,
			wantApps: []string{},
		},
		{
			name: "skips incomplete ios",
			configs: []*query.OIDCAppLinkConfig{
				{IOSTeamID: "TEAM"},
				{IOSBundleID: "com.example"},
			},
			wantApps: []string{},
		},
		{
			name: "aggregates and dedupes",
			configs: []*query.OIDCAppLinkConfig{
				{IOSTeamID: "TEAM1", IOSBundleID: "com.one"},
				{IOSTeamID: "TEAM2", IOSBundleID: "com.two"},
				{IOSTeamID: "TEAM1", IOSBundleID: "com.one"},
				{AndroidPackageName: "com.android", AndroidSHA256CertFingerprints: []string{"AA"}},
			},
			wantApps: []string{"TEAM1.com.one", "TEAM2.com.two"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := buildAppleAppSiteAssociation(tc.configs)
			assert.Equal(t, tc.wantApps, got.WebCredentials.Apps)
		})
	}
}

func TestBuildAssetLinks(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name    string
		configs []*query.OIDCAppLinkConfig
		want    []assetLink
	}{
		{
			name:    "empty",
			configs: nil,
			want:    []assetLink{},
		},
		{
			name: "skips incomplete android",
			configs: []*query.OIDCAppLinkConfig{
				{AndroidPackageName: "com.example"},
				{AndroidSHA256CertFingerprints: []string{"AA"}},
			},
			want: []assetLink{},
		},
		{
			name: "aggregates and normalizes fingerprints",
			configs: []*query.OIDCAppLinkConfig{
				{
					IOSTeamID:                     "TEAM",
					IOSBundleID:                   "com.ios",
					AndroidPackageName:            "com.one",
					AndroidSHA256CertFingerprints: []string{"aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899"},
				},
				{
					AndroidPackageName: "com.two",
					AndroidSHA256CertFingerprints: []string{
						"cc:dd:ee:ff:00:11:22:33:44:55:66:77:88:99:aa:bb:cc:dd:ee:ff:00:11:22:33:44:55:66:77:88:99:aa:bb",
						"eeff00112233445566778899aabbccddeeff00112233445566778899aabbccdd",
					},
				},
			},
			want: []assetLink{
				{
					Relation: []string{"delegate_permission/common.handle_all_urls", "delegate_permission/common.get_login_creds"},
					Target: assetLinkTarget{
						Namespace:   "android_app",
						PackageName: "com.one",
						SHA256CertFingerprints: []string{
							"AA:BB:CC:DD:EE:FF:00:11:22:33:44:55:66:77:88:99:AA:BB:CC:DD:EE:FF:00:11:22:33:44:55:66:77:88:99",
						},
					},
				},
				{
					Relation: []string{"delegate_permission/common.handle_all_urls", "delegate_permission/common.get_login_creds"},
					Target: assetLinkTarget{
						Namespace:   "android_app",
						PackageName: "com.two",
						SHA256CertFingerprints: []string{
							"CC:DD:EE:FF:00:11:22:33:44:55:66:77:88:99:AA:BB:CC:DD:EE:FF:00:11:22:33:44:55:66:77:88:99:AA:BB",
							"EE:FF:00:11:22:33:44:55:66:77:88:99:AA:BB:CC:DD:EE:FF:00:11:22:33:44:55:66:77:88:99:AA:BB:CC:DD",
						},
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, buildAssetLinks(t.Context(), tc.configs))
		})
	}
}

func TestNormalizeSHA256Fingerprint(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "already canonical",
			in:   "AA:BB:CC:DD:EE:FF:00:11:22:33:44:55:66:77:88:99:AA:BB:CC:DD:EE:FF:00:11:22:33:44:55:66:77:88:99",
			want: "AA:BB:CC:DD:EE:FF:00:11:22:33:44:55:66:77:88:99:AA:BB:CC:DD:EE:FF:00:11:22:33:44:55:66:77:88:99",
		},
		{
			name: "lowercase with colons",
			in:   "aa:bb:cc:dd:ee:ff:00:11:22:33:44:55:66:77:88:99:aa:bb:cc:dd:ee:ff:00:11:22:33:44:55:66:77:88:99",
			want: "AA:BB:CC:DD:EE:FF:00:11:22:33:44:55:66:77:88:99:AA:BB:CC:DD:EE:FF:00:11:22:33:44:55:66:77:88:99",
		},
		{
			name: "lowercase without separators",
			in:   "aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899",
			want: "AA:BB:CC:DD:EE:FF:00:11:22:33:44:55:66:77:88:99:AA:BB:CC:DD:EE:FF:00:11:22:33:44:55:66:77:88:99",
		},
		{
			name: "invalid continuous hex uppercased",
			in:   "aabb",
			want: "AABB",
		},
		{
			name: "colon form uppercased without reparse",
			in:   "aa:bb",
			want: "AA:BB",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, normalizeSHA256Fingerprint(t.Context(), tc.in))
		})
	}
}

type stubAppLinkQuerier struct {
	configs []*query.OIDCAppLinkConfig
	err     error
}

func (s stubAppLinkQuerier) SearchOIDCAppLinkConfigs(_ context.Context) ([]*query.OIDCAppLinkConfig, error) {
	return s.configs, s.err
}

func TestHandlerCacheControl(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name       string
		maxAge     time.Duration
		wantHeader string
	}{
		{
			name:       "default max-age",
			maxAge:     5 * time.Minute,
			wantHeader: "max-age=300, must-revalidate",
		},
		{
			name:       "no-store when zero",
			maxAge:     0,
			wantHeader: "no-store",
		},
		{
			name:       "no-store when negative",
			maxAge:     -time.Minute,
			wantHeader: "no-store",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			h := NewHandler(stubAppLinkQuerier{}, Config{AppLinksCacheControlMaxAge: tc.maxAge})

			for _, path := range []string{AppleAppSiteAssociationPath, AssetLinksPath} {
				req := httptest.NewRequest(http.MethodGet, path, nil)
				rec := httptest.NewRecorder()
				h.ServeHTTP(rec, req)
				require.Equal(t, http.StatusOK, rec.Code)
				cacheControl := rec.Header().Get(http_util.CacheControl)
				assert.Equal(t, tc.wantHeader, cacheControl, path)
				assert.NotContains(t, cacheControl, "stale-while-revalidate")
				// "public" must never be set: it opts the response into shared caches
				// (e.g. a CDN), whose independent edges can pin diverging copies.
				assert.NotContains(t, cacheControl, "public", path)
			}
		})
	}
}

func TestHandlerHeadOmitsBody(t *testing.T) {
	t.Parallel()

	h := NewHandler(stubAppLinkQuerier{}, Config{AppLinksCacheControlMaxAge: 5 * time.Minute})
	for _, path := range []string{AppleAppSiteAssociationPath, AssetLinksPath} {
		t.Run(path, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(http.MethodHead, path, nil)
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			require.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
			assert.Equal(t, "max-age=300, must-revalidate", rec.Header().Get(http_util.CacheControl))
			assert.Empty(t, rec.Body.Bytes())
		})
	}
}
