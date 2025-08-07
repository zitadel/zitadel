package admin

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/grpc/object"
	policy_grpc "github.com/zitadel/zitadel/internal/api/grpc/policy"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
)

func (s *Server) GetLockoutPolicy(ctx context.Context, req *admin_pb.GetLockoutPolicyRequest) (*admin_pb.GetLockoutPolicyResponse, error) {
	policy, err := s.query.DefaultLockoutPolicy(ctx)
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
