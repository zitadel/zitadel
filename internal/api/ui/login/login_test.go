package login

import (
	"context"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/feature"
)

func TestConfig_defaultBaseURL(t *testing.T) {
	t.Parallel()

	config := &DefaultPaths{BasePath: &url.URL{Path: "/basepath"}}

	baseCustomURI, err := url.Parse("https://custom")
	require.NoError(t, err)

	tt := []struct {
		name     string
		inputCtx context.Context
		http.DomainCtx
		expected string
	}{
		{
			name:     "LoginV2 not required",
			inputCtx: authz.NewMockContext("instance1", "org1", "user1"),
			expected: "",
		},
		{
			name: "LoginV2 required, no custom BaseURI",
			inputCtx: http.WithDomainContext(
				authz.NewMockContext("instance1", "org1", "user1",
					authz.WithMockFeatures(feature.Features{LoginV2: feature.LoginV2{Required: true}}),
				),
				&http.DomainCtx{Protocol: "https", PublicHost: "origin"},
			),
			expected: "https://origin/basepath",
		},
		{
			name: "LoginV2 required, absolute custom BaseURI",
			inputCtx: http.WithDomainContext(
				authz.NewMockContext("instance1", "org1", "user1",
					authz.WithMockFeatures(feature.Features{LoginV2: feature.LoginV2{Required: true, BaseURI: baseCustomURI}}),
				),
				&http.DomainCtx{Protocol: "https", PublicHost: "origin"},
			),
			expected: "https://custom",
		},
		{
			name: "LoginV2 required, custom BaseURI empty string",
			inputCtx: http.WithDomainContext(
				authz.NewMockContext("instance1", "org1", "user1",
					authz.WithMockFeatures(feature.Features{LoginV2: feature.LoginV2{Required: true, BaseURI: &url.URL{}}}),
				),
				&http.DomainCtx{Protocol: "https", PublicHost: "origin"},
			),
			expected: "https://origin/basepath",
		},
		{
			name: "LoginV2 required, relative custom BaseURI",
			inputCtx: http.WithDomainContext(
				authz.NewMockContext("instance1", "org1", "user1",
					authz.WithMockFeatures(feature.Features{LoginV2: feature.LoginV2{Required: true, BaseURI: &url.URL{Path: "/custom"}}}),
				),
				&http.DomainCtx{Protocol: "https", PublicHost: "origin"},
			),
			expected: "https://origin/custom",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := config.defaultBaseURL(tc.inputCtx)
			if tc.expected == "" {
				assert.Nil(t, result)
				return
			}
			assert.Equal(t, tc.expected, result.String())
		})
	}
}

func TestConfig_DefaultEmailCodeURLTemplate(t *testing.T) {
	t.Parallel()

	tt := []struct {
		testName                 string
		inputCtx                 context.Context
		expectedEmailURLTemplate string
	}{
		{
			testName: "when base path is empty should return empty email url template",
			inputCtx: http.WithDomainContext(
				authz.NewMockContext("instance1", "org1", "user1",
					authz.WithMockFeatures(feature.Features{LoginV2: feature.LoginV2{Required: false, BaseURI: &url.URL{}}}),
				),
				&http.DomainCtx{Protocol: "https", PublicHost: "origin"},
			),
			expectedEmailURLTemplate: "",
		},
		{
			testName: "when base path is not empty should return expected url template",
			inputCtx: http.WithDomainContext(
				authz.NewMockContext("instance1", "org1", "user1",
					authz.WithMockFeatures(feature.Features{LoginV2: feature.LoginV2{Required: true, BaseURI: &url.URL{}}}),
				),
				&http.DomainCtx{Protocol: "https", PublicHost: "origin"},
			),
			expectedEmailURLTemplate: "https://origin/basepath/email-code-path",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// Given
			c := &DefaultPaths{
				BasePath:      &url.URL{Path: "/basepath"},
				EmailCodePath: &url.URL{Path: "/email-code-path"},
			}

			// Test
			res := c.DefaultEmailCodeURLTemplate(tc.inputCtx)

			// Verify
			assert.Equal(t, tc.expectedEmailURLTemplate, res)
		})
	}
}

func TestConfig_DefaultPasswordSetURLTemplate(t *testing.T) {
	t.Parallel()

	tt := []struct {
		testName                 string
		inputCtx                 context.Context
		expectedEmailURLTemplate string
	}{
		{
			testName: "when base path is empty should return empty email url template",
			inputCtx: http.WithDomainContext(
				authz.NewMockContext("instance1", "org1", "user1",
					authz.WithMockFeatures(feature.Features{LoginV2: feature.LoginV2{Required: false, BaseURI: &url.URL{}}}),
				),
				&http.DomainCtx{Protocol: "https", PublicHost: "origin"},
			),
			expectedEmailURLTemplate: "",
		},
		{
			testName: "when base path is not empty should return expected url template",
			inputCtx: http.WithDomainContext(
				authz.NewMockContext("instance1", "org1", "user1",
					authz.WithMockFeatures(feature.Features{LoginV2: feature.LoginV2{Required: true, BaseURI: &url.URL{}}}),
				),
				&http.DomainCtx{Protocol: "https", PublicHost: "origin"},
			),
			expectedEmailURLTemplate: "https://origin/basepath/password-set-path",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// Given
			c := &DefaultPaths{
				BasePath:        &url.URL{Path: "/basepath"},
				PasswordSetPath: &url.URL{Path: "/password-set-path"},
			}

			// Test
			res := c.DefaultPasswordSetURLTemplate(tc.inputCtx)

			// Verify
			assert.Equal(t, tc.expectedEmailURLTemplate, res)
		})
	}
}

