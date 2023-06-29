//go:build integration

package oidc_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v2/pkg/client/rp"
	"github.com/zitadel/oidc/v2/pkg/oidc"

	http_util "github.com/zitadel/zitadel/internal/api/http"
	oidc_internal "github.com/zitadel/zitadel/internal/api/oidc"
	"github.com/zitadel/zitadel/internal/integration"
	app_pb "github.com/zitadel/zitadel/pkg/grpc/app"
	mgmt "github.com/zitadel/zitadel/pkg/grpc/management"
	session "github.com/zitadel/zitadel/pkg/grpc/session/v2alpha"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2alpha"
)

var (
	CTX               context.Context
	Tester            *integration.Tester
	Client            session.SessionServiceClient
	User              *user.AddHumanUserResponse
	GenericOAuthIDPID string
)

const (
	redirectURI = "oidc_integration_test://callback"
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, errCtx, cancel := integration.Contexts(5 * time.Minute)
		defer cancel()

		Tester = integration.NewTester(ctx)
		defer Tester.Done()
		Client = Tester.Client.SessionV2

		CTX, _ = Tester.WithSystemAuthorization(ctx, integration.OrgOwner), errCtx
		return m.Run()
	}())
}

func createClient(t testing.TB) string {
	project, err := Tester.Client.Mgmt.AddProject(CTX, &mgmt.AddProjectRequest{
		Name: fmt.Sprintf("project-%d", time.Now().UnixNano()),
	})
	require.NoError(t, err)

	app, err := Tester.Client.Mgmt.AddOIDCApp(CTX, &mgmt.AddOIDCAppRequest{
		ProjectId:                project.Id,
		Name:                     fmt.Sprintf("app-%d", time.Now().UnixNano()),
		RedirectUris:             []string{redirectURI},
		ResponseTypes:            []app_pb.OIDCResponseType{app_pb.OIDCResponseType_OIDC_RESPONSE_TYPE_CODE},
		GrantTypes:               []app_pb.OIDCGrantType{app_pb.OIDCGrantType_OIDC_GRANT_TYPE_AUTHORIZATION_CODE},
		AppType:                  app_pb.OIDCAppType_OIDC_APP_TYPE_NATIVE,
		AuthMethodType:           app_pb.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_NONE,
		PostLogoutRedirectUris:   nil,
		Version:                  app_pb.OIDCVersion_OIDC_VERSION_1_0,
		DevMode:                  false,
		AccessTokenType:          app_pb.OIDCTokenType_OIDC_TOKEN_TYPE_JWT,
		AccessTokenRoleAssertion: false,
		IdTokenRoleAssertion:     false,
		IdTokenUserinfoAssertion: false,
		ClockSkew:                nil,
		AdditionalOrigins:        nil,
		SkipNativeAppSuccessPage: false,
	})
	require.NoError(t, err)

	return app.GetClientId()
}

func createAuthRequest(t testing.TB, clientID string) string {
	issuer := http_util.BuildHTTP(Tester.Config.ExternalDomain, Tester.Config.Port, Tester.Config.ExternalSecure)
	provider, err := rp.NewRelyingPartyOIDC(issuer, clientID, "", redirectURI, []string{oidc.ScopeOpenID})
	require.NoError(t, err)
	codeVerifier := "codeVerifier"
	codeChallenge := oidc.NewSHACodeChallenge(codeVerifier)
	authURL := rp.AuthURL("state", provider, rp.WithCodeChallenge(codeChallenge))

	req, err := http.NewRequest(http.MethodGet, authURL, nil)
	require.NoError(t, err)
	req.Header.Set(oidc_internal.LoginClientHeader, "something")

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	loc, err := resp.Location()
	require.NoError(t, err)

	prefixWithHost := issuer + Tester.Config.OIDC.DefaultLoginURLV2
	require.Truef(t, strings.HasPrefix(loc.String(), prefixWithHost),
		"login location has not prefix %s, but is %s", prefixWithHost, loc.String())

	return strings.TrimPrefix(loc.String(), prefixWithHost)
}

func TestOPStorage_CreateAuthRequest(t *testing.T) {
	clientID := createClient(t)

	id := createAuthRequest(t, clientID)
	require.NotEmpty(t, id)
}
