package admin

import (
	"context"

	"github.com/caos/zitadel/internal/api/grpc/object"
	policy_grpc "github.com/caos/zitadel/internal/api/grpc/policy"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func (s *Server) GetDefaultLabelPolicy(ctx context.Context, req *admin_pb.GetDefaultLabelPolicyRequest) (*admin_pb.GetDefaultLabelPolicyResponse, error) {
	policy, err := s.iam.GetDefaultLabelPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetDefaultLabelPolicyResponse{Policy: policy_grpc.ModelLabelPolicyToPb(policy)}, nil
}

func (s *Server) UpdateDefaultLabelPolicy(ctx context.Context, req *admin_pb.UpdateDefaultLabelPolicyRequest) (*admin_pb.UpdateDefaultLabelPolicyResponse, error) {
	policy, err := s.command.ChangeDefaultLabelPolicy(ctx, updateDefaultLabelPolicyToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdateDefaultLabelPolicyResponse{
		Details: object.ToDetailsPb(
			policy.Sequence,
			policy.CreationDate,
			policy.ChangeDate,
			policy.ResourceOwner,
		),
	}, nil
}
