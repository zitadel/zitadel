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

	config := &Config{
		DefaultBasePath: "/basepath",
	}

	baseCustomURI, err := url.Parse("https://custom")
	require.Nil(t, err)

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
			name: "LoginV2 required, custom BaseURI",
			inputCtx: http.WithDomainContext(
				authz.NewMockContext("instance1", "org1", "user1",
					authz.WithMockFeatures(feature.Features{LoginV2: feature.LoginV2{Required: true, BaseURI: baseCustomURI}}),
				),
				&http.DomainCtx{Protocol: "https", PublicHost: "origin"},
			),

			expected: "https://custom/basepath",
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
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := config.defaultBaseURL(tc.inputCtx)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestConfig_DefaultEmailCodeURLTemplate(t *testing.T) {
	t.Parallel()

	// Given
	ctx := http.WithDomainContext(
		authz.NewMockContext("instance1", "org1", "user1",
			authz.WithMockFeatures(feature.Features{LoginV2: feature.LoginV2{Required: true, BaseURI: &url.URL{}}}),
		),
		&http.DomainCtx{Protocol: "https", PublicHost: "origin"},
	)
	c := &Config{
		DefaultBasePath: "/basepath",
		DefaultPaths: &DefaultPaths{
			EmailCodePath: "/email-code-path",
		},
	}

	// Test
	res := c.DefaultEmailCodeURLTemplate(ctx)

	// Verify
	assert.Equal(t, "https://origin/basepath/email-code-path", res)
}

func TestConfig_DefaultPasswordSetURLTemplate(t *testing.T) {
	t.Parallel()

	// Given
	ctx := http.WithDomainContext(
		authz.NewMockContext("instance1", "org1", "user1",
			authz.WithMockFeatures(feature.Features{LoginV2: feature.LoginV2{Required: true, BaseURI: &url.URL{}}}),
		),
		&http.DomainCtx{Protocol: "https", PublicHost: "origin"},
	)
	c := &Config{
		DefaultBasePath: "/basepath",
		DefaultPaths: &DefaultPaths{
			PasswordSetPath: "/password-set-path",
		},
	}

	// Test
	res := c.DefaultPasswordSetURLTemplate(ctx)

	// Verify
	assert.Equal(t, "https://origin/basepath/password-set-path", res)
}
