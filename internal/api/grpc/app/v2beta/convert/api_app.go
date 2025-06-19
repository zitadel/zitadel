package convert

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/query"
	app "github.com/zitadel/zitadel/pkg/grpc/app/v2beta"
)

func CreateAPIApplicationRequestToDomain(name, projectID string, app *app.CreateAPIApplicationRequest) *domain.APIApp {
	return &domain.APIApp{
		ObjectRoot: models.ObjectRoot{
			AggregateID: projectID,
		},
		AppName:        name,
		AuthMethodType: apiAuthMethodTypeToDomain(app.GetAuthMethodType()),
	}
}

func PatchAPIApplicationConfigurationRequestToDomain(appID, projectID string, app *app.UpdateAPIApplicationConfigurationRequest) *domain.APIApp {
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