func TestConfig_DefaultPasskeySetURLTemplate(t *testing.T) {
	t.Parallel()

	tt := []struct {
		testName                 string
		inputCtx                 context.Context
		expectedEmailURLTemplate string
	}{
		{
			testName: "when base path is empty should return empty email url template",
			inputCtx: http.WithDomainContext(
				authz.NewMockContext("instance1", "org1", "user1",
					authz.WithMockFeatures(feature.Features{LoginV2: feature.LoginV2{Required: false, BaseURI: &url.URL{}}}),
				),
				&http.DomainCtx{Protocol: "https", PublicHost: "origin"},
			),
			expectedEmailURLTemplate: "",
		},
		{
			testName: "when base path is not empty should return expected url template",
			inputCtx: http.WithDomainContext(
				authz.NewMockContext("instance1", "org1", "user1",
					authz.WithMockFeatures(feature.Features{LoginV2: feature.LoginV2{Required: true, BaseURI: &url.URL{}}}),
				),
				&http.DomainCtx{Protocol: "https", PublicHost: "origin"},
			),
			expectedEmailURLTemplate: "https://origin/basepath/passkey-set-path",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// Given
			c := &DefaultPaths{
				BasePath:       &url.URL{Path: "/basepath"},
				PasskeySetPath: &url.URL{Path: "/passkey-set-path"},
			}

			// Test
			res := c.DefaultPasskeySetURLTemplate(tc.inputCtx)

			// Verify
			assert.Equal(t, tc.expectedEmailURLTemplate, res)
		})
	}
}

func TestConfig_DefaultDomainClaimedURLTemplate(t *testing.T) {
	t.Parallel()

	tt := []struct {
		testName                 string
		inputCtx                 context.Context
		expectedEmailURLTemplate string
	}{
		{
			testName: "when base path is empty should return empty email url template",
			inputCtx: http.WithDomainContext(
				authz.NewMockContext("instance1", "org1", "user1",
					authz.WithMockFeatures(feature.Features{LoginV2: feature.LoginV2{Required: false, BaseURI: &url.URL{}}}),
				),
				&http.DomainCtx{Protocol: "https", PublicHost: "origin"},
			),
			expectedEmailURLTemplate: "",
		},
		{
			testName: "when base path is not empty should return expected url template",
			inputCtx: http.WithDomainContext(
				authz.NewMockContext("instance1", "org1", "user1",
					authz.WithMockFeatures(feature.Features{LoginV2: feature.LoginV2{Required: true, BaseURI: &url.URL{}}}),
				),
				&http.DomainCtx{Protocol: "https", PublicHost: "origin"},
			),
			expectedEmailURLTemplate: "https://origin/basepath/loginname",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// Given
			c := &DefaultPaths{
				BasePath:          &url.URL{Path: "/basepath"},
				DomainClaimedPath: &url.URL{Path: "/loginname"},
			}

			// Test
			res := c.DefaultDomainClaimedURLTemplate(tc.inputCtx)

			// Verify
			assert.Equal(t, tc.expectedEmailURLTemplate, res)
		})
	}
}

func Test_mergeURLs(t *testing.T) {
	type args struct {
		base *url.URL
		path *url.URL
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "both base and path are nil, return empty string",
			args: args{
				base: nil,
				path: nil,
			},
			want: "",
		},
		{
			name: "base is nil, path is set, return path",
			args: args{
				base: nil,
				path: &url.URL{Path: "/path"},
			},
			want: "/path",
		},
		{
			name: "base is set, path is nil, return base",
			args: args{
				base: &url.URL{Path: "/base"},
				path: nil,
			},
			want: "/base",
		},
		{
			name: "both base and path are set, return merged path",
			args: args{
				base: &url.URL{Path: "/base"},
				path: &url.URL{Path: "/path"},
			},
			want: "/base/path",
		},
		{
			name: "both base and path are set incl. placeholders, return merged path",
			args: args{
				base: &url.URL{Path: "/base", RawQuery: "query={{.Placeholder}}"},
				path: &url.URL{Path: "/path", RawQuery: "query={{.Placeholder1}}&query2={{.Placeholder2}}"},
			},
			want: "/base/path?query={{.Placeholder}}&query={{.Placeholder1}}&query2={{.Placeholder2}}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mergeURLs(tt.args.base, tt.args.path)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_mergeQueries(t *testing.T) {
	type args struct {
		base url.Values
		path url.Values
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "both base and path are empty, return empty string",
			args: args{
				base: map[string][]string{},
				path: map[string][]string{},
			},
			want: "",
		},
		{
			name: "base is empty, path is set, return path",
			args: args{
				base: map[string][]string{},
				path: map[string][]string{
					"query": {"{{.Placeholder1}}", "{{.Placeholder2}}"},
				},
			},
			want: "query={{.Placeholder1}}&query={{.Placeholder2}}",
		},
		{
			name: "base is set, path is empty, return base",
			args: args{
				base: map[string][]string{
					"query": {"{{.Placeholder}}"},
				},
				path: map[string][]string{},
			},
			want: "query={{.Placeholder}}",
		},
		{
			name: "both base and path are set, return merged query",
			args: args{
				base: map[string][]string{
					"query":  {"{{.Placeholder}}"},
					"query2": {"{{.Placeholder}}"},
				},
				path: map[string][]string{
					"query":  {"{{.Placeholder1}}"},
					"query2": {"{{.Placeholder}}", "{{.Placeholder2}}"},
				},
			},
			want: "query={{.Placeholder}}&query={{.Placeholder1}}&query2={{.Placeholder}}&query2={{.Placeholder2}}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, mergeQueries(tt.args.base, tt.args.path), "mergeQueries(%v, %v)", tt.args.base, tt.args.path)
		})
	}
}
