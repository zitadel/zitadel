package admin

import (
	"context"

	"github.com/caos/zitadel/internal/api/grpc/object"
	policy_grpc "github.com/caos/zitadel/internal/api/grpc/policy"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func (s *Server) GetPasswordComplexityPolicy(ctx context.Context, _ *admin_pb.GetPasswordComplexityPolicyRequest) (*admin_pb.GetPasswordComplexityPolicyResponse, error) {
	policy, err := s.iam.GetDefaultPasswordComplexityPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetPasswordComplexityPolicyResponse{Policy: policy_grpc.ModelPasswordComplexityPolicyToPb(policy)}, nil
}

func (s *Server) UpdatePasswordComplexityPolicy(ctx context.Context, req *admin_pb.UpdatePasswordComplexityPolicyRequest) (*admin_pb.UpdatePasswordComplexityPolicyResponse, error) {
	result, err := s.command.ChangeDefaultPasswordComplexityPolicy(ctx, UpdatePasswordComplexityPolicyToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdatePasswordComplexityPolicyResponse{
		Details: object.ToDetailsPb(
			result.Sequence,
			result.ChangeDate,
			result.ResourceOwner,
		),
	}, nil
}
