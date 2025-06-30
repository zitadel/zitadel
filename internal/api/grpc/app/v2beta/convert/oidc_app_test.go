package convert

import (
	"net/url"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/query"
	app "github.com/zitadel/zitadel/pkg/grpc/app/v2beta"
)

func TestCreateOIDCAppRequestToDomain(t *testing.T) {
	t.Parallel()

	tt := []struct {
		testName  string
		projectID string
		req       *app.CreateOIDCApplicationRequest

		expectedModel *domain.OIDCApp
		expectedError error
	}{
		{
			testName:  "unparsable login version 2 URL",
			projectID: "pid",
			req: &app.CreateOIDCApplicationRequest{
				LoginVersion: &app.LoginVersion{Version: &app.LoginVersion_LoginV2{
					LoginV2: &app.LoginV2{BaseUri: gu.Ptr("%+o")}},
				},
			},
			expectedModel: nil,
			expectedError: &url.Error{
				URL: "%+o",
				Op:  "parse",
				Err: url.EscapeError("%+o"),
			},
		},
		{
			testName:  "all fields set",
			projectID: "project1",
			req: &app.CreateOIDCApplicationRequest{
				RedirectUris:             []string{"https://redirect"},
				ResponseTypes:            []app.OIDCResponseType{app.OIDCResponseType_OIDC_RESPONSE_TYPE_CODE},
				GrantTypes:               []app.OIDCGrantType{app.OIDCGrantType_OIDC_GRANT_TYPE_AUTHORIZATION_CODE},
				AppType:                  app.OIDCAppType_OIDC_APP_TYPE_WEB,
				AuthMethodType:           app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_BASIC,
				PostLogoutRedirectUris:   []string{"https://logout"},
				DevMode:                  true,
				AccessTokenType:          app.OIDCTokenType_OIDC_TOKEN_TYPE_BEARER,
				AccessTokenRoleAssertion: true,
				IdTokenRoleAssertion:     true,
				IdTokenUserinfoAssertion: true,
				ClockSkew:                durationpb.New(5 * time.Second),
				AdditionalOrigins:        []string{"https://origin"},
				SkipNativeAppSuccessPage: true,
				BackChannelLogoutUri:     "https://backchannel",
				LoginVersion: &app.LoginVersion{Version: &app.LoginVersion_LoginV2{LoginV2: &app.LoginV2{
					BaseUri: gu.Ptr("https://login"),
				}}},
			},
			expectedModel: &domain.OIDCApp{
				ObjectRoot:               models.ObjectRoot{AggregateID: "project1"},
				AppName:                  "all fields set",
				OIDCVersion:              gu.Ptr(domain.OIDCVersionV1),
				RedirectUris:             []string{"https://redirect"},
				ResponseTypes:            []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
				GrantTypes:               []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
				ApplicationType:          gu.Ptr(domain.OIDCApplicationTypeWeb),
				AuthMethodType:           gu.Ptr(domain.OIDCAuthMethodTypeBasic),
				PostLogoutRedirectUris:   []string{"https://logout"},
				DevMode:                  gu.Ptr(true),
				AccessTokenType:          gu.Ptr(domain.OIDCTokenTypeBearer),
				AccessTokenRoleAssertion: gu.Ptr(true),
				IDTokenRoleAssertion:     gu.Ptr(true),
				IDTokenUserinfoAssertion: gu.Ptr(true),
				ClockSkew:                gu.Ptr(5 * time.Second),
				AdditionalOrigins:        []string{"https://origin"},
				SkipNativeAppSuccessPage: gu.Ptr(true),
				BackChannelLogoutURI:     gu.Ptr("https://backchannel"),
				LoginVersion:             gu.Ptr(domain.LoginVersion2),
				LoginBaseURI:             gu.Ptr("https://login"),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// When
			res, err := CreateOIDCAppRequestToDomain(tc.testName, tc.projectID, tc.req)

			// Then
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedModel, res)
		})
	}
}

