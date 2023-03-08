package management

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	idp_grpc "github.com/zitadel/zitadel/internal/api/grpc/idp"
	object_pb "github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/query"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
)

func (s *Server) GetOrgIDPByID(ctx context.Context, req *mgmt_pb.GetOrgIDPByIDRequest) (*mgmt_pb.GetOrgIDPByIDResponse, error) {
	idp, err := s.query.IDPByIDAndResourceOwner(ctx, true, req.Id, authz.GetCtxData(ctx).OrgID, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetOrgIDPByIDResponse{Idp: idp_grpc.ModelIDPViewToPb(idp)}, nil
}

func (s *Server) ListOrgIDPs(ctx context.Context, req *mgmt_pb.ListOrgIDPsRequest) (*mgmt_pb.ListOrgIDPsResponse, error) {
	queries, err := listIDPsToModel(ctx, req)
	if err != nil {
		return nil, err
	}
	resp, err := s.query.IDPs(ctx, queries, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListOrgIDPsResponse{
		Result:  idp_grpc.IDPViewsToPb(resp.IDPs),
		Details: object_pb.ToListDetails(resp.Count, resp.Sequence, resp.Timestamp),
	}, nil
}

func (s *Server) AddOrgOIDCIDP(ctx context.Context, req *mgmt_pb.AddOrgOIDCIDPRequest) (*mgmt_pb.AddOrgOIDCIDPResponse, error) {
	config, err := s.command.AddIDPConfig(ctx, AddOIDCIDPRequestToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddOrgOIDCIDPResponse{
		IdpId: config.IDPConfigID,
		Details: object_pb.AddToDetailsPb(
			config.Sequence,
			config.ChangeDate,
			config.ResourceOwner,
		),
	}, nil
}

func (s *Server) AddOrgJWTIDP(ctx context.Context, req *mgmt_pb.AddOrgJWTIDPRequest) (*mgmt_pb.AddOrgJWTIDPResponse, error) {
	config, err := s.command.AddIDPConfig(ctx, AddJWTIDPRequestToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddOrgJWTIDPResponse{
		IdpId: config.IDPConfigID,
		Details: object_pb.AddToDetailsPb(
			config.Sequence,
			config.ChangeDate,
			config.ResourceOwner,
		),
	}, nil
}

func (s *Server) DeactivateOrgIDP(ctx context.Context, req *mgmt_pb.DeactivateOrgIDPRequest) (*mgmt_pb.DeactivateOrgIDPResponse, error) {
	objectDetails, err := s.command.DeactivateIDPConfig(ctx, req.IdpId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.DeactivateOrgIDPResponse{Details: object_pb.DomainToChangeDetailsPb(objectDetails)}, nil
}

func (s *Server) ReactivateOrgIDP(ctx context.Context, req *mgmt_pb.ReactivateOrgIDPRequest) (*mgmt_pb.ReactivateOrgIDPResponse, error) {
	objectDetails, err := s.command.ReactivateIDPConfig(ctx, req.IdpId, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ReactivateOrgIDPResponse{Details: object_pb.DomainToChangeDetailsPb(objectDetails)}, nil
}

func (s *Server) RemoveOrgIDP(ctx context.Context, req *mgmt_pb.RemoveOrgIDPRequest) (*mgmt_pb.RemoveOrgIDPResponse, error) {
	idp, err := s.query.IDPByIDAndResourceOwner(ctx, true, req.IdpId, authz.GetCtxData(ctx).OrgID, true)
	if err != nil {
		return nil, err
	}
	idpQuery, err := query.NewIDPUserLinkIDPIDSearchQuery(req.IdpId)
	if err != nil {
		return nil, err
	}
	userLinks, err := s.query.IDPUserLinks(ctx, &query.IDPUserLinksSearchQuery{
		Queries: []query.SearchQuery{idpQuery},
	}, true)
	if err != nil {
		return nil, err
	}
	_, err = s.command.RemoveIDPConfig(ctx, req.IdpId, authz.GetCtxData(ctx).OrgID, idp != nil, userLinksToDomain(userLinks.Links)...)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveOrgIDPResponse{}, nil
}

func (s *Server) UpdateOrgIDP(ctx context.Context, req *mgmt_pb.UpdateOrgIDPRequest) (*mgmt_pb.UpdateOrgIDPResponse, error) {
	config, err := s.command.ChangeIDPConfig(ctx, updateIDPToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateOrgIDPResponse{
		Details: object_pb.ChangeToDetailsPb(
			config.Sequence,
			config.ChangeDate,
			config.ResourceOwner,
		),
	}, nil
}

func (s *Server) UpdateOrgIDPOIDCConfig(ctx context.Context, req *mgmt_pb.UpdateOrgIDPOIDCConfigRequest) (*mgmt_pb.UpdateOrgIDPOIDCConfigResponse, error) {
	config, err := s.command.ChangeIDPOIDCConfig(ctx, updateOIDCConfigToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateOrgIDPOIDCConfigResponse{
		Details: object_pb.ChangeToDetailsPb(
			config.Sequence,
			config.ChangeDate,
			config.ResourceOwner,
		),
	}, nil
}

func (s *Server) UpdateOrgIDPJWTConfig(ctx context.Context, req *mgmt_pb.UpdateOrgIDPJWTConfigRequest) (*mgmt_pb.UpdateOrgIDPJWTConfigResponse, error) {
	config, err := s.command.ChangeIDPJWTConfig(ctx, updateJWTConfigToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateOrgIDPJWTConfigResponse{
		Details: object_pb.ChangeToDetailsPb(
			config.Sequence,
			config.ChangeDate,
			config.ResourceOwner,
		),
	}, nil
}

func (s *Server) GetProviderByID(ctx context.Context, req *mgmt_pb.GetProviderByIDRequest) (*mgmt_pb.GetProviderByIDResponse, error) {
	orgIDQuery, err := query.NewIDPTemplateResourceOwnerSearchQuery(authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	idp, err := s.query.IDPTemplateByID(ctx, true, req.Id, false, orgIDQuery)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetProviderByIDResponse{Idp: idp_grpc.ProviderToPb(idp)}, nil
}

func (s *Server) ListProviders(ctx context.Context, req *mgmt_pb.ListProvidersRequest) (*mgmt_pb.ListProvidersResponse, error) {
	queries, err := listProvidersToQuery(ctx, req)
	if err != nil {
		return nil, err
	}
	resp, err := s.query.IDPTemplates(ctx, queries, false)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListProvidersResponse{
		Result:  idp_grpc.ProvidersToPb(resp.Templates),
		Details: object_pb.ToListDetails(resp.Count, resp.Sequence, resp.Timestamp),
	}, nil
}

func (s *Server) AddGenericOAuthProvider(ctx context.Context, req *mgmt_pb.AddGenericOAuthProviderRequest) (*mgmt_pb.AddGenericOAuthProviderResponse, error) {
	id, details, err := s.command.AddOrgGenericOAuthProvider(ctx, authz.GetCtxData(ctx).OrgID, addGenericOAuthProviderToCommand(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddGenericOAuthProviderResponse{
		Id:      id,
		Details: object_pb.DomainToAddDetailsPb(details),
	}, nil
}

func (s *Server) UpdateGenericOAuthProvider(ctx context.Context, req *mgmt_pb.UpdateGenericOAuthProviderRequest) (*mgmt_pb.UpdateGenericOAuthProviderResponse, error) {
	details, err := s.command.UpdateOrgGenericOAuthProvider(ctx, authz.GetCtxData(ctx).OrgID, req.Id, updateGenericOAuthProviderToCommand(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateGenericOAuthProviderResponse{
		Details: object_pb.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) AddGenericOIDCProvider(ctx context.Context, req *mgmt_pb.AddGenericOIDCProviderRequest) (*mgmt_pb.AddGenericOIDCProviderResponse, error) {
	id, details, err := s.command.AddOrgGenericOIDCProvider(ctx, authz.GetCtxData(ctx).OrgID, addGenericOIDCProviderToCommand(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddGenericOIDCProviderResponse{
		Id:      id,
		Details: object_pb.DomainToAddDetailsPb(details),
	}, nil
}

func (s *Server) UpdateGenericOIDCProvider(ctx context.Context, req *mgmt_pb.UpdateGenericOIDCProviderRequest) (*mgmt_pb.UpdateGenericOIDCProviderResponse, error) {
	details, err := s.command.UpdateOrgGenericOIDCProvider(ctx, authz.GetCtxData(ctx).OrgID, req.Id, updateGenericOIDCProviderToCommand(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateGenericOIDCProviderResponse{
		Details: object_pb.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) AddJWTProvider(ctx context.Context, req *mgmt_pb.AddJWTProviderRequest) (*mgmt_pb.AddJWTProviderResponse, error) {
	id, details, err := s.command.AddOrgJWTProvider(ctx, authz.GetCtxData(ctx).OrgID, addJWTProviderToCommand(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddJWTProviderResponse{
		Id:      id,
		Details: object_pb.DomainToAddDetailsPb(details),
	}, nil
}

func (s *Server) UpdateJWTProvider(ctx context.Context, req *mgmt_pb.UpdateJWTProviderRequest) (*mgmt_pb.UpdateJWTProviderResponse, error) {
	details, err := s.command.UpdateOrgJWTProvider(ctx, authz.GetCtxData(ctx).OrgID, req.Id, updateJWTProviderToCommand(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateJWTProviderResponse{
		Details: object_pb.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) AddGitHubProvider(ctx context.Context, req *mgmt_pb.AddGitHubProviderRequest) (*mgmt_pb.AddGitHubProviderResponse, error) {
	id, details, err := s.command.AddOrgGitHubProvider(ctx, authz.GetCtxData(ctx).OrgID, addGitHubProviderToCommand(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddGitHubProviderResponse{
		Id:      id,
		Details: object_pb.DomainToAddDetailsPb(details),
	}, nil
}

func (s *Server) UpdateGitHubProvider(ctx context.Context, req *mgmt_pb.UpdateGitHubProviderRequest) (*mgmt_pb.UpdateGitHubProviderResponse, error) {
	details, err := s.command.UpdateOrgGitHubProvider(ctx, authz.GetCtxData(ctx).OrgID, req.Id, updateGitHubProviderToCommand(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateGitHubProviderResponse{
		Details: object_pb.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) AddGitHubEnterpriseServerProvider(ctx context.Context, req *mgmt_pb.AddGitHubEnterpriseServerProviderRequest) (*mgmt_pb.AddGitHubEnterpriseServerProviderResponse, error) {
	id, details, err := s.command.AddOrgGitHubEnterpriseProvider(ctx, authz.GetCtxData(ctx).OrgID, addGitHubEnterpriseProviderToCommand(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddGitHubEnterpriseServerProviderResponse{
		Id:      id,
		Details: object_pb.DomainToAddDetailsPb(details),
	}, nil
}

func (s *Server) UpdateGitHubEnterpriseServerProvider(ctx context.Context, req *mgmt_pb.UpdateGitHubEnterpriseServerProviderRequest) (*mgmt_pb.UpdateGitHubEnterpriseServerProviderResponse, error) {
	details, err := s.command.UpdateOrgGitHubEnterpriseProvider(ctx, authz.GetCtxData(ctx).OrgID, req.Id, updateGitHubEnterpriseProviderToCommand(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateGitHubEnterpriseServerProviderResponse{
		Details: object_pb.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) AddGoogleProvider(ctx context.Context, req *mgmt_pb.AddGoogleProviderRequest) (*mgmt_pb.AddGoogleProviderResponse, error) {
	id, details, err := s.command.AddOrgGoogleProvider(ctx, authz.GetCtxData(ctx).OrgID, addGoogleProviderToCommand(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddGoogleProviderResponse{
		Id:      id,
		Details: object_pb.DomainToAddDetailsPb(details),
	}, nil
}

func (s *Server) UpdateGoogleProvider(ctx context.Context, req *mgmt_pb.UpdateGoogleProviderRequest) (*mgmt_pb.UpdateGoogleProviderResponse, error) {
	details, err := s.command.UpdateOrgGoogleProvider(ctx, authz.GetCtxData(ctx).OrgID, req.Id, updateGoogleProviderToCommand(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateGoogleProviderResponse{
		Details: object_pb.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) AddLDAPProvider(ctx context.Context, req *mgmt_pb.AddLDAPProviderRequest) (*mgmt_pb.AddLDAPProviderResponse, error) {
	id, details, err := s.command.AddOrgLDAPProvider(ctx, authz.GetCtxData(ctx).OrgID, addLDAPProviderToCommand(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddLDAPProviderResponse{
		Id:      id,
		Details: object_pb.DomainToAddDetailsPb(details),
	}, nil
}

func (s *Server) UpdateLDAPProvider(ctx context.Context, req *mgmt_pb.UpdateLDAPProviderRequest) (*mgmt_pb.UpdateLDAPProviderResponse, error) {
	details, err := s.command.UpdateOrgLDAPProvider(ctx, authz.GetCtxData(ctx).OrgID, req.Id, updateLDAPProviderToCommand(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateLDAPProviderResponse{
		Details: object_pb.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) DeleteProvider(ctx context.Context, req *mgmt_pb.DeleteProviderRequest) (*mgmt_pb.DeleteProviderResponse, error) {
	details, err := s.command.DeleteOrgProvider(ctx, authz.GetCtxData(ctx).OrgID, req.Id)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.DeleteProviderResponse{
		Details: object_pb.DomainToChangeDetailsPb(details),
	}, nil
}
