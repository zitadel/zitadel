package management

import (
	"context"
	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/pkg/grpc/management"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) GetPasswordLockoutPolicy(ctx context.Context, _ *empty.Empty) (*management.PasswordLockoutPolicyView, error) {
	result, err := s.org.GetPasswordLockoutPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return passwordLockoutPolicyViewFromModel(result), nil
}

func (s *Server) GetDefaultPasswordLockoutPolicy(ctx context.Context, _ *empty.Empty) (*management.PasswordLockoutPolicyView, error) {
	result, err := s.org.GetDefaultPasswordLockoutPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return passwordLockoutPolicyViewFromModel(result), nil
}

func (s *Server) CreatePasswordLockoutPolicy(ctx context.Context, policy *management.PasswordLockoutPolicyRequest) (*management.PasswordLockoutPolicy, error) {
	result, err := s.command.AddPasswordLockoutPolicy(ctx, passwordLockoutPolicyRequestToDomain(ctx, policy))
	if err != nil {
		return nil, err
	}
	return passwordLockoutPolicyFromDomain(result), nil
}

func (s *Server) UpdatePasswordLockoutPolicy(ctx context.Context, policy *management.PasswordLockoutPolicyRequest) (*management.PasswordLockoutPolicy, error) {
	result, err := s.command.ChangePasswordLockoutPolicy(ctx, passwordLockoutPolicyRequestToDomain(ctx, policy))
	if err != nil {
		return nil, err
	}
	return passwordLockoutPolicyFromDomain(result), nil
}

func (s *Server) RemovePasswordLockoutPolicy(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	err := s.command.RemovePasswordLockoutPolicy(ctx, authz.GetCtxData(ctx).OrgID)
	return &empty.Empty{}, err
}
