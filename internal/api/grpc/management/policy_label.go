package management

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/object"
	policy_grpc "github.com/caos/zitadel/internal/api/grpc/policy"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) GetLabelPolicy(ctx context.Context, req *mgmt_pb.GetLabelPolicyRequest) (*mgmt_pb.GetLabelPolicyResponse, error) {
	policy, err := s.org.GetLabelPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetLabelPolicyResponse{Policy: policy_grpc.ModelLabelPolicyToPb(policy)}, nil
}

func (s *Server) GetDefaultLabelPolicy(ctx context.Context, req *mgmt_pb.GetDefaultLabelPolicyRequest) (*mgmt_pb.GetDefaultLabelPolicyResponse, error) {
	policy, err := s.org.GetDefaultLabelPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetDefaultLabelPolicyResponse{Policy: policy_grpc.ModelLabelPolicyToPb(policy)}, nil
}

func (s *Server) AddCustomLabelPolicy(ctx context.Context, req *mgmt_pb.AddCustomLabelPolicyRequest) (*mgmt_pb.AddCustomLabelPolicyResponse, error) {
	policy, err := s.command.AddLabelPolicy(ctx, authz.GetCtxData(ctx).OrgID, addLabelPolicyToDomain(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddCustomLabelPolicyResponse{
		Details: object.AddToDetailsPb(
			policy.Sequence,
			policy.ChangeDate,
			policy.ResourceOwner,
		),
	}, nil
}

func (s *Server) UpdateCustomLabelPolicy(ctx context.Context, req *mgmt_pb.UpdateCustomLabelPolicyRequest) (*mgmt_pb.UpdateCustomLabelPolicyResponse, error) {
	policy, err := s.command.ChangeLabelPolicy(ctx, authz.GetCtxData(ctx).OrgID, updateLabelPolicyToDomain(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateCustomLabelPolicyResponse{
		Details: object.ChangeToDetailsPb(
			policy.Sequence,
			policy.ChangeDate,
			policy.ResourceOwner,
		),
	}, nil
}

func (s *Server) ResetLabelPolicyToDefault(ctx context.Context, req *mgmt_pb.ResetLabelPolicyToDefaultRequest) (*mgmt_pb.ResetLabelPolicyToDefaultResponse, error) {
	objectDetails, err := s.command.RemoveLabelPolicy(ctx, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ResetLabelPolicyToDefaultResponse{
		Details: object.DomainToChangeDetailsPb(objectDetails),
	}, nil
}
