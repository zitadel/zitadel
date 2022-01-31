package management

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
	features_grpc "github.com/caos/zitadel/v2/internal/api/grpc/features"
)

func (s *Server) GetFeatures(ctx context.Context, req *mgmt_pb.GetFeaturesRequest) (*mgmt_pb.GetFeaturesResponse, error) {
	features, err := s.query.FeaturesByOrgID(ctx, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetFeaturesResponse{
		Features: features_grpc.ModelFeaturesToPb(features),
	}, nil
}
