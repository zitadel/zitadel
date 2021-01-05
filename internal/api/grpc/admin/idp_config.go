package admin

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"

	"github.com/caos/zitadel/pkg/grpc/admin"
)

func (s *Server) IdpByID(ctx context.Context, id *admin.IdpID) (*admin.IdpView, error) {
	config, err := s.query.DefaultIDPConfigByID(ctx, id.Id)
	if err != nil {
		return nil, err
	}
	return idpViewFromDomain(config), nil
}

func (s *Server) CreateOidcIdp(ctx context.Context, oidcIdpConfig *admin.OidcIdpConfigCreate) (*admin.Idp, error) {
	config, err := s.command.AddDefaultIDPConfig(ctx, createOIDCIDPToDomain(oidcIdpConfig))
	if err != nil {
		return nil, err
	}
	return idpFromDomain(config), nil
}

func (s *Server) UpdateIdpConfig(ctx context.Context, idpConfig *admin.IdpUpdate) (*admin.Idp, error) {
	config, err := s.command.ChangeDefaultIDPConfig(ctx, updateIdpToDomain(idpConfig))
	if err != nil {
		return nil, err
	}
	return idpFromDomain(config), nil
}

func (s *Server) DeactivateIdpConfig(ctx context.Context, id *admin.IdpID) (*admin.Idp, error) {
	config, err := s.command.DeactivateDefaultIDPConfig(ctx, id.Id)
	if err != nil {
		return nil, err
	}
	return idpFromDomain(config), nil
}

func (s *Server) ReactivateIdpConfig(ctx context.Context, id *admin.IdpID) (*admin.Idp, error) {
	config, err := s.command.ReactivateDefaultIDPConfig(ctx, id.Id)
	if err != nil {
		return nil, err
	}
	return idpFromDomain(config), nil
}

//TODO: Change To V2
func (s *Server) RemoveIdpConfig(ctx context.Context, id *admin.IdpID) (*empty.Empty, error) {
	err := s.iam.RemoveIDPConfig(ctx, id.Id)
	return &empty.Empty{}, err
}

func (s *Server) UpdateOidcIdpConfig(ctx context.Context, request *admin.OidcIdpConfigUpdate) (*admin.OidcIdpConfig, error) {
	config, err := s.command.ChangeDefaultIDPOIDCConfig(ctx, updateOIDCIDPToDomain(request))
	if err != nil {
		return nil, err
	}
	return oidcIDPConfigFromDomain(config), nil
}

func (s *Server) SearchIdps(ctx context.Context, request *admin.IdpSearchRequest) (*admin.IdpSearchResponse, error) {
	response, err := s.iam.SearchIDPConfigs(ctx, idpConfigSearchRequestToModel(request))
	if err != nil {
		return nil, err
	}
	return idpConfigSearchResponseFromModel(response), nil
}
