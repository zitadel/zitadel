package management

import (
	"encoding/json"

	"github.com/caos/logging"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/model"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/pkg/management/grpc"
)

func appFromModel(app *proj_model.Application) *grpc.Application {
	creationDate, err := ptypes.TimestampProto(app.CreationDate)
	logging.Log("GRPC-iejs3").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(app.ChangeDate)
	logging.Log("GRPC-di7rw").OnError(err).Debug("unable to parse timestamp")

	return &grpc.Application{
		Id:           app.AppID,
		State:        appStateFromModel(app.State),
		CreationDate: creationDate,
		ChangeDate:   changeDate,
		Name:         app.Name,
		Sequence:     app.Sequence,
		AppConfig:    appConfigFromModel(app),
	}
}

func appConfigFromModel(app *proj_model.Application) grpc.AppConfig {
	if app.Type == proj_model.AppTypeOIDC {
		return &grpc.Application_OidcConfig{
			OidcConfig: oidcConfigFromModel(app.OIDCConfig),
		}
	}
	return nil
}

func oidcConfigFromModel(config *proj_model.OIDCConfig) *grpc.OIDCConfig {
	return &grpc.OIDCConfig{
		RedirectUris:           config.RedirectUris,
		ResponseTypes:          oidcResponseTypesFromModel(config.ResponseTypes),
		GrantTypes:             oidcGrantTypesFromModel(config.GrantTypes),
		ApplicationType:        oidcApplicationTypeFromModel(config.ApplicationType),
		ClientId:               config.ClientID,
		ClientSecret:           config.ClientSecretString,
		AuthMethodType:         oidcAuthMethodTypeFromModel(config.AuthMethodType),
		PostLogoutRedirectUris: config.PostLogoutRedirectUris,
	}
}

func oidcConfigFromApplicationViewModel(app *proj_model.ApplicationView) *grpc.OIDCConfig {
	return &grpc.OIDCConfig{
		RedirectUris:           app.OIDCRedirectUris,
		ResponseTypes:          oidcResponseTypesFromModel(app.OIDCResponseTypes),
		GrantTypes:             oidcGrantTypesFromModel(app.OIDCGrantTypes),
		ApplicationType:        oidcApplicationTypeFromModel(app.OIDCApplicationType),
		ClientId:               app.OIDCClientID,
		AuthMethodType:         oidcAuthMethodTypeFromModel(app.OIDCAuthMethodType),
		PostLogoutRedirectUris: app.OIDCPostLogoutRedirectUris,
	}
}

func oidcAppCreateToModel(app *grpc.OIDCApplicationCreate) *proj_model.Application {
	return &proj_model.Application{
		ObjectRoot: models.ObjectRoot{
			AggregateID: app.ProjectId,
		},
		Name: app.Name,
		Type: proj_model.AppTypeOIDC,
		OIDCConfig: &proj_model.OIDCConfig{
			RedirectUris:           app.RedirectUris,
			ResponseTypes:          oidcResponseTypesToModel(app.ResponseTypes),
			GrantTypes:             oidcGrantTypesToModel(app.GrantTypes),
			ApplicationType:        oidcApplicationTypeToModel(app.ApplicationType),
			AuthMethodType:         oidcAuthMethodTypeToModel(app.AuthMethodType),
			PostLogoutRedirectUris: app.PostLogoutRedirectUris,
		},
	}
}

func appUpdateToModel(app *grpc.ApplicationUpdate) *proj_model.Application {
	return &proj_model.Application{
		ObjectRoot: models.ObjectRoot{
			AggregateID: app.ProjectId,
		},
		AppID: app.Id,
		Name:  app.Name,
	}
}

func oidcConfigUpdateToModel(app *grpc.OIDCConfigUpdate) *proj_model.OIDCConfig {
	return &proj_model.OIDCConfig{
		ObjectRoot: models.ObjectRoot{
			AggregateID: app.ProjectId,
		},
		AppID:                  app.ApplicationId,
		RedirectUris:           app.RedirectUris,
		ResponseTypes:          oidcResponseTypesToModel(app.ResponseTypes),
		GrantTypes:             oidcGrantTypesToModel(app.GrantTypes),
		ApplicationType:        oidcApplicationTypeToModel(app.ApplicationType),
		AuthMethodType:         oidcAuthMethodTypeToModel(app.AuthMethodType),
		PostLogoutRedirectUris: app.PostLogoutRedirectUris,
	}
}

func applicationSearchRequestsToModel(request *grpc.ApplicationSearchRequest) *proj_model.ApplicationSearchRequest {
	return &proj_model.ApplicationSearchRequest{
		Offset:  request.Offset,
		Limit:   request.Limit,
		Queries: applicationSearchQueriesToModel(request.ProjectId, request.Queries),
	}
}

func applicationSearchQueriesToModel(projectID string, queries []*grpc.ApplicationSearchQuery) []*proj_model.ApplicationSearchQuery {
	converted := make([]*proj_model.ApplicationSearchQuery, len(queries)+1)
	for i, q := range queries {
		converted[i] = applicationSearchQueryToModel(q)
	}
	converted[len(queries)] = &proj_model.ApplicationSearchQuery{Key: proj_model.AppSearchKeyProjectID, Method: model.SearchMethodEquals, Value: projectID}

	return converted
}