func TestUpdateOIDCAppConfigRequestToDomain(t *testing.T) {
	t.Parallel()

	tt := []struct {
		testName string

		appID     string
		projectID string
		req       *app.UpdateOIDCApplicationConfigurationRequest

		expectedModel *domain.OIDCApp
		expectedError error
	}{
		{
			testName:  "unparsable login version 2 URL",
			appID:     "app1",
			projectID: "pid",
			req: &app.UpdateOIDCApplicationConfigurationRequest{
				LoginVersion: &app.LoginVersion{Version: &app.LoginVersion_LoginV2{
					LoginV2: &app.LoginV2{BaseUri: gu.Ptr("%+o")},
				}},
			},
			expectedModel: nil,
			expectedError: &url.Error{
				URL: "%+o",
				Op:  "parse",
				Err: url.EscapeError("%+o"),
			},
		},
		{
			testName:  "successful Update",
			appID:     "app1",
			projectID: "proj1",
			req: &app.UpdateOIDCApplicationConfigurationRequest{
				RedirectUris:             []string{"https://redirect"},
				ResponseTypes:            []app.OIDCResponseType{app.OIDCResponseType_OIDC_RESPONSE_TYPE_CODE},
				GrantTypes:               []app.OIDCGrantType{app.OIDCGrantType_OIDC_GRANT_TYPE_AUTHORIZATION_CODE},
				AppType:                  gu.Ptr(app.OIDCAppType_OIDC_APP_TYPE_WEB),
				AuthMethodType:           gu.Ptr(app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_BASIC),
				PostLogoutRedirectUris:   []string{"https://logout"},
				DevMode:                  gu.Ptr(true),
				AccessTokenType:          gu.Ptr(app.OIDCTokenType_OIDC_TOKEN_TYPE_BEARER),
				AccessTokenRoleAssertion: gu.Ptr(true),
				IdTokenRoleAssertion:     gu.Ptr(true),
				IdTokenUserinfoAssertion: gu.Ptr(true),
				ClockSkew:                durationpb.New(5 * time.Second),
				AdditionalOrigins:        []string{"https://origin"},
				SkipNativeAppSuccessPage: gu.Ptr(true),
				BackChannelLogoutUri:     gu.Ptr("https://backchannel"),
				LoginVersion: &app.LoginVersion{Version: &app.LoginVersion_LoginV2{
					LoginV2: &app.LoginV2{BaseUri: gu.Ptr("https://login")},
				}},
			},
			expectedModel: &domain.OIDCApp{
				ObjectRoot:               models.ObjectRoot{AggregateID: "proj1"},
				AppID:                    "app1",
				RedirectUris:             []string{"https://redirect"},
				ResponseTypes:            []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
				GrantTypes:               []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
				ApplicationType:          gu.Ptr(domain.OIDCApplicationTypeWeb),
				AuthMethodType:           gu.Ptr(domain.OIDCAuthMethodTypeBasic),
				PostLogoutRedirectUris:   []string{"https://logout"},
				DevMode:                  gu.Ptr(true),
				AccessTokenType:          gu.Ptr(domain.OIDCTokenTypeBearer),
				AccessTokenRoleAssertion: gu.Ptr(true),
				IDTokenRoleAssertion:     gu.Ptr(true),
				IDTokenUserinfoAssertion: gu.Ptr(true),
				ClockSkew:                gu.Ptr(5 * time.Second),
				AdditionalOrigins:        []string{"https://origin"},
				SkipNativeAppSuccessPage: gu.Ptr(true),
				BackChannelLogoutURI:     gu.Ptr("https://backchannel"),
				LoginVersion:             gu.Ptr(domain.LoginVersion2),
				LoginBaseURI:             gu.Ptr("https://login"),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// When
			got, err := UpdateOIDCAppConfigRequestToDomain(tc.appID, tc.projectID, tc.req)

			// Then
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedModel, got)
		})
	}
}

