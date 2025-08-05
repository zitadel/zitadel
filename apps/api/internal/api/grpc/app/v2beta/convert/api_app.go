package convert

import (
	"strings"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/query"
	app "github.com/zitadel/zitadel/pkg/grpc/app/v2beta"
)

func CreateAPIApplicationRequestToDomain(name, projectID, appID string, app *app.CreateAPIApplicationRequest) *domain.APIApp {
	return &domain.APIApp{
		ObjectRoot: models.ObjectRoot{
			AggregateID: projectID,
		},
		AppName:        name,
		AppID:          appID,
		AuthMethodType: apiAuthMethodTypeToDomain(app.GetAuthMethodType()),
	}
}

func UpdateAPIApplicationConfigurationRequestToDomain(appID, projectID string, app *app.UpdateAPIApplicationConfigurationRequest) *domain.APIApp {
	return &domain.APIApp{
		ObjectRoot: models.ObjectRoot{
			AggregateID: projectID,
		},
		AppID:          appID,
		AuthMethodType: apiAuthMethodTypeToDomain(app.GetAuthMethodType()),
	}
}

func appAPIConfigToPb(apiApp *query.APIApp) app.ApplicationConfig {
	return &app.Application_ApiConfig{
		ApiConfig: &app.APIConfig{
			ClientId:       apiApp.ClientID,
			AuthMethodType: apiAuthMethodTypeToPb(apiApp.AuthMethodType),
		},
	}
}

func apiAuthMethodTypeToDomain(authType app.APIAuthMethodType) domain.APIAuthMethodType {
	switch authType {
	case app.APIAuthMethodType_API_AUTH_METHOD_TYPE_BASIC:
		return domain.APIAuthMethodTypeBasic
	case app.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT:
		return domain.APIAuthMethodTypePrivateKeyJWT
	default:
		return domain.APIAuthMethodTypeBasic
	}
}

func apiAuthMethodTypeToPb(methodType domain.APIAuthMethodType) app.APIAuthMethodType {
	switch methodType {
	case domain.APIAuthMethodTypeBasic:
		return app.APIAuthMethodType_API_AUTH_METHOD_TYPE_BASIC
	case domain.APIAuthMethodTypePrivateKeyJWT:
		return app.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT
	default:
		return app.APIAuthMethodType_API_AUTH_METHOD_TYPE_BASIC
	}
}

func GetApplicationKeyQueriesRequestToDomain(orgID, projectID, appID string) ([]query.SearchQuery, error) {
	var searchQueries []query.SearchQuery

	orgID, projectID, appID = strings.TrimSpace(orgID), strings.TrimSpace(projectID), strings.TrimSpace(appID)

	if orgID != "" {
		resourceOwner, err := query.NewAuthNKeyResourceOwnerQuery(orgID)
		if err != nil {
			return nil, err
		}

		searchQueries = append(searchQueries, resourceOwner)
	}

	if projectID != "" {
		aggregateID, err := query.NewAuthNKeyAggregateIDQuery(projectID)
		if err != nil {
			return nil, err
		}

		searchQueries = append(searchQueries, aggregateID)
	}

	if appID != "" {
		objectID, err := query.NewAuthNKeyObjectIDQuery(appID)

		if err != nil {
			return nil, err
		}

		searchQueries = append(searchQueries, objectID)
	}

	return searchQueries, nil
}
