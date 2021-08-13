package admin

import (
	"context"

	"github.com/caos/zitadel/internal/api/grpc/object"
	policy_grpc "github.com/caos/zitadel/internal/api/grpc/policy"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func (s *Server) GetLockoutPolicy(ctx context.Context, req *admin_pb.GetLockoutPolicyRequest) (*admin_pb.GetLockoutPolicyResponse, error) {
	policy, err := s.iam.GetDefaultLockoutPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetLockoutPolicyResponse{Policy: policy_grpc.ModelLockoutPolicyToPb(policy)}, nil
}

func (s *Server) UpdateLockoutPolicy(ctx context.Context, req *admin_pb.UpdateLockoutPolicyRequest) (*admin_pb.UpdateLockoutPolicyResponse, error) {
	policy, err := s.command.ChangeDefaultLockoutPolicy(ctx, UpdateLockoutPolicyToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdateLockoutPolicyResponse{
		Details: object.ChangeToDetailsPb(
			policy.Sequence,
			policy.ChangeDate,
			policy.ResourceOwner,
		),
	}, nil
}
