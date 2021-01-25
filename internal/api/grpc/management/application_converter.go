package management

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/v2/domain"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/caos/logging"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/model"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/pkg/grpc/management"
	"github.com/caos/zitadel/pkg/grpc/message"
)

func appFromDomain(app *domain.Application) *management.Application {
	return &management.Application{
		Id:         app.AppID,
		State:      appStateFromDomain(app.State),
		ChangeDate: timestamppb.New(app.ChangeDate),
		Name:       app.Name,
		Sequence:   app.Sequence,
	}
}

func oidcAppFromDomain(app *domain.OIDCApp) *management.Application {
	return &management.Application{
		Id:         app.AppID,
		State:      appStateFromDomain(app.State),
		ChangeDate: timestamppb.New(app.ChangeDate),
		Name:       app.AppName,
		Sequence:   app.Sequence,
		AppConfig:  oidcAppConfigFromDomain(app),
	}
}

func oidcAppConfigFromDomain(app *domain.OIDCApp) management.AppConfig {
	return &management.Application_OidcConfig{
		OidcConfig: oidcConfigFromDomain(app),
	}
}

func oidcConfigFromDomain(config *domain.OIDCApp) *management.OIDCConfig {
	return &management.OIDCConfig{
		RedirectUris:             config.RedirectUris,
		ResponseTypes:            oidcResponseTypesFromDomain(config.ResponseTypes),
		GrantTypes:               oidcGrantTypesFromDomain(config.GrantTypes),
		ApplicationType:          oidcApplicationTypeFromDomain(config.ApplicationType),
		ClientId:                 config.ClientID,
		ClientSecret:             config.ClientSecretString,
		AuthMethodType:           oidcAuthMethodTypeFromDomain(config.AuthMethodType),
		PostLogoutRedirectUris:   config.PostLogoutRedirectUris,
		Version:                  oidcVersionFromDomain(config.OIDCVersion),
		NoneCompliant:            config.Compliance.NoneCompliant,
		ComplianceProblems:       complianceProblemsToLocalizedMessages(config.Compliance.Problems),
		DevMode:                  config.DevMode,
		AccessTokenType:          oidcTokenTypeFromDomain(config.AccessTokenType),
		AccessTokenRoleAssertion: config.AccessTokenRoleAssertion,
		IdTokenRoleAssertion:     config.IDTokenRoleAssertion,
		IdTokenUserinfoAssertion: config.IDTokenUserinfoAssertion,
		ClockSkew:                durationpb.New(config.ClockSkew),
	}
}

func oidcConfigFromModel(config *proj_model.OIDCConfig) *management.OIDCConfig {
	return &management.OIDCConfig{
		RedirectUris:             config.RedirectUris,
		ResponseTypes:            oidcResponseTypesFromModel(config.ResponseTypes),
		GrantTypes:               oidcGrantTypesFromModel(config.GrantTypes),
		ApplicationType:          oidcApplicationTypeFromModel(config.ApplicationType),
		ClientId:                 config.ClientID,
		ClientSecret:             config.ClientSecretString,
		AuthMethodType:           oidcAuthMethodTypeFromModel(config.AuthMethodType),
		PostLogoutRedirectUris:   config.PostLogoutRedirectUris,
		Version:                  oidcVersionFromModel(config.OIDCVersion),
		NoneCompliant:            config.Compliance.NoneCompliant,
		ComplianceProblems:       complianceProblemsToLocalizedMessages(config.Compliance.Problems),
		DevMode:                  config.DevMode,
		AccessTokenType:          oidcTokenTypeFromModel(config.AccessTokenType),
		AccessTokenRoleAssertion: config.AccessTokenRoleAssertion,
		IdTokenRoleAssertion:     config.IDTokenRoleAssertion,
		IdTokenUserinfoAssertion: config.IDTokenUserinfoAssertion,
		ClockSkew:                durationpb.New(config.ClockSkew),
	}
}

