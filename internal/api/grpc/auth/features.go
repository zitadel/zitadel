package auth

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	auth_pb "github.com/caos/zitadel/pkg/grpc/auth"
)

func (s *Server) ListMyZitadelFeatures(ctx context.Context, _ *auth_pb.ListMyZitadelFeaturesRequest) (*auth_pb.ListMyZitadelFeaturesResponse, error) {
	features, err := s.repo.GetOrgFeatures(ctx, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &auth_pb.ListMyZitadelFeaturesResponse{
		Result: features.FeatureList(),
	}, nil
}
