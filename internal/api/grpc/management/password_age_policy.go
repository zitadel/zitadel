package management

import (
	"context"
	"github.com/caos/zitadel/pkg/grpc/management"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) GetPasswordAgePolicy(ctx context.Context, _ *empty.Empty) (*management.PasswordAgePolicyView, error) {
	result, err := s.org.GetPasswordAgePolicy(ctx)
	if err != nil {
		return nil, err
	}
	return passwordAgePolicyViewFromModel(result), nil
}

func (s *Server) GetDefaultPasswordAgePolicy(ctx context.Context, _ *empty.Empty) (*management.PasswordAgePolicyView, error) {
	result, err := s.org.GetDefaultPasswordAgePolicy(ctx)
	if err != nil {
		return nil, err
	}
	return passwordAgePolicyViewFromModel(result), nil
}

func (s *Server) CreatePasswordAgePolicy(ctx context.Context, policy *management.PasswordAgePolicyAdd) (*management.PasswordAgePolicy, error) {
	result, err := s.org.AddPasswordAgePolicy(ctx, passwordAgePolicyAddToModel(policy))
	if err != nil {
		return nil, err
	}
	return passwordAgePolicyFromModel(result), nil
}

func (s *Server) UpdatePasswordAgePolicy(ctx context.Context, policy *management.PasswordAgePolicy) (*management.PasswordAgePolicy, error) {
	result, err := s.org.ChangePasswordAgePolicy(ctx, passwordAgePolicyToModel(policy))
	if err != nil {
		return nil, err
	}
	return passwordAgePolicyFromModel(result), nil
}

func (s *Server) RemovePasswordAgePolicy(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	err := s.org.RemovePasswordAgePolicy(ctx)
	return &empty.Empty{}, err
}
