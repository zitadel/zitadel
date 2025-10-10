package integration

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

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
	user_v2 "github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func (i *Instance) CreateOIDCClientLoginVersion(ctx context.Context, redirectURI, logoutRedirectURI, projectID string, appType app.OIDCAppType, authMethod app.OIDCAuthMethodType, devMode bool, loginVersion *app.LoginVersion, grantTypes ...app.OIDCGrantType) (*management.AddOIDCAppResponse, error) {
	if len(grantTypes) == 0 {
		grantTypes = []app.OIDCGrantType{app.OIDCGrantType_OIDC_GRANT_TYPE_AUTHORIZATION_CODE, app.OIDCGrantType_OIDC_GRANT_TYPE_REFRESH_TOKEN}
	}
	resp, err := i.Client.Mgmt.AddOIDCApp(ctx, &management.AddOIDCAppRequest{
		ProjectId:                projectID,
		Name:                     fmt.Sprintf("app-%d", time.Now().UnixNano()),
		RedirectUris:             []string{redirectURI},
		ResponseTypes:            []app.OIDCResponseType{app.OIDCResponseType_OIDC_RESPONSE_TYPE_CODE},
		GrantTypes:               grantTypes,
		AppType:                  appType,
		AuthMethodType:           authMethod,
		PostLogoutRedirectUris:   []string{logoutRedirectURI},
		Version:                  app.OIDCVersion_OIDC_VERSION_1_0,
		DevMode:                  devMode,
		AccessTokenType:          app.OIDCTokenType_OIDC_TOKEN_TYPE_JWT,
		AccessTokenRoleAssertion: false,
		IdTokenRoleAssertion:     false,
		IdTokenUserinfoAssertion: false,
		ClockSkew:                nil,
		AdditionalOrigins:        nil,
		SkipNativeAppSuccessPage: false,
		LoginVersion:             loginVersion,
	})
	if err != nil {
		return nil, err
	}
	return resp, await(func() error {
		_, err := i.Client.Mgmt.GetProjectByID(ctx, &management.GetProjectByIDRequest{
			Id: projectID,
		})
		if err != nil {
			return err
		}
		_, err = i.Client.Mgmt.GetAppByID(ctx, &management.GetAppByIDRequest{
			ProjectId: projectID,
			AppId:     resp.GetAppId(),
		})
		return err
	})
}

func (i *Instance) CreateOIDCClient(ctx context.Context, redirectURI, logoutRedirectURI, projectID string, appType app.OIDCAppType, authMethod app.OIDCAuthMethodType, devMode bool, grantTypes ...app.OIDCGrantType) (*management.AddOIDCAppResponse, error) {
	return i.CreateOIDCClientLoginVersion(ctx, redirectURI, logoutRedirectURI, projectID, appType, authMethod, devMode, nil, grantTypes...)
}

func (i *Instance) CreateOIDCNativeClient(ctx context.Context, redirectURI, logoutRedirectURI, projectID string, devMode bool) (*management.AddOIDCAppResponse, error) {
	return i.CreateOIDCClient(ctx, redirectURI, logoutRedirectURI, projectID, app.OIDCAppType_OIDC_APP_TYPE_NATIVE, app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_NONE, devMode)
}

func (i *Instance) CreateOIDCWebClientBasic(ctx context.Context, redirectURI, logoutRedirectURI, projectID string) (*management.AddOIDCAppResponse, error) {
	return i.CreateOIDCClient(ctx, redirectURI, logoutRedirectURI, projectID, app.OIDCAppType_OIDC_APP_TYPE_WEB, app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_BASIC, false)
}

func (i *Instance) CreateOIDCWebClientJWT(ctx context.Context, redirectURI, logoutRedirectURI, projectID string, grantTypes ...app.OIDCGrantType) (client *management.AddOIDCAppResponse, keyData []byte, err error) {
	client, err = i.CreateOIDCClient(ctx, redirectURI, logoutRedirectURI, projectID, app.OIDCAppType_OIDC_APP_TYPE_WEB, app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT, false, grantTypes...)
	if err != nil {
		return nil, nil, err
	}
	key, err := i.Client.Mgmt.AddAppKey(ctx, &management.AddAppKeyRequest{
		ProjectId:      projectID,
		AppId:          client.GetAppId(),
		Type:           authn.KeyType_KEY_TYPE_JSON,
		ExpirationDate: timestamppb.New(time.Now().Add(time.Hour)),
	})
	if err != nil {
		return nil, nil, err
	}
	mustAwait(func() error {
		_, err := i.Client.Mgmt.GetAppByID(ctx, &management.GetAppByIDRequest{
			ProjectId: projectID,
			AppId:     client.GetAppId(),
		})
		return err
	})

	return client, key.GetKeyDetails(), nil
}

