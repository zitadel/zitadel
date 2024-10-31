package project

import (
	"google.golang.org/protobuf/types/known/durationpb"

	object_grpc "github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	app_pb "github.com/zitadel/zitadel/pkg/grpc/app"
	message_pb "github.com/zitadel/zitadel/pkg/grpc/message"
)

func AppsToPb(apps []*query.App) []*app_pb.App {
	a := make([]*app_pb.App, len(apps))
	for i, app := range apps {
		a[i] = AppToPb(app)
	}
	return a
}

func AppToPb(app *query.App) *app_pb.App {
	return &app_pb.App{
		Id:      app.ID,
		Details: object_grpc.ToViewDetailsPb(app.Sequence, app.CreationDate, app.ChangeDate, app.ResourceOwner),
		State:   AppStateToPb(app.State),
		Name:    app.Name,
		Config:  AppConfigToPb(app),
	}
}

func AppConfigToPb(app *query.App) app_pb.AppConfig {
	if app.OIDCConfig != nil {
		return AppOIDCConfigToPb(app.OIDCConfig)
	}
	if app.SAMLConfig != nil {
		return AppSAMLConfigToPb(app.SAMLConfig)
	}
	return AppAPIConfigToPb(app.APIConfig)
}

func AppOIDCConfigToPb(app *query.OIDCApp) *app_pb.App_OidcConfig {
	return &app_pb.App_OidcConfig{
		OidcConfig: &app_pb.OIDCConfig{
			RedirectUris:             app.RedirectURIs,
			ResponseTypes:            OIDCResponseTypesFromModel(app.ResponseTypes),
			GrantTypes:               OIDCGrantTypesFromModel(app.GrantTypes),
			AppType:                  OIDCApplicationTypeToPb(app.AppType),
			ClientId:                 app.ClientID,
			AuthMethodType:           OIDCAuthMethodTypeToPb(app.AuthMethodType),
			PostLogoutRedirectUris:   app.PostLogoutRedirectURIs,
			Version:                  OIDCVersionToPb(domain.OIDCVersion(app.Version)),
			NoneCompliant:            len(app.ComplianceProblems) != 0,
			ComplianceProblems:       ComplianceProblemsToLocalizedMessages(app.ComplianceProblems),
			DevMode:                  app.IsDevMode,
			AccessTokenType:          oidcTokenTypeToPb(app.AccessTokenType),
			AccessTokenRoleAssertion: app.AssertAccessTokenRole,
			IdTokenRoleAssertion:     app.AssertIDTokenRole,
			IdTokenUserinfoAssertion: app.AssertIDTokenUserinfo,
			ClockSkew:                durationpb.New(app.ClockSkew),
			AdditionalOrigins:        app.AdditionalOrigins,
			AllowedOrigins:           app.AllowedOrigins,
			SkipNativeAppSuccessPage: app.SkipNativeAppSuccessPage,
			BackChannelLogoutUri:     app.BackChannelLogoutURI,
		},
	}
}

func AppSAMLConfigToPb(app *query.SAMLApp) app_pb.AppConfig {
	return &app_pb.App_SamlConfig{
		SamlConfig: &app_pb.SAMLConfig{
			Metadata: &app_pb.SAMLConfig_MetadataXml{MetadataXml: app.Metadata},
		},
	}
}

func AppAPIConfigToPb(app *query.APIApp) app_pb.AppConfig {
	return &app_pb.App_ApiConfig{
		ApiConfig: &app_pb.APIConfig{
			ClientId:       app.ClientID,
			AuthMethodType: APIAuthMethodeTypeToPb(app.AuthMethodType),
		},
	}
}

func AppStateToPb(state domain.AppState) app_pb.AppState {
	switch state {
	case domain.AppStateActive:
		return app_pb.AppState_APP_STATE_ACTIVE
	case domain.AppStateInactive:
		return app_pb.AppState_APP_STATE_INACTIVE
	default:
		return app_pb.AppState_APP_STATE_UNSPECIFIED
	}
}

func OIDCResponseTypesFromModel(responseTypes []domain.OIDCResponseType) []app_pb.OIDCResponseType {
	oidcResponseTypes := make([]app_pb.OIDCResponseType, len(responseTypes))
	for i, responseType := range responseTypes {
		switch responseType {
		case domain.OIDCResponseTypeCode:
			oidcResponseTypes[i] = app_pb.OIDCResponseType_OIDC_RESPONSE_TYPE_CODE
		case domain.OIDCResponseTypeIDToken:
			oidcResponseTypes[i] = app_pb.OIDCResponseType_OIDC_RESPONSE_TYPE_ID_TOKEN
		case domain.OIDCResponseTypeIDTokenToken:
			oidcResponseTypes[i] = app_pb.OIDCResponseType_OIDC_RESPONSE_TYPE_ID_TOKEN_TOKEN
		}
	}
	return oidcResponseTypes
}