func oidcConfigFromApplicationViewModel(app *proj_model.ApplicationView) *management.OIDCConfig {
	return &management.OIDCConfig{
		RedirectUris:             app.OIDCRedirectUris,
		ResponseTypes:            oidcResponseTypesFromModel(app.OIDCResponseTypes),
		GrantTypes:               oidcGrantTypesFromModel(app.OIDCGrantTypes),
		ApplicationType:          oidcApplicationTypeFromModel(app.OIDCApplicationType),
		ClientId:                 app.OIDCClientID,
		AuthMethodType:           oidcAuthMethodTypeFromModel(app.OIDCAuthMethodType),
		PostLogoutRedirectUris:   app.OIDCPostLogoutRedirectUris,
		Version:                  oidcVersionFromModel(app.OIDCVersion),
		NoneCompliant:            app.NoneCompliant,
		ComplianceProblems:       complianceProblemsToLocalizedMessages(app.ComplianceProblems),
		DevMode:                  app.DevMode,
		AccessTokenType:          oidcTokenTypeFromModel(app.AccessTokenType),
		AccessTokenRoleAssertion: app.AccessTokenRoleAssertion,
		IdTokenRoleAssertion:     app.IDTokenRoleAssertion,
		IdTokenUserinfoAssertion: app.IDTokenUserinfoAssertion,
		ClockSkew:                durationpb.New(app.ClockSkew),
	}
}

func complianceProblemsToLocalizedMessages(problems []string) []*message.LocalizedMessage {
	converted := make([]*message.LocalizedMessage, len(problems))
	for i, p := range problems {
		converted[i] = message.NewLocalizedMessage(p)
	}
	return converted

}

func oidcAppCreateToDomain(app *management.OIDCApplicationCreate) *domain.OIDCApp {
	return &domain.OIDCApp{
		ObjectRoot: models.ObjectRoot{
			AggregateID: app.ProjectId,
		},
		AppName:                  app.Name,
		OIDCVersion:              oidcVersionToDomain(app.Version),
		RedirectUris:             app.RedirectUris,
		ResponseTypes:            oidcResponseTypesToDomain(app.ResponseTypes),
		GrantTypes:               oidcGrantTypesToDomain(app.GrantTypes),
		ApplicationType:          oidcApplicationTypeToDomain(app.ApplicationType),
		AuthMethodType:           oidcAuthMethodTypeToDomain(app.AuthMethodType),
		PostLogoutRedirectUris:   app.PostLogoutRedirectUris,
		DevMode:                  app.DevMode,
		AccessTokenType:          oidcTokenTypeToDomain(app.AccessTokenType),
		AccessTokenRoleAssertion: app.AccessTokenRoleAssertion,
		IDTokenRoleAssertion:     app.IdTokenRoleAssertion,
		IDTokenUserinfoAssertion: app.IdTokenUserinfoAssertion,
		ClockSkew:                app.ClockSkew.AsDuration(),
	}
}

func appUpdateToDomain(app *management.ApplicationUpdate) *domain.Application {
	return &domain.Application{
		ObjectRoot: models.ObjectRoot{
			AggregateID: app.ProjectId,
		},
		AppID: app.Id,
		Name:  app.Name,
	}
}

func oidcConfigUpdateToModel(app *management.OIDCConfigUpdate) *proj_model.OIDCConfig {
	return &proj_model.OIDCConfig{
		ObjectRoot: models.ObjectRoot{
			AggregateID: app.ProjectId,
		},
		AppID:                    app.ApplicationId,
		RedirectUris:             app.RedirectUris,
		ResponseTypes:            oidcResponseTypesToModel(app.ResponseTypes),
		GrantTypes:               oidcGrantTypesToModel(app.GrantTypes),
		ApplicationType:          oidcApplicationTypeToModel(app.ApplicationType),
		AuthMethodType:           oidcAuthMethodTypeToModel(app.AuthMethodType),
		PostLogoutRedirectUris:   app.PostLogoutRedirectUris,
		DevMode:                  app.DevMode,
		AccessTokenType:          oidcTokenTypeToModel(app.AccessTokenType),
		AccessTokenRoleAssertion: app.AccessTokenRoleAssertion,
		IDTokenRoleAssertion:     app.IdTokenRoleAssertion,
		IDTokenUserinfoAssertion: app.IdTokenUserinfoAssertion,
		ClockSkew:                app.ClockSkew.AsDuration(),
	}
}

