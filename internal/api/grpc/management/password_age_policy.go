package management

import (
	"context"
	"github.com/caos/zitadel/internal/api/authz"
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

func (s *Server) CreatePasswordAgePolicy(ctx context.Context, policy *management.PasswordAgePolicyRequest) (*management.PasswordAgePolicy, error) {
	result, err := s.command.AddPasswordAgePolicy(ctx, authz.GetCtxData(ctx).OrgID, passwordAgePolicyRequestToDomain(ctx, policy))
	if err != nil {
		return nil, err
	}
	return passwordAgePolicyFromDomain(result), nil
}

func (s *Server) UpdatePasswordAgePolicy(ctx context.Context, policy *management.PasswordAgePolicyRequest) (*management.PasswordAgePolicy, error) {
	result, err := s.command.ChangePasswordAgePolicy(ctx, authz.GetCtxData(ctx).OrgID, passwordAgePolicyRequestToDomain(ctx, policy))
	if err != nil {
		return nil, err
	}
	return passwordAgePolicyFromDomain(result), nil
}

func (s *Server) RemovePasswordAgePolicy(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	err := s.command.RemovePasswordAgePolicy(ctx, authz.GetCtxData(ctx).OrgID)
	return &empty.Empty{}, err
}
