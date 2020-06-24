package management

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/caos/zitadel/pkg/management/grpc"
)

func (s *Server) SearchApplications(ctx context.Context, in *grpc.ApplicationSearchRequest) (*grpc.ApplicationSearchResponse, error) {
	response, err := s.project.SearchApplications(ctx, applicationSearchRequestsToModel(in))
	if err != nil {
		return nil, err
	}
	return applicationSearchResponseFromModel(response), nil
}

func (s *Server) ApplicationByID(ctx context.Context, in *grpc.ApplicationID) (*grpc.ApplicationView, error) {
	app, err := s.project.ApplicationByID(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return applicationViewFromModel(app), nil
}

func (s *Server) CreateOIDCApplication(ctx context.Context, in *grpc.OIDCApplicationCreate) (*grpc.Application, error) {
	app, err := s.project.AddApplication(ctx, oidcAppCreateToModel(in))
	if err != nil {
		return nil, err
	}
	return appFromModel(app), nil
}
func (s *Server) UpdateApplication(ctx context.Context, in *grpc.ApplicationUpdate) (*grpc.Application, error) {
	app, err := s.project.ChangeApplication(ctx, appUpdateToModel(in))
	if err != nil {
		return nil, err
	}
	return appFromModel(app), nil
}
func (s *Server) DeactivateApplication(ctx context.Context, in *grpc.ApplicationID) (*grpc.Application, error) {
	app, err := s.project.DeactivateApplication(ctx, in.ProjectId, in.Id)
	if err != nil {
		return nil, err
	}
	return appFromModel(app), nil
}
func (s *Server) ReactivateApplication(ctx context.Context, in *grpc.ApplicationID) (*grpc.Application, error) {
	app, err := s.project.ReactivateApplication(ctx, in.ProjectId, in.Id)
	if err != nil {
		return nil, err
	}
	return appFromModel(app), nil
}

func (s *Server) RemoveApplication(ctx context.Context, in *grpc.ApplicationID) (*empty.Empty, error) {
	err := s.project.RemoveApplication(ctx, in.ProjectId, in.Id)
	return &empty.Empty{}, err
}

func (s *Server) UpdateApplicationOIDCConfig(ctx context.Context, in *grpc.OIDCConfigUpdate) (*grpc.OIDCConfig, error) {
	config, err := s.project.ChangeOIDCConfig(ctx, oidcConfigUpdateToModel(in))
	if err != nil {
		return nil, err
	}
	return oidcConfigFromModel(config), nil
}

func (s *Server) RegenerateOIDCClientSecret(ctx context.Context, in *grpc.ApplicationID) (*grpc.ClientSecret, error) {
	config, err := s.project.ChangeOIDConfigSecret(ctx, in.ProjectId, in.Id)
	if err != nil {
		return nil, err
	}
	return &grpc.ClientSecret{ClientSecret: config.ClientSecretString}, nil
}

func (s *Server) ApplicationChanges(ctx context.Context, changesRequest *grpc.ChangeRequest) (*grpc.Changes, error) {
	response, err := s.project.ApplicationChanges(ctx, changesRequest.Id, changesRequest.SecId, 0, 0)
	if err != nil {
		return nil, err
	}
	return appChangesToResponse(response, changesRequest.GetSequenceOffset(), changesRequest.GetLimit()), nil
}
