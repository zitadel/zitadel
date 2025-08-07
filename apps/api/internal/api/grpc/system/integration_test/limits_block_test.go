//go:build integration

package system_test

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/admin"
	"github.com/zitadel/zitadel/pkg/grpc/system"
)

func TestServer_Limits_Block(t *testing.T) {
	isoInstance := integration.NewInstance(CTX)
	iamOwnerCtx := isoInstance.WithAuthorization(CTX, integration.UserTypeIAMOwner)
	tests := []*test{
		publicAPIBlockingTest(isoInstance.Domain),
		{
			name: "mutating API",
			testGrpc: func(tt assert.TestingT, expectBlocked bool) {
				randomGrpcIdpName := randomString("idp-grpc", 5)
				_, err := isoInstance.Client.Admin.AddGitHubProvider(iamOwnerCtx, &admin.AddGitHubProviderRequest{
					Name:         randomGrpcIdpName,
					ClientId:     "client-id",
					ClientSecret: "client-secret",
				})
				assertGrpcError(tt, err, expectBlocked)
			},
			testHttp: func(tt assert.TestingT) (*http.Request, error, func(assert.TestingT, *http.Response, bool)) {
				randomHttpIdpName := randomString("idp-http", 5)
				req, err := http.NewRequestWithContext(
					CTX,
					"POST",
					fmt.Sprintf("http://%s/admin/v1/idps/github", net.JoinHostPort(isoInstance.Domain, "8080")),
					strings.NewReader(`{
	"name": "`+randomHttpIdpName+`",
	"clientId": "client-id",
	"clientSecret": "client-secret"
}`),
				)
				if err != nil {
					return nil, err, nil
				}
				req.Header.Set("Authorization", isoInstance.BearerToken(iamOwnerCtx))
				return req, nil, func(ttt assert.TestingT, response *http.Response, expectBlocked bool) {
					assertLimitResponse(ttt, response, expectBlocked)
					assertSetLimitingCookie(ttt, response, expectBlocked)
				}
			},
		}, {
			name: "discovery",
			testHttp: func(tt assert.TestingT) (*http.Request, error, func(assert.TestingT, *http.Response, bool)) {
				req, err := http.NewRequestWithContext(
					CTX,
					"GET",
					fmt.Sprintf("http://%s/.well-known/openid-configuration", net.JoinHostPort(isoInstance.Domain, "8080")),
					nil,
				)
				return req, err, func(ttt assert.TestingT, response *http.Response, expectBlocked bool) {
					assertLimitResponse(ttt, response, expectBlocked)
					assertSetLimitingCookie(ttt, response, expectBlocked)
				}
			},
		}, {
			name: "login",
			testHttp: func(tt assert.TestingT) (*http.Request, error, func(assert.TestingT, *http.Response, bool)) {
				req, err := http.NewRequestWithContext(
					CTX,
					"GET",
					fmt.Sprintf("http://%s/ui/login/login/externalidp/callback", net.JoinHostPort(isoInstance.Domain, "8080")),
					nil,
				)
				return req, err, func(ttt assert.TestingT, response *http.Response, expectBlocked bool) {
					// the login paths should return a redirect if the instance is blocked
					if expectBlocked {
						assert.Equal(ttt, http.StatusFound, response.StatusCode)
					} else {
						assertLimitResponse(ttt, response, false)
					}
					assertSetLimitingCookie(ttt, response, expectBlocked)
				}
			},
		}, {
			name: "console",
			testHttp: func(tt assert.TestingT) (*http.Request, error, func(assert.TestingT, *http.Response, bool)) {
				req, err := http.NewRequestWithContext(
					CTX,
					"GET",
					fmt.Sprintf("http://%s/ui/console/", net.JoinHostPort(isoInstance.Domain, "8080")),
					nil,
				)
				return req, err, func(ttt assert.TestingT, response *http.Response, expectBlocked bool) {
					// the console is not blocked so we can render a link to an instance management portal.
					// A CDN can cache these assets easily
					// We also don't care about a cookie because the environment.json already takes care of that.
					assertLimitResponse(ttt, response, false)
				}
			},
		}, {
			name: "environment.json",
			testHttp: func(tt assert.TestingT) (*http.Request, error, func(assert.TestingT, *http.Response, bool)) {
				req, err := http.NewRequestWithContext(
					CTX,
					"GET",
					fmt.Sprintf("http://%s/ui/console/assets/environment.json", net.JoinHostPort(isoInstance.Domain, "8080")),
					nil,
				)
				return req, err, func(ttt assert.TestingT, response *http.Response, expectBlocked bool) {
					// the environment.json should always return successfully
					assertLimitResponse(ttt, response, false)
					assertSetLimitingCookie(ttt, response, expectBlocked)
					body, err := io.ReadAll(response.Body)
					assert.NoError(ttt, err)
					var compFunc assert.ComparisonAssertionFunc = assert.NotContains
					if expectBlocked {
						compFunc = assert.Contains
					}
					compFunc(ttt, string(body), `"exhausted":true`)
				}
			},
		}}
	_, err := integration.SystemClient().SetLimits(CTX, &system.SetLimitsRequest{
		InstanceId: isoInstance.ID(),
		Block:      gu.Ptr(true),
	})
	assert.NoError(t, err)
	// The following call ensures that an undefined bool is not deserialized to false
	_, err = integration.SystemClient().SetLimits(CTX, &system.SetLimitsRequest{
		InstanceId:        isoInstance.ID(),
		AuditLogRetention: durationpb.New(time.Hour),
	})
	assert.NoError(t, err)
	for _, tt := range tests {
		var isFirst bool
		t.Run(tt.name+" with blocking", func(t *testing.T) {
			isFirst = isFirst || !t.Skipped()
			testBlockingAPI(t, tt, true, isFirst)
		})
	}
	_, err = integration.SystemClient().SetLimits(CTX, &system.SetLimitsRequest{
		InstanceId: isoInstance.ID(),
		Block:      gu.Ptr(false),
	})
	require.NoError(t, err)
	for _, tt := range tests {
		var isFirst bool
		t.Run(tt.name+" without blocking", func(t *testing.T) {
			isFirst = isFirst || !t.Skipped()
			testBlockingAPI(t, tt, false, isFirst)
		})
	}
}

