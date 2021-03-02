package management

import (
	"context"
	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/object"
	policy_grpc "github.com/caos/zitadel/internal/api/grpc/policy"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) GetPasswordComplexityPolicy(ctx context.Context, req *mgmt_pb.GetPasswordComplexityPolicyRequest) (*mgmt_pb.GetPasswordComplexityPolicyResponse, error) {
	policy, err := s.org.GetPasswordComplexityPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetPasswordComplexityPolicyResponse{Policy: policy_grpc.ModelPasswordComplexityPolicyToPb(policy)}, nil
}

func (s *Server) GetDefaultPasswordComplexityPolicy(ctx context.Context, req *mgmt_pb.GetDefaultPasswordComplexityPolicyRequest) (*mgmt_pb.GetDefaultPasswordComplexityPolicyResponse, error) {
	policy, err := s.org.GetDefaultPasswordComplexityPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetDefaultPasswordComplexityPolicyResponse{Policy: policy_grpc.ModelPasswordComplexityPolicyToPb(policy)}, nil
}

func (s *Server) AddCustomPasswordComplexityPolicy(ctx context.Context, req *mgmt_pb.AddCustomPasswordComplexityPolicyRequest) (*mgmt_pb.AddCustomPasswordComplexityPolicyResponse, error) {
	result, err := s.command.AddPasswordComplexityPolicy(ctx, authz.GetCtxData(ctx).OrgID, AddPasswordComplexityPolicyToDomain(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddCustomPasswordComplexityPolicyResponse{
		Details: object.ToDetailsPb(
			result.Sequence,
			result.ChangeDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) UpdateCustomPasswordComplexityPolicy(ctx context.Context, req *mgmt_pb.UpdateCustomPasswordComplexityPolicyRequest) (*mgmt_pb.UpdateCustomPasswordComplexityPolicyResponse, error) {
	result, err := s.command.ChangePasswordComplexityPolicy(ctx, authz.GetCtxData(ctx).OrgID, UpdatePasswordComplexityPolicyToDomain(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateCustomPasswordComplexityPolicyResponse{
		Details: object.ToDetailsPb(
			result.Sequence,
			result.ChangeDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) ResetPasswordComplexityPolicyToDefault(ctx context.Context, req *mgmt_pb.ResetPasswordComplexityPolicyToDefaultRequest) (*mgmt_pb.ResetPasswordComplexityPolicyToDefaultResponse, error) {
	objectDetails, err := s.command.RemovePasswordComplexityPolicy(ctx, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ResetPasswordComplexityPolicyToDefaultResponse{
		Details: object.DomainToDetailsPb(objectDetails),
	}, nil
}
