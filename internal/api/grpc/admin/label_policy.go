package admin

import (
	"context"

	"github.com/caos/zitadel/internal/api/grpc/object"
	policy_grpc "github.com/caos/zitadel/internal/api/grpc/policy"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
)

func (s *Server) GetLabelPolicy(ctx context.Context, req *admin_pb.GetLabelPolicyRequest) (*admin_pb.GetLabelPolicyResponse, error) {
	policy, err := s.iam.GetDefaultLabelPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetLabelPolicyResponse{Policy: policy_grpc.ModelLabelPolicyToPb(policy)}, nil
}

func (s *Server) UpdateLabelPolicy(ctx context.Context, req *admin_pb.UpdateLabelPolicyRequest) (*admin_pb.UpdateLabelPolicyResponse, error) {
	policy, err := s.command.ChangeDefaultLabelPolicy(ctx, updateLabelPolicyToDomain(req))
	if err != nil {
		return nil, err
	}
	return &admin_pb.UpdateLabelPolicyResponse{
		Details: object.ChangeToDetailsPb(
			policy.Sequence,
			policy.ChangeDate,
			policy.ResourceOwner,
		),
	}, nil
}
