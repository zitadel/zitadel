//go:build integration

package admin_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/admin"
)

func TestServer_Restrictions_DisallowPublicOrgRegistration(t *testing.T) {
	instance := integration.NewInstance(CTX)
	regOrgUrl, err := url.Parse("http://" + instance.Domain + ":8080/ui/login/register/org")
	require.NoError(t, err)
	// The CSRF cookie must be sent with every request.
	// We can simulate a browser session using a cookie jar.
	jar, err := cookiejar.New(nil)
	require.NoError(t, err)
	browserSession := &http.Client{Jar: jar}
	var csrfToken string
	iamOwnerCtx := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)

	t.Run("public org registration is allowed by default", func(tt *testing.T) {
		csrfToken = awaitPubOrgRegAllowed(tt, iamOwnerCtx, instance.Client, browserSession, regOrgUrl)
	})
	t.Run("disallowing public org registration disables the endpoints", func(tt *testing.T) {
		_, err = instance.Client.Admin.SetRestrictions(iamOwnerCtx, &admin.SetRestrictionsRequest{DisallowPublicOrgRegistration: gu.Ptr(true)})
		require.NoError(tt, err)
		awaitPubOrgRegDisallowed(tt, iamOwnerCtx, instance.Client, browserSession, regOrgUrl, csrfToken)
	})
	t.Run("allowing public org registration again re-enables the endpoints", func(tt *testing.T) {
		_, err = instance.Client.Admin.SetRestrictions(iamOwnerCtx, &admin.SetRestrictionsRequest{DisallowPublicOrgRegistration: gu.Ptr(false)})
		require.NoError(tt, err)
		awaitPubOrgRegAllowed(tt, iamOwnerCtx, instance.Client, browserSession, regOrgUrl)
	})
}

// awaitPubOrgRegAllowed doesn't accept a CSRF token, as we expected it to always produce a new one
func awaitPubOrgRegAllowed(t *testing.T, ctx context.Context, cc *integration.Client, client *http.Client, parsedURL *url.URL) string {
	csrfToken := awaitGetSSRGetResponse(t, ctx, client, parsedURL, http.StatusOK)
	awaitPostFormResponse(t, ctx, client, parsedURL, http.StatusOK, csrfToken)
	restrictions, err := cc.Admin.GetRestrictions(ctx, &admin.GetRestrictionsRequest{})
	require.NoError(t, err)
	require.False(t, restrictions.DisallowPublicOrgRegistration)
	return csrfToken
}

// awaitPubOrgRegDisallowed accepts an old CSRF token, as we don't expect to get a CSRF token from the GET request anymore
func awaitPubOrgRegDisallowed(t *testing.T, ctx context.Context, cc *integration.Client, client *http.Client, parsedURL *url.URL, reuseOldCSRFToken string) {
	awaitGetSSRGetResponse(t, ctx, client, parsedURL, http.StatusNotFound)
	awaitPostFormResponse(t, ctx, client, parsedURL, http.StatusConflict, reuseOldCSRFToken)
	restrictions, err := cc.Admin.GetRestrictions(ctx, &admin.GetRestrictionsRequest{})
	require.NoError(t, err)
	require.True(t, restrictions.DisallowPublicOrgRegistration)
}

// awaitGetSSRGetResponse cuts the CSRF token from the response body if it exists
func awaitGetSSRGetResponse(t *testing.T, ctx context.Context, client *http.Client, parsedURL *url.URL, expectCode int) string {
	var csrfToken []byte
	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
	require.EventuallyWithT(t,
		func(tt *assert.CollectT) {
			resp, err := client.Get(parsedURL.String())
			require.NoError(tt, err)
			body, err := io.ReadAll(resp.Body)
			require.NoError(tt, err)
			searchField := `<input type="hidden" name="gorilla.csrf.Token" value="`
			_, after, hasCsrfToken := bytes.Cut(body, []byte(searchField))
			if hasCsrfToken {
				csrfToken, _, _ = bytes.Cut(after, []byte(`">`))
			}
			assert.Equal(tt, resp.StatusCode, expectCode)
		}, retryDuration, tick, "awaiting successful get SSR get response failed",
	)
	return string(csrfToken)
}

// awaitPostFormResponse needs a valid CSRF token to make it to the actual endpoint implementation and get the expected status code
func awaitPostFormResponse(t *testing.T, ctx context.Context, client *http.Client, parsedURL *url.URL, expectCode int, csrfToken string) {
	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
	require.EventuallyWithT(t,
		func(tt *assert.CollectT) {
			resp, err := client.PostForm(parsedURL.String(), url.Values{
				"gorilla.csrf.Token": {csrfToken},
			})
			require.NoError(tt, err)
			assert.Equal(tt, resp.StatusCode, expectCode)
		}, retryDuration, tick, "awaiting successful Post Form failed",
	)
}
