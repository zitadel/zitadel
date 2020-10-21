package management

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"

	"github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) IdpByID(ctx context.Context, id *management.IdpID) (*management.IdpView, error) {
	config, err := s.org.IDPConfigByID(ctx, id.Id)
	if err != nil {
		return nil, err
	}
	return idpViewFromModel(config), nil
}

func (s *Server) CreateOidcIdp(ctx context.Context, oidcIdpConfig *management.OidcIdpConfigCreate) (*management.Idp, error) {
	config, err := s.org.AddOIDCIDPConfig(ctx, createOidcIdpToModel(oidcIdpConfig))
	if err != nil {
		return nil, err
	}
	return idpFromModel(config), nil
}

func (s *Server) UpdateIdpConfig(ctx context.Context, idpConfig *management.IdpUpdate) (*management.Idp, error) {
	config, err := s.org.ChangeIDPConfig(ctx, updateIdpToModel(idpConfig))
	if err != nil {
		return nil, err
	}
	return idpFromModel(config), nil
}

func (s *Server) DeactivateIdpConfig(ctx context.Context, id *management.IdpID) (*management.Idp, error) {
	config, err := s.org.DeactivateIDPConfig(ctx, id.Id)
	if err != nil {
		return nil, err
	}
	return idpFromModel(config), nil
}

func (s *Server) ReactivateIdpConfig(ctx context.Context, id *management.IdpID) (*management.Idp, error) {
	config, err := s.org.ReactivateIDPConfig(ctx, id.Id)
	if err != nil {
		return nil, err
	}
	return idpFromModel(config), nil
}

func (s *Server) RemoveIdpConfig(ctx context.Context, id *management.IdpID) (*empty.Empty, error) {
	err := s.org.RemoveIDPConfig(ctx, id.Id)
	return &empty.Empty{}, err
}

func (s *Server) UpdateOidcIdpConfig(ctx context.Context, request *management.OidcIdpConfigUpdate) (*management.OidcIdpConfig, error) {
	config, err := s.org.ChangeOIDCIDPConfig(ctx, updateOidcIdpToModel(request))
	if err != nil {
		return nil, err
	}
	return oidcIdpConfigFromModel(config), nil
}

func (s *Server) SearchIdps(ctx context.Context, request *management.IdpSearchRequest) (*management.IdpSearchResponse, error) {
	searchRequest, err := idpConfigSearchRequestToModel(request)
	if err != nil {
		return nil, err
	}
	response, err := s.org.SearchIDPConfigs(ctx, searchRequest)
	if err != nil {
		return nil, err
	}
	return idpConfigSearchResponseFromModel(response), nil
}
