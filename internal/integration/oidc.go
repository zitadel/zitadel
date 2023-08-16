package integration

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/zitadel/oidc/v2/pkg/client"
	"github.com/zitadel/oidc/v2/pkg/client/rp"
	"github.com/zitadel/oidc/v2/pkg/client/rs"
	"github.com/zitadel/oidc/v2/pkg/oidc"

	http_util "github.com/zitadel/zitadel/internal/api/http"
	oidc_internal "github.com/zitadel/zitadel/internal/api/oidc"
	"github.com/zitadel/zitadel/pkg/grpc/app"
	"github.com/zitadel/zitadel/pkg/grpc/management"
)

func (s *Tester) CreateOIDCNativeClient(ctx context.Context, redirectURI, logoutRedirectURI, projectID string) (*management.AddOIDCAppResponse, error) {
	return s.Client.Mgmt.AddOIDCApp(ctx, &management.AddOIDCAppRequest{
		ProjectId:                projectID,
		Name:                     fmt.Sprintf("app-%d", time.Now().UnixNano()),
		RedirectUris:             []string{redirectURI},
		ResponseTypes:            []app.OIDCResponseType{app.OIDCResponseType_OIDC_RESPONSE_TYPE_CODE},
		GrantTypes:               []app.OIDCGrantType{app.OIDCGrantType_OIDC_GRANT_TYPE_AUTHORIZATION_CODE, app.OIDCGrantType_OIDC_GRANT_TYPE_REFRESH_TOKEN},
		AppType:                  app.OIDCAppType_OIDC_APP_TYPE_NATIVE,
		AuthMethodType:           app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_NONE,
		PostLogoutRedirectUris:   []string{logoutRedirectURI},
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
		DevMode:                  true,
		AccessTokenType:          app.OIDCTokenType_OIDC_TOKEN_TYPE_JWT,
		AccessTokenRoleAssertion: false,
		IdTokenRoleAssertion:     false,
		IdTokenUserinfoAssertion: false,
		ClockSkew:                nil,
		AdditionalOrigins:        nil,
		SkipNativeAppSuccessPage: false,
	})
}

func (s *Tester) CreateProject(ctx context.Context) (*management.AddProjectResponse, error) {
	return s.Client.Mgmt.AddProject(ctx, &management.AddProjectRequest{
		Name: fmt.Sprintf("project-%d", time.Now().UnixNano()),
	})
}

func (s *Tester) CreateAPIClient(ctx context.Context, projectID string) (*management.AddAPIAppResponse, error) {
	return s.Client.Mgmt.AddAPIApp(ctx, &management.AddAPIAppRequest{
		ProjectId:      projectID,
		Name:           fmt.Sprintf("api-%d", time.Now().UnixNano()),
		AuthMethodType: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT,
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

func (s *Tester) CreateOIDCAuthRequestImplicit(clientID, loginClient, redirectURI string, scope ...string) (authRequestID string, err error) {
	provider, err := s.CreateRelyingParty(clientID, redirectURI, scope...)
	if err != nil {
		return "", err
	}

	authURL := rp.AuthURL("state", provider)

	// implicit is not natively supported so let's just overwrite the response type
	parsed, _ := url.Parse(authURL)
	queries := parsed.Query()
	queries.Set("response_type", string(oidc.ResponseTypeIDToken))
	parsed.RawQuery = queries.Encode()
	authURL = parsed.String()

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

func (s *Tester) OIDCIssuer() string {
	return http_util.BuildHTTP(s.Config.ExternalDomain, s.Config.Port, s.Config.ExternalSecure)
}

func (s *Tester) CreateRelyingParty(clientID, redirectURI string, scope ...string) (rp.RelyingParty, error) {
	if len(scope) == 0 {
		scope = []string{oidc.ScopeOpenID}
	}
	loginClient := &http.Client{Transport: &loginRoundTripper{http.DefaultTransport}}
	return rp.NewRelyingPartyOIDC(s.OIDCIssuer(), clientID, "", redirectURI, scope, rp.WithHTTPClient(loginClient))
}

type loginRoundTripper struct {
	http.RoundTripper
}

func (c *loginRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set(oidc_internal.LoginClientHeader, LoginUser)
	return c.RoundTripper.RoundTrip(req)
}

func (s *Tester) CreateResourceServer(keyFileData []byte) (rs.ResourceServer, error) {
	keyFile, err := client.ConfigFromKeyFileData(keyFileData)
	if err != nil {
		return nil, err
	}
	return rs.NewResourceServerJWTProfile(s.OIDCIssuer(), keyFile.ClientID, keyFile.KeyID, []byte(keyFile.Key))
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
