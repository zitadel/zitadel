package management

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
	config, err := s.command.AddDefaultIDPConfig(ctx, addOIDCIDPRequestToDomain(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddOrgOIDCIDPResponse{
		IdpId: config.AggregateID,
		Details: object_pb.AddToDetailsPb(
			config.Sequence,
			config.ChangeDate,
			config.ResourceOwner,
		),
	}, nil
}
func (s *Server) DeactivateOrgIDP(ctx context.Context, req *mgmt_pb.DeactivateOrgIDPRequest) (*mgmt_pb.DeactivateOrgIDPResponse, error) {
	objectDetails, err := s.command.DeactivateDefaultIDPConfig(ctx, req.IdpId)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.DeactivateOrgIDPResponse{Details: object_pb.DomainToChangeDetailsPb(objectDetails)}, nil
}
func (s *Server) ReactivateOrgIDP(ctx context.Context, req *mgmt_pb.ReactivateOrgIDPRequest) (*mgmt_pb.ReactivateOrgIDPResponse, error) {
	objectDetails, err := s.command.ReactivateDefaultIDPConfig(ctx, req.IdpId)
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
	_, err = s.command.RemoveDefaultIDPConfig(ctx, req.IdpId, idpProviderViewsToDomain(idpProviders), externalIDPViewsToDomain(externalIDPs)...)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.RemoveOrgIDPResponse{}, nil
}
func (s *Server) UpdateOrgIDP(ctx context.Context, req *mgmt_pb.UpdateOrgIDPRequest) (*mgmt_pb.UpdateOrgIDPResponse, error) {
	config, err := s.command.ChangeDefaultIDPConfig(ctx, updateIDPToDomain(req))
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
	return nil, status.Errorf(codes.Unimplemented, "method UpdateOrgIDPOIDCConfig not implemented")
}
