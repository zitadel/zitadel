//go:build integration

package system_test

import (
	"fmt"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/admin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/zitadel/pkg/grpc/system"
)

func TestServer_Limits_Block(t *testing.T) {
	domain, instanceID, iamOwnerCtx := Tester.UseIsolatedInstance(t, CTX, SystemCTX)
	type test struct {
		name     string
		testHttp func(tt assert.TestingT) (req *http.Request, err error, assertResponse func(response *http.Response, withBlocking bool))
		testGrpc func(tt assert.TestingT, withBlocking bool)
	}
	tests := []test{{
		name: "public API",
		testGrpc: func(tt assert.TestingT, withBlocking bool) {
			_, err := Tester.Client.Admin.Healthz(CTX, &admin.HealthzRequest{})
			assertGrpcError(tt, err, withBlocking)
		},
		testHttp: func(tt assert.TestingT) (*http.Request, error, func(*http.Response, bool)) {
			req, err := http.NewRequestWithContext(
				CTX,
				"GET",
				fmt.Sprintf("http://%s/admin/v1/healthz", net.JoinHostPort(domain, "8080")),
				nil,
			)
			return req, err, func(response *http.Response, withBlocking bool) {
				assertLimitResponse(tt, response, withBlocking)
				assertSetLimitingCookie(tt, response, withBlocking)
			}
		},
	}, {
		name: "mutating API",
		testGrpc: func(tt assert.TestingT, withBlocking bool) {
			randomGrpcIdpName := randomString("idp-grpc", 5)
			_, err := Tester.Client.Admin.AddGitHubProvider(iamOwnerCtx, &admin.AddGitHubProviderRequest{
				Name:         randomGrpcIdpName,
				ClientId:     "client-id",
				ClientSecret: "client-secret",
			})
			assertGrpcError(tt, err, withBlocking)
			//nolint:contextcheck
			idpExists := idpExistsCondition(tt, instanceID, randomGrpcIdpName)
			if withBlocking {
				// We ensure that the idp really is not created
				assert.Neverf(tt, idpExists, 5*time.Second, 1*time.Second, "idp should never be created")
			} else {
				assert.Eventuallyf(tt, idpExists, 5*time.Second, 1*time.Second, "idp should be created")
			}
		},
		testHttp: func(tt assert.TestingT) (*http.Request, error, func(*http.Response, bool)) {
			randomHttpIdpName := randomString("idp-http", 5)
			req, err := http.NewRequestWithContext(
				CTX,
				"POST",
				fmt.Sprintf("http://%s/admin/v1/idps/github", net.JoinHostPort(domain, "8080")),
				strings.NewReader(`{
	"name": "`+randomHttpIdpName+`",
	"clientId": "client-id",
	"clientSecret": "client-secret"
}`),
			)
			if err != nil {
				return nil, err, nil
			}
			req.Header.Set("Authorization", Tester.BearerToken(iamOwnerCtx))
			return req, nil, func(response *http.Response, withBlocking bool) {
				assertLimitResponse(tt, response, withBlocking)
				assertSetLimitingCookie(tt, response, withBlocking)
			}
		},
	}, {
		name: "discovery",
		testHttp: func(tt assert.TestingT) (*http.Request, error, func(*http.Response, bool)) {
			req, err := http.NewRequestWithContext(
				CTX,
				"GET",
				fmt.Sprintf("http://%s/.well-known/openid-configuration", net.JoinHostPort(domain, "8080")),
				nil,
			)
			return req, err, func(response *http.Response, withBlocking bool) {
				assertLimitResponse(tt, response, withBlocking)
				assertSetLimitingCookie(tt, response, withBlocking)
			}
		},
	}, {
		name: "login",
		testHttp: func(tt assert.TestingT) (*http.Request, error, func(*http.Response, bool)) {
			req, err := http.NewRequestWithContext(
				CTX,
				"GET",
				fmt.Sprintf("http://%s/ui/login/login/externalidp/callback", net.JoinHostPort(domain, "8080")),
				nil,
			)
			return req, err, func(response *http.Response, withBlocking bool) {
				// the login paths should return a redirect if the instance is blocked
				if withBlocking {
					assert.GreaterOrEqual(tt, response.StatusCode, http.StatusMultipleChoices)
					assert.LessOrEqual(tt, response.StatusCode, http.StatusPermanentRedirect)
				} else {
					assertLimitResponse(tt, response, false)
				}
				assertSetLimitingCookie(tt, response, withBlocking)
			}
		},
	}, {
		name: "console",
		testHttp: func(tt assert.TestingT) (*http.Request, error, func(*http.Response, bool)) {
			req, err := http.NewRequestWithContext(
				CTX,
				"GET",
				fmt.Sprintf("http://%s/ui/console/", net.JoinHostPort(domain, "8080")),
				nil,
			)
			return req, err, func(response *http.Response, withBlocking bool) {
				// the console is not blocked so we can render a link to an instance management portal.
				// A CDN can cache these assets easily
				// We also don't care about a cookie because the environment.json already takes care of that.
				assertLimitResponse(tt, response, false)
			}
		},
	}, {
		name: "environment.json",
		testHttp: func(tt assert.TestingT) (*http.Request, error, func(*http.Response, bool)) {
			req, err := http.NewRequestWithContext(
				CTX,
				"GET",
				fmt.Sprintf("http://%s/ui/console/assets/environment.json", net.JoinHostPort(domain, "8080")),
				nil,
			)
			return req, err, func(response *http.Response, withBlocking bool) {
				// the environment.json should always return successfully
				assertLimitResponse(tt, response, false)
				assertSetLimitingCookie(tt, response, withBlocking)
				body, err := io.ReadAll(response.Body)
				assert.NoError(tt, err)
				var compFunc assert.ComparisonAssertionFunc = assert.NotContains
				if withBlocking {
					compFunc = assert.Contains
				}
				compFunc(tt, string(body), `"exhausted":true`)
			}
		},
	}}
	runTest := func(t *testing.T, tt test, withBlocking bool, isFirst bool) {
		req, err, assertResponse := tt.testHttp(t)
		require.NoError(t, err)
		testHTTP := func() {
			resp, err := (&http.Client{
				// Don't follow redirects
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}).Do(req)
			defer func() {
				require.NoError(t, resp.Body.Close())
			}()
			require.NoError(t, err)
			assertResponse(resp, withBlocking)
		}
		if isFirst {
			// limits are eventually consistent, so we need to wait for the blocking to be set on the first test
			assert.EventuallyWithT(t, func(c *assert.CollectT) {
				testHTTP()
			}, 5*time.Second, time.Second, "wait for blocking to be set")
		} else {
			testHTTP()
		}
		if tt.testGrpc != nil {
			tt.testGrpc(t, withBlocking)
		}
	}
	_, err := Tester.Client.System.SetLimits(SystemCTX, &system.SetLimitsRequest{
		InstanceId: instanceID,
		Block:      wrapperspb.Bool(true),
	})
	require.NoError(t, err)
	for _, tt := range tests {
		var isFirst bool
		t.Run(tt.name+" with blocking", func(t *testing.T) {
			isFirst = isFirst || !t.Skipped()
			runTest(t, tt, true, isFirst)
		})
	}
	_, err = Tester.Client.System.SetLimits(SystemCTX, &system.SetLimitsRequest{
		InstanceId: instanceID,
		Block:      wrapperspb.Bool(false),
	})
	require.NoError(t, err)
	for _, tt := range tests {
		var isFirst bool
		t.Run(tt.name+" without blocking", func(t *testing.T) {
			isFirst = isFirst || !t.Skipped()
			runTest(t, tt, false, isFirst)
		})
	}
}

