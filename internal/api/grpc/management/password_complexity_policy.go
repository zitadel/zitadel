package management

import (
	"context"
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
	result, err := s.org.AddPasswordComplexityPolicy(ctx, passwordComplexityPolicyRequestToModel(policy))
	if err != nil {
		return nil, err
	}
	return passwordComplexityPolicyFromModel(result), nil
}

func (s *Server) UpdatePasswordComplexityPolicy(ctx context.Context, policy *management.PasswordComplexityPolicyRequest) (*management.PasswordComplexityPolicy, error) {
	result, err := s.org.ChangePasswordComplexityPolicy(ctx, passwordComplexityPolicyRequestToModel(policy))
	if err != nil {
		return nil, err
	}
	return passwordComplexityPolicyFromModel(result), nil
}

func (s *Server) RemovePasswordComplexityPolicy(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	err := s.org.RemovePasswordComplexityPolicy(ctx)
	return &empty.Empty{}, err
}
