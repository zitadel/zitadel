package convert

import (
	"net/url"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	app "github.com/zitadel/zitadel/pkg/grpc/app/v2beta"
)

func CreateOIDCAppRequestToDomain(name, projectID string, req *app.CreateOIDCApplicationRequest) (*domain.OIDCApp, error) {
	loginVersion, loginBaseURI, err := LoginVersionToDomain(req.GetLoginVersion())
	if err != nil {
		return nil, err
	}
	return &domain.OIDCApp{
		ObjectRoot: models.ObjectRoot{
			AggregateID: projectID,
		},
		AppName:                  name,
		OIDCVersion:              OIDCVersionToDomain(req.GetVersion()),
		RedirectUris:             req.GetRedirectUris(),
		ResponseTypes:            OIDCResponseTypesToDomain(req.GetResponseTypes()),
		GrantTypes:               OIDCGrantTypesToDomain(req.GetGrantTypes()),
		ApplicationType:          OIDCApplicationTypeToDomain(req.GetAppType()),
		AuthMethodType:           OIDCAuthMethodTypeToDomain(req.GetAuthMethodType()),
		PostLogoutRedirectUris:   req.GetPostLogoutRedirectUris(),
		DevMode:                  req.GetDevMode(),
		AccessTokenType:          OIDCTokenTypeToDomain(req.GetAccessTokenType()),
		AccessTokenRoleAssertion: req.GetAccessTokenRoleAssertion(),
		IDTokenRoleAssertion:     req.GetIdTokenRoleAssertion(),
		IDTokenUserinfoAssertion: req.GetIdTokenUserinfoAssertion(),
		ClockSkew:                req.GetClockSkew().AsDuration(),
		AdditionalOrigins:        req.GetAdditionalOrigins(),
		SkipNativeAppSuccessPage: req.GetSkipNativeAppSuccessPage(),
		BackChannelLogoutURI:     req.GetBackChannelLogoutUri(),
		LoginVersion:             loginVersion,
		LoginBaseURI:             loginBaseURI,
	}, nil
}

func LoginVersionToDomain(version *app.LoginVersion) (domain.LoginVersion, string, error) {
	switch v := version.GetVersion().(type) {
	case nil:
		return domain.LoginVersionUnspecified, "", nil
	case *app.LoginVersion_LoginV1:
		return domain.LoginVersion1, "", nil
	case *app.LoginVersion_LoginV2:
		_, err := url.Parse(v.LoginV2.GetBaseUri())
		return domain.LoginVersion2, v.LoginV2.GetBaseUri(), err
	default:
		return domain.LoginVersionUnspecified, "", nil
	}
}

func OIDCVersionToDomain(version app.OIDCVersion) domain.OIDCVersion {
	switch version {
	case app.OIDCVersion_OIDC_VERSION_1_0:
		return domain.OIDCVersionV1
	}
	return domain.OIDCVersionV1
}

func OIDCResponseTypesToDomain(responseTypes []app.OIDCResponseType) []domain.OIDCResponseType {
	if len(responseTypes) == 0 {
		return []domain.OIDCResponseType{domain.OIDCResponseTypeCode}
	}
	oidcResponseTypes := make([]domain.OIDCResponseType, len(responseTypes))
	for i, responseType := range responseTypes {
		switch responseType {
		case app.OIDCResponseType_OIDC_RESPONSE_TYPE_CODE:
			oidcResponseTypes[i] = domain.OIDCResponseTypeCode
		case app.OIDCResponseType_OIDC_RESPONSE_TYPE_ID_TOKEN:
			oidcResponseTypes[i] = domain.OIDCResponseTypeIDToken
		case app.OIDCResponseType_OIDC_RESPONSE_TYPE_ID_TOKEN_TOKEN:
			oidcResponseTypes[i] = domain.OIDCResponseTypeIDTokenToken
		}
	}
	return oidcResponseTypes
}

func OIDCGrantTypesToDomain(grantTypes []app.OIDCGrantType) []domain.OIDCGrantType {
	if len(grantTypes) == 0 {
		return []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode}
	}
	oidcGrantTypes := make([]domain.OIDCGrantType, len(grantTypes))
	for i, grantType := range grantTypes {
		switch grantType {
		case app.OIDCGrantType_OIDC_GRANT_TYPE_AUTHORIZATION_CODE:
			oidcGrantTypes[i] = domain.OIDCGrantTypeAuthorizationCode
		case app.OIDCGrantType_OIDC_GRANT_TYPE_IMPLICIT:
			oidcGrantTypes[i] = domain.OIDCGrantTypeImplicit
		case app.OIDCGrantType_OIDC_GRANT_TYPE_REFRESH_TOKEN:
			oidcGrantTypes[i] = domain.OIDCGrantTypeRefreshToken
		case app.OIDCGrantType_OIDC_GRANT_TYPE_DEVICE_CODE:
			oidcGrantTypes[i] = domain.OIDCGrantTypeDeviceCode
		case app.OIDCGrantType_OIDC_GRANT_TYPE_TOKEN_EXCHANGE:
			oidcGrantTypes[i] = domain.OIDCGrantTypeTokenExchange
		}
	}
	return oidcGrantTypes
}

func OIDCApplicationTypeToDomain(appType app.OIDCAppType) domain.OIDCApplicationType {
	switch appType {
	case app.OIDCAppType_OIDC_APP_TYPE_WEB:
		return domain.OIDCApplicationTypeWeb
	case app.OIDCAppType_OIDC_APP_TYPE_USER_AGENT:
		return domain.OIDCApplicationTypeUserAgent
	case app.OIDCAppType_OIDC_APP_TYPE_NATIVE:
		return domain.OIDCApplicationTypeNative
	}
	return domain.OIDCApplicationTypeWeb
}

func OIDCAuthMethodTypeToDomain(authType app.OIDCAuthMethodType) domain.OIDCAuthMethodType {
	switch authType {
	case app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_BASIC:
		return domain.OIDCAuthMethodTypeBasic
	case app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_POST:
		return domain.OIDCAuthMethodTypePost
	case app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_NONE:
		return domain.OIDCAuthMethodTypeNone
	case app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT:
		return domain.OIDCAuthMethodTypePrivateKeyJWT
	default:
		return domain.OIDCAuthMethodTypeBasic
	}
}

func OIDCTokenTypeToDomain(tokenType app.OIDCTokenType) domain.OIDCTokenType {
	switch tokenType {
	case app.OIDCTokenType_OIDC_TOKEN_TYPE_BEARER:
		return domain.OIDCTokenTypeBearer
	case app.OIDCTokenType_OIDC_TOKEN_TYPE_JWT:
		return domain.OIDCTokenTypeJWT
	default:
		return domain.OIDCTokenTypeBearer
	}
}

func ComplianceProblemsToLocalizedMessages(complianceProblems []string) []*app.CreateOIDCApplicationResponse_OIDCApplicationLocalizedMessage {
	converted := make([]*app.CreateOIDCApplicationResponse_OIDCApplicationLocalizedMessage, len(complianceProblems))
	for i, p := range complianceProblems {
		converted[i] = &app.CreateOIDCApplicationResponse_OIDCApplicationLocalizedMessage{Key: p}
	}

	return converted
}
