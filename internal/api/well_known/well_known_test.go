package well_known

import (
	"testing"

	"github.com/stretchr/testify/assert"

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
			name: "aggregates entries",
			configs: []*query.OIDCAppLinkConfig{
				{
					IOSTeamID:                     "TEAM",
					IOSBundleID:                   "com.ios",
					AndroidPackageName:            "com.one",
					AndroidSHA256CertFingerprints: []string{"AA:BB"},
				},
				{
					AndroidPackageName:            "com.two",
					AndroidSHA256CertFingerprints: []string{"CC:DD", "EE:FF"},
				},
			},
			want: []assetLink{
				{
					Relation: []string{"delegate_permission/common.get_login_creds"},
					Target: assetLinkTarget{
						Namespace:              "android_app",
						PackageName:            "com.one",
						SHA256CertFingerprints: []string{"AA:BB"},
					},
				},
				{
					Relation: []string{"delegate_permission/common.get_login_creds"},
					Target: assetLinkTarget{
						Namespace:              "android_app",
						PackageName:            "com.two",
						SHA256CertFingerprints: []string{"CC:DD", "EE:FF"},
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, buildAssetLinks(tc.configs))
		})
	}
}