type test struct {
	name     string
	testHttp func(t assert.TestingT) (req *http.Request, err error, assertResponse func(t assert.TestingT, response *http.Response, expectBlocked bool))
	testGrpc func(t assert.TestingT, expectBlocked bool)
}

func testBlockingAPI(t *testing.T, tt *test, expectBlocked bool, isFirst bool) {
	req, err, assertResponse := tt.testHttp(t)
	require.NoError(t, err)
	testHTTP := func(t require.TestingT) {
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
		assertResponse(t, resp, expectBlocked)
	}
	if isFirst {
		// limits are eventually consistent, so we need to wait for the blocking to be set on the first test
		assert.EventuallyWithT(t, func(c *assert.CollectT) {
			testHTTP(c)
		}, time.Minute, time.Second, "wait for blocking to be set")
	} else {
		testHTTP(t)
	}
	if tt.testGrpc != nil {
		tt.testGrpc(t, expectBlocked)
	}
}

func publicAPIBlockingTest(domain string) *test {
	return &test{
		name: "public API",
		testGrpc: func(tt assert.TestingT, expectBlocked bool) {
			conn, err := grpc.DialContext(CTX, net.JoinHostPort(domain, "8080"),
				grpc.WithBlock(),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			)
			assert.NoError(tt, err)
			_, err = admin.NewAdminServiceClient(conn).Healthz(CTX, &admin.HealthzRequest{})
			assertGrpcError(tt, err, expectBlocked)
		},
		testHttp: func(tt assert.TestingT) (*http.Request, error, func(assert.TestingT, *http.Response, bool)) {
			req, err := http.NewRequestWithContext(
				CTX,
				"GET",
				fmt.Sprintf("http://%s/admin/v1/healthz", net.JoinHostPort(domain, "8080")),
				nil,
			)
			return req, err, func(ttt assert.TestingT, response *http.Response, expectBlocked bool) {
				assertLimitResponse(ttt, response, expectBlocked)
				assertSetLimitingCookie(ttt, response, expectBlocked)
			}
		},
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
	assert.Fail(t, "cookie not found")
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
