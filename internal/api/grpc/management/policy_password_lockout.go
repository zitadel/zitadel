package management

import (
	"context"
	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/object"
	policy_grpc "github.com/caos/zitadel/internal/api/grpc/policy"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) GetPasswordLockoutPolicy(ctx context.Context, req *mgmt_pb.GetPasswordLockoutPolicyRequest) (*mgmt_pb.GetPasswordLockoutPolicyResponse, error) {
	policy, err := s.org.GetPasswordLockoutPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetPasswordLockoutPolicyResponse{Policy: policy_grpc.ModelPasswordLockoutPolicyToPb(policy), IsDefault: policy.Default}, nil
}

func (s *Server) GetDefaultPasswordLockoutPolicy(ctx context.Context, req *mgmt_pb.GetDefaultPasswordLockoutPolicyRequest) (*mgmt_pb.GetDefaultPasswordLockoutPolicyResponse, error) {
	policy, err := s.org.GetDefaultPasswordLockoutPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetDefaultPasswordLockoutPolicyResponse{Policy: policy_grpc.ModelPasswordLockoutPolicyToPb(policy)}, nil
}

func (s *Server) AddCustomPasswordLockoutPolicy(ctx context.Context, req *mgmt_pb.AddCustomPasswordLockoutPolicyRequest) (*mgmt_pb.AddCustomPasswordLockoutPolicyResponse, error) {
	policy, err := s.command.AddPasswordLockoutPolicy(ctx, authz.GetCtxData(ctx).OrgID, AddPasswordLockoutPolicyToDomain(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddCustomPasswordLockoutPolicyResponse{
		Details: object.AddToDetailsPb(
			policy.Sequence,
			policy.ChangeDate,
			policy.ResourceOwner,
		),
	}, nil
}

func (s *Server) UpdateCustomPasswordLockoutPolicy(ctx context.Context, req *mgmt_pb.UpdateCustomPasswordLockoutPolicyRequest) (*mgmt_pb.UpdateCustomPasswordLockoutPolicyResponse, error) {
	policy, err := s.command.ChangePasswordLockoutPolicy(ctx, authz.GetCtxData(ctx).OrgID, UpdatePasswordLockoutPolicyToDomain(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateCustomPasswordLockoutPolicyResponse{
		Details: object.ChangeToDetailsPb(
			policy.Sequence,
			policy.ChangeDate,
			policy.ResourceOwner,
		),
	}, nil
}

func (s *Server) ResetPasswordLockoutPolicyToDefault(ctx context.Context, req *mgmt_pb.ResetPasswordLockoutPolicyToDefaultRequest) (*mgmt_pb.ResetPasswordLockoutPolicyToDefaultResponse, error) {
	objectDetails, err := s.command.RemovePasswordComplexityPolicy(ctx, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ResetPasswordLockoutPolicyToDefaultResponse{
		Details: object.DomainToChangeDetailsPb(objectDetails),
	}, nil
}
