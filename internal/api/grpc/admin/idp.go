package admin

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	idp_grpc "github.com/zitadel/zitadel/internal/api/grpc/idp"
	object_pb "github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/query"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
)

func (s *Server) GetIDPByID(ctx context.Context, req *admin_pb.GetIDPByIDRequest) (*admin_pb.GetIDPByIDResponse, error) {
	idp, err := s.query.IDPByIDAndResourceOwner(ctx, true, req.Id, authz.GetInstance(ctx).InstanceID(), false)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetIDPByIDResponse{Idp: idp_grpc.IDPViewToPb(idp)}, nil
}

func (s *Server) ListIDPs(ctx context.Context, req *admin_pb.ListIDPsRequest) (*admin_pb.ListIDPsResponse, error) {
	queries, err := listIDPsToModel(authz.GetInstance(ctx).InstanceID(), req)
	if err != nil {
		return nil, err
	}
	resp, err := s.query.IDPs(ctx, queries, false)
	if err != nil {
		return nil, err
	}
	return &admin_pb.ListIDPsResponse{
		Result:  idp_grpc.IDPViewsToPb(resp.IDPs),
		Details: object_pb.ToListDetails(resp.Count, resp.Sequence, resp.Timestamp),
	}, nil
}

func (s *Server) AddOIDCIDP(ctx context.Context, req *admin_pb.AddOIDCIDPRequest) (*admin_pb.AddOIDCIDPResponse, error) {
	config, err := s.command.AddDefaultIDPConfig(ctx, addOIDCIDPRequestToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.AddOIDCIDPResponse{
		IdpId: config.IDPConfigID,
		Details: object_pb.AddToDetailsPb(
			config.Sequence,
			config.ChangeDate,
			config.ResourceOwner,
		),
	}, nil
}

func (s *Server) AddJWTIDP(ctx context.Context, req *admin_pb.AddJWTIDPRequest) (*admin_pb.AddJWTIDPResponse, error) {
	config, err := s.command.AddDefaultIDPConfig(ctx, addJWTIDPRequestToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.AddJWTIDPResponse{
		IdpId: config.IDPConfigID,
		Details: object_pb.AddToDetailsPb(
			config.Sequence,
			config.ChangeDate,
			config.ResourceOwner,
		),
	}, nil
}

func (s *Server) UpdateIDP(ctx context.Context, req *admin_pb.UpdateIDPRequest) (*admin_pb.UpdateIDPResponse, error) {
	config, err := s.command.ChangeDefaultIDPConfig(ctx, updateIDPToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdateIDPResponse{
		Details: object_pb.ChangeToDetailsPb(
			config.Sequence,
			config.ChangeDate,
			config.ResourceOwner,
		),
	}, nil
}

func (s *Server) DeactivateIDP(ctx context.Context, req *admin_pb.DeactivateIDPRequest) (*admin_pb.DeactivateIDPResponse, error) {
	objectDetails, err := s.command.DeactivateDefaultIDPConfig(ctx, req.IdpId)
	if err != nil {
		return nil, err
	}
	return &admin_pb.DeactivateIDPResponse{Details: object_pb.DomainToChangeDetailsPb(objectDetails)}, nil
}

func (s *Server) ReactivateIDP(ctx context.Context, req *admin_pb.ReactivateIDPRequest) (*admin_pb.ReactivateIDPResponse, error) {
	objectDetails, err := s.command.ReactivateDefaultIDPConfig(ctx, req.IdpId)
	if err != nil {
		return nil, err
	}
	return &admin_pb.ReactivateIDPResponse{Details: object_pb.DomainToChangeDetailsPb(objectDetails)}, nil
}

func (s *Server) RemoveIDP(ctx context.Context, req *admin_pb.RemoveIDPRequest) (*admin_pb.RemoveIDPResponse, error) {
	providerQuery, err := query.NewIDPIDSearchQuery(req.IdpId)
	if err != nil {
		return nil, err
	}
	idps, err := s.query.IDPs(ctx, &query.IDPSearchQueries{
		Queries: []query.SearchQuery{providerQuery},
	}, true)
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

	objectDetails, err := s.command.RemoveDefaultIDPConfig(ctx, req.IdpId, idpsToDomain(idps.IDPs), idpUserLinksToDomain(userLinks.Links)...)
	if err != nil {
		return nil, err
	}
	return &admin_pb.RemoveIDPResponse{Details: object_pb.DomainToChangeDetailsPb(objectDetails)}, nil
}

func (s *Server) UpdateIDPOIDCConfig(ctx context.Context, req *admin_pb.UpdateIDPOIDCConfigRequest) (*admin_pb.UpdateIDPOIDCConfigResponse, error) {
	config, err := s.command.ChangeDefaultIDPOIDCConfig(ctx, updateOIDCConfigToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdateIDPOIDCConfigResponse{
		Details: object_pb.ChangeToDetailsPb(
			config.Sequence,
			config.ChangeDate,
			config.ResourceOwner,
		),
	}, nil
}

func (s *Server) UpdateIDPJWTConfig(ctx context.Context, req *admin_pb.UpdateIDPJWTConfigRequest) (*admin_pb.UpdateIDPJWTConfigResponse, error) {
	config, err := s.command.ChangeDefaultIDPJWTConfig(ctx, updateJWTConfigToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdateIDPJWTConfigResponse{
		Details: object_pb.ChangeToDetailsPb(
			config.Sequence,
			config.ChangeDate,
			config.ResourceOwner,
		),
	}, nil
}

func (s *Server) GetProviderByID(ctx context.Context, req *admin_pb.GetProviderByIDRequest) (*admin_pb.GetProviderByIDResponse, error) {
	idp, err := s.query.IDPTemplateByIDAndResourceOwner(ctx, true, req.Id, authz.GetInstance(ctx).InstanceID(), false)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetProviderByIDResponse{Idp: idp_grpc.ProviderToPb(idp)}, nil
}

func (s *Server) ListProviders(ctx context.Context, req *admin_pb.ListProvidersRequest) (*admin_pb.ListProvidersResponse, error) {
	queries, err := listProvidersToQuery(authz.GetInstance(ctx).InstanceID(), req)
	if err != nil {
		return nil, err
	}
	resp, err := s.query.IDPTemplates(ctx, queries, false)
	if err != nil {
		return nil, err
	}
	return &admin_pb.ListProvidersResponse{
		Result:  idp_grpc.ProvidersToPb(resp.Templates),
		Details: object_pb.ToListDetails(resp.Count, resp.Sequence, resp.Timestamp),
	}, nil
}

func (s *Server) AddGenericOAuthProvider(ctx context.Context, req *admin_pb.AddGenericOAuthProviderRequest) (*admin_pb.AddGenericOAuthProviderResponse, error) {
	id, details, err := s.command.AddOrgGenericOAuthProvider(ctx, authz.GetCtxData(ctx).OrgID, addGenericOAuthProviderToCommand(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.AddGenericOAuthProviderResponse{
		Id:      id,
		Details: object_pb.DomainToAddDetailsPb(details),
	}, nil
}

func (s *Server) UpdateGenericOAuthProvider(ctx context.Context, req *admin_pb.UpdateGenericOAuthProviderRequest) (*admin_pb.UpdateGenericOAuthProviderResponse, error) {
	details, err := s.command.UpdateOrgGenericOAuthProvider(ctx, authz.GetCtxData(ctx).OrgID, req.Id, updateGenericOAuthProviderToCommand(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdateGenericOAuthProviderResponse{
		Details: object_pb.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) AddGitHubProvider(ctx context.Context, req *admin_pb.AddGitHubProviderRequest) (*admin_pb.AddGitHubProviderResponse, error) {
	id, details, err := s.command.AddOrgGitHubProvider(ctx, authz.GetCtxData(ctx).OrgID, addGitHubProviderToCommand(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.AddGitHubProviderResponse{
		Id:      id,
		Details: object_pb.DomainToAddDetailsPb(details),
	}, nil
}

func (s *Server) UpdateGitHubProvider(ctx context.Context, req *admin_pb.UpdateGitHubProviderRequest) (*admin_pb.UpdateGitHubProviderResponse, error) {
	details, err := s.command.UpdateOrgGitHubProvider(ctx, authz.GetCtxData(ctx).OrgID, req.Id, updateGitHubProviderToCommand(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdateGitHubProviderResponse{
		Details: object_pb.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) AddGitHubEnterpriseServerProvider(ctx context.Context, req *admin_pb.AddGitHubEnterpriseServerProviderRequest) (*admin_pb.AddGitHubEnterpriseServerProviderResponse, error) {
	id, details, err := s.command.AddOrgGitHubEnterpriseProvider(ctx, authz.GetCtxData(ctx).OrgID, addGitHubEnterpriseProviderToCommand(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.AddGitHubEnterpriseServerProviderResponse{
		Id:      id,
		Details: object_pb.DomainToAddDetailsPb(details),
	}, nil
}

func (s *Server) UpdateGitHubEnterpriseServerProvider(ctx context.Context, req *admin_pb.UpdateGitHubEnterpriseServerProviderRequest) (*admin_pb.UpdateGitHubEnterpriseServerProviderResponse, error) {
	details, err := s.command.UpdateOrgGitHubEnterpriseProvider(ctx, authz.GetCtxData(ctx).OrgID, req.Id, updateGitHubEnterpriseProviderToCommand(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdateGitHubEnterpriseServerProviderResponse{
		Details: object_pb.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) AddGoogleProvider(ctx context.Context, req *admin_pb.AddGoogleProviderRequest) (*admin_pb.AddGoogleProviderResponse, error) {
	id, details, err := s.command.AddOrgGoogleProvider(ctx, authz.GetCtxData(ctx).OrgID, addGoogleProviderToCommand(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.AddGoogleProviderResponse{
		Id:      id,
		Details: object_pb.DomainToAddDetailsPb(details),
	}, nil
}

func (s *Server) UpdateGoogleProvider(ctx context.Context, req *admin_pb.UpdateGoogleProviderRequest) (*admin_pb.UpdateGoogleProviderResponse, error) {
	details, err := s.command.UpdateOrgGoogleProvider(ctx, authz.GetCtxData(ctx).OrgID, req.Id, updateGoogleProviderToCommand(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdateGoogleProviderResponse{
		Details: object_pb.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) AddLDAPProvider(ctx context.Context, req *admin_pb.AddLDAPProviderRequest) (*admin_pb.AddLDAPProviderResponse, error) {
	id, details, err := s.command.AddInstanceLDAPProvider(ctx, addLDAPProviderToCommand(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.AddLDAPProviderResponse{
		Id:      id,
		Details: object_pb.DomainToAddDetailsPb(details),
	}, nil
}

func (s *Server) UpdateLDAPProvider(ctx context.Context, req *admin_pb.UpdateLDAPProviderRequest) (*admin_pb.UpdateLDAPProviderResponse, error) {
	details, err := s.command.UpdateInstanceLDAPProvider(ctx, req.Id, updateLDAPProviderToCommand(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdateLDAPProviderResponse{
		Details: object_pb.DomainToChangeDetailsPb(details),
	}, nil
}

func (s *Server) DeleteProvider(ctx context.Context, req *admin_pb.DeleteProviderRequest) (*admin_pb.DeleteProviderResponse, error) {
	details, err := s.command.DeleteInstanceProvider(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &admin_pb.DeleteProviderResponse{
		Details: object_pb.DomainToChangeDetailsPb(details),
	}, nil
}
