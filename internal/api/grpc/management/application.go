package management

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) SearchApplications(ctx context.Context, in *management.ApplicationSearchRequest) (*management.ApplicationSearchResponse, error) {
	response, err := s.project.SearchApplications(ctx, applicationSearchRequestsToModel(in))
	if err != nil {
		return nil, err
	}
	return applicationSearchResponseFromModel(response), nil
}

func (s *Server) ApplicationByID(ctx context.Context, in *management.ApplicationID) (*management.ApplicationView, error) {
	app, err := s.project.ApplicationByID(ctx, in.ProjectId, in.Id)
	if err != nil {
		return nil, err
	}
	return applicationViewFromModel(app), nil
}

func (s *Server) CreateOIDCApplication(ctx context.Context, in *management.OIDCApplicationCreate) (*management.Application, error) {
	app, err := s.project.AddApplication(ctx, oidcAppCreateToModel(in))
	if err != nil {
		return nil, err
	}
	return appFromModel(app), nil
}
func (s *Server) UpdateApplication(ctx context.Context, in *management.ApplicationUpdate) (*management.Application, error) {
	app, err := s.project.ChangeApplication(ctx, appUpdateToModel(in))
	if err != nil {
		return nil, err
	}
	return appFromModel(app), nil
}
func (s *Server) DeactivateApplication(ctx context.Context, in *management.ApplicationID) (*management.Application, error) {
	app, err := s.project.DeactivateApplication(ctx, in.ProjectId, in.Id)
	if err != nil {
		return nil, err
	}
	return appFromModel(app), nil
}
func (s *Server) ReactivateApplication(ctx context.Context, in *management.ApplicationID) (*management.Application, error) {
	app, err := s.project.ReactivateApplication(ctx, in.ProjectId, in.Id)
	if err != nil {
		return nil, err
	}
	return appFromModel(app), nil
}

func (s *Server) RemoveApplication(ctx context.Context, in *management.ApplicationID) (*empty.Empty, error) {
	err := s.project.RemoveApplication(ctx, in.ProjectId, in.Id)
	return &empty.Empty{}, err
}

func (s *Server) UpdateApplicationOIDCConfig(ctx context.Context, in *management.OIDCConfigUpdate) (*management.OIDCConfig, error) {
	config, err := s.project.ChangeOIDCConfig(ctx, oidcConfigUpdateToModel(in))
	if err != nil {
		return nil, err
	}
	return oidcConfigFromModel(config), nil
}

func (s *Server) RegenerateOIDCClientSecret(ctx context.Context, in *management.ApplicationID) (*management.ClientSecret, error) {
	config, err := s.project.ChangeOIDConfigSecret(ctx, in.ProjectId, in.Id)
	if err != nil {
		return nil, err
	}
	return &management.ClientSecret{ClientSecret: config.ClientSecretString}, nil
}

func (s *Server) ApplicationChanges(ctx context.Context, changesRequest *management.ChangeRequest) (*management.Changes, error) {
	response, err := s.project.ApplicationChanges(ctx, changesRequest.Id, changesRequest.SecId, changesRequest.SequenceOffset, changesRequest.Limit, changesRequest.Asc)
	if err != nil {
		return nil, err
	}
	return appChangesToResponse(response, changesRequest.GetSequenceOffset(), changesRequest.GetLimit()), nil
}

func (s *Server) SearchApplicationKeys(ctx context.Context, req *management.ApplicationKeySearchRequest) (*management.ApplicationKeySearchResponse, error) {
	result, err := s.project.SearchApplicationKeys(ctx, applicationKeySearchRequestToModel(req))
	if err != nil {
		return nil, err
	}
	return applicationKeySearchResponseFromModel(result), nil
}

func (s *Server) GetApplicationKey(ctx context.Context, req *management.ApplicationKeyIDRequest) (*management.ApplicationKeyView, error) {
	key, err := s.project.GetApplicationKey(ctx, req.ProjectId, req.ApplicationId, req.KeyId)
	if err != nil {
		return nil, err
	}
	return applicationKeyViewFromModel(key), nil
}

func (s *Server) AddApplicationKey(ctx context.Context, req *management.AddApplicationKeyRequest) (*management.AddApplicationKeyResponse, error) {
	key, err := s.project.AddApplicationKey(ctx, addApplicationKeyToModel(req))
	if err != nil {
		return nil, err
	}
	return addApplicationKeyFromModel(key), nil
}

func (s *Server) DeleteApplicationKey(ctx context.Context, req *management.ApplicationKeyIDRequest) (*empty.Empty, error) {
	err := s.project.RemoveApplicationKey(ctx, req.ProjectId, req.ApplicationId, req.KeyId)
	return &empty.Empty{}, err
}
