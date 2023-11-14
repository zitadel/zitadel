//go:build integration

package admin_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/pkg/grpc/admin"
)

func TestServer_Limits_DisallowPublicOrgRegistration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	domain, _, iamOwnerCtx := Tester.UseIsolatedInstance(ctx, SystemCTX)
	regOrgUrl, err := url.Parse("http://" + domain + ":8080/ui/login/register/org")
	require.NoError(t, err)
	jar, err := cookiejar.New(nil)
	require.NoError(t, err)
	browserSession := &http.Client{Jar: jar}
	// Default should be allowed
	awaitGetResponse(t, iamOwnerCtx, browserSession, regOrgUrl, http.StatusOK)
	_, err = Tester.Client.Admin.SetInstanceLimits(iamOwnerCtx, &admin.SetInstanceLimitsRequest{DisallowPublicOrgRegistration: true})
	require.NoError(t, err)
	awaitGetResponse(t, iamOwnerCtx, browserSession, regOrgUrl, http.StatusNotFound)
	_, err = Tester.Client.Admin.ResetInstanceLimits(iamOwnerCtx, &admin.ResetInstanceLimitsRequest{})
	require.NoError(t, err)
	awaitGetResponse(t, iamOwnerCtx, browserSession, regOrgUrl, http.StatusOK)
}

func awaitGetResponse(t *testing.T, ctx context.Context, parsedURL *url.URL, expectCode int) {
	await(t, ctx, func() bool {
		resp, err := http.Get(parsedURL.String())
		require.NoError(t, err)
		return resp.StatusCode == expectCode
	})
}

func awaitPostSuccess(t *testing.T, ctx context.Context, parsedURL *url.URL) {
	await(t, ctx, func() bool {
		resp, err := http.PostForm(parsedURL.String(), url.Values{})
		require.NoError(t, err)
		return resp.StatusCode == http.StatusOK
	})
}

func await(t *testing.T, ctx context.Context, cb func() bool) {
	deadline, ok := ctx.Deadline()
	require.True(t, ok, "context must have deadline")
	require.Eventuallyf(
		t,
		func() bool {
			defer func() {
				require.Nil(t, recover(), "panic in awaitHTTPStatusCode")
			}()
			return cb()
		},
		time.Until(deadline),
		100*time.Millisecond,
		"org registration response not received",
	)
}

func awaitHTTPStatusCode(t *testing.T, ctx context.Context, client *http.Client, method string, url string, expectCode int) {
	deadline, ok := ctx.Deadline()
	require.True(t, ok, "context must have deadline")
	require.Eventuallyf(
		t,
		func() bool {
			defer func() {
				require.Nil(t, recover(), "panic in awaitHTTPStatusCode")
			}()
			req, err := http.NewRequest(method, url, nil)
			require.NoError(t, err)
			resp, err := client.Do(req)
			require.NoError(t, err)
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			fmt.Println(string(body))
			return resp.StatusCode == expectCode
		},
		time.Until(deadline),
		100*time.Millisecond,
		"org registration response not received",
	)
}