func applicationSearchRequestsToModel(request *management.ApplicationSearchRequest) *proj_model.ApplicationSearchRequest {
	return &proj_model.ApplicationSearchRequest{
		Offset:  request.Offset,
		Limit:   request.Limit,
		Queries: applicationSearchQueriesToModel(request.ProjectId, request.Queries),
	}
}

func applicationSearchQueriesToModel(projectID string, queries []*management.ApplicationSearchQuery) []*proj_model.ApplicationSearchQuery {
	converted := make([]*proj_model.ApplicationSearchQuery, len(queries)+1)
	for i, q := range queries {
		converted[i] = applicationSearchQueryToModel(q)
	}
	converted[len(queries)] = &proj_model.ApplicationSearchQuery{Key: proj_model.AppSearchKeyProjectID, Method: model.SearchMethodEquals, Value: projectID}

	return converted
}

func applicationSearchQueryToModel(query *management.ApplicationSearchQuery) *proj_model.ApplicationSearchQuery {
	return &proj_model.ApplicationSearchQuery{
		Key:    applicationSearchKeyToModel(query.Key),
		Method: searchMethodToModel(query.Method),
		Value:  query.Value,
	}
}

func applicationSearchKeyToModel(key management.ApplicationSearchKey) proj_model.AppSearchKey {
	switch key {
	case management.ApplicationSearchKey_APPLICATIONSEARCHKEY_APP_NAME:
		return proj_model.AppSearchKeyName
	default:
		return proj_model.AppSearchKeyUnspecified
	}
}

func applicationSearchResponseFromModel(response *proj_model.ApplicationSearchResponse) *management.ApplicationSearchResponse {
	timestamp, err := ptypes.TimestampProto(response.Timestamp)
	logging.Log("GRPC-Lp06f").OnError(err).Debug("unable to parse timestamp")
	return &management.ApplicationSearchResponse{
		Offset:            response.Offset,
		Limit:             response.Limit,
		TotalResult:       response.TotalResult,
		Result:            applicationViewsFromModel(response.Result),
		ProcessedSequence: response.Sequence,
		ViewTimestamp:     timestamp,
	}
}

func applicationViewsFromModel(apps []*proj_model.ApplicationView) []*management.ApplicationView {
	converted := make([]*management.ApplicationView, len(apps))
	for i, app := range apps {
		converted[i] = applicationViewFromModel(app)
	}
	return converted
}

func applicationViewFromModel(application *proj_model.ApplicationView) *management.ApplicationView {
	creationDate, err := ptypes.TimestampProto(application.CreationDate)
	logging.Log("GRPC-lo9sw").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(application.ChangeDate)
	logging.Log("GRPC-8uwsd").OnError(err).Debug("unable to parse timestamp")

	converted := &management.ApplicationView{
		Id:           application.ID,
		State:        appStateFromModel(application.State),
		CreationDate: creationDate,
		ChangeDate:   changeDate,
		Name:         application.Name,
		Sequence:     application.Sequence,
	}
	if application.IsOIDC {
		converted.AppConfig = &management.ApplicationView_OidcConfig{
			OidcConfig: oidcConfigFromApplicationViewModel(application),
		}
	}
	return converted
}

func appStateFromDomain(state domain.AppState) management.AppState {
	switch state {
	case domain.AppStateActive:
		return management.AppState_APPSTATE_ACTIVE
	case domain.AppStateInactive:
		return management.AppState_APPSTATE_INACTIVE
	default:
		return management.AppState_APPSTATE_UNSPECIFIED
	}
}

func appStateFromModel(state proj_model.AppState) management.AppState {
	switch state {
	case proj_model.AppStateActive:
		return management.AppState_APPSTATE_ACTIVE
	case proj_model.AppStateInactive:
		return management.AppState_APPSTATE_INACTIVE
	default:
		return management.AppState_APPSTATE_UNSPECIFIED
	}
}

