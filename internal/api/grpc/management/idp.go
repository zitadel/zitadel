package management

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) GetOrgIDPByID(ctx context.Context, req *mgmt_pb.GetOrgIDPByIDRequest) (*mgmt_pb.GetOrgIDPByIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOrgIDPByID not implemented")
}
func (s *Server) ListOrgIDPs(ctx context.Context, req *mgmt_pb.ListOrgIDPsRequest) (*mgmt_pb.ListOrgIDPsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListOrgIDPs not implemented")
}
func (s *Server) AddOrgOIDCIDP(ctx context.Context, req *mgmt_pb.AddOrgOIDCIDPRequest) (*mgmt_pb.AddOrgOIDCIDPResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddOrgOIDCIDP not implemented")
}
func (s *Server) DeactivateOrgIDP(ctx context.Context, req *mgmt_pb.DeactivateOrgIDPRequest) (*mgmt_pb.DeactivateOrgIDPResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeactivateOrgIDP not implemented")
}
func (s *Server) ReactivateOrgIDP(ctx context.Context, req *mgmt_pb.ReactivateOrgIDPRequest) (*mgmt_pb.ReactivateOrgIDPResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReactivateOrgIDP not implemented")
}
func (s *Server) RemoveOrgIDP(ctx context.Context, req *mgmt_pb.RemoveOrgIDPRequest) (*mgmt_pb.RemoveOrgIDPResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveOrgIDP not implemented")
}
func (s *Server) UpdateOrgIDP(ctx context.Context, req *mgmt_pb.UpdateOrgIDPRequest) (*mgmt_pb.UpdateOrgIDPResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateOrgIDP not implemented")
}
func (s *Server) UpdateOrgIDPOIDCConfig(ctx context.Context, req *mgmt_pb.UpdateOrgIDPOIDCConfigRequest) (*mgmt_pb.UpdateOrgIDPOIDCConfigResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateOrgIDPOIDCConfig not implemented")
}
