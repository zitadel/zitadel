package admin

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	idp_grpc "github.com/caos/zitadel/internal/api/grpc/idp"
	object_pb "github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/query"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func (s *Server) GetIDPByID(ctx context.Context, req *admin_pb.GetIDPByIDRequest) (*admin_pb.GetIDPByIDResponse, error) {
	idp, err := s.query.IDPByIDAndResourceOwner(ctx, req.Id, authz.GetInstance(ctx).InstanceID())
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
	resp, err := s.query.IDPs(ctx, queries)
	if err != nil {
		return nil, err
	}
	return &admin_pb.ListIDPsResponse{
		Result:  idp_grpc.IDPViewsToPb(resp.IDPs),
		Details: object_pb.ToListDetails(resp.Count, resp.Sequence, resp.Timestamp),
	}, nil
}

func (s *Server) AddOIDCIDP(ctx context.Context, req *admin_pb.AddOIDCIDPRequest) (*admin_pb.AddOIDCIDPResponse, error) {
	config, err := s.command.AddDefaultIDPConfig(ctx, authz.GetInstance(ctx).InstanceID(), addOIDCIDPRequestToDomain(req))
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
	config, err := s.command.AddDefaultIDPConfig(ctx, authz.GetInstance(ctx).InstanceID(), addJWTIDPRequestToDomain(req))
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
	config, err := s.command.ChangeDefaultIDPConfig(ctx, authz.GetInstance(ctx).InstanceID(), updateIDPToDomain(req))
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
	objectDetails, err := s.command.DeactivateDefaultIDPConfig(ctx, authz.GetInstance(ctx).InstanceID(), req.IdpId)
	if err != nil {
		return nil, err
	}
	return &admin_pb.DeactivateIDPResponse{Details: object_pb.DomainToChangeDetailsPb(objectDetails)}, nil
}

func (s *Server) ReactivateIDP(ctx context.Context, req *admin_pb.ReactivateIDPRequest) (*admin_pb.ReactivateIDPResponse, error) {
	objectDetails, err := s.command.ReactivateDefaultIDPConfig(ctx, authz.GetInstance(ctx).InstanceID(), req.IdpId)
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
	})
	if err != nil {
		return nil, err
	}

	idpQuery, err := query.NewIDPUserLinkIDPIDSearchQuery(req.IdpId)
	if err != nil {
		return nil, err
	}
	userLinks, err := s.query.IDPUserLinks(ctx, &query.IDPUserLinksSearchQuery{
		Queries: []query.SearchQuery{idpQuery},
	})
	if err != nil {
		return nil, err
	}

	objectDetails, err := s.command.RemoveDefaultIDPConfig(ctx, authz.GetInstance(ctx).InstanceID(), req.IdpId, idpsToDomain(idps.IDPs), idpUserLinksToDomain(userLinks.Links)...)
	if err != nil {
		return nil, err
	}
	return &admin_pb.RemoveIDPResponse{Details: object_pb.DomainToChangeDetailsPb(objectDetails)}, nil
}

func (s *Server) UpdateIDPOIDCConfig(ctx context.Context, req *admin_pb.UpdateIDPOIDCConfigRequest) (*admin_pb.UpdateIDPOIDCConfigResponse, error) {
	config, err := s.command.ChangeDefaultIDPOIDCConfig(ctx, authz.GetInstance(ctx).InstanceID(), updateOIDCConfigToDomain(req))
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
	config, err := s.command.ChangeDefaultIDPJWTConfig(ctx, authz.GetInstance(ctx).InstanceID(), updateJWTConfigToDomain(req))
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
