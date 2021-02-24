package admin

import (
	"context"
	"github.com/caos/zitadel/pkg/grpc/admin"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) GetDefaultPasswordComplexityPolicy(ctx context.Context, _ *empty.Empty) (*admin.DefaultPasswordComplexityPolicyView, error) {
	result, err := s.iam.GetDefaultPasswordComplexityPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return passwordComplexityPolicyViewFromModel(result), nil
}

func (s *Server) UpdateDefaultPasswordComplexityPolicy(ctx context.Context, policy *admin.DefaultPasswordComplexityPolicyRequest) (*admin.DefaultPasswordComplexityPolicy, error) {
	result, err := s.command.ChangeDefaultPasswordComplexityPolicy(ctx, passwordComplexityPolicyToDomain(policy))
	if err != nil {
		return nil, err
	}
	return passwordComplexityPolicyFromDomain(result), nil
}
