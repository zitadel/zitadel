package admin

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/object"
	policy_grpc "github.com/caos/zitadel/internal/api/grpc/policy"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func (s *Server) GetPasswordAgePolicy(ctx context.Context, req *admin_pb.GetPasswordAgePolicyRequest) (*admin_pb.GetPasswordAgePolicyResponse, error) {
	policy, err := s.query.DefaultPasswordAgePolicy(ctx)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetPasswordAgePolicyResponse{
		Policy: policy_grpc.ModelPasswordAgePolicyToPb(policy),
	}, nil
}

func (s *Server) UpdatePasswordAgePolicy(ctx context.Context, req *admin_pb.UpdatePasswordAgePolicyRequest) (*admin_pb.UpdatePasswordAgePolicyResponse, error) {
	result, err := s.command.ChangeDefaultPasswordAgePolicy(ctx, authz.GetInstance(ctx).ID, UpdatePasswordAgePolicyToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdatePasswordAgePolicyResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.ChangeDate,
			result.ResourceOwner,
		),
	}, nil
}
