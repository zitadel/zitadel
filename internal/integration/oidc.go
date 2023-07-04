package integration

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/zitadel/oidc/v2/pkg/client/rp"
	"github.com/zitadel/oidc/v2/pkg/oidc"

	http_util "github.com/zitadel/zitadel/internal/api/http"
	oidc_internal "github.com/zitadel/zitadel/internal/api/oidc"
	"github.com/zitadel/zitadel/pkg/grpc/app"
	"github.com/zitadel/zitadel/pkg/grpc/management"
)

func (s *Tester) CreateOIDCNativeClient(ctx context.Context, redirectURI string) (*management.AddOIDCAppResponse, error) {
	project, err := s.Client.Mgmt.AddProject(ctx, &management.AddProjectRequest{
		Name: fmt.Sprintf("project-%d", time.Now().UnixNano()),
	})
	if err != nil {
		return nil, err
	}
	return s.Client.Mgmt.AddOIDCApp(ctx, &management.AddOIDCAppRequest{
		ProjectId:                project.GetId(),
		Name:                     fmt.Sprintf("app-%d", time.Now().UnixNano()),
		RedirectUris:             []string{redirectURI},
		ResponseTypes:            []app.OIDCResponseType{app.OIDCResponseType_OIDC_RESPONSE_TYPE_CODE},
		GrantTypes:               []app.OIDCGrantType{app.OIDCGrantType_OIDC_GRANT_TYPE_AUTHORIZATION_CODE, app.OIDCGrantType_OIDC_GRANT_TYPE_REFRESH_TOKEN},
		AppType:                  app.OIDCAppType_OIDC_APP_TYPE_NATIVE,
		AuthMethodType:           app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_NONE,
		PostLogoutRedirectUris:   nil,
		Version:                  app.OIDCVersion_OIDC_VERSION_1_0,
		DevMode:                  false,
		AccessTokenType:          app.OIDCTokenType_OIDC_TOKEN_TYPE_JWT,
		AccessTokenRoleAssertion: false,
		IdTokenRoleAssertion:     false,
		IdTokenUserinfoAssertion: false,
		ClockSkew:                nil,
		AdditionalOrigins:        nil,
		SkipNativeAppSuccessPage: false,
	})
}

func (s *Tester) CreateOIDCImplicitFlowClient(ctx context.Context, redirectURI string) (*management.AddOIDCAppResponse, error) {
	project, err := s.Client.Mgmt.AddProject(ctx, &management.AddProjectRequest{
		Name: fmt.Sprintf("project-%d", time.Now().UnixNano()),
	})
	if err != nil {
		return nil, err
	}
	return s.Client.Mgmt.AddOIDCApp(ctx, &management.AddOIDCAppRequest{
		ProjectId:                project.GetId(),
		Name:                     fmt.Sprintf("app-%d", time.Now().UnixNano()),
		RedirectUris:             []string{redirectURI},
		ResponseTypes:            []app.OIDCResponseType{app.OIDCResponseType_OIDC_RESPONSE_TYPE_ID_TOKEN_TOKEN},
		GrantTypes:               []app.OIDCGrantType{app.OIDCGrantType_OIDC_GRANT_TYPE_IMPLICIT},
		AppType:                  app.OIDCAppType_OIDC_APP_TYPE_USER_AGENT,
		AuthMethodType:           app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_NONE,
		PostLogoutRedirectUris:   nil,
		Version:                  app.OIDCVersion_OIDC_VERSION_1_0,
		DevMode:                  false,
		AccessTokenType:          app.OIDCTokenType_OIDC_TOKEN_TYPE_JWT,
		AccessTokenRoleAssertion: false,
		IdTokenRoleAssertion:     false,
		IdTokenUserinfoAssertion: false,
		ClockSkew:                nil,
		AdditionalOrigins:        nil,
		SkipNativeAppSuccessPage: false,
	})
}

func (s *Tester) CreateOIDCAuthRequest(clientID, loginClient, redirectURI string, scope ...string) (authRequestID string, err error) {
	provider, err := s.CreateRelyingParty(clientID, redirectURI, scope...)
	if err != nil {
		return "", err
	}

	codeVerifier := "codeVerifier"
	codeChallenge := oidc.NewSHACodeChallenge(codeVerifier)
	authURL := rp.AuthURL("state", provider, rp.WithCodeChallenge(codeChallenge))

	loc, err := CheckRedirect(authURL, map[string]string{oidc_internal.LoginClientHeader: loginClient})
	if err != nil {
		return "", err
	}

	prefixWithHost := provider.Issuer() + s.Config.OIDC.DefaultLoginURLV2
	if !strings.HasPrefix(loc.String(), prefixWithHost) {
		return "", fmt.Errorf("login location has not prefix %s, but is %s", prefixWithHost, loc.String())
	}
	return strings.TrimPrefix(loc.String(), prefixWithHost), nil
}

func (s *Tester) CreateRelyingParty(clientID, redirectURI string, scope ...string) (rp.RelyingParty, error) {
	issuer := http_util.BuildHTTP(s.Config.ExternalDomain, s.Config.Port, s.Config.ExternalSecure)
	if len(scope) == 0 {
		scope = []string{oidc.ScopeOpenID}
	}
	return rp.NewRelyingPartyOIDC(issuer, clientID, "", redirectURI, scope)
}

func CheckRedirect(url string, headers map[string]string) (*url.URL, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp.Location()
}
