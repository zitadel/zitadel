package management

import (
	"time"

	authn_grpc "github.com/caos/zitadel/internal/api/grpc/authn"
	"github.com/caos/zitadel/internal/api/grpc/object"
	app_grpc "github.com/caos/zitadel/internal/api/grpc/project"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	key_model "github.com/caos/zitadel/internal/key/model"
	proj_model "github.com/caos/zitadel/internal/project/model"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func ListAppsRequestToModel(req *mgmt_pb.ListAppsRequest) (*proj_model.ApplicationSearchRequest, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries, err := app_grpc.AppQueriesToModel(req.Queries)
	if err != nil {
		return nil, err
	}
	queries = append(queries, &proj_model.ApplicationSearchQuery{
		Key:    proj_model.AppSearchKeyProjectID,
		Method: domain.SearchMethodEquals,
		Value:  req.ProjectId,
	})
	return &proj_model.ApplicationSearchRequest{
		Offset: offset,
		Limit:  limit,
		Asc:    asc,
		//SortingColumn: //TODO: sorting
		Queries: queries,
	}, nil
}

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

func UpdateOIDCAppConfigRequestToDomain(app *mgmt_pb.UpdateOIDCAppConfigRequest) *domain.OIDCApp {
	return &domain.OIDCApp{
		ObjectRoot: models.ObjectRoot{
			AggregateID: app.ProjectId,
		},
		AppID:                    app.AppId,
		RedirectUris:             app.RedirectUris,
		ResponseTypes:            app_grpc.OIDCResponseTypesToDomain(app.ResponseTypes),
		GrantTypes:               app_grpc.OIDCGrantTypesToDomain(app.GrantTypes),
		ApplicationType:          app_grpc.OIDCApplicationTypeToDomain(app.AppType),
		AuthMethodType:           app_grpc.OIDCAuthMethodTypeToDomain(app.AuthMethodType),
		PostLogoutRedirectUris:   app.PostLogoutRedirectUris,
		DevMode:                  app.DevMode,
		AccessTokenType:          app_grpc.OIDCTokenTypeToDomain(app.AccessTokenType),
		AccessTokenRoleAssertion: app.AccessTokenRoleAssertion,
		IDTokenRoleAssertion:     app.IdTokenRoleAssertion,
		IDTokenUserinfoAssertion: app.IdTokenUserinfoAssertion,
		ClockSkew:                app.ClockSkew.AsDuration(),
	}
}

func UpdateAPIAppConfigRequestToDomain(app *mgmt_pb.UpdateAPIAppConfigRequest) *domain.APIApp {
	return &domain.APIApp{
		ObjectRoot: models.ObjectRoot{
			AggregateID: app.ProjectId,
		},
		AppID:          app.AppId,
		AuthMethodType: app_grpc.APIAuthMethodTypeToDomain(app.AuthMethodType),
	}
}

func AddAPIClientKeyRequestToDomain(key *mgmt_pb.AddAppKeyRequest) *domain.ApplicationKey {
	expirationDate := time.Time{}
	if key.ExpirationDate != nil {
		expirationDate = key.ExpirationDate.AsTime()
	}

	return &domain.ApplicationKey{
		ObjectRoot: models.ObjectRoot{
			AggregateID: key.ProjectId,
		},
		ExpirationDate: expirationDate,
		Type:           authn_grpc.KeyTypeToDomain(key.Type),
		ApplicationID:  key.AppId,
	}
}

func ListAPIClientKeysRequestToModel(req *mgmt_pb.ListAppKeysRequest) (*key_model.AuthNKeySearchRequest, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries := make([]*key_model.AuthNKeySearchQuery, 2)
	queries = append(queries, &key_model.AuthNKeySearchQuery{
		Key:    key_model.AuthNKeyObjectID,
		Method: domain.SearchMethodEquals,
		Value:  req.AppId,
	})
	return &key_model.AuthNKeySearchRequest{
		Offset: offset,
		Limit:  limit,
		Asc:    asc,
		//SortingColumn: //TODO: sorting
		Queries: queries,
	}, nil
}
