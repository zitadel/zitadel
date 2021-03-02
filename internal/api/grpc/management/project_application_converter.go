package management

import (
	app_grpc "github.com/caos/zitadel/internal/api/grpc/project"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func AddOIDCAppRequestToDomain(req *mgmt_pb.AddOIDCAppRequest) *domain.OIDCApp {
	return &domain.OIDCApp{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.ProjectId,
		},
		AppName:                  req.Name,
		OIDCVersion:              app_grpc.OIDCVersionToDomain(req.Version),
		RedirectUris:             req.RedirectUris,
		ResponseTypes:            app_grpc.OIDCResponseTypesToDomain(req.ResponseTypes),
		GrantTypes:               app_grpc.OIDCGrantTypesToDomain(req.GrantTypes),
		ApplicationType:          app_grpc.OIDCApplicationTypeToDomain(req.AppType),
		AuthMethodType:           app_grpc.OIDCAuthMethodTypeToDomain(req.AuthMethodType),
		PostLogoutRedirectUris:   req.PostLogoutRedirectUris,
		DevMode:                  req.DevMode,
		AccessTokenType:          app_grpc.OIDCTokenTypeToDomain(req.AccessTokenType),
		AccessTokenRoleAssertion: req.AccessTokenRoleAssertion,
		IDTokenRoleAssertion:     req.IdTokenRoleAssertion,
		IDTokenUserinfoAssertion: req.IdTokenUserinfoAssertion,
		ClockSkew:                req.ClockSkew.AsDuration(),
	}
}

func AddAPIAppRequestToDomain(app *mgmt_pb.AddAPIAppRequest) *domain.APIApp {
	return &domain.APIApp{
		ObjectRoot: models.ObjectRoot{
			AggregateID: app.ProjectId,
		},
		AppName:        app.Name,
		AuthMethodType: app_grpc.APIAuthMethodTypeToDomain(app.AuthMethodType),
	}
}

func UpdateAppRequestToDomain(app *mgmt_pb.UpdateAppRequest) domain.Application {
	return &domain.ChangeApp{
		AppID:   app.AppId,
		AppName: app.Name,
	}
}