func oidcResponseTypesToDomain(responseTypes []management.OIDCResponseType) []domain.OIDCResponseType {
	if responseTypes == nil || len(responseTypes) == 0 {
		return []domain.OIDCResponseType{domain.OIDCResponseTypeCode}
	}
	oidcResponseTypes := make([]domain.OIDCResponseType, len(responseTypes))

	for i, responseType := range responseTypes {
		switch responseType {
		case management.OIDCResponseType_OIDCRESPONSETYPE_CODE:
			oidcResponseTypes[i] = domain.OIDCResponseTypeCode
		case management.OIDCResponseType_OIDCRESPONSETYPE_ID_TOKEN:
			oidcResponseTypes[i] = domain.OIDCResponseTypeIDToken
		case management.OIDCResponseType_OIDCRESPONSETYPE_ID_TOKEN_TOKEN:
			oidcResponseTypes[i] = domain.OIDCResponseTypeIDTokenToken
		}
	}

	return oidcResponseTypes
}

func oidcResponseTypesToModel(responseTypes []management.OIDCResponseType) []proj_model.OIDCResponseType {
	if responseTypes == nil || len(responseTypes) == 0 {
		return []proj_model.OIDCResponseType{proj_model.OIDCResponseTypeCode}
	}
	oidcResponseTypes := make([]proj_model.OIDCResponseType, len(responseTypes))

	for i, responseType := range responseTypes {
		switch responseType {
		case management.OIDCResponseType_OIDCRESPONSETYPE_CODE:
			oidcResponseTypes[i] = proj_model.OIDCResponseTypeCode
		case management.OIDCResponseType_OIDCRESPONSETYPE_ID_TOKEN:
			oidcResponseTypes[i] = proj_model.OIDCResponseTypeIDToken
		case management.OIDCResponseType_OIDCRESPONSETYPE_ID_TOKEN_TOKEN:
			oidcResponseTypes[i] = proj_model.OIDCResponseTypeIDTokenToken
		}
	}

	return oidcResponseTypes
}

func oidcResponseTypesFromDomain(responseTypes []domain.OIDCResponseType) []management.OIDCResponseType {
	oidcResponseTypes := make([]management.OIDCResponseType, len(responseTypes))

	for i, responseType := range responseTypes {
		switch responseType {
		case domain.OIDCResponseTypeCode:
			oidcResponseTypes[i] = management.OIDCResponseType_OIDCRESPONSETYPE_CODE
		case domain.OIDCResponseTypeIDToken:
			oidcResponseTypes[i] = management.OIDCResponseType_OIDCRESPONSETYPE_ID_TOKEN
		case domain.OIDCResponseTypeIDTokenToken:
			oidcResponseTypes[i] = management.OIDCResponseType_OIDCRESPONSETYPE_ID_TOKEN_TOKEN
		}
	}

	return oidcResponseTypes
}
func oidcResponseTypesFromModel(responseTypes []proj_model.OIDCResponseType) []management.OIDCResponseType {
	oidcResponseTypes := make([]management.OIDCResponseType, len(responseTypes))

	for i, responseType := range responseTypes {
		switch responseType {
		case proj_model.OIDCResponseTypeCode:
			oidcResponseTypes[i] = management.OIDCResponseType_OIDCRESPONSETYPE_CODE
		case proj_model.OIDCResponseTypeIDToken:
			oidcResponseTypes[i] = management.OIDCResponseType_OIDCRESPONSETYPE_ID_TOKEN
		case proj_model.OIDCResponseTypeIDTokenToken:
			oidcResponseTypes[i] = management.OIDCResponseType_OIDCRESPONSETYPE_ID_TOKEN_TOKEN
		}
	}

	return oidcResponseTypes
}

