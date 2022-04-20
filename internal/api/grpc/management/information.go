package management

import (
	"context"

	"github.com/caos/oidc/pkg/oidc"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/http"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) Healthz(context.Context, *mgmt_pb.HealthzRequest) (*mgmt_pb.HealthzResponse, error) {
	return &mgmt_pb.HealthzResponse{}, nil
}

func (s *Server) GetOIDCInformation(ctx context.Context, _ *mgmt_pb.GetOIDCInformationRequest) (*mgmt_pb.GetOIDCInformationResponse, error) {
	issuer := http.BuildHTTP(authz.GetInstance(ctx).RequestedDomain(), s.externalPort, s.externalSecure) + s.issuerPath
	return &mgmt_pb.GetOIDCInformationResponse{
		Issuer:            issuer,
		DiscoveryEndpoint: issuer + oidc.DiscoveryEndpoint,
	}, nil
}
