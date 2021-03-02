package management

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/caos/zitadel/internal/api/authz"
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
	return nil, status.Errorf(codes.Unimplemented, "method ListApps not implemented")
}

func (s *Server) ListAppChanges(ctx context.Context, req *mgmt_pb.ListAppChangesRequest) (*mgmt_pb.ListAppChangesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAppChanges not implemented")
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
	_, err := s.command.ChangeApplication(ctx, req.ProjectId, UpdateAppRequestToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateAppResponse{
		//TODO: details
		//Details: object_grpc.ToDetailsPb(
		//	app.Sequence,
		//	app.ChangeDate,
		//	app.ResourceOwner,
		//),
	}, nil
}
func (s *Server) UpdateOIDCAppConfig(ctx context.Context, req *mgmt_pb.UpdateOIDCAppConfigRequest) (*mgmt_pb.UpdateOIDCAppConfigResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateOIDCAppConfig not implemented")
}
func (s *Server) UpdateAPIAppConfig(ctx context.Context, req *mgmt_pb.UpdateAPIAppConfigRequest) (*mgmt_pb.UpdateAPIAppConfigResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateAPIAppConfig not implemented")
}
func (s *Server) DeactivateApp(ctx context.Context, req *mgmt_pb.DeactivateAppRequest) (*mgmt_pb.DeactivateAppResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeactivateApp not implemented")
}
func (s *Server) ReactivateApp(ctx context.Context, req *mgmt_pb.ReactivateAppRequest) (*mgmt_pb.ReactivateAppResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReactivateApp not implemented")
}
func (s *Server) RemoveApp(ctx context.Context, req *mgmt_pb.RemoveAppRequest) (*mgmt_pb.RemoveAppResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveApp not implemented")
}
func (s *Server) RegenerateOIDCClientSecret(ctx context.Context, req *mgmt_pb.RegenerateOIDCClientSecretRequest) (*mgmt_pb.RegenerateOIDCClientSecretResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegenerateOIDCClientSecret not implemented")
}
func (s *Server) RegenerateAPIClientSecret(ctx context.Context, req *mgmt_pb.RegenerateAPIClientSecretRequest) (*mgmt_pb.RegenerateAPIClientSecretResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegenerateAPIClientSecret not implemented")
}
func (s *Server) GetAPIClientKey(ctx context.Context, req *mgmt_pb.GetAPIClientKeyRequest) (*mgmt_pb.GetAPIClientKeyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAPIClientKey not implemented")
}
func (s *Server) ListAPIClientKeys(ctx context.Context, req *mgmt_pb.ListAPIClientKeysRequest) (*mgmt_pb.ListAPIClientKeysResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAPIClientKeys not implemented")
}
func (s *Server) AddAPIClientKey(ctx context.Context, req *mgmt_pb.AddAPIClientKeyRequest) (*mgmt_pb.AddAPIClientKeyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddAPIClientKey not implemented")
}
func (s *Server) RemoveAPIClientKey(ctx context.Context, req *mgmt_pb.RemoveAPIClientKeyRequest) (*mgmt_pb.RemoveAPIClientKeyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveAPIClientKey not implemented")
}
