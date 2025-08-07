package convert

import (
	"github.com/muhlemmer/gu"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/query"
	app "github.com/zitadel/zitadel/pkg/grpc/app/v2beta"
)

func CreateOIDCAppRequestToDomain(name, projectID string, req *app.CreateOIDCApplicationRequest) (*domain.OIDCApp, error) {
	loginVersion, loginBaseURI, err := loginVersionToDomain(req.GetLoginVersion())
	if err != nil {
		return nil, err
	}
	return &domain.OIDCApp{
		ObjectRoot: models.ObjectRoot{
			AggregateID: projectID,
		},
		AppName:                  name,
		OIDCVersion:              gu.Ptr(domain.OIDCVersionV1),
		RedirectUris:             req.GetRedirectUris(),
		ResponseTypes:            oidcResponseTypesToDomain(req.GetResponseTypes()),
		GrantTypes:               oidcGrantTypesToDomain(req.GetGrantTypes()),
		ApplicationType:          gu.Ptr(oidcApplicationTypeToDomain(req.GetAppType())),
		AuthMethodType:           gu.Ptr(oidcAuthMethodTypeToDomain(req.GetAuthMethodType())),
		PostLogoutRedirectUris:   req.GetPostLogoutRedirectUris(),
		DevMode:                  &req.DevMode,
		AccessTokenType:          gu.Ptr(oidcTokenTypeToDomain(req.GetAccessTokenType())),
		AccessTokenRoleAssertion: gu.Ptr(req.GetAccessTokenRoleAssertion()),
		IDTokenRoleAssertion:     gu.Ptr(req.GetIdTokenRoleAssertion()),
		IDTokenUserinfoAssertion: gu.Ptr(req.GetIdTokenUserinfoAssertion()),
		ClockSkew:                gu.Ptr(req.GetClockSkew().AsDuration()),
		AdditionalOrigins:        req.GetAdditionalOrigins(),
		SkipNativeAppSuccessPage: gu.Ptr(req.GetSkipNativeAppSuccessPage()),
		BackChannelLogoutURI:     gu.Ptr(req.GetBackChannelLogoutUri()),
		LoginVersion:             loginVersion,
		LoginBaseURI:             loginBaseURI,
	}, nil
}

func UpdateOIDCAppConfigRequestToDomain(appID, projectID string, app *app.UpdateOIDCApplicationConfigurationRequest) (*domain.OIDCApp, error) {
	loginVersion, loginBaseURI, err := loginVersionToDomain(app.GetLoginVersion())
	if err != nil {
		return nil, err
	}
	return &domain.OIDCApp{
		ObjectRoot: models.ObjectRoot{
			AggregateID: projectID,
		},
		AppID:                    appID,
		RedirectUris:             app.RedirectUris,
		ResponseTypes:            oidcResponseTypesToDomain(app.ResponseTypes),
		GrantTypes:               oidcGrantTypesToDomain(app.GrantTypes),
		ApplicationType:          oidcApplicationTypeToDomainPtr(app.AppType),
		AuthMethodType:           oidcAuthMethodTypeToDomainPtr(app.AuthMethodType),
		PostLogoutRedirectUris:   app.PostLogoutRedirectUris,
		DevMode:                  app.DevMode,
		AccessTokenType:          oidcTokenTypeToDomainPtr(app.AccessTokenType),
		AccessTokenRoleAssertion: app.AccessTokenRoleAssertion,
		IDTokenRoleAssertion:     app.IdTokenRoleAssertion,
		IDTokenUserinfoAssertion: app.IdTokenUserinfoAssertion,
		ClockSkew:                gu.Ptr(app.GetClockSkew().AsDuration()),
		AdditionalOrigins:        app.AdditionalOrigins,
		SkipNativeAppSuccessPage: app.SkipNativeAppSuccessPage,
		BackChannelLogoutURI:     app.BackChannelLogoutUri,
		LoginVersion:             loginVersion,
		LoginBaseURI:             loginBaseURI,
	}, nil
}