func TestOIDCResponseTypesToDomain(t *testing.T) {
	t.Parallel()

	tt := []struct {
		testName          string
		inputResponseType []app.OIDCResponseType
		expectedResponse  []domain.OIDCResponseType
	}{
		{
			testName:          "empty response types",
			inputResponseType: []app.OIDCResponseType{},
			expectedResponse:  []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
		},
		{
			testName: "all response types",
			inputResponseType: []app.OIDCResponseType{
				app.OIDCResponseType_OIDC_RESPONSE_TYPE_UNSPECIFIED,
				app.OIDCResponseType_OIDC_RESPONSE_TYPE_CODE,
				app.OIDCResponseType_OIDC_RESPONSE_TYPE_ID_TOKEN,
				app.OIDCResponseType_OIDC_RESPONSE_TYPE_ID_TOKEN_TOKEN,
			},
			expectedResponse: []domain.OIDCResponseType{
				domain.OIDCResponseTypeUnspecified,
				domain.OIDCResponseTypeCode,
				domain.OIDCResponseTypeIDToken,
				domain.OIDCResponseTypeIDTokenToken,
			},
		},
		{
			testName: "single response type",
			inputResponseType: []app.OIDCResponseType{
				app.OIDCResponseType_OIDC_RESPONSE_TYPE_CODE,
			},
			expectedResponse: []domain.OIDCResponseType{
				domain.OIDCResponseTypeCode,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// When
			res := oidcResponseTypesToDomain(tc.inputResponseType)

			// Then
			assert.Equal(t, tc.expectedResponse, res)
		})
	}
}

func TestOIDCGrantTypesToDomain(t *testing.T) {
	t.Parallel()

	tt := []struct {
		testName       string
		inputGrantType []app.OIDCGrantType
		expectedGrants []domain.OIDCGrantType
	}{
		{
			testName:       "empty grant types",
			inputGrantType: []app.OIDCGrantType{},
			expectedGrants: []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
		},
		{
			testName: "all grant types",
			inputGrantType: []app.OIDCGrantType{
				app.OIDCGrantType_OIDC_GRANT_TYPE_AUTHORIZATION_CODE,
				app.OIDCGrantType_OIDC_GRANT_TYPE_IMPLICIT,
				app.OIDCGrantType_OIDC_GRANT_TYPE_REFRESH_TOKEN,
				app.OIDCGrantType_OIDC_GRANT_TYPE_DEVICE_CODE,
				app.OIDCGrantType_OIDC_GRANT_TYPE_TOKEN_EXCHANGE,
			},
			expectedGrants: []domain.OIDCGrantType{
				domain.OIDCGrantTypeAuthorizationCode,
				domain.OIDCGrantTypeImplicit,
				domain.OIDCGrantTypeRefreshToken,
				domain.OIDCGrantTypeDeviceCode,
				domain.OIDCGrantTypeTokenExchange,
			},
		},
		{
			testName: "single grant type",
			inputGrantType: []app.OIDCGrantType{
				app.OIDCGrantType_OIDC_GRANT_TYPE_AUTHORIZATION_CODE,
			},
			expectedGrants: []domain.OIDCGrantType{
				domain.OIDCGrantTypeAuthorizationCode,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// When
			res := oidcGrantTypesToDomain(tc.inputGrantType)

			// Then
			assert.Equal(t, tc.expectedGrants, res)
		})
	}
}