func OIDCResponseTypesToDomain(responseTypes []app_pb.OIDCResponseType) []domain.OIDCResponseType {
	if len(responseTypes) == 0 {
		return []domain.OIDCResponseType{domain.OIDCResponseTypeCode}
	}
	oidcResponseTypes := make([]domain.OIDCResponseType, len(responseTypes))
	for i, responseType := range responseTypes {
		switch responseType {
		case app_pb.OIDCResponseType_OIDC_RESPONSE_TYPE_CODE:
			oidcResponseTypes[i] = domain.OIDCResponseTypeCode
		case app_pb.OIDCResponseType_OIDC_RESPONSE_TYPE_ID_TOKEN:
			oidcResponseTypes[i] = domain.OIDCResponseTypeIDToken
		case app_pb.OIDCResponseType_OIDC_RESPONSE_TYPE_ID_TOKEN_TOKEN:
			oidcResponseTypes[i] = domain.OIDCResponseTypeIDTokenToken
		}
	}
	return oidcResponseTypes
}

func OIDCGrantTypesFromModel(grantTypes []domain.OIDCGrantType) []app_pb.OIDCGrantType {
	oidcGrantTypes := make([]app_pb.OIDCGrantType, len(grantTypes))
	for i, grantType := range grantTypes {
		switch grantType {
		case domain.OIDCGrantTypeAuthorizationCode:
			oidcGrantTypes[i] = app_pb.OIDCGrantType_OIDC_GRANT_TYPE_AUTHORIZATION_CODE
		case domain.OIDCGrantTypeImplicit:
			oidcGrantTypes[i] = app_pb.OIDCGrantType_OIDC_GRANT_TYPE_IMPLICIT
		case domain.OIDCGrantTypeRefreshToken:
			oidcGrantTypes[i] = app_pb.OIDCGrantType_OIDC_GRANT_TYPE_REFRESH_TOKEN
		case domain.OIDCGrantTypeDeviceCode:
			oidcGrantTypes[i] = app_pb.OIDCGrantType_OIDC_GRANT_TYPE_DEVICE_CODE
		case domain.OIDCGrantTypeTokenExchange:
			oidcGrantTypes[i] = app_pb.OIDCGrantType_OIDC_GRANT_TYPE_TOKEN_EXCHANGE
		}
	}
	return oidcGrantTypes
}

func OIDCGrantTypesToDomain(grantTypes []app_pb.OIDCGrantType) []domain.OIDCGrantType {
	if len(grantTypes) == 0 {
		return []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode}
	}
	oidcGrantTypes := make([]domain.OIDCGrantType, len(grantTypes))
	for i, grantType := range grantTypes {
		switch grantType {
		case app_pb.OIDCGrantType_OIDC_GRANT_TYPE_AUTHORIZATION_CODE:
			oidcGrantTypes[i] = domain.OIDCGrantTypeAuthorizationCode
		case app_pb.OIDCGrantType_OIDC_GRANT_TYPE_IMPLICIT:
			oidcGrantTypes[i] = domain.OIDCGrantTypeImplicit
		case app_pb.OIDCGrantType_OIDC_GRANT_TYPE_REFRESH_TOKEN:
			oidcGrantTypes[i] = domain.OIDCGrantTypeRefreshToken
		case app_pb.OIDCGrantType_OIDC_GRANT_TYPE_DEVICE_CODE:
			oidcGrantTypes[i] = domain.OIDCGrantTypeDeviceCode
		case app_pb.OIDCGrantType_OIDC_GRANT_TYPE_TOKEN_EXCHANGE:
			oidcGrantTypes[i] = domain.OIDCGrantTypeTokenExchange
		}
	}
	return oidcGrantTypes
}

func OIDCApplicationTypeToPb(appType domain.OIDCApplicationType) app_pb.OIDCAppType {
	switch appType {
	case domain.OIDCApplicationTypeWeb:
		return app_pb.OIDCAppType_OIDC_APP_TYPE_WEB
	case domain.OIDCApplicationTypeUserAgent:
		return app_pb.OIDCAppType_OIDC_APP_TYPE_USER_AGENT
	case domain.OIDCApplicationTypeNative:
		return app_pb.OIDCAppType_OIDC_APP_TYPE_NATIVE
	default:
		return app_pb.OIDCAppType_OIDC_APP_TYPE_WEB
	}
}

func OIDCApplicationTypeToDomain(appType app_pb.OIDCAppType) domain.OIDCApplicationType {
	switch appType {
	case app_pb.OIDCAppType_OIDC_APP_TYPE_WEB:
		return domain.OIDCApplicationTypeWeb
	case app_pb.OIDCAppType_OIDC_APP_TYPE_USER_AGENT:
		return domain.OIDCApplicationTypeUserAgent
	case app_pb.OIDCAppType_OIDC_APP_TYPE_NATIVE:
		return domain.OIDCApplicationTypeNative
	}
	return domain.OIDCApplicationTypeWeb
}

