package admin

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"

	"github.com/caos/zitadel/pkg/grpc/admin"
)

func (s *Server) IdpByID(ctx context.Context, id *admin.IdpID) (*admin.IdpView, error) {
	config, err := s.iam.IDPConfigByID(ctx, id.Id)
	if err != nil {
		return nil, err
	}
	return idpViewFromModel(config), nil
}

func (s *Server) CreateOidcIdp(ctx context.Context, oidcIdpConfig *admin.OidcIdpConfigCreate) (*admin.Idp, error) {
	config, err := s.iam.AddOIDCIDPConfig(ctx, createOidcIdpToModel(oidcIdpConfig))
	if err != nil {
		return nil, err
	}
	return idpFromModel(config), nil
}

func (s *Server) UpdateIdpConfig(ctx context.Context, idpConfig *admin.IdpUpdate) (*admin.Idp, error) {
	config, err := s.iam.ChangeIDPConfig(ctx, updateIdpToModel(idpConfig))
	if err != nil {
		return nil, err
	}
	return idpFromModel(config), nil
}

func (s *Server) DeactivateIdpConfig(ctx context.Context, id *admin.IdpID) (*admin.Idp, error) {
	config, err := s.iam.DeactivateIDPConfig(ctx, id.Id)
	if err != nil {
		return nil, err
	}
	return idpFromModel(config), nil
}

func (s *Server) ReactivateIdpConfig(ctx context.Context, id *admin.IdpID) (*admin.Idp, error) {
	config, err := s.iam.ReactivateIDPConfig(ctx, id.Id)
	if err != nil {
		return nil, err
	}
	return idpFromModel(config), nil
}

func (s *Server) RemoveIdpConfig(ctx context.Context, id *admin.IdpID) (*empty.Empty, error) {
	err := s.iam.RemoveIDPConfig(ctx, id.Id)
	return &empty.Empty{}, err
}

func (s *Server) UpdateOidcIdpConfig(ctx context.Context, request *admin.OidcIdpConfigUpdate) (*admin.OidcIdpConfig, error) {
	config, err := s.iam.ChangeOidcIDPConfig(ctx, updateOidcIdpToModel(request))
	if err != nil {
		return nil, err
	}
	return oidcIdpConfigFromModel(config), nil
}

func (s *Server) SearchIdps(ctx context.Context, request *admin.IdpSearchRequest) (*admin.IdpSearchResponse, error) {
	response, err := s.iam.SearchIDPConfigs(ctx, idpConfigSearchRequestToModel(request))
	if err != nil {
		return nil, err
	}
	return idpConfigSearchResponseFromModel(response), nil
}
