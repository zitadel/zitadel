package admin

import (
	"context"

	"github.com/caos/zitadel/pkg/grpc/admin"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) GetDefaultLabelPolicy(ctx context.Context, _ *empty.Empty) (*admin.DefaultLabelPolicyView, error) {
	result, err := s.iam.GetDefaultLabelPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return labelPolicyViewFromModel(result), nil
}

func (s *Server) UpdateDefaultLabelPolicy(ctx context.Context, policy *admin.DefaultLabelPolicyUpdate) (*admin.DefaultLabelPolicy, error) {
	result, err := s.iam.ChangeDefaultLabelPolicy(ctx, labelPolicyToModel(policy))
	if err != nil {
		return nil, err
	}
	return labelPolicyFromModel(result), nil
}
