package convert

import (
	"github.com/muhlemmer/gu"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/application/v2"
)

func CreateOIDCAppRequestToDomain(name, appID, projectID string, req *application.CreateOIDCApplicationRequest) (*domain.OIDCApp, error) {
	loginVersion, loginBaseURI, err := loginVersionToDomain(req.GetLoginVersion())
	if err != nil {
		return nil, err
	}
	return &domain.OIDCApp{
		ObjectRoot: models.ObjectRoot{
			AggregateID: projectID,
		},
		AppID:                    appID,
		AppName:                  name,
		OIDCVersion:              gu.Ptr(domain.OIDCVersionV1),
		RedirectUris:             req.GetRedirectUris(),
		ResponseTypes:            oidcResponseTypesToDomain(req.GetResponseTypes()),
		GrantTypes:               oidcGrantTypesToDomain(req.GetGrantTypes()),
		ApplicationType:          gu.Ptr(oidcApplicationTypeToDomain(req.GetApplicationType())),
		AuthMethodType:           gu.Ptr(oidcAuthMethodTypeToDomain(req.GetAuthMethodType())),
		PostLogoutRedirectUris:   req.GetPostLogoutRedirectUris(),
		DevMode:                  &req.DevelopmentMode,
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

func UpdateOIDCAppConfigRequestToDomain(appID, projectID string, app *application.UpdateOIDCApplicationConfigurationRequest) (*domain.OIDCApp, error) {
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
		ApplicationType:          oidcApplicationTypeToDomainPtr(app.ApplicationType),
		AuthMethodType:           oidcAuthMethodTypeToDomainPtr(app.AuthMethodType),
		PostLogoutRedirectUris:   app.PostLogoutRedirectUris,
		DevMode:                  app.DevelopmentMode,
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

func oidcResponseTypesToDomain(responseTypes []application.OIDCResponseType) []domain.OIDCResponseType {
	if len(responseTypes) == 0 {
		return []domain.OIDCResponseType{domain.OIDCResponseTypeCode}
	}
	oidcResponseTypes := make([]domain.OIDCResponseType, len(responseTypes))
	for i, responseType := range responseTypes {
		switch responseType {
		case application.OIDCResponseType_OIDC_RESPONSE_TYPE_UNSPECIFIED:
			oidcResponseTypes[i] = domain.OIDCResponseTypeUnspecified
		case application.OIDCResponseType_OIDC_RESPONSE_TYPE_CODE:
			oidcResponseTypes[i] = domain.OIDCResponseTypeCode
		case application.OIDCResponseType_OIDC_RESPONSE_TYPE_ID_TOKEN:
			oidcResponseTypes[i] = domain.OIDCResponseTypeIDToken
		case application.OIDCResponseType_OIDC_RESPONSE_TYPE_ID_TOKEN_TOKEN:
			oidcResponseTypes[i] = domain.OIDCResponseTypeIDTokenToken
		}
	}
	return oidcResponseTypes
}

func oidcGrantTypesToDomain(grantTypes []application.OIDCGrantType) []domain.OIDCGrantType {
	if len(grantTypes) == 0 {
		return []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode}
	}
	oidcGrantTypes := make([]domain.OIDCGrantType, len(grantTypes))
	for i, grantType := range grantTypes {
		switch grantType {
		case application.OIDCGrantType_OIDC_GRANT_TYPE_AUTHORIZATION_CODE:
			oidcGrantTypes[i] = domain.OIDCGrantTypeAuthorizationCode
		case application.OIDCGrantType_OIDC_GRANT_TYPE_IMPLICIT:
			oidcGrantTypes[i] = domain.OIDCGrantTypeImplicit
		case application.OIDCGrantType_OIDC_GRANT_TYPE_REFRESH_TOKEN:
			oidcGrantTypes[i] = domain.OIDCGrantTypeRefreshToken
		case application.OIDCGrantType_OIDC_GRANT_TYPE_DEVICE_CODE:
			oidcGrantTypes[i] = domain.OIDCGrantTypeDeviceCode
		case application.OIDCGrantType_OIDC_GRANT_TYPE_TOKEN_EXCHANGE:
			oidcGrantTypes[i] = domain.OIDCGrantTypeTokenExchange
		}
	}
	return oidcGrantTypes
}

func oidcApplicationTypeToDomainPtr(appType *application.OIDCApplicationType) *domain.OIDCApplicationType {
	if appType == nil {
		return nil
	}

	res := oidcApplicationTypeToDomain(*appType)
	return &res
}

func oidcApplicationTypeToDomain(appType application.OIDCApplicationType) domain.OIDCApplicationType {
	switch appType {
	case application.OIDCApplicationType_OIDC_APP_TYPE_WEB:
		return domain.OIDCApplicationTypeWeb
	case application.OIDCApplicationType_OIDC_APP_TYPE_USER_AGENT:
		return domain.OIDCApplicationTypeUserAgent
	case application.OIDCApplicationType_OIDC_APP_TYPE_NATIVE:
		return domain.OIDCApplicationTypeNative
	default:
		return domain.OIDCApplicationTypeWeb
	}
}

func oidcAuthMethodTypeToDomainPtr(authType *application.OIDCAuthMethodType) *domain.OIDCAuthMethodType {
	if authType == nil {
		return nil
	}

	res := oidcAuthMethodTypeToDomain(*authType)
	return &res
}

func oidcAuthMethodTypeToDomain(authType application.OIDCAuthMethodType) domain.OIDCAuthMethodType {
	switch authType {
	case application.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_BASIC:
		return domain.OIDCAuthMethodTypeBasic
	case application.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_POST:
		return domain.OIDCAuthMethodTypePost
	case application.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_NONE:
		return domain.OIDCAuthMethodTypeNone
	case application.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT:
		return domain.OIDCAuthMethodTypePrivateKeyJWT
	default:
		return domain.OIDCAuthMethodTypeBasic
	}
}

func oidcTokenTypeToDomainPtr(tokenType *application.OIDCTokenType) *domain.OIDCTokenType {
	if tokenType == nil {
		return nil
	}

	res := oidcTokenTypeToDomain(*tokenType)
	return &res
}

func oidcTokenTypeToDomain(tokenType application.OIDCTokenType) domain.OIDCTokenType {
	switch tokenType {
	case application.OIDCTokenType_OIDC_TOKEN_TYPE_BEARER:
		return domain.OIDCTokenTypeBearer
	case application.OIDCTokenType_OIDC_TOKEN_TYPE_JWT:
		return domain.OIDCTokenTypeJWT
	default:
		return domain.OIDCTokenTypeBearer
	}
}

func ComplianceProblemsToLocalizedMessages(complianceProblems []string) []*application.OIDCLocalizedMessage {
	converted := make([]*application.OIDCLocalizedMessage, len(complianceProblems))
	for i, p := range complianceProblems {
		converted[i] = &application.OIDCLocalizedMessage{Key: p}
	}

	return converted
}

func appOIDCConfigToPb(oidcApp *query.OIDCApp) *application.Application_OidcConfiguration {
	return &application.Application_OidcConfiguration{
		OidcConfiguration: &application.OIDCConfiguration{
			RedirectUris:             oidcApp.RedirectURIs,
			ResponseTypes:            oidcResponseTypesFromModel(oidcApp.ResponseTypes),
			GrantTypes:               oidcGrantTypesFromModel(oidcApp.GrantTypes),
			ApplicationType:          oidcApplicationTypeToPb(oidcApp.AppType),
			ClientId:                 oidcApp.ClientID,
			AuthMethodType:           oidcAuthMethodTypeToPb(oidcApp.AuthMethodType),
			PostLogoutRedirectUris:   oidcApp.PostLogoutRedirectURIs,
			Version:                  application.OIDCVersion_OIDC_VERSION_1_0,
			NonCompliant:             len(oidcApp.ComplianceProblems) != 0,
			ComplianceProblems:       ComplianceProblemsToLocalizedMessages(oidcApp.ComplianceProblems),
			DevelopmentMode:          oidcApp.IsDevMode,
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

func oidcResponseTypesFromModel(responseTypes []domain.OIDCResponseType) []application.OIDCResponseType {
	oidcResponseTypes := make([]application.OIDCResponseType, len(responseTypes))
	for i, responseType := range responseTypes {
		switch responseType {
		case domain.OIDCResponseTypeUnspecified:
			oidcResponseTypes[i] = application.OIDCResponseType_OIDC_RESPONSE_TYPE_UNSPECIFIED
		case domain.OIDCResponseTypeCode:
			oidcResponseTypes[i] = application.OIDCResponseType_OIDC_RESPONSE_TYPE_CODE
		case domain.OIDCResponseTypeIDToken:
			oidcResponseTypes[i] = application.OIDCResponseType_OIDC_RESPONSE_TYPE_ID_TOKEN
		case domain.OIDCResponseTypeIDTokenToken:
			oidcResponseTypes[i] = application.OIDCResponseType_OIDC_RESPONSE_TYPE_ID_TOKEN_TOKEN
		}
	}
	return oidcResponseTypes
}

func oidcGrantTypesFromModel(grantTypes []domain.OIDCGrantType) []application.OIDCGrantType {
	oidcGrantTypes := make([]application.OIDCGrantType, len(grantTypes))
	for i, grantType := range grantTypes {
		switch grantType {
		case domain.OIDCGrantTypeAuthorizationCode:
			oidcGrantTypes[i] = application.OIDCGrantType_OIDC_GRANT_TYPE_AUTHORIZATION_CODE
		case domain.OIDCGrantTypeImplicit:
			oidcGrantTypes[i] = application.OIDCGrantType_OIDC_GRANT_TYPE_IMPLICIT
		case domain.OIDCGrantTypeRefreshToken:
			oidcGrantTypes[i] = application.OIDCGrantType_OIDC_GRANT_TYPE_REFRESH_TOKEN
		case domain.OIDCGrantTypeDeviceCode:
			oidcGrantTypes[i] = application.OIDCGrantType_OIDC_GRANT_TYPE_DEVICE_CODE
		case domain.OIDCGrantTypeTokenExchange:
			oidcGrantTypes[i] = application.OIDCGrantType_OIDC_GRANT_TYPE_TOKEN_EXCHANGE
		}
	}
	return oidcGrantTypes
}

func oidcApplicationTypeToPb(appType domain.OIDCApplicationType) application.OIDCApplicationType {
	switch appType {
	case domain.OIDCApplicationTypeWeb:
		return application.OIDCApplicationType_OIDC_APP_TYPE_WEB
	case domain.OIDCApplicationTypeUserAgent:
		return application.OIDCApplicationType_OIDC_APP_TYPE_USER_AGENT
	case domain.OIDCApplicationTypeNative:
		return application.OIDCApplicationType_OIDC_APP_TYPE_NATIVE
	default:
		return application.OIDCApplicationType_OIDC_APP_TYPE_WEB
	}
}

func oidcAuthMethodTypeToPb(authType domain.OIDCAuthMethodType) application.OIDCAuthMethodType {
	switch authType {
	case domain.OIDCAuthMethodTypeBasic:
		return application.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_BASIC
	case domain.OIDCAuthMethodTypePost:
		return application.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_POST
	case domain.OIDCAuthMethodTypeNone:
		return application.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_NONE
	case domain.OIDCAuthMethodTypePrivateKeyJWT:
		return application.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT
	default:
		return application.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_BASIC
	}
}

func oidcTokenTypeToPb(tokenType domain.OIDCTokenType) application.OIDCTokenType {
	switch tokenType {
	case domain.OIDCTokenTypeBearer:
		return application.OIDCTokenType_OIDC_TOKEN_TYPE_BEARER
	case domain.OIDCTokenTypeJWT:
		return application.OIDCTokenType_OIDC_TOKEN_TYPE_JWT
	default:
		return application.OIDCTokenType_OIDC_TOKEN_TYPE_BEARER
	}
}
