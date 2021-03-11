package management

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	authn_grpc "github.com/caos/zitadel/internal/api/grpc/authn"
	change_grpc "github.com/caos/zitadel/internal/api/grpc/change"
	object_grpc "github.com/caos/zitadel/internal/api/grpc/object"
	project_grpc "github.com/caos/zitadel/internal/api/grpc/project"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) GetAppByID(ctx context.Context, req *mgmt_pb.GetAppByIDRequest) (*mgmt_pb.GetAppByIDResponse, error) {
	app, err := s.project.ApplicationByID(ctx, req.ProjectId, req.AppId)
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
	domains, err := s.project.SearchApplications(ctx, queries)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListAppsResponse{
		Result: project_grpc.AppsToPb(domains.Result),
		Details: object_grpc.ToListDetails(
			domains.TotalResult,
			domains.Sequence,
			domains.Timestamp,
		),
	}, nil
}

func (s *Server) ListAppChanges(ctx context.Context, req *mgmt_pb.ListAppChangesRequest) (*mgmt_pb.ListAppChangesResponse, error) {
	offset, limit, asc := object_grpc.ListQueryToModel(req.Query)
	res, err := s.project.ApplicationChanges(ctx, req.ProjectId, req.AppId, offset, limit, asc)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListAppChangesResponse{
		Result: change_grpc.AppChangesToPb(res.Changes),
	}, nil
}

func (s *Server) AddOIDCApp(ctx context.Context, req *mgmt_pb.AddOIDCAppRequest) (*mgmt_pb.AddOIDCAppResponse, error) {
	app, err := s.command.AddOIDCApplication(ctx, AddOIDCAppRequestToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddOIDCAppResponse{
		AppId:              app.AppID,
		Details:            object_grpc.ToDetailsPb(app.Sequence, app.ChangeDate, app.ResourceOwner),
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
		Details:      object_grpc.ToDetailsPb(app.Sequence, app.ChangeDate, app.ResourceOwner),
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
		Details: object_grpc.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) UpdateOIDCAppConfig(ctx context.Context, req *mgmt_pb.UpdateOIDCAppConfigRequest) (*mgmt_pb.UpdateOIDCAppConfigResponse, error) {
	config, err := s.command.ChangeOIDCApplication(ctx, UpdateOIDCAppConfigRequestToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateOIDCAppConfigResponse{
		Details: object_grpc.ToDetailsPb(
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
		Details: object_grpc.ToDetailsPb(
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
		Details: object_grpc.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) ReactivateApp(ctx context.Context, req *mgmt_pb.ReactivateAppRequest) (*mgmt_pb.ReactivateAppResponse, error) {
	details, err := s.command.ReactivateApplication(ctx, req.ProjectId, req.AppId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ReactivateAppResponse{
		Details: object_grpc.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) RemoveApp(ctx context.Context, req *mgmt_pb.RemoveAppRequest) (*mgmt_pb.RemoveAppResponse, error) {
	details, err := s.command.RemoveApplication(ctx, req.ProjectId, req.AppId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveAppResponse{
		Details: object_grpc.DomainToDetailsPb(details),
	}, nil
}

func (s *Server) RegenerateOIDCClientSecret(ctx context.Context, req *mgmt_pb.RegenerateOIDCClientSecretRequest) (*mgmt_pb.RegenerateOIDCClientSecretResponse, error) {
	config, err := s.command.ChangeOIDCApplicationSecret(ctx, req.ProjectId, req.AppId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RegenerateOIDCClientSecretResponse{
		ClientSecret: config.ClientSecretString,
		Details: object_grpc.ToDetailsPb(
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
		Details: object_grpc.ToDetailsPb(
			config.Sequence,
			config.ChangeDate,
			config.ResourceOwner,
		),
	}, nil
}

func (s *Server) GetAppKey(ctx context.Context, req *mgmt_pb.GetAppKeyRequest) (*mgmt_pb.GetAppKeyResponse, error) {
	key, err := s.project.GetClientKey(ctx, req.ProjectId, req.AppId, req.KeyId)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetAppKeyResponse{
		Key: authn_grpc.KeyToPb(key),
	}, nil
}

func (s *Server) ListAppKeys(ctx context.Context, req *mgmt_pb.ListAppKeysRequest) (*mgmt_pb.ListAppKeysResponse, error) {
	queries, err := ListAPIClientKeysRequestToModel(req)
	if err != nil {
		return nil, err
	}
	domains, err := s.project.SearchClientKeys(ctx, queries)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListAppKeysResponse{
		Result: authn_grpc.KeyViewsToPb(domains.Result),
		Details: object_grpc.ToListDetails(
			domains.TotalResult,
			domains.Sequence,
			domains.Timestamp,
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
		Details:    object_grpc.ToDetailsPb(key.Sequence, key.ChangeDate, key.ResourceOwner),
		KeyDetails: keyDetails,
	}, nil
}

func (s *Server) RemoveAppKey(ctx context.Context, req *mgmt_pb.RemoveAppKeyRequest) (*mgmt_pb.RemoveAppKeyResponse, error) {
	details, err := s.command.RemoveApplicationKey(ctx, req.ProjectId, req.AppId, req.KeyId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveAppKeyResponse{
		Details: object_grpc.DomainToDetailsPb(details),
	}, nil
}
