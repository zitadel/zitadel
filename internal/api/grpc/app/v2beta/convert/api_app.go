package convert

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
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

func PatchAPIApplicationConfigurationRequestToDomain(appID, projectID string, app *app.PatchAPIApplicationConfigurationRequest) *domain.APIApp {
	return &domain.APIApp{
		ObjectRoot: models.ObjectRoot{
			AggregateID: projectID,
		},
		AppID:          appID,
		AuthMethodType: apiAuthMethodTypeToDomain(app.GetAuthMethodType()),
	}
}
