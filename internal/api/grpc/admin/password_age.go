package admin

import (
	"context"

	"github.com/caos/zitadel/internal/api/grpc/object"
	policy_grpc "github.com/caos/zitadel/internal/api/grpc/policy"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func (s *Server) GetDefaultPasswordAgePolicy(ctx context.Context, req *admin_pb.GetDefaultPasswordAgePolicyRequest) (*admin_pb.GetDefaultPasswordAgePolicyResponse, error) {
	policy, err := s.iam.GetDefaultPasswordAgePolicy(ctx)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetDefaultPasswordAgePolicyResponse{
		Policy: policy_grpc.ModelPasswordAgePolicyToPb(policy),
	}, nil
}

func (s *Server) UpdateDefaultPasswordAgePolicy(ctx context.Context, req *admin_pb.UpdateDefaultPasswordAgePolicyRequest) (*admin_pb.UpdateDefaultPasswordAgePolicyResponse, error) {
	result, err := s.command.ChangeDefaultPasswordAgePolicy(ctx, UpdateDefaultPasswordAgePolicyToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdateDefaultPasswordAgePolicyResponse{
		Details: object.ToDetailsPb(
			result.Sequence,
			result.CreationDate,
			result.ChangeDate,
			result.ResourceOwner,
		),
	}, nil
}