func TestOIDCApplicationTypeToDomain(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		appType  app.OIDCAppType
		expected domain.OIDCApplicationType
	}{
		{
			name:     "web type",
			appType:  app.OIDCAppType_OIDC_APP_TYPE_WEB,
			expected: domain.OIDCApplicationTypeWeb,
		},
		{
			name:     "user agent type",
			appType:  app.OIDCAppType_OIDC_APP_TYPE_USER_AGENT,
			expected: domain.OIDCApplicationTypeUserAgent,
		},
		{
			name:     "native type",
			appType:  app.OIDCAppType_OIDC_APP_TYPE_NATIVE,
			expected: domain.OIDCApplicationTypeNative,
		},
		{
			name:     "unspecified type defaults to web",
			expected: domain.OIDCApplicationTypeWeb,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// When
			result := oidcApplicationTypeToDomain(tc.appType)

			// Then
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestOIDCAuthMethodTypeToDomain(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name             string
		authType         app.OIDCAuthMethodType
		expectedResponse domain.OIDCAuthMethodType
	}{
		{
			name:             "basic auth type",
			authType:         app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_BASIC,
			expectedResponse: domain.OIDCAuthMethodTypeBasic,
		},
		{
			name:             "post auth type",
			authType:         app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_POST,
			expectedResponse: domain.OIDCAuthMethodTypePost,
		},
		{
			name:             "none auth type",
			authType:         app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_NONE,
			expectedResponse: domain.OIDCAuthMethodTypeNone,
		},
		{
			name:             "private key jwt auth type",
			authType:         app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT,
			expectedResponse: domain.OIDCAuthMethodTypePrivateKeyJWT,
		},
		{
			name:             "unspecified auth type defaults to basic",
			expectedResponse: domain.OIDCAuthMethodTypeBasic,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// When
			res := oidcAuthMethodTypeToDomain(tc.authType)

			// Then
			assert.Equal(t, tc.expectedResponse, res)
		})
	}
}

