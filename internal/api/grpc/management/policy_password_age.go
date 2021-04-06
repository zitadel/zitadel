package management

import (
	"context"
	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/grpc/object"
	policy_grpc "github.com/caos/zitadel/internal/api/grpc/policy"
	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) GetPasswordAgePolicy(ctx context.Context, req *mgmt_pb.GetPasswordAgePolicyRequest) (*mgmt_pb.GetPasswordAgePolicyResponse, error) {
	policy, err := s.org.GetPasswordAgePolicy(ctx)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetPasswordAgePolicyResponse{
		Policy: policy_grpc.ModelPasswordAgePolicyToPb(policy),
	}, nil
}

func (s *Server) GetDefaultPasswordAgePolicy(ctx context.Context, req *mgmt_pb.GetDefaultPasswordAgePolicyRequest) (*mgmt_pb.GetDefaultPasswordAgePolicyResponse, error) {
	policy, err := s.org.GetDefaultPasswordAgePolicy(ctx)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetDefaultPasswordAgePolicyResponse{
		Policy: policy_grpc.ModelPasswordAgePolicyToPb(policy),
	}, nil
}

func (s *Server) AddCustomPasswordAgePolicy(ctx context.Context, req *mgmt_pb.AddCustomPasswordAgePolicyRequest) (*mgmt_pb.AddCustomPasswordAgePolicyResponse, error) {
	result, err := s.command.AddPasswordAgePolicy(ctx, authz.GetCtxData(ctx).OrgID, AddPasswordAgePolicyToDomain(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.AddCustomPasswordAgePolicyResponse{
		Details: object.AddToDetailsPb(
			result.Sequence,
			result.ChangeDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) UpdateCustomPasswordAgePolicy(ctx context.Context, req *mgmt_pb.UpdateCustomPasswordAgePolicyRequest) (*mgmt_pb.UpdateCustomPasswordAgePolicyResponse, error) {
	result, err := s.command.ChangePasswordAgePolicy(ctx, authz.GetCtxData(ctx).OrgID, UpdatePasswordAgePolicyToDomain(req))
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.UpdateCustomPasswordAgePolicyResponse{
		Details: object.ChangeToDetailsPb(
			result.Sequence,
			result.ChangeDate,
			result.ResourceOwner,
		),
	}, nil
}

func (s *Server) ResetPasswordAgePolicyToDefault(ctx context.Context, req *mgmt_pb.ResetPasswordAgePolicyToDefaultRequest) (*mgmt_pb.ResetPasswordAgePolicyToDefaultResponse, error) {
	objectDetails, err := s.command.RemovePasswordAgePolicy(ctx, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.ResetPasswordAgePolicyToDefaultResponse{
		Details: object.DomainToChangeDetailsPb(objectDetails),
	}, nil
}
