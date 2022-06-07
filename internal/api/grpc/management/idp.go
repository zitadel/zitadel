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
	idp, err := s.query.IDPByIDAndResourceOwner(ctx, true, req.Id, authz.GetCtxData(ctx).OrgID)
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
	resp, err := s.query.IDPs(ctx, queries)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListOrgIDPsResponse{
		Result:  idp_grpc.IDPViewsToPb(resp.IDPs),
		Details: object_pb.ToListDetails(resp.Count, resp.Sequence, resp.Timestamp),
	}, nil
}

func (s *Server) AddOrgOIDCIDP(ctx context.Context, req *mgmt_pb.AddOrgOIDCIDPRequest) (*mgmt_pb.AddOrgOIDCIDPResponse, error) {
	config, err := s.command.AddIDPConfig(ctx, addOIDCIDPRequestToDomain(req), authz.GetCtxData(ctx).OrgID)
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
	config, err := s.command.AddIDPConfig(ctx, addJWTIDPRequestToDomain(req), authz.GetCtxData(ctx).OrgID)
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
	idp, err := s.query.IDPByIDAndResourceOwner(ctx, false, req.IdpId, authz.GetCtxData(ctx).OrgID)
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
