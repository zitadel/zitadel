package auth

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	auth_pb "github.com/caos/zitadel/pkg/grpc/auth"
	policy_grpc "github.com/caos/zitadel/v2/internal/api/grpc/policy"
)

func (s *Server) GetMyPasswordComplexityPolicy(ctx context.Context, _ *auth_pb.GetMyPasswordComplexityPolicyRequest) (*auth_pb.GetMyPasswordComplexityPolicyResponse, error) {
	policy, err := s.query.PasswordComplexityPolicyByOrg(ctx, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &auth_pb.GetMyPasswordComplexityPolicyResponse{Policy: policy_grpc.ModelPasswordComplexityPolicyToPb(policy)}, nil
}
