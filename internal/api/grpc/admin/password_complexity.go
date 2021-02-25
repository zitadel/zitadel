package admin

import (
	"context"

	"github.com/caos/zitadel/internal/api/grpc/object"
	policy_grpc "github.com/caos/zitadel/internal/api/grpc/policy"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func (s *Server) GetDefaultPasswordComplexityPolicy(ctx context.Context, _ *admin_pb.GetDefaultPasswordComplexityPolicyRequest) (*admin_pb.GetDefaultPasswordComplexityPolicyResponse, error) {
	policy, err := s.iam.GetDefaultPasswordComplexityPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetDefaultPasswordComplexityPolicyResponse{Policy: policy_grpc.ModelPasswordComplexityPolicyToPb(policy)}, nil
}

func (s *Server) UpdateDefaultPasswordComplexityPolicy(ctx context.Context, req *admin_pb.UpdateDefaultPasswordComplexityPolicyRequest) (*admin_pb.UpdateDefaultPasswordComplexityPolicyResponse, error) {
	result, err := s.command.ChangeDefaultPasswordComplexityPolicy(ctx, UpdateDefaultPasswordComplexityPolicyToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdateDefaultPasswordComplexityPolicyResponse{
		Details: object.ToDetailsPb(
			result.Sequence,
			result.CreationDate,
			result.ChangeDate,
			result.ResourceOwner,
		),
	}, nil
}
