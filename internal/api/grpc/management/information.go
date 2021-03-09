package management

import (
	"context"

	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) Healthz(context.Context, *mgmt_pb.HealthzRequest) (*mgmt_pb.HealthzResponse, error) {
	return &mgmt_pb.HealthzResponse{}, nil
}

func (s *Server) GetOIDCInformation(ctx context.Context, req *mgmt_pb.GetOIDCInformationRequest) (*mgmt_pb.GetOIDCInformationResponse, error) {
	return &mgmt_pb.GetOIDCInformationResponse{
		Issuer:            s.systemDefaults.ZitadelDocs.Issuer,
		DiscoveryEndpoint: s.systemDefaults.ZitadelDocs.DiscoveryEndpoint,
	}, nil
}
