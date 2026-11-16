package test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v3/pkg/oidc"

	"github.com/zitadel/zitadel/internal/idp"
	"github.com/zitadel/zitadel/internal/idp/providers/oauth"
)

// AssertRedirectAuthWithPKCE asserts that the actual auth URL contains a S256 code challenge
// matching the code verifier of the persistent parameters and otherwise equals the expected one,
// since the code challenge is random.
func AssertRedirectAuthWithPKCE(t *testing.T, expected, actual idp.Auth, persistentParameters map[string]any) {
	t.Helper()

	expectedRedirect, ok := expected.(*idp.RedirectAuth)
	require.True(t, ok)
	actualRedirect, ok := actual.(*idp.RedirectAuth)
	require.True(t, ok)

	expectedURL, err := url.Parse(expectedRedirect.RedirectURL)
	require.NoError(t, err)
	actualURL, err := url.Parse(actualRedirect.RedirectURL)
	require.NoError(t, err)

	verifier, _ := persistentParameters[oauth.CodeVerifier].(string)
	require.NotEmpty(t, verifier)

	actualQuery := actualURL.Query()
	assert.Equal(t, oidc.NewSHACodeChallenge(verifier), actualQuery.Get("code_challenge"))
	assert.Equal(t, "S256", actualQuery.Get("code_challenge_method"))

	actualQuery.Del("code_challenge")
	actualQuery.Del("code_challenge_method")

	assert.Equal(t, expectedURL.Scheme, actualURL.Scheme)
	assert.Equal(t, expectedURL.Host, actualURL.Host)
	assert.Equal(t, expectedURL.Path, actualURL.Path)
	assert.Equal(t, expectedURL.Query(), actualQuery)
}