func oidcGrantTypesToDomain(grantTypes []management.OIDCGrantType) []domain.OIDCGrantType {
	if grantTypes == nil || len(grantTypes) == 0 {
		return []domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode}
	}
	oidcGrantTypes := make([]domain.OIDCGrantType, len(grantTypes))

	for i, grantType := range grantTypes {
		switch grantType {
		case management.OIDCGrantType_OIDCGRANTTYPE_AUTHORIZATION_CODE:
			oidcGrantTypes[i] = domain.OIDCGrantTypeAuthorizationCode
		case management.OIDCGrantType_OIDCGRANTTYPE_IMPLICIT:
			oidcGrantTypes[i] = domain.OIDCGrantTypeImplicit
		case management.OIDCGrantType_OIDCGRANTTYPE_REFRESH_TOKEN:
			oidcGrantTypes[i] = domain.OIDCGrantTypeRefreshToken
		}
	}
	return oidcGrantTypes
}

func oidcGrantTypesToModel(grantTypes []management.OIDCGrantType) []proj_model.OIDCGrantType {
	if grantTypes == nil || len(grantTypes) == 0 {
		return []proj_model.OIDCGrantType{proj_model.OIDCGrantTypeAuthorizationCode}
	}
	oidcGrantTypes := make([]proj_model.OIDCGrantType, len(grantTypes))

	for i, grantType := range grantTypes {
		switch grantType {
		case management.OIDCGrantType_OIDCGRANTTYPE_AUTHORIZATION_CODE:
			oidcGrantTypes[i] = proj_model.OIDCGrantTypeAuthorizationCode
		case management.OIDCGrantType_OIDCGRANTTYPE_IMPLICIT:
			oidcGrantTypes[i] = proj_model.OIDCGrantTypeImplicit
		case management.OIDCGrantType_OIDCGRANTTYPE_REFRESH_TOKEN:
			oidcGrantTypes[i] = proj_model.OIDCGrantTypeRefreshToken
		}
	}
	return oidcGrantTypes
}

func oidcGrantTypesFromDomain(grantTypes []domain.OIDCGrantType) []management.OIDCGrantType {
	oidcGrantTypes := make([]management.OIDCGrantType, len(grantTypes))

	for i, grantType := range grantTypes {
		switch grantType {
		case domain.OIDCGrantTypeAuthorizationCode:
			oidcGrantTypes[i] = management.OIDCGrantType_OIDCGRANTTYPE_AUTHORIZATION_CODE
		case domain.OIDCGrantTypeImplicit:
			oidcGrantTypes[i] = management.OIDCGrantType_OIDCGRANTTYPE_IMPLICIT
		case domain.OIDCGrantTypeRefreshToken:
			oidcGrantTypes[i] = management.OIDCGrantType_OIDCGRANTTYPE_REFRESH_TOKEN
		}
	}
	return oidcGrantTypes
}

func oidcGrantTypesFromModel(grantTypes []proj_model.OIDCGrantType) []management.OIDCGrantType {
	oidcGrantTypes := make([]management.OIDCGrantType, len(grantTypes))

	for i, grantType := range grantTypes {
		switch grantType {
		case proj_model.OIDCGrantTypeAuthorizationCode:
			oidcGrantTypes[i] = management.OIDCGrantType_OIDCGRANTTYPE_AUTHORIZATION_CODE
		case proj_model.OIDCGrantTypeImplicit:
			oidcGrantTypes[i] = management.OIDCGrantType_OIDCGRANTTYPE_IMPLICIT
		case proj_model.OIDCGrantTypeRefreshToken:
			oidcGrantTypes[i] = management.OIDCGrantType_OIDCGRANTTYPE_REFRESH_TOKEN
		}
	}
	return oidcGrantTypes
}

func oidcApplicationTypeToDomain(appType management.OIDCApplicationType) domain.OIDCApplicationType {
	switch appType {
	case management.OIDCApplicationType_OIDCAPPLICATIONTYPE_WEB:
		return domain.OIDCApplicationTypeWeb
	case management.OIDCApplicationType_OIDCAPPLICATIONTYPE_USER_AGENT:
		return domain.OIDCApplicationTypeUserAgent
	case management.OIDCApplicationType_OIDCAPPLICATIONTYPE_NATIVE:
		return domain.OIDCApplicationTypeNative
	}
	return domain.OIDCApplicationTypeWeb
}
func oidcApplicationTypeToModel(appType management.OIDCApplicationType) proj_model.OIDCApplicationType {
	switch appType {
	case management.OIDCApplicationType_OIDCAPPLICATIONTYPE_WEB:
		return proj_model.OIDCApplicationTypeWeb
	case management.OIDCApplicationType_OIDCAPPLICATIONTYPE_USER_AGENT:
		return proj_model.OIDCApplicationTypeUserAgent
	case management.OIDCApplicationType_OIDCAPPLICATIONTYPE_NATIVE:
		return proj_model.OIDCApplicationTypeNative
	}
	return proj_model.OIDCApplicationTypeWeb
}

