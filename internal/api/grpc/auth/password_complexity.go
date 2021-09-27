package auth

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	policy_grpc "github.com/caos/zitadel/internal/api/grpc/policy"
	auth_pb "github.com/caos/zitadel/pkg/grpc/auth"
)

func (s *Server) GetMyPasswordComplexityPolicy(ctx context.Context, _ *auth_pb.GetMyPasswordComplexityPolicyRequest) (*auth_pb.GetMyPasswordComplexityPolicyResponse, error) {
	policy, err := s.query.MyPasswordComplexityPolicy(ctx, authz.GetCtxData(ctx).UserID)
	if err != nil {
		return nil, err
	}
	return &auth_pb.GetMyPasswordComplexityPolicyResponse{Policy: policy_grpc.ModelPasswordComplexityPolicyToPb(policy)}, nil
}
