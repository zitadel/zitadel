package management

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	idp_grpc "github.com/caos/zitadel/internal/api/grpc/idp"
	object_pb "github.com/caos/zitadel/internal/api/grpc/object"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) GetOrgIDPByID(ctx context.Context, req *mgmt_pb.GetOrgIDPByIDRequest) (*mgmt_pb.GetOrgIDPByIDResponse, error) {
	idp, err := s.org.IDPConfigByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetOrgIDPByIDResponse{Idp: idp_grpc.ModelIDPViewToPb(idp)}, nil
}
func (s *Server) ListOrgIDPs(ctx context.Context, req *mgmt_pb.ListOrgIDPsRequest) (*mgmt_pb.ListOrgIDPsResponse, error) {
	resp, err := s.org.SearchIDPConfigs(ctx, listIDPsToModel(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ListOrgIDPsResponse{
		Result:  idp_grpc.IDPViewsToPb(resp.Result),
		Details: object_pb.ToListDetails(resp.TotalResult, resp.Sequence, resp.Timestamp),
	}, nil
}
func (s *Server) AddOrgOIDCIDP(ctx context.Context, req *mgmt_pb.AddOrgOIDCIDPRequest) (*mgmt_pb.AddOrgOIDCIDPResponse, error) {
	id, details, err := s.command.AddIDPConfig(ctx, addOIDCIDPRequestToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddOrgOIDCIDPResponse{
		IdpId:   id,
		Details: object_pb.DomainToChangeDetailsPb(details),
	}, nil
}
func (s *Server) AddOrgAuthConnectorIDP(ctx context.Context, req *mgmt_pb.AddOrgAuthConnectorIDPRequest) (*mgmt_pb.AddOrgAuthConnectorIDPResponse, error) {
	id, details, err := s.command.AddIDPConfig(ctx, addAuthConnectorIDPRequestToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddOrgAuthConnectorIDPResponse{
		IdpId:   id,
		Details: object_pb.DomainToChangeDetailsPb(details),
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
	idpProviders, err := s.org.GetIDPProvidersByIDPConfigID(ctx, authz.GetCtxData(ctx).OrgID, req.IdpId)
	if err != nil {
		return nil, err
	}
	externalIDPs, err := s.user.ExternalIDPsByIDPConfigID(ctx, req.IdpId)
	if err != nil {
		return nil, err
	}
	_, err = s.command.RemoveIDPConfig(ctx, req.IdpId, authz.GetCtxData(ctx).OrgID, len(idpProviders) > 0, externalIDPViewsToDomain(externalIDPs)...)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveOrgIDPResponse{}, nil
}
func (s *Server) UpdateOrgIDP(ctx context.Context, req *mgmt_pb.UpdateOrgIDPRequest) (*mgmt_pb.UpdateOrgIDPResponse, error) {
	details, err := s.command.ChangeIDPConfig(ctx, updateIDPToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateOrgIDPResponse{
		Details: object_pb.DomainToChangeDetailsPb(details),
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

func (s *Server) UpdateOrgIDPAuthConnectorConfig(ctx context.Context, req *mgmt_pb.UpdateOrgIDPAuthConnectorConfigRequest) (*mgmt_pb.UpdateOrgIDPAuthConnectorConfigResponse, error) {
	details, err := s.command.ChangeIDPAuthConnectorConfig(ctx, updateAuthConnectorConfigToDomain(req), authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateOrgIDPAuthConnectorConfigResponse{
		Details: object_pb.DomainToChangeDetailsPb(details),
	}, nil
}