func TestOIDCTokenTypeToDomain(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name         string
		tokenType    app.OIDCTokenType
		expectedType domain.OIDCTokenType
	}{
		{
			name:         "bearer token type",
			tokenType:    app.OIDCTokenType_OIDC_TOKEN_TYPE_BEARER,
			expectedType: domain.OIDCTokenTypeBearer,
		},
		{
			name:         "jwt token type",
			tokenType:    app.OIDCTokenType_OIDC_TOKEN_TYPE_JWT,
			expectedType: domain.OIDCTokenTypeJWT,
		},
		{
			name:         "unspecified defaults to bearer",
			expectedType: domain.OIDCTokenTypeBearer,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// When
			result := oidcTokenTypeToDomain(tc.tokenType)

			// Then
			assert.Equal(t, tc.expectedType, result)
		})
	}
}
func TestAppOIDCConfigToPb(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		input    *query.OIDCApp
		expected *app.Application_OidcConfig
	}{
		{
			name:  "empty config",
			input: &query.OIDCApp{},
			expected: &app.Application_OidcConfig{
				OidcConfig: &app.OIDCConfig{
					Version:            app.OIDCVersion_OIDC_VERSION_1_0,
					ComplianceProblems: []*app.OIDCLocalizedMessage{},
					ClockSkew:          durationpb.New(0),
					ResponseTypes:      []app.OIDCResponseType{},
					GrantTypes:         []app.OIDCGrantType{},
				},
			},
		},
		{
			name: "full config",
			input: &query.OIDCApp{
				RedirectURIs:             []string{"https://example.com/callback"},
				ResponseTypes:            []domain.OIDCResponseType{domain.OIDCResponseTypeCode},
				GrantTypes:               []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
				AppType:                  domain.OIDCApplicationTypeWeb,
				ClientID:                 "client123",
				AuthMethodType:           domain.OIDCAuthMethodTypeBasic,
				PostLogoutRedirectURIs:   []string{"https://example.com/logout"},
				ComplianceProblems:       []string{"problem1", "problem2"},
				IsDevMode:                true,
				AccessTokenType:          domain.OIDCTokenTypeBearer,
				AssertAccessTokenRole:    true,
				AssertIDTokenRole:        true,
				AssertIDTokenUserinfo:    true,
				ClockSkew:                5 * time.Second,
				AdditionalOrigins:        []string{"https://app.example.com"},
				AllowedOrigins:           []string{"https://allowed.example.com"},
				SkipNativeAppSuccessPage: true,
				BackChannelLogoutURI:     "https://example.com/backchannel",
				LoginVersion:             domain.LoginVersion2,
				LoginBaseURI:             gu.Ptr("https://login.example.com"),
			},
			expected: &app.Application_OidcConfig{
				OidcConfig: &app.OIDCConfig{
					RedirectUris:           []string{"https://example.com/callback"},
					ResponseTypes:          []app.OIDCResponseType{app.OIDCResponseType_OIDC_RESPONSE_TYPE_CODE},
					GrantTypes:             []app.OIDCGrantType{app.OIDCGrantType_OIDC_GRANT_TYPE_AUTHORIZATION_CODE},
					AppType:                app.OIDCAppType_OIDC_APP_TYPE_WEB,
					ClientId:               "client123",
					AuthMethodType:         app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_BASIC,
					PostLogoutRedirectUris: []string{"https://example.com/logout"},
					Version:                app.OIDCVersion_OIDC_VERSION_1_0,
					NoneCompliant:          true,
					ComplianceProblems: []*app.OIDCLocalizedMessage{
						{Key: "problem1"},
						{Key: "problem2"},
					},
					DevMode:                  true,
					AccessTokenType:          app.OIDCTokenType_OIDC_TOKEN_TYPE_BEARER,
					AccessTokenRoleAssertion: true,
					IdTokenRoleAssertion:     true,
					IdTokenUserinfoAssertion: true,
					ClockSkew:                durationpb.New(5 * time.Second),
					AdditionalOrigins:        []string{"https://app.example.com"},
					AllowedOrigins:           []string{"https://allowed.example.com"},
					SkipNativeAppSuccessPage: true,
					BackChannelLogoutUri:     "https://example.com/backchannel",
					LoginVersion: &app.LoginVersion{
						Version: &app.LoginVersion_LoginV2{
							LoginV2: &app.LoginV2{
								BaseUri: gu.Ptr("https://login.example.com"),
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tt {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// When
			result := appOIDCConfigToPb(tt.input)

			// Then
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestOIDCResponseTypesFromModel(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name          string
		responseTypes []domain.OIDCResponseType
		expected      []app.OIDCResponseType
	}{
		{
			name:          "empty response types",
			responseTypes: []domain.OIDCResponseType{},
			expected:      []app.OIDCResponseType{},
		},
		{
			name: "all response types",
			responseTypes: []domain.OIDCResponseType{
				domain.OIDCResponseTypeUnspecified,
				domain.OIDCResponseTypeCode,
				domain.OIDCResponseTypeIDToken,
				domain.OIDCResponseTypeIDTokenToken,
			},
			expected: []app.OIDCResponseType{
				app.OIDCResponseType_OIDC_RESPONSE_TYPE_UNSPECIFIED,
				app.OIDCResponseType_OIDC_RESPONSE_TYPE_CODE,
				app.OIDCResponseType_OIDC_RESPONSE_TYPE_ID_TOKEN,
				app.OIDCResponseType_OIDC_RESPONSE_TYPE_ID_TOKEN_TOKEN,
			},
		},
		{
			name: "single response type",
			responseTypes: []domain.OIDCResponseType{
				domain.OIDCResponseTypeCode,
			},
			expected: []app.OIDCResponseType{
				app.OIDCResponseType_OIDC_RESPONSE_TYPE_CODE,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// When
			result := oidcResponseTypesFromModel(tc.responseTypes)

			// Then
			assert.Equal(t, tc.expected, result)
		})
	}
}
func TestOIDCGrantTypesFromModel(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name       string
		grantTypes []domain.OIDCGrantType
		expected   []app.OIDCGrantType
	}{
		{
			name:       "empty grant types",
			grantTypes: []domain.OIDCGrantType{},
			expected:   []app.OIDCGrantType{},
		},
		{
			name: "all grant types",
			grantTypes: []domain.OIDCGrantType{
				domain.OIDCGrantTypeAuthorizationCode,
				domain.OIDCGrantTypeImplicit,
				domain.OIDCGrantTypeRefreshToken,
				domain.OIDCGrantTypeDeviceCode,
				domain.OIDCGrantTypeTokenExchange,
			},
			expected: []app.OIDCGrantType{
				app.OIDCGrantType_OIDC_GRANT_TYPE_AUTHORIZATION_CODE,
				app.OIDCGrantType_OIDC_GRANT_TYPE_IMPLICIT,
				app.OIDCGrantType_OIDC_GRANT_TYPE_REFRESH_TOKEN,
				app.OIDCGrantType_OIDC_GRANT_TYPE_DEVICE_CODE,
				app.OIDCGrantType_OIDC_GRANT_TYPE_TOKEN_EXCHANGE,
			},
		},
		{
			name: "single grant type",
			grantTypes: []domain.OIDCGrantType{
				domain.OIDCGrantTypeAuthorizationCode,
			},
			expected: []app.OIDCGrantType{
				app.OIDCGrantType_OIDC_GRANT_TYPE_AUTHORIZATION_CODE,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// When
			result := oidcGrantTypesFromModel(tc.grantTypes)

			// Then
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestOIDCApplicationTypeToPb(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		appType  domain.OIDCApplicationType
		expected app.OIDCAppType
	}{
		{
			name:     "web type",
			appType:  domain.OIDCApplicationTypeWeb,
			expected: app.OIDCAppType_OIDC_APP_TYPE_WEB,
		},
		{
			name:     "user agent type",
			appType:  domain.OIDCApplicationTypeUserAgent,
			expected: app.OIDCAppType_OIDC_APP_TYPE_USER_AGENT,
		},
		{
			name:     "native type",
			appType:  domain.OIDCApplicationTypeNative,
			expected: app.OIDCAppType_OIDC_APP_TYPE_NATIVE,
		},
		{
			name:     "unspecified type defaults to web",
			appType:  domain.OIDCApplicationType(999), // Invalid value to trigger default case
			expected: app.OIDCAppType_OIDC_APP_TYPE_WEB,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// When
			result := oidcApplicationTypeToPb(tc.appType)

			// Then
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestOIDCAuthMethodTypeToPb(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		authType domain.OIDCAuthMethodType
		expected app.OIDCAuthMethodType
	}{
		{
			name:     "basic auth type",
			authType: domain.OIDCAuthMethodTypeBasic,
			expected: app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_BASIC,
		},
		{
			name:     "post auth type",
			authType: domain.OIDCAuthMethodTypePost,
			expected: app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_POST,
		},
		{
			name:     "none auth type",
			authType: domain.OIDCAuthMethodTypeNone,
			expected: app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_NONE,
		},
		{
			name:     "private key jwt auth type",
			authType: domain.OIDCAuthMethodTypePrivateKeyJWT,
			expected: app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT,
		},
		{
			name:     "unknown auth type defaults to basic",
			authType: domain.OIDCAuthMethodType(999),
			expected: app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_BASIC,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// When
			result := oidcAuthMethodTypeToPb(tc.authType)

			// Then
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestOIDCTokenTypeToPb(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name      string
		tokenType domain.OIDCTokenType
		expected  app.OIDCTokenType
	}{
		{
			name:      "bearer token type",
			tokenType: domain.OIDCTokenTypeBearer,
			expected:  app.OIDCTokenType_OIDC_TOKEN_TYPE_BEARER,
		},
		{
			name:      "jwt token type",
			tokenType: domain.OIDCTokenTypeJWT,
			expected:  app.OIDCTokenType_OIDC_TOKEN_TYPE_JWT,
		},
		{
			name:      "unknown token type defaults to bearer",
			tokenType: domain.OIDCTokenType(999), // Invalid value to trigger default case
			expected:  app.OIDCTokenType_OIDC_TOKEN_TYPE_BEARER,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// When
			result := oidcTokenTypeToPb(tc.tokenType)

			// Then
			assert.Equal(t, tc.expected, result)
		})
	}
}