// If expectSet is true, we expect the cookie to be set
// If expectSet is false, we expect the cookie to be deleted
func assertSetLimitingCookie(t assert.TestingT, response *http.Response, expectSet bool) {
	for _, cookie := range response.Cookies() {
		if cookie.Name == "zitadel.quota.exhausted" {
			if expectSet {
				assert.Greater(t, cookie.MaxAge, 0)
			} else {
				assert.LessOrEqual(t, cookie.MaxAge, 0)
			}
			return
		}
	}
	assert.FailNow(t, "cookie not found")
}

func assertGrpcError(t assert.TestingT, err error, expectBlocked bool) {
	if expectBlocked {
		assert.Equal(t, codes.ResourceExhausted, status.Convert(err).Code())
		return
	}
	assert.NoError(t, err)
}

func assertLimitResponse(t assert.TestingT, response *http.Response, expectBlocked bool) {
	if expectBlocked {
		assert.Equal(t, http.StatusTooManyRequests, response.StatusCode)
		return
	}
	assert.GreaterOrEqual(t, response.StatusCode, 200)
	assert.Less(t, response.StatusCode, 300)
}

func idpExistsCondition(t assert.TestingT, instanceID, idpName string) func() bool {
	return func() bool {
		nameQuery, err := query.NewIDPTemplateNameSearchQuery(query.TextEquals, idpName)
		assert.NoError(t, err)
		instanceQuery, err := query.NewIDPTemplateResourceOwnerSearchQuery(instanceID)
		assert.NoError(t, err)
		idps, err := Tester.Queries.IDPTemplates(authz.WithInstanceID(CTX, instanceID), &query.IDPTemplateSearchQueries{
			Queries: []query.SearchQuery{
				instanceQuery,
				nameQuery,
			},
		}, false)
		assert.NoError(t, err)
		return len(idps.Templates) > 0
	}
}
