package management

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	authn_grpc "github.com/zitadel/zitadel/internal/api/grpc/authn"
	change_grpc "github.com/zitadel/zitadel/internal/api/grpc/change"
	object_grpc "github.com/zitadel/zitadel/internal/api/grpc/object"
	project_grpc "github.com/zitadel/zitadel/internal/api/grpc/project"
	"github.com/zitadel/zitadel/internal/query"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
)

func (s *Server) GetAppByID(ctx context.Context, req *mgmt_pb.GetAppByIDRequest) (*mgmt_pb.GetAppByIDResponse, error) {
	app, err := s.query.AppByProjectAndAppID(ctx, true, req.ProjectId, req.AppId)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetAppByIDResponse{
		App: project_grpc.AppToPb(app),
	}, nil
}

func (s *Server) ListApps(ctx context.Context, req *mgmt_pb.ListAppsRequest) (*mgmt_pb.ListAppsResponse, error) {
	queries, err := ListAppsRequestToModel(req)
	if err != nil {
		return nil, err
	}
	apps, err := s.query.SearchApps(ctx, queries)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListAppsResponse{
		Result: project_grpc.AppsToPb(apps.Apps),
		Details: object_grpc.ToListDetails(
			apps.Count,
			apps.Sequence,
			apps.Timestamp,
		),
	}, nil
}