func applicationSearchQueryToModel(query *grpc.ApplicationSearchQuery) *proj_model.ApplicationSearchQuery {
	return &proj_model.ApplicationSearchQuery{
		Key:    applicationSearchKeyToModel(query.Key),
		Method: searchMethodToModel(query.Method),
		Value:  query.Value,
	}
}

func applicationSearchKeyToModel(key grpc.ApplicationSearchKey) proj_model.AppSearchKey {
	switch key {
	case grpc.ApplicationSearchKey_APPLICATIONSEARCHKEY_APP_NAME:
		return proj_model.AppSearchKeyName
	default:
		return proj_model.AppSearchKeyUnspecified
	}
}

func applicationSearchResponseFromModel(response *proj_model.ApplicationSearchResponse) *grpc.ApplicationSearchResponse {
	return &grpc.ApplicationSearchResponse{
		Offset:      response.Offset,
		Limit:       response.Limit,
		TotalResult: response.TotalResult,
		Result:      applicationViewsFromModel(response.Result),
	}
}

func applicationViewsFromModel(apps []*proj_model.ApplicationView) []*grpc.ApplicationView {
	converted := make([]*grpc.ApplicationView, len(apps))
	for i, app := range apps {
		converted[i] = applicationViewFromModel(app)
	}
	return converted
}

func applicationViewFromModel(application *proj_model.ApplicationView) *grpc.ApplicationView {
	creationDate, err := ptypes.TimestampProto(application.CreationDate)
	logging.Log("GRPC-lo9sw").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(application.ChangeDate)
	logging.Log("GRPC-8uwsd").OnError(err).Debug("unable to parse timestamp")

	converted := &grpc.ApplicationView{
		Id:           application.ID,
		State:        appStateFromModel(application.State),
		CreationDate: creationDate,
		ChangeDate:   changeDate,
		Name:         application.Name,
		Sequence:     application.Sequence,
	}
	if application.IsOIDC {
		converted.AppConfig = &grpc.ApplicationView_OidcConfig{
			OidcConfig: oidcConfigFromApplicationViewModel(application),
		}
	}
	return converted
}

func appStateFromModel(state proj_model.AppState) grpc.AppState {
	switch state {
	case proj_model.AppStateActive:
		return grpc.AppState_APPSTATE_ACTIVE
	case proj_model.AppStateInactive:
		return grpc.AppState_APPSTATE_INACTIVE
	default:
		return grpc.AppState_APPSTATE_UNSPECIFIED
	}
}

func oidcResponseTypesToModel(responseTypes []grpc.OIDCResponseType) []proj_model.OIDCResponseType {
	if responseTypes == nil || len(responseTypes) == 0 {
		return []proj_model.OIDCResponseType{proj_model.OIDCResponseTypeCode}
	}
	oidcResponseTypes := make([]proj_model.OIDCResponseType, len(responseTypes))

	for i, responseType := range responseTypes {
		switch responseType {
		case grpc.OIDCResponseType_OIDCRESPONSETYPE_CODE:
			oidcResponseTypes[i] = proj_model.OIDCResponseTypeCode
		case grpc.OIDCResponseType_OIDCRESPONSETYPE_ID_TOKEN:
			oidcResponseTypes[i] = proj_model.OIDCResponseTypeIDToken
		case grpc.OIDCResponseType_OIDCRESPONSETYPE_TOKEN:
			oidcResponseTypes[i] = proj_model.OIDCResponseTypeToken
		}
	}

	return oidcResponseTypes
}

func oidcResponseTypesFromModel(responseTypes []proj_model.OIDCResponseType) []grpc.OIDCResponseType {
	oidcResponseTypes := make([]grpc.OIDCResponseType, len(responseTypes))

	for i, responseType := range responseTypes {
		switch responseType {
		case proj_model.OIDCResponseTypeCode:
			oidcResponseTypes[i] = grpc.OIDCResponseType_OIDCRESPONSETYPE_CODE
		case proj_model.OIDCResponseTypeIDToken:
			oidcResponseTypes[i] = grpc.OIDCResponseType_OIDCRESPONSETYPE_ID_TOKEN
		case proj_model.OIDCResponseTypeToken:
			oidcResponseTypes[i] = grpc.OIDCResponseType_OIDCRESPONSETYPE_TOKEN
		}
	}

	return oidcResponseTypes
}

func oidcGrantTypesToModel(grantTypes []grpc.OIDCGrantType) []proj_model.OIDCGrantType {
	if grantTypes == nil || len(grantTypes) == 0 {
		return []proj_model.OIDCGrantType{proj_model.OIDCGrantTypeAuthorizationCode}
	}
	oidcGrantTypes := make([]proj_model.OIDCGrantType, len(grantTypes))

	for i, grantType := range grantTypes {
		switch grantType {
		case grpc.OIDCGrantType_OIDCGRANTTYPE_AUTHORIZATION_CODE:
			oidcGrantTypes[i] = proj_model.OIDCGrantTypeAuthorizationCode
		case grpc.OIDCGrantType_OIDCGRANTTYPE_IMPLICIT:
			oidcGrantTypes[i] = proj_model.OIDCGrantTypeImplicit
		case grpc.OIDCGrantType_OIDCGRANTTYPE_REFRESH_TOKEN:
			oidcGrantTypes[i] = proj_model.OIDCGrantTypeRefreshToken
		}
	}
	return oidcGrantTypes
}

