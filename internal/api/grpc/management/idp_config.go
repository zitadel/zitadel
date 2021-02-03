package management

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/caos/zitadel/internal/api/authz"
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
	config, err := s.command.AddIDPConfig(ctx, createOidcIdpToDomain(oidcIdpConfig))
	if err != nil {
		return nil, err
	}
	return idpFromDomain(config), nil
}

func (s *Server) UpdateIdpConfig(ctx context.Context, idpConfig *management.IdpUpdate) (*management.Idp, error) {
	config, err := s.command.ChangeIDPConfig(ctx, updateIdpToDomain(ctx, idpConfig))
	if err != nil {
		return nil, err
	}
	return idpFromDomain(config), nil
}

func (s *Server) DeactivateIdpConfig(ctx context.Context, id *management.IdpID) (*empty.Empty, error) {
	err := s.command.DeactivateIDPConfig(ctx, id.Id, authz.GetCtxData(ctx).OrgID)
	return &empty.Empty{}, err
}

func (s *Server) ReactivateIdpConfig(ctx context.Context, id *management.IdpID) (*empty.Empty, error) {
	err := s.command.ReactivateIDPConfig(ctx, id.Id, authz.GetCtxData(ctx).OrgID)
	return &empty.Empty{}, err
}

func (s *Server) RemoveIdpConfig(ctx context.Context, id *management.IdpID) (*empty.Empty, error) {
	externalIdps, err := s.user.ExternalIDPsByIDPConfigID(ctx, id.Id)
	if err != nil {
		return &empty.Empty{}, err
	}
	providers, err := s.org.GetIDPProvidersByIDPConfigID(ctx, authz.GetCtxData(ctx).OrgID, id.Id)
	if err != nil {
		return &empty.Empty{}, err
	}
	err = s.command.RemoveIDPConfig(ctx, id.Id, authz.GetCtxData(ctx).OrgID, len(providers) > 0, externalIDPViewsToDomain(externalIdps)...)
	return &empty.Empty{}, err
}

func (s *Server) UpdateOidcIdpConfig(ctx context.Context, request *management.OidcIdpConfigUpdate) (*management.OidcIdpConfig, error) {
	config, err := s.command.ChangeIDPOIDCConfig(ctx, updateOidcIdpToDomain(ctx, request))
	if err != nil {
		return nil, err
	}
	return oidcIdpConfigFromDomain(config), nil
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
