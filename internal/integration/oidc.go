package integration

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/zitadel/oidc/v3/pkg/client"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/client/rs"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"google.golang.org/protobuf/types/known/timestamppb"

	http_util "github.com/zitadel/zitadel/internal/api/http"
	oidc_internal "github.com/zitadel/zitadel/internal/api/oidc"
	"github.com/zitadel/zitadel/pkg/grpc/app"
	"github.com/zitadel/zitadel/pkg/grpc/authn"
	"github.com/zitadel/zitadel/pkg/grpc/management"
	"github.com/zitadel/zitadel/pkg/grpc/user"
)

func (s *Tester) CreateOIDCClient(ctx context.Context, redirectURI, logoutRedirectURI, projectID string, appType app.OIDCAppType, authMethod app.OIDCAuthMethodType) (*management.AddOIDCAppResponse, error) {
	return s.Client.Mgmt.AddOIDCApp(ctx, &management.AddOIDCAppRequest{
		ProjectId:                projectID,
		Name:                     fmt.Sprintf("app-%d", time.Now().UnixNano()),
		RedirectUris:             []string{redirectURI},
		ResponseTypes:            []app.OIDCResponseType{app.OIDCResponseType_OIDC_RESPONSE_TYPE_CODE},
		GrantTypes:               []app.OIDCGrantType{app.OIDCGrantType_OIDC_GRANT_TYPE_AUTHORIZATION_CODE, app.OIDCGrantType_OIDC_GRANT_TYPE_REFRESH_TOKEN},
		AppType:                  appType,
		AuthMethodType:           authMethod,
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

func (s *Tester) CreateOIDCNativeClient(ctx context.Context, redirectURI, logoutRedirectURI, projectID string) (*management.AddOIDCAppResponse, error) {
	return s.CreateOIDCClient(ctx, redirectURI, logoutRedirectURI, projectID, app.OIDCAppType_OIDC_APP_TYPE_NATIVE, app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_NONE)
}

func (s *Tester) CreateOIDCWebClientBasic(ctx context.Context, redirectURI, logoutRedirectURI, projectID string) (*management.AddOIDCAppResponse, error) {
	return s.CreateOIDCClient(ctx, redirectURI, logoutRedirectURI, projectID, app.OIDCAppType_OIDC_APP_TYPE_WEB, app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_BASIC)
}

func (s *Tester) CreateOIDCWebClientJWT(ctx context.Context, redirectURI, logoutRedirectURI, projectID string) (client *management.AddOIDCAppResponse, keyData []byte, err error) {
	client, err = s.CreateOIDCClient(ctx, redirectURI, logoutRedirectURI, projectID, app.OIDCAppType_OIDC_APP_TYPE_WEB, app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT)
	if err != nil {
		return nil, nil, err
	}
	key, err := s.Client.Mgmt.AddAppKey(ctx, &management.AddAppKeyRequest{
		ProjectId:      projectID,
		AppId:          client.GetAppId(),
		Type:           authn.KeyType_KEY_TYPE_JSON,
		ExpirationDate: timestamppb.New(time.Now().Add(time.Hour)),
	})
	if err != nil {
		return nil, nil, err
	}
	return client, key.GetKeyDetails(), nil
}

func (s *Tester) CreateOIDCInactivateClient(ctx context.Context, redirectURI, logoutRedirectURI, projectID string) (*management.AddOIDCAppResponse, error) {
	client, err := s.CreateOIDCNativeClient(ctx, redirectURI, logoutRedirectURI, projectID)
	if err != nil {
		return nil, err
	}
	_, err = s.Client.Mgmt.DeactivateApp(ctx, &management.DeactivateAppRequest{
		ProjectId: projectID,
		AppId:     client.GetAppId(),
	})
	if err != nil {
		return nil, err
	}
	return client, err
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

const CodeVerifier = "codeVerifier"

func (s *Tester) CreateOIDCAuthRequest(ctx context.Context, clientID, loginClient, redirectURI string, scope ...string) (authRequestID string, err error) {
	provider, err := s.CreateRelyingParty(ctx, clientID, redirectURI, scope...)
	if err != nil {
		return "", err
	}
	codeChallenge := oidc.NewSHACodeChallenge(CodeVerifier)
	authURL := rp.AuthURL("state", provider, rp.WithCodeChallenge(codeChallenge))

	req, err := GetRequest(authURL, map[string]string{oidc_internal.LoginClientHeader: loginClient})
	if err != nil {
		return "", err
	}

	loc, err := CheckRedirect(req)
	if err != nil {
		return "", err
	}

	prefixWithHost := provider.Issuer() + s.Config.OIDC.DefaultLoginURLV2
	if !strings.HasPrefix(loc.String(), prefixWithHost) {
		return "", fmt.Errorf("login location has not prefix %s, but is %s", prefixWithHost, loc.String())
	}
	return strings.TrimPrefix(loc.String(), prefixWithHost), nil
}

func (s *Tester) CreateOIDCAuthRequestImplicit(ctx context.Context, clientID, loginClient, redirectURI string, scope ...string) (authRequestID string, err error) {
	provider, err := s.CreateRelyingParty(ctx, clientID, redirectURI, scope...)
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

	req, err := GetRequest(authURL, map[string]string{oidc_internal.LoginClientHeader: loginClient})
	if err != nil {
		return "", err
	}

	loc, err := CheckRedirect(req)
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

func (s *Tester) CreateRelyingParty(ctx context.Context, clientID, redirectURI string, scope ...string) (rp.RelyingParty, error) {
	if len(scope) == 0 {
		scope = []string{oidc.ScopeOpenID}
	}
	loginClient := &http.Client{Transport: &loginRoundTripper{http.DefaultTransport}}
	return rp.NewRelyingPartyOIDC(ctx, s.OIDCIssuer(), clientID, "", redirectURI, scope, rp.WithHTTPClient(loginClient))
}

type loginRoundTripper struct {
	http.RoundTripper
}

func (c *loginRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set(oidc_internal.LoginClientHeader, LoginUser)
	return c.RoundTripper.RoundTrip(req)
}

func (s *Tester) CreateResourceServer(ctx context.Context, keyFileData []byte) (rs.ResourceServer, error) {
	keyFile, err := client.ConfigFromKeyFileData(keyFileData)
	if err != nil {
		return nil, err
	}
	return rs.NewResourceServerJWTProfile(ctx, s.OIDCIssuer(), keyFile.ClientID, keyFile.KeyID, []byte(keyFile.Key))
}

func GetRequest(url string, headers map[string]string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	return req, nil
}

func CheckRedirect(req *http.Request) (*url.URL, error) {
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

func (s *Tester) CreateOIDCCredentialsClient(ctx context.Context) (string, string, error) {
	name := gofakeit.Username()
	user, err := s.Client.Mgmt.AddMachineUser(ctx, &management.AddMachineUserRequest{
		Name:            name,
		UserName:        name,
		AccessTokenType: user.AccessTokenType_ACCESS_TOKEN_TYPE_JWT,
	})
	if err != nil {
		return "", "", err
	}
	secret, err := s.Client.Mgmt.GenerateMachineSecret(ctx, &management.GenerateMachineSecretRequest{
		UserId: user.GetUserId(),
	})
	if err != nil {
		return "", "", err
	}
	return secret.GetClientId(), secret.GetClientSecret(), nil
}