func oidcGrantTypesFromModel(grantTypes []proj_model.OIDCGrantType) []grpc.OIDCGrantType {
	oidcGrantTypes := make([]grpc.OIDCGrantType, len(grantTypes))

	for i, grantType := range grantTypes {
		switch grantType {
		case proj_model.OIDCGrantTypeAuthorizationCode:
			oidcGrantTypes[i] = grpc.OIDCGrantType_OIDCGRANTTYPE_AUTHORIZATION_CODE
		case proj_model.OIDCGrantTypeImplicit:
			oidcGrantTypes[i] = grpc.OIDCGrantType_OIDCGRANTTYPE_IMPLICIT
		case proj_model.OIDCGrantTypeRefreshToken:
			oidcGrantTypes[i] = grpc.OIDCGrantType_OIDCGRANTTYPE_REFRESH_TOKEN
		}
	}
	return oidcGrantTypes
}

func oidcApplicationTypeToModel(appType grpc.OIDCApplicationType) proj_model.OIDCApplicationType {
	switch appType {
	case grpc.OIDCApplicationType_OIDCAPPLICATIONTYPE_WEB:
		return proj_model.OIDCApplicationTypeWeb
	case grpc.OIDCApplicationType_OIDCAPPLICATIONTYPE_USER_AGENT:
		return proj_model.OIDCApplicationTypeUserAgent
	case grpc.OIDCApplicationType_OIDCAPPLICATIONTYPE_NATIVE:
		return proj_model.OIDCApplicationTypeNative
	}
	return proj_model.OIDCApplicationTypeWeb
}

func oidcApplicationTypeFromModel(appType proj_model.OIDCApplicationType) grpc.OIDCApplicationType {
	switch appType {
	case proj_model.OIDCApplicationTypeWeb:
		return grpc.OIDCApplicationType_OIDCAPPLICATIONTYPE_WEB
	case proj_model.OIDCApplicationTypeUserAgent:
		return grpc.OIDCApplicationType_OIDCAPPLICATIONTYPE_USER_AGENT
	case proj_model.OIDCApplicationTypeNative:
		return grpc.OIDCApplicationType_OIDCAPPLICATIONTYPE_NATIVE
	default:
		return grpc.OIDCApplicationType_OIDCAPPLICATIONTYPE_WEB
	}
}

func oidcAuthMethodTypeToModel(authType grpc.OIDCAuthMethodType) proj_model.OIDCAuthMethodType {
	switch authType {
	case grpc.OIDCAuthMethodType_OIDCAUTHMETHODTYPE_BASIC:
		return proj_model.OIDCAuthMethodTypeBasic
	case grpc.OIDCAuthMethodType_OIDCAUTHMETHODTYPE_POST:
		return proj_model.OIDCAuthMethodTypePost
	case grpc.OIDCAuthMethodType_OIDCAUTHMETHODTYPE_NONE:
		return proj_model.OIDCAuthMethodTypeNone
	default:
		return proj_model.OIDCAuthMethodTypeBasic
	}
}

func oidcAuthMethodTypeFromModel(authType proj_model.OIDCAuthMethodType) grpc.OIDCAuthMethodType {
	switch authType {
	case proj_model.OIDCAuthMethodTypeBasic:
		return grpc.OIDCAuthMethodType_OIDCAUTHMETHODTYPE_BASIC
	case proj_model.OIDCAuthMethodTypePost:
		return grpc.OIDCAuthMethodType_OIDCAUTHMETHODTYPE_POST
	case proj_model.OIDCAuthMethodTypeNone:
		return grpc.OIDCAuthMethodType_OIDCAUTHMETHODTYPE_NONE
	default:
		return grpc.OIDCAuthMethodType_OIDCAUTHMETHODTYPE_BASIC
	}
}

func appChangesToResponse(response *proj_model.ApplicationChanges, offset uint64, limit uint64) (_ *grpc.Changes) {
	return &grpc.Changes{
		Limit:   limit,
		Offset:  offset,
		Changes: appChangesToMgtAPI(response),
	}
}

func appChangesToMgtAPI(changes *proj_model.ApplicationChanges) (_ []*grpc.Change) {
	result := make([]*grpc.Change, len(changes.Changes))

	for i, change := range changes.Changes {
		b, err := json.Marshal(change.Data)
		data := &structpb.Struct{}
		err = protojson.Unmarshal(b, data)
		if err != nil {
		}
		result[i] = &grpc.Change{
			ChangeDate: change.ChangeDate,
			EventType:  change.EventType,
			Sequence:   change.Sequence,
			Editor:     change.ModifierName,
			EditorId:   change.ModifierId,
			Data:       data,
		}
	}

	return result
}
