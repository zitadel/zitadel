package management

import (
	"context"
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

func (s *Server) CreatePasswordLockoutPolicy(ctx context.Context, policy *management.PasswordLockoutPolicyAdd) (*management.PasswordLockoutPolicy, error) {
	result, err := s.org.AddPasswordLockoutPolicy(ctx, passwordLockoutPolicyAddToModel(policy))
	if err != nil {
		return nil, err
	}
	return passwordLockoutPolicyFromModel(result), nil
}

func (s *Server) UpdatePasswordLockoutPolicy(ctx context.Context, policy *management.PasswordLockoutPolicy) (*management.PasswordLockoutPolicy, error) {
	result, err := s.org.ChangePasswordLockoutPolicy(ctx, passwordLockoutPolicyToModel(policy))
	if err != nil {
		return nil, err
	}
	return passwordLockoutPolicyFromModel(result), nil
}

func (s *Server) RemovePasswordLockoutPolicy(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	err := s.org.RemovePasswordLockoutPolicy(ctx)
	return &empty.Empty{}, err
}
