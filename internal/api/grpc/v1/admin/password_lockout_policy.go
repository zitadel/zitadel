package admin

import (
	"context"
	"github.com/caos/zitadel/pkg/grpc/admin"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) GetDefaultPasswordLockoutPolicy(ctx context.Context, _ *empty.Empty) (*admin.DefaultPasswordLockoutPolicyView, error) {
	result, err := s.iam.GetDefaultPasswordLockoutPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return passwordLockoutPolicyViewFromModel(result), nil
}

func (s *Server) UpdateDefaultPasswordLockoutPolicy(ctx context.Context, policy *admin.DefaultPasswordLockoutPolicyRequest) (*admin.DefaultPasswordLockoutPolicy, error) {
	result, err := s.command.ChangeDefaultPasswordLockoutPolicy(ctx, passwordLockoutPolicyToDomain(policy))
	if err != nil {
		return nil, err
	}
	return passwordLockoutPolicyFromDomain(result), nil
}
