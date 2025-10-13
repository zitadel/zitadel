package convert

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/application/v2"
)

func CreateAPIApplicationRequestToDomain(name, projectID, appID string, app *application.CreateAPIApplicationRequest) *domain.APIApp {
	return &domain.APIApp{
		ObjectRoot: models.ObjectRoot{
			AggregateID: projectID,
		},
		AppName:        name,
		AppID:          appID,
		AuthMethodType: apiAuthMethodTypeToDomain(app.GetAuthMethodType()),
	}
}

func UpdateAPIApplicationConfigurationRequestToDomain(appID, projectID string, app *application.UpdateAPIApplicationConfigurationRequest) *domain.APIApp {
	return &domain.APIApp{
		ObjectRoot: models.ObjectRoot{
			AggregateID: projectID,
		},
		AppID:          appID,
		AuthMethodType: apiAuthMethodTypeToDomain(app.GetAuthMethodType()),
	}
}

func appAPIConfigToPb(apiApp *query.APIApp) application.IsApplicationConfiguration {
	return &application.Application_ApiConfiguration{
		ApiConfiguration: &application.APIConfiguration{
			ClientId:       apiApp.ClientID,
			AuthMethodType: apiAuthMethodTypeToPb(apiApp.AuthMethodType),
		},
	}
}

func apiAuthMethodTypeToDomain(authType application.APIAuthMethodType) domain.APIAuthMethodType {
	switch authType {
	case application.APIAuthMethodType_API_AUTH_METHOD_TYPE_BASIC:
		return domain.APIAuthMethodTypeBasic
	case application.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT:
		return domain.APIAuthMethodTypePrivateKeyJWT
	default:
		return domain.APIAuthMethodTypeBasic
	}
}

func apiAuthMethodTypeToPb(methodType domain.APIAuthMethodType) application.APIAuthMethodType {
	switch methodType {
	case domain.APIAuthMethodTypeBasic:
		return application.APIAuthMethodType_API_AUTH_METHOD_TYPE_BASIC
	case domain.APIAuthMethodTypePrivateKeyJWT:
		return application.APIAuthMethodType_API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT
	default:
		return application.APIAuthMethodType_API_AUTH_METHOD_TYPE_BASIC
	}
}
