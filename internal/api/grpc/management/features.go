package management

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	features_grpc "github.com/caos/zitadel/internal/api/grpc/features"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) GetFeatures(ctx context.Context, req *mgmt_pb.GetFeaturesRequest) (*mgmt_pb.GetFeaturesResponse, error) {
	features, err := s.features.GetOrgFeatures(ctx, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetFeaturesResponse{
		Features: features_grpc.FeaturesFromModel(features),
	}, nil
}