func (i *Instance) CreateOIDCInactivateClient(ctx context.Context, redirectURI, logoutRedirectURI, projectID string) (*management.AddOIDCAppResponse, error) {
	client, err := i.CreateOIDCNativeClient(ctx, redirectURI, logoutRedirectURI, projectID, false)
	if err != nil {
		return nil, err
	}
	_, err = i.Client.Mgmt.DeactivateApp(ctx, &management.DeactivateAppRequest{
		ProjectId: projectID,
		AppId:     client.GetAppId(),
	})
	if err != nil {
		return nil, err
	}
	return client, err
}

func (i *Instance) CreateOIDCInactivateProjectClient(ctx context.Context, redirectURI, logoutRedirectURI, projectID string) (*management.AddOIDCAppResponse, error) {
	client, err := i.CreateOIDCNativeClient(ctx, redirectURI, logoutRedirectURI, projectID, false)
	if err != nil {
		return nil, err
	}
	_, err = i.Client.Mgmt.DeactivateProject(ctx, &management.DeactivateProjectRequest{
		Id: projectID,
	})
	if err != nil {
		return nil, err
	}
	return client, err
}

func (i *Instance) CreateOIDCImplicitFlowClient(ctx context.Context, t *testing.T, redirectURI string, loginVersion *app.LoginVersion) (*management.AddOIDCAppResponse, error) {
	project := i.CreateProject(ctx, t, "", ProjectName(), false, false)

	resp, err := i.Client.Mgmt.AddOIDCApp(ctx, &management.AddOIDCAppRequest{
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
		LoginVersion:             loginVersion,
	})
	if err != nil {
		return nil, err
	}
	return resp, await(func() error {
		_, err := i.Client.Mgmt.GetProjectByID(ctx, &management.GetProjectByIDRequest{
			Id: project.GetId(),
		})
		if err != nil {
			return err
		}
		_, err = i.Client.Mgmt.GetAppByID(ctx, &management.GetAppByIDRequest{
			ProjectId: project.GetId(),
			AppId:     resp.GetAppId(),
		})
		return err
	})
}

func (i *Instance) CreateOIDCTokenExchangeClient(ctx context.Context, t *testing.T) (client *management.AddOIDCAppResponse, keyData []byte, err error) {
	project := i.CreateProject(ctx, t, "", ProjectName(), false, false)
	return i.CreateOIDCWebClientJWT(ctx, "", "", project.GetId(), app.OIDCGrantType_OIDC_GRANT_TYPE_TOKEN_EXCHANGE, app.OIDCGrantType_OIDC_GRANT_TYPE_AUTHORIZATION_CODE, app.OIDCGrantType_OIDC_GRANT_TYPE_REFRESH_TOKEN)
}

func (i *Instance) CreateAPIClientJWT(ctx context.Context, projectID string) (*management.AddAPIAppResponse, error) {
	return i.Client.Mgmt.AddAPIApp(ctx, &management.AddAPIAppRequest{
		ProjectId:      projectID,
		Name:           fmt.Sprintf("api-%d", time.Now().UnixNano()),
		AuthMethodType: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT,
	})
}

func (i *Instance) CreateAPIClientBasic(ctx context.Context, projectID string) (*management.AddAPIAppResponse, error) {
	return i.Client.Mgmt.AddAPIApp(ctx, &management.AddAPIAppRequest{
		ProjectId:      projectID,
		Name:           fmt.Sprintf("api-%d", time.Now().UnixNano()),
		AuthMethodType: app.APIAuthMethodType_API_AUTH_METHOD_TYPE_BASIC,
	})
}

const CodeVerifier = "codeVerifier"

func (i *Instance) CreateOIDCAuthRequest(ctx context.Context, clientID, loginClient, redirectURI string, scope ...string) (now time.Time, authRequestID string, err error) {
	return i.CreateOIDCAuthRequestWithDomain(ctx, i.Domain, clientID, loginClient, redirectURI, scope...)
}

func (i *Instance) CreateOIDCAuthRequestWithoutLoginClientHeader(ctx context.Context, clientID, redirectURI, loginBaseURI string, scope ...string) (now time.Time, authRequestID string, err error) {
	return i.createOIDCAuthRequestWithDomain(ctx, i.Domain, clientID, redirectURI, "", loginBaseURI, scope...)
}

func (i *Instance) CreateOIDCAuthRequestWithDomain(ctx context.Context, domain, clientID, loginClient, redirectURI string, scope ...string) (now time.Time, authRequestID string, err error) {
	return i.createOIDCAuthRequestWithDomain(ctx, domain, clientID, redirectURI, loginClient, "", scope...)
}

