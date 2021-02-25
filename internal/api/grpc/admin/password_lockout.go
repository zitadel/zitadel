package admin

import (
	"context"

	"github.com/caos/zitadel/internal/api/grpc/object"
	policy_grpc "github.com/caos/zitadel/internal/api/grpc/policy"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func (s *Server) GetDefaultPasswordLockoutPolicy(ctx context.Context, req *admin_pb.GetDefaultPasswordLockoutPolicyRequest) (*admin_pb.GetDefaultPasswordLockoutPolicyResponse, error) {
	policy, err := s.iam.GetDefaultPasswordLockoutPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetDefaultPasswordLockoutPolicyResponse{Policy: policy_grpc.ModelPasswordLockoutPolicyToPb(policy)}, nil
}

func (s *Server) UpdateDefaultPasswordLockoutPolicy(ctx context.Context, req *admin_pb.UpdateDefaultPasswordLockoutPolicyRequest) (*admin_pb.UpdateDefaultPasswordLockoutPolicyResponse, error) {
	policy, err := s.command.ChangeDefaultPasswordLockoutPolicy(ctx, UpdateDefaultPasswordLockoutPolicyToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdateDefaultPasswordLockoutPolicyResponse{
		Details: object.ToDetailsPb(
			policy.Sequence,
			policy.CreationDate,
			policy.ChangeDate,
			policy.ResourceOwner,
		),
	}, nil
}
