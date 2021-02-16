package management

import (
	"context"
	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/pkg/grpc/management"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) GetPasswordComplexityPolicy(ctx context.Context, _ *empty.Empty) (*management.PasswordComplexityPolicyView, error) {
	result, err := s.org.GetPasswordComplexityPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return passwordComplexityPolicyViewFromModel(result), nil
}

func (s *Server) GetDefaultPasswordComplexityPolicy(ctx context.Context, _ *empty.Empty) (*management.PasswordComplexityPolicyView, error) {
	result, err := s.org.GetDefaultPasswordComplexityPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return passwordComplexityPolicyViewFromModel(result), nil
}

func (s *Server) CreatePasswordComplexityPolicy(ctx context.Context, policy *management.PasswordComplexityPolicyRequest) (*management.PasswordComplexityPolicy, error) {
	result, err := s.command.AddPasswordComplexityPolicy(ctx, authz.GetCtxData(ctx).OrgID, passwordComplexityPolicyRequestToDomain(ctx, policy))
	if err != nil {
		return nil, err
	}
	return passwordComplexityPolicyFromDomain(result), nil
}

func (s *Server) UpdatePasswordComplexityPolicy(ctx context.Context, policy *management.PasswordComplexityPolicyRequest) (*management.PasswordComplexityPolicy, error) {
	result, err := s.command.ChangePasswordComplexityPolicy(ctx, authz.GetCtxData(ctx).OrgID, passwordComplexityPolicyRequestToDomain(ctx, policy))
	if err != nil {
		return nil, err
	}
	return passwordComplexityPolicyFromDomain(result), nil
}

func (s *Server) RemovePasswordComplexityPolicy(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	err := s.command.RemovePasswordComplexityPolicy(ctx, authz.GetCtxData(ctx).OrgID)
	return &empty.Empty{}, err
}
