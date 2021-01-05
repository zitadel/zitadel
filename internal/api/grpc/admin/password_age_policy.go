package admin

import (
	"context"
	"github.com/caos/zitadel/pkg/grpc/admin"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) GetDefaultPasswordAgePolicy(ctx context.Context, _ *empty.Empty) (*admin.DefaultPasswordAgePolicyView, error) {
	result, err := s.iam.GetDefaultPasswordAgePolicy(ctx)
	if err != nil {
		return nil, err
	}
	return passwordAgePolicyViewFromModel(result), nil
}

func (s *Server) UpdateDefaultPasswordAgePolicy(ctx context.Context, policy *admin.DefaultPasswordAgePolicyRequest) (*admin.DefaultPasswordAgePolicy, error) {
	result, err := s.command.ChangeDefaultPasswordAgePolicy(ctx, passwordAgePolicyToDomain(policy))
	if err != nil {
		return nil, err
	}
	return passwordAgePolicyFromDomain(result), nil
}
