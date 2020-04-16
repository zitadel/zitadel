package grpc

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
)

func (s *Server) SearchApplications(ctx context.Context, in *ApplicationSearchRequest) (*ApplicationSearchResponse, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-yW23f", "Not implemented")
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

func (s *Server) UpdateApplicationOIDCConfig(ctx context.Context, in *OIDCConfigUpdate) (*OIDCConfig, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-xm56g", "Not implemented")
}
func (s *Server) RegenerateOIDCClientSecret(ctx context.Context, in *ApplicationID) (*ClientSecret, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-dlwp3", "Not implemented")
}

func (s *Server) ApplicationChanges(ctx context.Context, changesRequest *ChangeRequest) (*Changes, error) {
	return nil, errors.ThrowUnimplemented(nil, "GRPC-due45", "Not implemented")
}