func oidcVersionToDomain(version management.OIDCVersion) domain.OIDCVersion {
	switch version {
	case management.OIDCVersion_OIDCV1_0:
		return domain.OIDCVersionV1
	}
	return domain.OIDCVersionV1
}

func oidcApplicationTypeFromDomain(appType domain.OIDCApplicationType) management.OIDCApplicationType {
	switch appType {
	case domain.OIDCApplicationTypeWeb:
		return management.OIDCApplicationType_OIDCAPPLICATIONTYPE_WEB
	case domain.OIDCApplicationTypeUserAgent:
		return management.OIDCApplicationType_OIDCAPPLICATIONTYPE_USER_AGENT
	case domain.OIDCApplicationTypeNative:
		return management.OIDCApplicationType_OIDCAPPLICATIONTYPE_NATIVE
	default:
		return management.OIDCApplicationType_OIDCAPPLICATIONTYPE_WEB
	}
}

func oidcApplicationTypeFromModel(appType proj_model.OIDCApplicationType) management.OIDCApplicationType {
	switch appType {
	case proj_model.OIDCApplicationTypeWeb:
		return management.OIDCApplicationType_OIDCAPPLICATIONTYPE_WEB
	case proj_model.OIDCApplicationTypeUserAgent:
		return management.OIDCApplicationType_OIDCAPPLICATIONTYPE_USER_AGENT
	case proj_model.OIDCApplicationTypeNative:
		return management.OIDCApplicationType_OIDCAPPLICATIONTYPE_NATIVE
	default:
		return management.OIDCApplicationType_OIDCAPPLICATIONTYPE_WEB
	}
}

func oidcAuthMethodTypeToDomain(authType management.OIDCAuthMethodType) domain.OIDCAuthMethodType {
	switch authType {
	case management.OIDCAuthMethodType_OIDCAUTHMETHODTYPE_BASIC:
		return domain.OIDCAuthMethodTypeBasic
	case management.OIDCAuthMethodType_OIDCAUTHMETHODTYPE_POST:
		return domain.OIDCAuthMethodTypePost
	case management.OIDCAuthMethodType_OIDCAUTHMETHODTYPE_NONE:
		return domain.OIDCAuthMethodTypeNone
	default:
		return domain.OIDCAuthMethodTypeBasic
	}
}

func oidcAuthMethodTypeToModel(authType management.OIDCAuthMethodType) proj_model.OIDCAuthMethodType {
	switch authType {
	case management.OIDCAuthMethodType_OIDCAUTHMETHODTYPE_BASIC:
		return proj_model.OIDCAuthMethodTypeBasic
	case management.OIDCAuthMethodType_OIDCAUTHMETHODTYPE_POST:
		return proj_model.OIDCAuthMethodTypePost
	case management.OIDCAuthMethodType_OIDCAUTHMETHODTYPE_NONE:
		return proj_model.OIDCAuthMethodTypeNone
	default:
		return proj_model.OIDCAuthMethodTypeBasic
	}
}

func oidcAuthMethodTypeFromDomain(authType domain.OIDCAuthMethodType) management.OIDCAuthMethodType {
	switch authType {
	case domain.OIDCAuthMethodTypeBasic:
		return management.OIDCAuthMethodType_OIDCAUTHMETHODTYPE_BASIC
	case domain.OIDCAuthMethodTypePost:
		return management.OIDCAuthMethodType_OIDCAUTHMETHODTYPE_POST
	case domain.OIDCAuthMethodTypeNone:
		return management.OIDCAuthMethodType_OIDCAUTHMETHODTYPE_NONE
	default:
		return management.OIDCAuthMethodType_OIDCAUTHMETHODTYPE_BASIC
	}
}