func OIDCAuthMethodTypeToPb(authType domain.OIDCAuthMethodType) app_pb.OIDCAuthMethodType {
	switch authType {
	case domain.OIDCAuthMethodTypeBasic:
		return app_pb.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_BASIC
	case domain.OIDCAuthMethodTypePost:
		return app_pb.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_POST
	case domain.OIDCAuthMethodTypeNone:
		return app_pb.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_NONE
	case domain.OIDCAuthMethodTypePrivateKeyJWT:
		return app_pb.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT
	default:
		return app_pb.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_BASIC
	}
}

func OIDCAuthMethodTypeToDomain(authType app_pb.OIDCAuthMethodType) domain.OIDCAuthMethodType {
	switch authType {
	case app_pb.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_BASIC:
		return domain.OIDCAuthMethodTypeBasic
	case app_pb.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_POST:
		return domain.OIDCAuthMethodTypePost
	case app_pb.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_NONE:
		return domain.OIDCAuthMethodTypeNone
	case app_pb.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT:
		return domain.OIDCAuthMethodTypePrivateKeyJWT
	default:
		return domain.OIDCAuthMethodTypeBasic
	}
}

func OIDCVersionToPb(version domain.OIDCVersion) app_pb.OIDCVersion {
	switch version {
	case domain.OIDCVersionV1:
		return app_pb.OIDCVersion_OIDC_VERSION_1_0
	}
	return app_pb.OIDCVersion_OIDC_VERSION_1_0
}

func OIDCVersionToDomain(version app_pb.OIDCVersion) domain.OIDCVersion {
	switch version {
	case app_pb.OIDCVersion_OIDC_VERSION_1_0:
		return domain.OIDCVersionV1
	}
	return domain.OIDCVersionV1
}

func oidcTokenTypeToPb(tokenType domain.OIDCTokenType) app_pb.OIDCTokenType {
	switch tokenType {
	case domain.OIDCTokenTypeBearer:
		return app_pb.OIDCTokenType_OIDC_TOKEN_TYPE_BEARER
	case domain.OIDCTokenTypeJWT:
		return app_pb.OIDCTokenType_OIDC_TOKEN_TYPE_JWT
	default:
		return app_pb.OIDCTokenType_OIDC_TOKEN_TYPE_BEARER
	}
}

func OIDCTokenTypeToDomain(tokenType app_pb.OIDCTokenType) domain.OIDCTokenType {
	switch tokenType {
	case app_pb.OIDCTokenType_OIDC_TOKEN_TYPE_BEARER:
		return domain.OIDCTokenTypeBearer
	case app_pb.OIDCTokenType_OIDC_TOKEN_TYPE_JWT:
		return domain.OIDCTokenTypeJWT
	default:
		return domain.OIDCTokenTypeBearer
	}
}

func ComplianceProblemsToLocalizedMessages(problems []string) []*message_pb.LocalizedMessage {
	converted := make([]*message_pb.LocalizedMessage, len(problems))
	for i, p := range problems {
		converted[i] = message_pb.NewLocalizedMessage(p)
	}
	return converted

}

func APIAuthMethodeTypeToPb(methodType domain.APIAuthMethodType) app_pb.APIAuthMethodType {
	switch methodType {
	case domain.APIAuthMethodTypeBasic:
		return app_pb.APIAuthMethodType_API_AUTH_METHOD_TYPE_BASIC
	case domain.APIAuthMethodTypePrivateKeyJWT:
		return app_pb.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT
	default:
		return app_pb.APIAuthMethodType_API_AUTH_METHOD_TYPE_BASIC
	}
}

func APIAuthMethodTypeToDomain(authType app_pb.APIAuthMethodType) domain.APIAuthMethodType {
	switch authType {
	case app_pb.APIAuthMethodType_API_AUTH_METHOD_TYPE_BASIC:
		return domain.APIAuthMethodTypeBasic
	case app_pb.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT:
		return domain.APIAuthMethodTypePrivateKeyJWT
	default:
		return domain.APIAuthMethodTypeBasic
	}
}

func AppQueriesToModel(queries []*app_pb.AppQuery) (q []query.SearchQuery, err error) {
	q = make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = AppQueryToModel(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func AppQueryToModel(appQuery *app_pb.AppQuery) (query.SearchQuery, error) {
	switch q := appQuery.Query.(type) {
	case *app_pb.AppQuery_NameQuery:
		return query.NewAppNameSearchQuery(object_grpc.TextMethodToQuery(q.NameQuery.Method), q.NameQuery.Name)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "APP-Add46", "List.Query.Invalid")
	}
}
