package grpc

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) SearchApplications(ctx context.Context, in *ApplicationSearchRequest) (*ApplicationSearchResponse, error) {
	response, err := s.project.SearchApplications(ctx, applicationSearchRequestsToModel(in))
	if err != nil {
		return nil, err
	}
	return applicationSearchResponseFromModel(response), nil
}

func (s *Server) ApplicationByID(ctx context.Context, in *ApplicationID) (*Application, error) {
	app, err := s.project.ApplicationByID(ctx, in.ProjectId, in.Id)
	if err != nil {
		return nil, err
	}
	return appFromModel(app), nil
}

func (s *Server) CreateOIDCApplication(ctx context.Context, in *OIDCApplicationCreate) (*Application, error) {
	app, err := s.project.AddApplication(ctx, oidcAppCreateToModel(in))
	if err != nil {
		return nil, err
	}
	return appFromModel(app), nil
}
func (s *Server) UpdateApplication(ctx context.Context, in *ApplicationUpdate) (*Application, error) {
	app, err := s.project.ChangeApplication(ctx, appUpdateToModel(in))
	if err != nil {
		return nil, err
	}
	return appFromModel(app), nil
}
func (s *Server) DeactivateApplication(ctx context.Context, in *ApplicationID) (*Application, error) {
	app, err := s.project.DeactivateApplication(ctx, in.ProjectId, in.Id)
	if err != nil {
		return nil, err
	}
	return appFromModel(app), nil
}
func (s *Server) ReactivateApplication(ctx context.Context, in *ApplicationID) (*Application, error) {
	app, err := s.project.ReactivateApplication(ctx, in.ProjectId, in.Id)
	if err != nil {
		return nil, err
	}
	return appFromModel(app), nil
}

func (s *Server) RemoveApplication(ctx context.Context, in *ApplicationID) (*empty.Empty, error) {
	err := s.project.RemoveApplication(ctx, in.ProjectId, in.Id)
	return &empty.Empty{}, err
}

func (s *Server) UpdateApplicationOIDCConfig(ctx context.Context, in *OIDCConfigUpdate) (*OIDCConfig, error) {
	config, err := s.project.ChangeOIDCConfig(ctx, oidcConfigUpdateToModel(in))
	if err != nil {
		return nil, err
	}
	return oidcConfigFromModel(config), nil
}

func (s *Server) RegenerateOIDCClientSecret(ctx context.Context, in *ApplicationID) (*ClientSecret, error) {
	config, err := s.project.ChangeOIDConfigSecret(ctx, in.ProjectId, in.Id)
	if err != nil {
		return nil, err
	}
	return &ClientSecret{ClientSecret: config.ClientSecretString}, nil
}

func (s *Server) ApplicationChanges(ctx context.Context, changesRequest *ChangeRequest) (*Changes, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-due45", "Not implemented")
}
