package management

import (
	"context"
	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/object"
	policy_grpc "github.com/caos/zitadel/internal/api/grpc/policy"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) GetLockoutPolicy(ctx context.Context, req *mgmt_pb.GetLockoutPolicyRequest) (*mgmt_pb.GetLockoutPolicyResponse, error) {
	policy, err := s.org.GetLockoutPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetLockoutPolicyResponse{Policy: policy_grpc.ModelLockoutPolicyToPb(policy), IsDefault: policy.Default}, nil
}

func (s *Server) GetDefaultLockoutPolicy(ctx context.Context, req *mgmt_pb.GetDefaultLockoutPolicyRequest) (*mgmt_pb.GetDefaultLockoutPolicyResponse, error) {
	policy, err := s.org.GetDefaultLockoutPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetDefaultLockoutPolicyResponse{Policy: policy_grpc.ModelLockoutPolicyToPb(policy)}, nil
}

func (s *Server) AddCustomLockoutPolicy(ctx context.Context, req *mgmt_pb.AddCustomLockoutPolicyRequest) (*mgmt_pb.AddCustomLockoutPolicyResponse, error) {
	policy, err := s.command.AddLockoutPolicy(ctx, authz.GetCtxData(ctx).OrgID, AddLockoutPolicyToDomain(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddCustomLockoutPolicyResponse{
		Details: object.AddToDetailsPb(
			policy.Sequence,
			policy.ChangeDate,
			policy.ResourceOwner,
		),
	}, nil
}

func (s *Server) UpdateCustomLockoutPolicy(ctx context.Context, req *mgmt_pb.UpdateCustomLockoutPolicyRequest) (*mgmt_pb.UpdateCustomLockoutPolicyResponse, error) {
	policy, err := s.command.ChangeLockoutPolicy(ctx, authz.GetCtxData(ctx).OrgID, UpdateLockoutPolicyToDomain(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateCustomLockoutPolicyResponse{
		Details: object.ChangeToDetailsPb(
			policy.Sequence,
			policy.ChangeDate,
			policy.ResourceOwner,
		),
	}, nil
}

func (s *Server) ResetLockoutPolicyToDefault(ctx context.Context, req *mgmt_pb.ResetLockoutPolicyToDefaultRequest) (*mgmt_pb.ResetLockoutPolicyToDefaultResponse, error) {
	objectDetails, err := s.command.RemovePasswordComplexityPolicy(ctx, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ResetLockoutPolicyToDefaultResponse{
		Details: object.DomainToChangeDetailsPb(objectDetails),
	}, nil
}
