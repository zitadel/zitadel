package auth

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	policy_grpc "github.com/zitadel/zitadel/internal/api/grpc/policy"
	auth_pb "github.com/zitadel/zitadel/pkg/grpc/auth"
)

func (s *Server) GetMyPasswordComplexityPolicy(ctx context.Context, _ *auth_pb.GetMyPasswordComplexityPolicyRequest) (*auth_pb.GetMyPasswordComplexityPolicyResponse, error) {
	policy, err := s.query.PasswordComplexityPolicyByOrg(ctx, true, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &auth_pb.GetMyPasswordComplexityPolicyResponse{Policy: policy_grpc.ModelPasswordComplexityPolicyToPb(policy)}, nil
}