func (s *Server) ListAppChanges(ctx context.Context, req *mgmt_pb.ListAppChangesRequest) (*mgmt_pb.ListAppChangesResponse, error) {
	sequence, limit, asc := change_grpc.ChangeQueryToQuery(req.Query)
	features, err := s.query.FeaturesByOrgID(ctx, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	res, err := s.query.ApplicationChanges(ctx, req.ProjectId, req.AppId, sequence, limit, asc, features.AuditLogRetention)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListAppChangesResponse{
		Result: change_grpc.ChangesToPb(res.Changes, s.assetAPIPrefix),
	}, nil
}

func (s *Server) AddOIDCApp(ctx context.Context, req *mgmt_pb.AddOIDCAppRequest) (*mgmt_pb.AddOIDCAppResponse, error) {
	app, err := s.command.AddOIDCApplication(ctx, AddOIDCAppRequestToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddOIDCAppResponse{
		AppId:              app.AppID,
		Details:            object_grpc.AddToDetailsPb(app.Sequence, app.ChangeDate, app.ResourceOwner),
		ClientId:           app.ClientID,
		ClientSecret:       app.ClientSecretString,
		NoneCompliant:      app.Compliance.NoneCompliant,
		ComplianceProblems: project_grpc.ComplianceProblemsToLocalizedMessages(app.Compliance.Problems),
	}, nil
}

func (s *Server) AddAPIApp(ctx context.Context, req *mgmt_pb.AddAPIAppRequest) (*mgmt_pb.AddAPIAppResponse, error) {
	app, err := s.command.AddAPIApplication(ctx, AddAPIAppRequestToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddAPIAppResponse{
		AppId:        app.AppID,
		Details:      object_grpc.AddToDetailsPb(app.Sequence, app.ChangeDate, app.ResourceOwner),
		ClientId:     app.ClientID,
		ClientSecret: app.ClientSecretString,
	}, nil
}

func (s *Server) UpdateApp(ctx context.Context, req *mgmt_pb.UpdateAppRequest) (*mgmt_pb.UpdateAppResponse, error) {
	details, err := s.command.ChangeApplication(ctx, req.ProjectId, UpdateAppRequestToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateAppResponse{
		Details: object_grpc.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) UpdateOIDCAppConfig(ctx context.Context, req *mgmt_pb.UpdateOIDCAppConfigRequest) (*mgmt_pb.UpdateOIDCAppConfigResponse, error) {
	config, err := s.command.ChangeOIDCApplication(ctx, UpdateOIDCAppConfigRequestToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateOIDCAppConfigResponse{
		Details: object_grpc.ChangeToDetailsPb(
			config.Sequence,
			config.ChangeDate,
			config.ResourceOwner,
		),
	}, nil
}

func (s *Server) UpdateAPIAppConfig(ctx context.Context, req *mgmt_pb.UpdateAPIAppConfigRequest) (*mgmt_pb.UpdateAPIAppConfigResponse, error) {
	config, err := s.command.ChangeAPIApplication(ctx, UpdateAPIAppConfigRequestToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateAPIAppConfigResponse{
		Details: object_grpc.ChangeToDetailsPb(
			config.Sequence,
			config.ChangeDate,
			config.ResourceOwner,
		),
	}, nil
}

func (s *Server) DeactivateApp(ctx context.Context, req *mgmt_pb.DeactivateAppRequest) (*mgmt_pb.DeactivateAppResponse, error) {
	details, err := s.command.DeactivateApplication(ctx, req.ProjectId, req.AppId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.DeactivateAppResponse{
		Details: object_grpc.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) ReactivateApp(ctx context.Context, req *mgmt_pb.ReactivateAppRequest) (*mgmt_pb.ReactivateAppResponse, error) {
	details, err := s.command.ReactivateApplication(ctx, req.ProjectId, req.AppId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ReactivateAppResponse{
		Details: object_grpc.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) RemoveApp(ctx context.Context, req *mgmt_pb.RemoveAppRequest) (*mgmt_pb.RemoveAppResponse, error) {
	details, err := s.command.RemoveApplication(ctx, req.ProjectId, req.AppId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveAppResponse{
		Details: object_grpc.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) RegenerateOIDCClientSecret(ctx context.Context, req *mgmt_pb.RegenerateOIDCClientSecretRequest) (*mgmt_pb.RegenerateOIDCClientSecretResponse, error) {
	config, err := s.command.ChangeOIDCApplicationSecret(ctx, req.ProjectId, req.AppId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RegenerateOIDCClientSecretResponse{
		ClientSecret: config.ClientSecretString,
		Details: object_grpc.ChangeToDetailsPb(
			config.Sequence,
			config.ChangeDate,
			config.ResourceOwner,
		),
	}, nil
}

func (s *Server) RegenerateAPIClientSecret(ctx context.Context, req *mgmt_pb.RegenerateAPIClientSecretRequest) (*mgmt_pb.RegenerateAPIClientSecretResponse, error) {
	config, err := s.command.ChangeAPIApplicationSecret(ctx, req.ProjectId, req.AppId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RegenerateAPIClientSecretResponse{
		ClientSecret: config.ClientSecretString,
		Details: object_grpc.ChangeToDetailsPb(
			config.Sequence,
			config.ChangeDate,
			config.ResourceOwner,
		),
	}, nil
}

func (s *Server) GetAppKey(ctx context.Context, req *mgmt_pb.GetAppKeyRequest) (*mgmt_pb.GetAppKeyResponse, error) {
	resourceOwner, err := query.NewAuthNKeyResourceOwnerQuery(authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	aggregateID, err := query.NewAuthNKeyAggregateIDQuery(req.ProjectId)
	if err != nil {
		return nil, err
	}
	objectID, err := query.NewAuthNKeyObjectIDQuery(req.AppId)
	if err != nil {
		return nil, err
	}
	key, err := s.query.GetAuthNKeyByID(ctx, true, req.KeyId, resourceOwner, aggregateID, objectID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetAppKeyResponse{
		Key: authn_grpc.KeyToPb(key),
	}, nil
}

func (s *Server) ListAppKeys(ctx context.Context, req *mgmt_pb.ListAppKeysRequest) (*mgmt_pb.ListAppKeysResponse, error) {
	queries, err := ListAPIClientKeysRequestToQuery(ctx, req)
	if err != nil {
		return nil, err
	}
	keys, err := s.query.SearchAuthNKeys(ctx, queries)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListAppKeysResponse{
		Result: authn_grpc.KeysToPb(keys.AuthNKeys),
		Details: object_grpc.ToListDetails(
			keys.Count,
			keys.Sequence,
			keys.Timestamp,
		),
	}, nil
}

func (s *Server) AddAppKey(ctx context.Context, req *mgmt_pb.AddAppKeyRequest) (*mgmt_pb.AddAppKeyResponse, error) {
	key, err := s.command.AddApplicationKey(ctx, AddAPIClientKeyRequestToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	keyDetails, err := key.Detail()
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddAppKeyResponse{
		Id:         key.KeyID,
		Details:    object_grpc.AddToDetailsPb(key.Sequence, key.ChangeDate, key.ResourceOwner),
		KeyDetails: keyDetails,
	}, nil
}

func (s *Server) RemoveAppKey(ctx context.Context, req *mgmt_pb.RemoveAppKeyRequest) (*mgmt_pb.RemoveAppKeyResponse, error) {
	details, err := s.command.RemoveApplicationKey(ctx, req.ProjectId, req.AppId, req.KeyId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveAppKeyResponse{
		Details: object_grpc.DomainToChangeDetailsPb(details),
	}, nil
}
