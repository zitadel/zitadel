package admin

import (
	"context"

	"github.com/caos/zitadel/internal/api/grpc/object"
	policy_grpc "github.com/caos/zitadel/internal/api/grpc/policy"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func (s *Server) GetPasswordLockoutPolicy(ctx context.Context, req *admin_pb.GetPasswordLockoutPolicyRequest) (*admin_pb.GetPasswordLockoutPolicyResponse, error) {
	policy, err := s.iam.GetDefaultPasswordLockoutPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetPasswordLockoutPolicyResponse{Policy: policy_grpc.ModelPasswordLockoutPolicyToPb(policy)}, nil
}

func (s *Server) UpdatePasswordLockoutPolicy(ctx context.Context, req *admin_pb.UpdatePasswordLockoutPolicyRequest) (*admin_pb.UpdatePasswordLockoutPolicyResponse, error) {
	policy, err := s.command.ChangeDefaultPasswordLockoutPolicy(ctx, UpdatePasswordLockoutPolicyToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdatePasswordLockoutPolicyResponse{
		Details: object.ToDetailsPb(
			policy.Sequence,
			policy.ChangeDate,
			policy.ResourceOwner,
		),
	}, nil
}
