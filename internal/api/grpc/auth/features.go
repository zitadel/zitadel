package auth

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	auth_pb "github.com/zitadel/zitadel/pkg/grpc/auth"
)

func (s *Server) ListMyZitadelFeatures(ctx context.Context, _ *auth_pb.ListMyZitadelFeaturesRequest) (*auth_pb.ListMyZitadelFeaturesResponse, error) {
	features, err := s.query.FeaturesByOrgID(ctx, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &auth_pb.ListMyZitadelFeaturesResponse{
		Result: features.EnabledFeatureTypes(),
	}, nil
}
