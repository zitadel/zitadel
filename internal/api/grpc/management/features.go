package management

import (
	"context"

	features_grpc "github.com/caos/zitadel/internal/api/grpc/features"
	"github.com/caos/zitadel/internal/domain"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) GetOrgFeatures(ctx context.Context, req *mgmt_pb.GetFeaturesRequest) (*mgmt_pb.GetFeaturesResponse, error) {
	var features *domain.Features
	var err error
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetFeaturesResponse{
		Features: features_grpc.FeaturesToPb(features),
	}, nil
}
