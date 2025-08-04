package management

import (
	"context"
	"time"

	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/internal/api/authz"
	authn_grpc "github.com/zitadel/zitadel/internal/api/grpc/authn"
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	app_grpc "github.com/zitadel/zitadel/internal/api/grpc/project"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/query"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
)

func ListAppsRequestToModel(req *mgmt_pb.ListAppsRequest) (*query.AppSearchQueries, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	queries, err := app_grpc.AppQueriesToModel(req.Queries)
	if err != nil {
		return nil, err
	}
	projectQuery, err := query.NewAppProjectIDSearchQuery(req.ProjectId)
	if err != nil {
		return nil, err
	}
	queries = append(queries, projectQuery)
	return &query.AppSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
		//SortingColumn: //TODO: sorting
		Queries: queries,
	}, nil
}

func AddOIDCAppRequestToDomain(req *mgmt_pb.AddOIDCAppRequest) (*domain.OIDCApp, error) {
	loginVersion, loginBaseURI, err := app_grpc.LoginVersionToDomain(req.GetLoginVersion())
	if err != nil {
		return nil, err
	}
	return &domain.OIDCApp{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.ProjectId,
		},
		AppName:                  req.Name,
		OIDCVersion:              gu.Ptr(app_grpc.OIDCVersionToDomain(req.Version)),
		RedirectUris:             req.RedirectUris,
		ResponseTypes:            app_grpc.OIDCResponseTypesToDomain(req.ResponseTypes),
		GrantTypes:               app_grpc.OIDCGrantTypesToDomain(req.GrantTypes),
		ApplicationType:          gu.Ptr(app_grpc.OIDCApplicationTypeToDomain(req.AppType)),
		AuthMethodType:           gu.Ptr(app_grpc.OIDCAuthMethodTypeToDomain(req.AuthMethodType)),
		PostLogoutRedirectUris:   req.PostLogoutRedirectUris,
		DevMode:                  gu.Ptr(req.GetDevMode()),
		AccessTokenType:          gu.Ptr(app_grpc.OIDCTokenTypeToDomain(req.AccessTokenType)),
		AccessTokenRoleAssertion: gu.Ptr(req.GetAccessTokenRoleAssertion()),
		IDTokenRoleAssertion:     gu.Ptr(req.GetIdTokenRoleAssertion()),
		IDTokenUserinfoAssertion: gu.Ptr(req.GetIdTokenUserinfoAssertion()),
		ClockSkew:                gu.Ptr(req.GetClockSkew().AsDuration()),
		AdditionalOrigins:        req.AdditionalOrigins,
		SkipNativeAppSuccessPage: gu.Ptr(req.GetSkipNativeAppSuccessPage()),
		BackChannelLogoutURI:     gu.Ptr(req.GetBackChannelLogoutUri()),
		LoginVersion:             gu.Ptr(loginVersion),
		LoginBaseURI:             gu.Ptr(loginBaseURI),
	}, nil
}

func AddSAMLAppRequestToDomain(req *mgmt_pb.AddSAMLAppRequest) (*domain.SAMLApp, error) {
	loginVersion, loginBaseURI, err := app_grpc.LoginVersionToDomain(req.GetLoginVersion())
	if err != nil {
		return nil, err
	}
	return &domain.SAMLApp{
		ObjectRoot: models.ObjectRoot{
			AggregateID: req.ProjectId,
		},
		AppName:      req.Name,
		Metadata:     req.GetMetadataXml(),
		MetadataURL:  gu.Ptr(req.GetMetadataUrl()),
		LoginVersion: gu.Ptr(loginVersion),
		LoginBaseURI: gu.Ptr(loginBaseURI),
	}, nil
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

func UpdateOIDCAppConfigRequestToDomain(app *mgmt_pb.UpdateOIDCAppConfigRequest) (*domain.OIDCApp, error) {
	loginVersion, loginBaseURI, err := app_grpc.LoginVersionToDomain(app.GetLoginVersion())
	if err != nil {
		return nil, err
	}
	return &domain.OIDCApp{
		ObjectRoot: models.ObjectRoot{
			AggregateID: app.ProjectId,
		},
		AppID:                    app.AppId,
		RedirectUris:             app.RedirectUris,
		ResponseTypes:            app_grpc.OIDCResponseTypesToDomain(app.ResponseTypes),
		GrantTypes:               app_grpc.OIDCGrantTypesToDomain(app.GrantTypes),
		ApplicationType:          gu.Ptr(app_grpc.OIDCApplicationTypeToDomain(app.AppType)),
		AuthMethodType:           gu.Ptr(app_grpc.OIDCAuthMethodTypeToDomain(app.AuthMethodType)),
		PostLogoutRedirectUris:   app.PostLogoutRedirectUris,
		DevMode:                  gu.Ptr(app.GetDevMode()),
		AccessTokenType:          gu.Ptr(app_grpc.OIDCTokenTypeToDomain(app.AccessTokenType)),
		AccessTokenRoleAssertion: gu.Ptr(app.GetAccessTokenRoleAssertion()),
		IDTokenRoleAssertion:     gu.Ptr(app.GetIdTokenRoleAssertion()),
		IDTokenUserinfoAssertion: gu.Ptr(app.GetIdTokenUserinfoAssertion()),
		ClockSkew:                gu.Ptr(app.GetClockSkew().AsDuration()),
		AdditionalOrigins:        app.AdditionalOrigins,
		SkipNativeAppSuccessPage: gu.Ptr(app.GetSkipNativeAppSuccessPage()),
		BackChannelLogoutURI:     gu.Ptr(app.GetBackChannelLogoutUri()),
		LoginVersion:             gu.Ptr(loginVersion),
		LoginBaseURI:             gu.Ptr(loginBaseURI),
	}, nil
}

func UpdateSAMLAppConfigRequestToDomain(app *mgmt_pb.UpdateSAMLAppConfigRequest) (*domain.SAMLApp, error) {
	loginVersion, loginBaseURI, err := app_grpc.LoginVersionToDomain(app.GetLoginVersion())
	if err != nil {
		return nil, err
	}
	return &domain.SAMLApp{
		ObjectRoot: models.ObjectRoot{
			AggregateID: app.ProjectId,
		},
		AppID:        app.AppId,
		Metadata:     app.GetMetadataXml(),
		MetadataURL:  gu.Ptr(app.GetMetadataUrl()),
		LoginVersion: gu.Ptr(loginVersion),
		LoginBaseURI: gu.Ptr(loginBaseURI),
	}, nil
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

func ListAPIClientKeysRequestToQuery(ctx context.Context, req *mgmt_pb.ListAppKeysRequest) (*query.AuthNKeySearchQueries, error) {
	resourceOwner, err := query.NewAuthNKeyResourceOwnerQuery(authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	projectID, err := query.NewAuthNKeyAggregateIDQuery(req.ProjectId)
	if err != nil {
		return nil, err
	}
	appID, err := query.NewAuthNKeyObjectIDQuery(req.AppId)
	if err != nil {
		return nil, err
	}
	offset, limit, asc := object.ListQueryToModel(req.Query)
	return &query.AuthNKeySearchQueries{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
		Queries: []query.SearchQuery{
			resourceOwner,
			projectID,
			appID,
		},
	}, nil
}