func (i *Instance) createOIDCAuthRequestWithDomain(ctx context.Context, domain, clientID, redirectURI, loginClient, loginBaseURI string, scope ...string) (now time.Time, authRequestID string, err error) {
	provider, err := i.CreateRelyingPartyForDomain(ctx, domain, clientID, redirectURI, loginClient, scope...)
	if err != nil {
		return now, "", fmt.Errorf("create relying party: %w", err)
	}
	codeChallenge := oidc.NewSHACodeChallenge(CodeVerifier)
	authURL := rp.AuthURL("state", provider, rp.WithCodeChallenge(codeChallenge))

	var headers map[string]string
	if loginClient != "" {
		headers = map[string]string{oidc_internal.LoginClientHeader: loginClient}
	}
	req, err := GetRequest(authURL, headers)
	if err != nil {
		return now, "", fmt.Errorf("get request: %w", err)
	}

	now = time.Now()
	loc, err := CheckRedirect(req)
	if err != nil {
		return now, "", fmt.Errorf("check redirect: %w", err)
	}

	if loginBaseURI == "" {
		loginBaseURI = provider.Issuer() + i.Config.LoginURLV2
	}
	if !strings.HasPrefix(loc.String(), loginBaseURI) {
		return now, "", fmt.Errorf("login location has not prefix %s, but is %s", loginBaseURI, loc.String())
	}
	return now, strings.TrimPrefix(loc.String(), loginBaseURI), nil
}

func (i *Instance) CreateOIDCAuthRequestImplicitWithoutLoginClientHeader(ctx context.Context, clientID, redirectURI string, scope ...string) (authRequestID string, err error) {
	return i.createOIDCAuthRequestImplicit(ctx, clientID, redirectURI, nil, scope...)
}

func (i *Instance) CreateOIDCAuthRequestImplicit(ctx context.Context, clientID, loginClient, redirectURI string, scope ...string) (authRequestID string, err error) {
	return i.createOIDCAuthRequestImplicit(ctx, clientID, redirectURI, map[string]string{oidc_internal.LoginClientHeader: loginClient}, scope...)
}

func (i *Instance) createOIDCAuthRequestImplicit(ctx context.Context, clientID, redirectURI string, headers map[string]string, scope ...string) (authRequestID string, err error) {
	provider, err := i.CreateRelyingParty(ctx, clientID, redirectURI, scope...)
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

	req, err := GetRequest(authURL, headers)
	if err != nil {
		return "", err
	}

	loc, err := CheckRedirect(req)
	if err != nil {
		return "", err
	}

	prefixWithHost := provider.Issuer() + i.Config.LoginURLV2
	if !strings.HasPrefix(loc.String(), prefixWithHost) {
		return "", fmt.Errorf("login location has not prefix %s, but is %s", prefixWithHost, loc.String())
	}
	return strings.TrimPrefix(loc.String(), prefixWithHost), nil
}

func (i *Instance) OIDCIssuer() string {
	return http_util.BuildHTTP(i.Domain, i.Config.Port, i.Config.Secure)
}

func (i *Instance) CreateRelyingParty(ctx context.Context, clientID, redirectURI string, scope ...string) (rp.RelyingParty, error) {
	return i.CreateRelyingPartyForDomain(ctx, i.Domain, clientID, redirectURI, i.Users.Get(UserTypeLogin).Username, scope...)
}

func (i *Instance) CreateRelyingPartyWithoutLoginClientHeader(ctx context.Context, clientID, redirectURI string, scope ...string) (rp.RelyingParty, error) {
	return i.CreateRelyingPartyForDomain(ctx, i.Domain, clientID, redirectURI, "", scope...)
}

func (i *Instance) CreateRelyingPartyForDomain(ctx context.Context, domain, clientID, redirectURI, loginClientUsername string, scope ...string) (rp.RelyingParty, error) {
	if len(scope) == 0 {
		scope = []string{oidc.ScopeOpenID}
	}
	if loginClientUsername == "" {
		return rp.NewRelyingPartyOIDC(ctx, http_util.BuildHTTP(domain, i.Config.Port, i.Config.Secure), clientID, "", redirectURI, scope)
	}
	loginClient := &http.Client{Transport: &loginRoundTripper{http.DefaultTransport, loginClientUsername}}
	return rp.NewRelyingPartyOIDC(ctx, http_util.BuildHTTP(domain, i.Config.Port, i.Config.Secure), clientID, "", redirectURI, scope, rp.WithHTTPClient(loginClient))
}

type loginRoundTripper struct {
	http.RoundTripper
	loginUsername string
}