func oidcResponseTypesToDomain(responseTypes []app.OIDCResponseType) []domain.OIDCResponseType {
	if len(responseTypes) == 0 {
		return []domain.OIDCResponseType{domain.OIDCResponseTypeCode}
	}
	oidcResponseTypes := make([]domain.OIDCResponseType, len(responseTypes))
	for i, responseType := range responseTypes {
		switch responseType {
		case app.OIDCResponseType_OIDC_RESPONSE_TYPE_UNSPECIFIED:
			oidcResponseTypes[i] = domain.OIDCResponseTypeUnspecified
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

func oidcGrantTypesToDomain(grantTypes []app.OIDCGrantType) []domain.OIDCGrantType {
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

func oidcApplicationTypeToDomainPtr(appType *app.OIDCAppType) *domain.OIDCApplicationType {
	if appType == nil {
		return nil
	}

	res := oidcApplicationTypeToDomain(*appType)
	return &res
}

func oidcApplicationTypeToDomain(appType app.OIDCAppType) domain.OIDCApplicationType {
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

func oidcAuthMethodTypeToDomainPtr(authType *app.OIDCAuthMethodType) *domain.OIDCAuthMethodType {
	if authType == nil {
		return nil
	}

	res := oidcAuthMethodTypeToDomain(*authType)
	return &res
}

func oidcAuthMethodTypeToDomain(authType app.OIDCAuthMethodType) domain.OIDCAuthMethodType {
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

func oidcTokenTypeToDomainPtr(tokenType *app.OIDCTokenType) *domain.OIDCTokenType {
	if tokenType == nil {
		return nil
	}

	res := oidcTokenTypeToDomain(*tokenType)
	return &res
}

func oidcTokenTypeToDomain(tokenType app.OIDCTokenType) domain.OIDCTokenType {
	switch tokenType {
	case app.OIDCTokenType_OIDC_TOKEN_TYPE_BEARER:
		return domain.OIDCTokenTypeBearer
	case app.OIDCTokenType_OIDC_TOKEN_TYPE_JWT:
		return domain.OIDCTokenTypeJWT
	default:
		return domain.OIDCTokenTypeBearer
	}
}

func ComplianceProblemsToLocalizedMessages(complianceProblems []string) []*app.OIDCLocalizedMessage {
	converted := make([]*app.OIDCLocalizedMessage, len(complianceProblems))
	for i, p := range complianceProblems {
		converted[i] = &app.OIDCLocalizedMessage{Key: p}
	}

	return converted
}

func appOIDCConfigToPb(oidcApp *query.OIDCApp) *app.Application_OidcConfig {
	return &app.Application_OidcConfig{
		OidcConfig: &app.OIDCConfig{
			RedirectUris:             oidcApp.RedirectURIs,
			ResponseTypes:            oidcResponseTypesFromModel(oidcApp.ResponseTypes),
			GrantTypes:               oidcGrantTypesFromModel(oidcApp.GrantTypes),
			AppType:                  oidcApplicationTypeToPb(oidcApp.AppType),
			ClientId:                 oidcApp.ClientID,
			AuthMethodType:           oidcAuthMethodTypeToPb(oidcApp.AuthMethodType),
			PostLogoutRedirectUris:   oidcApp.PostLogoutRedirectURIs,
			Version:                  app.OIDCVersion_OIDC_VERSION_1_0,
			NoneCompliant:            len(oidcApp.ComplianceProblems) != 0,
			ComplianceProblems:       ComplianceProblemsToLocalizedMessages(oidcApp.ComplianceProblems),
			DevMode:                  oidcApp.IsDevMode,
			AccessTokenType:          oidcTokenTypeToPb(oidcApp.AccessTokenType),
			AccessTokenRoleAssertion: oidcApp.AssertAccessTokenRole,
			IdTokenRoleAssertion:     oidcApp.AssertIDTokenRole,
			IdTokenUserinfoAssertion: oidcApp.AssertIDTokenUserinfo,
			ClockSkew:                durationpb.New(oidcApp.ClockSkew),
			AdditionalOrigins:        oidcApp.AdditionalOrigins,
			AllowedOrigins:           oidcApp.AllowedOrigins,
			SkipNativeAppSuccessPage: oidcApp.SkipNativeAppSuccessPage,
			BackChannelLogoutUri:     oidcApp.BackChannelLogoutURI,
			LoginVersion:             loginVersionToPb(oidcApp.LoginVersion, oidcApp.LoginBaseURI),
		},
	}
}

func oidcResponseTypesFromModel(responseTypes []domain.OIDCResponseType) []app.OIDCResponseType {
	oidcResponseTypes := make([]app.OIDCResponseType, len(responseTypes))
	for i, responseType := range responseTypes {
		switch responseType {
		case domain.OIDCResponseTypeUnspecified:
			oidcResponseTypes[i] = app.OIDCResponseType_OIDC_RESPONSE_TYPE_UNSPECIFIED
		case domain.OIDCResponseTypeCode:
			oidcResponseTypes[i] = app.OIDCResponseType_OIDC_RESPONSE_TYPE_CODE
		case domain.OIDCResponseTypeIDToken:
			oidcResponseTypes[i] = app.OIDCResponseType_OIDC_RESPONSE_TYPE_ID_TOKEN
		case domain.OIDCResponseTypeIDTokenToken:
			oidcResponseTypes[i] = app.OIDCResponseType_OIDC_RESPONSE_TYPE_ID_TOKEN_TOKEN
		}
	}
	return oidcResponseTypes
}

func oidcGrantTypesFromModel(grantTypes []domain.OIDCGrantType) []app.OIDCGrantType {
	oidcGrantTypes := make([]app.OIDCGrantType, len(grantTypes))
	for i, grantType := range grantTypes {
		switch grantType {
		case domain.OIDCGrantTypeAuthorizationCode:
			oidcGrantTypes[i] = app.OIDCGrantType_OIDC_GRANT_TYPE_AUTHORIZATION_CODE
		case domain.OIDCGrantTypeImplicit:
			oidcGrantTypes[i] = app.OIDCGrantType_OIDC_GRANT_TYPE_IMPLICIT
		case domain.OIDCGrantTypeRefreshToken:
			oidcGrantTypes[i] = app.OIDCGrantType_OIDC_GRANT_TYPE_REFRESH_TOKEN
		case domain.OIDCGrantTypeDeviceCode:
			oidcGrantTypes[i] = app.OIDCGrantType_OIDC_GRANT_TYPE_DEVICE_CODE
		case domain.OIDCGrantTypeTokenExchange:
			oidcGrantTypes[i] = app.OIDCGrantType_OIDC_GRANT_TYPE_TOKEN_EXCHANGE
		}
	}
	return oidcGrantTypes
}

func oidcApplicationTypeToPb(appType domain.OIDCApplicationType) app.OIDCAppType {
	switch appType {
	case domain.OIDCApplicationTypeWeb:
		return app.OIDCAppType_OIDC_APP_TYPE_WEB
	case domain.OIDCApplicationTypeUserAgent:
		return app.OIDCAppType_OIDC_APP_TYPE_USER_AGENT
	case domain.OIDCApplicationTypeNative:
		return app.OIDCAppType_OIDC_APP_TYPE_NATIVE
	default:
		return app.OIDCAppType_OIDC_APP_TYPE_WEB
	}
}

func oidcAuthMethodTypeToPb(authType domain.OIDCAuthMethodType) app.OIDCAuthMethodType {
	switch authType {
	case domain.OIDCAuthMethodTypeBasic:
		return app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_BASIC
	case domain.OIDCAuthMethodTypePost:
		return app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_POST
	case domain.OIDCAuthMethodTypeNone:
		return app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_NONE
	case domain.OIDCAuthMethodTypePrivateKeyJWT:
		return app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT
	default:
		return app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_BASIC
	}
}

func oidcTokenTypeToPb(tokenType domain.OIDCTokenType) app.OIDCTokenType {
	switch tokenType {
	case domain.OIDCTokenTypeBearer:
		return app.OIDCTokenType_OIDC_TOKEN_TYPE_BEARER
	case domain.OIDCTokenTypeJWT:
		return app.OIDCTokenType_OIDC_TOKEN_TYPE_JWT
	default:
		return app.OIDCTokenType_OIDC_TOKEN_TYPE_BEARER
	}
}