func oidcAuthMethodTypeFromModel(authType proj_model.OIDCAuthMethodType) management.OIDCAuthMethodType {
	switch authType {
	case proj_model.OIDCAuthMethodTypeBasic:
		return management.OIDCAuthMethodType_OIDCAUTHMETHODTYPE_BASIC
	case proj_model.OIDCAuthMethodTypePost:
		return management.OIDCAuthMethodType_OIDCAUTHMETHODTYPE_POST
	case proj_model.OIDCAuthMethodTypeNone:
		return management.OIDCAuthMethodType_OIDCAUTHMETHODTYPE_NONE
	default:
		return management.OIDCAuthMethodType_OIDCAUTHMETHODTYPE_BASIC
	}
}

func oidcTokenTypeToDomain(tokenType management.OIDCTokenType) domain.OIDCTokenType {
	switch tokenType {
	case management.OIDCTokenType_OIDCTokenType_Bearer:
		return domain.OIDCTokenTypeBearer
	case management.OIDCTokenType_OIDCTokenType_JWT:
		return domain.OIDCTokenTypeJWT
	default:
		return domain.OIDCTokenTypeBearer
	}
}

func oidcTokenTypeToModel(tokenType management.OIDCTokenType) proj_model.OIDCTokenType {
	switch tokenType {
	case management.OIDCTokenType_OIDCTokenType_Bearer:
		return proj_model.OIDCTokenTypeBearer
	case management.OIDCTokenType_OIDCTokenType_JWT:
		return proj_model.OIDCTokenTypeJWT
	default:
		return proj_model.OIDCTokenTypeBearer
	}
}

func oidcTokenTypeFromDomain(tokenType domain.OIDCTokenType) management.OIDCTokenType {
	switch tokenType {
	case domain.OIDCTokenTypeBearer:
		return management.OIDCTokenType_OIDCTokenType_Bearer
	case domain.OIDCTokenTypeJWT:
		return management.OIDCTokenType_OIDCTokenType_JWT
	default:
		return management.OIDCTokenType_OIDCTokenType_Bearer
	}
}

func oidcTokenTypeFromModel(tokenType proj_model.OIDCTokenType) management.OIDCTokenType {
	switch tokenType {
	case proj_model.OIDCTokenTypeBearer:
		return management.OIDCTokenType_OIDCTokenType_Bearer
	case proj_model.OIDCTokenTypeJWT:
		return management.OIDCTokenType_OIDCTokenType_JWT
	default:
		return management.OIDCTokenType_OIDCTokenType_Bearer
	}
}

func oidcVersionFromDomain(version domain.OIDCVersion) management.OIDCVersion {
	switch version {
	case domain.OIDCVersionV1:
		return management.OIDCVersion_OIDCV1_0
	default:
		return management.OIDCVersion_OIDCV1_0
	}
}

func oidcVersionFromModel(version proj_model.OIDCVersion) management.OIDCVersion {
	switch version {
	case proj_model.OIDCVersionV1:
		return management.OIDCVersion_OIDCV1_0
	default:
		return management.OIDCVersion_OIDCV1_0
	}
}

func appChangesToResponse(response *proj_model.ApplicationChanges, offset uint64, limit uint64) (_ *management.Changes) {
	return &management.Changes{
		Limit:   limit,
		Offset:  offset,
		Changes: appChangesToMgtAPI(response),
	}
}

func appChangesToMgtAPI(changes *proj_model.ApplicationChanges) (_ []*management.Change) {
	result := make([]*management.Change, len(changes.Changes))

	for i, change := range changes.Changes {
		b, err := json.Marshal(change.Data)
		data := &structpb.Struct{}
		err = protojson.Unmarshal(b, data)
		if err != nil {
		}
		result[i] = &management.Change{
			ChangeDate: change.ChangeDate,
			EventType:  message.NewLocalizedEventType(change.EventType),
			Sequence:   change.Sequence,
			Editor:     change.ModifierName,
			EditorId:   change.ModifierId,
			Data:       data,
		}
	}

	return result
}
