//go:build integration

package admin_test

import (
	"bytes"
	"context"
	"github.com/muhlemmer/gu"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/pkg/grpc/admin"
)

func TestServer_Restrictions_DisallowPublicOrgRegistration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	domain, _, iamOwnerCtx := Tester.UseIsolatedInstance(ctx, SystemCTX)
	regOrgUrl, err := url.Parse("http://" + domain + ":8080/ui/login/register/org")
	require.NoError(t, err)
	// The CSRF cookie must be sent with every request.
	// We can simulate a browser session using a cookie jar.
	jar, err := cookiejar.New(nil)
	require.NoError(t, err)
	browserSession := &http.Client{Jar: jar}
	// Default should be allowed
	csrfToken := awaitAllowed(t, iamOwnerCtx, browserSession, regOrgUrl)
	_, err = Tester.Client.Admin.SetRestrictions(iamOwnerCtx, &admin.SetRestrictionsRequest{DisallowPublicOrgRegistration: gu.Ptr(true)})
	require.NoError(t, err)
	awaitDisallowed(t, iamOwnerCtx, browserSession, regOrgUrl, csrfToken)
	_, err = Tester.Client.Admin.SetRestrictions(iamOwnerCtx, &admin.SetRestrictionsRequest{DisallowPublicOrgRegistration: gu.Ptr(false)})
	require.NoError(t, err)
	awaitAllowed(t, iamOwnerCtx, browserSession, regOrgUrl)
}

// awaitAllowed doesn't accept a CSRF token, as we expected it to always produce a new one
func awaitAllowed(t *testing.T, ctx context.Context, client *http.Client, parsedURL *url.URL) string {
	csrfToken := awaitGetResponse(t, ctx, client, parsedURL, http.StatusOK)
	awaitPostFormResponse(t, ctx, client, parsedURL, http.StatusOK, csrfToken)
	restrictions, err := Tester.Client.Admin.GetRestrictions(ctx, &admin.GetRestrictionsRequest{})
	require.NoError(t, err)
	require.False(t, restrictions.DisallowPublicOrgRegistration)
	return csrfToken
}

// awaitDisallowed accepts an old CSRF token, as we don't expect to get a CSRF token from the GET request anymore
func awaitDisallowed(t *testing.T, ctx context.Context, client *http.Client, parsedURL *url.URL, reuseOldCSRFToken string) {
	awaitGetResponse(t, ctx, client, parsedURL, http.StatusNotFound)
	awaitPostFormResponse(t, ctx, client, parsedURL, http.StatusConflict, reuseOldCSRFToken)
	restrictions, err := Tester.Client.Admin.GetRestrictions(ctx, &admin.GetRestrictionsRequest{})
	require.NoError(t, err)
	require.True(t, restrictions.DisallowPublicOrgRegistration)
}

// awaitGetResponse cuts the CSRF token from the response body if it exists
func awaitGetResponse(t *testing.T, ctx context.Context, client *http.Client, parsedURL *url.URL, expectCode int) string {
	var csrfToken []byte
	await(t, ctx, func() bool {
		resp, err := client.Get(parsedURL.String())
		require.NoError(t, err)
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		searchField := `<input type="hidden" name="gorilla.csrf.Token" value="`
		_, after, hasCsrfToken := bytes.Cut(body, []byte(searchField))
		if hasCsrfToken {
			csrfToken, _, _ = bytes.Cut(after, []byte(`">`))
		}
		return resp.StatusCode == expectCode
	})
	return string(csrfToken)
}

// awaitPostFormResponse needs a valid CSRF token to make it to the actual endpoint implementation and get the expected status code
func awaitPostFormResponse(t *testing.T, ctx context.Context, client *http.Client, parsedURL *url.URL, expectCode int, csrfToken string) {
	await(t, ctx, func() bool {
		resp, err := client.PostForm(parsedURL.String(), url.Values{
			"gorilla.csrf.Token": {csrfToken},
		})
		require.NoError(t, err)
		return resp.StatusCode == expectCode

	})
}

func await(t *testing.T, ctx context.Context, cb func() bool) {
	deadline, ok := ctx.Deadline()
	require.True(t, ok, "context must have deadline")
	require.Eventuallyf(
		t,
		func() bool {
			defer func() {
				require.Nil(t, recover(), "panic in await callback")
			}()
			return cb()
		},
		time.Until(deadline),
		100*time.Millisecond,
		"awaiting successful callback failed",
	)
}