func (c *loginRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set(oidc_internal.LoginClientHeader, c.loginUsername)
	return c.RoundTripper.RoundTrip(req)
}

func (i *Instance) CreateResourceServerJWTProfile(ctx context.Context, keyFileData []byte) (rs.ResourceServer, error) {
	keyFile, err := client.ConfigFromKeyFileData(keyFileData)
	if err != nil {
		return nil, err
	}
	return rs.NewResourceServerJWTProfile(ctx, i.OIDCIssuer(), keyFile.ClientID, keyFile.KeyID, []byte(keyFile.Key))
}

func (i *Instance) CreateResourceServerClientCredentials(ctx context.Context, clientID, clientSecret string) (rs.ResourceServer, error) {
	return rs.NewResourceServerClientCredentials(ctx, i.OIDCIssuer(), clientID, clientSecret)
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

func CheckPost(url string, values url.Values) (*url.URL, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.PostForm(url, values)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp.Location()
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
	if resp.StatusCode < 300 || resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("check redirect unexpected status: %q; body: %q", resp.Status, body)
	}

	return resp.Location()
}

func (i *Instance) CreateOIDCCredentialsClient(ctx context.Context) (machine *management.AddMachineUserResponse, name, clientID, clientSecret string, err error) {
	name = Username()
	machine, err = i.Client.Mgmt.AddMachineUser(ctx, &management.AddMachineUserRequest{
		Name:            name,
		UserName:        name,
		AccessTokenType: user.AccessTokenType_ACCESS_TOKEN_TYPE_JWT,
	})
	if err != nil {
		return nil, "", "", "", err
	}
	secret, err := i.Client.Mgmt.GenerateMachineSecret(ctx, &management.GenerateMachineSecretRequest{
		UserId: machine.GetUserId(),
	})
	if err != nil {
		return nil, "", "", "", err
	}
	return machine, name, secret.GetClientId(), secret.GetClientSecret(), nil
}

func (i *Instance) CreateOIDCCredentialsClientInactive(ctx context.Context) (machine *management.AddMachineUserResponse, name, clientID, clientSecret string, err error) {
	name = Username()
	machine, err = i.Client.Mgmt.AddMachineUser(ctx, &management.AddMachineUserRequest{
		Name:            name,
		UserName:        name,
		AccessTokenType: user.AccessTokenType_ACCESS_TOKEN_TYPE_JWT,
	})
	if err != nil {
		return nil, "", "", "", err
	}
	secret, err := i.Client.Mgmt.GenerateMachineSecret(ctx, &management.GenerateMachineSecretRequest{
		UserId: machine.GetUserId(),
	})
	if err != nil {
		return nil, "", "", "", err
	}
	_, err = i.Client.UserV2.DeactivateUser(ctx, &user_v2.DeactivateUserRequest{
		UserId: machine.GetUserId(),
	})
	if err != nil {
		return nil, "", "", "", err
	}
	return machine, name, secret.GetClientId(), secret.GetClientSecret(), nil
}

func (i *Instance) CreateOIDCJWTProfileClient(ctx context.Context, keyLifetime time.Duration) (machine *management.AddMachineUserResponse, name string, keyData []byte, err error) {
	name = Username()
	machine, err = i.Client.Mgmt.AddMachineUser(ctx, &management.AddMachineUserRequest{
		Name:            name,
		UserName:        name,
		AccessTokenType: user.AccessTokenType_ACCESS_TOKEN_TYPE_JWT,
	})
	if err != nil {
		return nil, "", nil, err
	}
	keyResp, err := i.Client.Mgmt.AddMachineKey(ctx, &management.AddMachineKeyRequest{
		UserId:         machine.GetUserId(),
		Type:           authn.KeyType_KEY_TYPE_JSON,
		ExpirationDate: timestamppb.New(time.Now().Add(keyLifetime)),
	})
	if err != nil {
		return nil, "", nil, err
	}
	mustAwait(func() error {
		_, err := i.Client.Mgmt.GetMachineKeyByIDs(ctx, &management.GetMachineKeyByIDsRequest{
			UserId: machine.GetUserId(),
			KeyId:  keyResp.GetKeyId(),
		})
		return err
	})

	return machine, name, keyResp.GetKeyDetails(), nil
}

func (i *Instance) CreateDeviceAuthorizationRequest(ctx context.Context, clientID string, scopes ...string) (*oidc.DeviceAuthorizationResponse, error) {
	provider, err := i.CreateRelyingParty(ctx, clientID, "", scopes...)
	if err != nil {
		return nil, err
	}
	return rp.DeviceAuthorization(ctx, scopes, provider, nil)
}
