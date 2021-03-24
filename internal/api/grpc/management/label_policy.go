package management

import (
	"context"

	"github.com/caos/zitadel/pkg/grpc/management"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) GetLabelPolicy(ctx context.Context, _ *empty.Empty) (*management.LabelPolicyView, error) {
	result, err := s.org.GetLabelPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return labelPolicyViewFromModel(result), nil
}

func (s *Server) GetDefaultLabelPolicy(ctx context.Context, _ *empty.Empty) (*management.LabelPolicyView, error) {
	result, err := s.org.GetLabelPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return labelPolicyViewFromModel(result), nil
}

func (s *Server) CreateLabelPolicy(ctx context.Context, policy *management.LabelPolicyRequest) (*management.LabelPolicy, error) {
	result, err := s.org.AddLabelPolicy(ctx, labelPolicyRequestToModel(policy))
	if err != nil {
		return nil, err
	}
	return labelPolicyFromModel(result), nil
}

func (s *Server) UpdateLabelPolicy(ctx context.Context, policy *management.LabelPolicyRequest) (*management.LabelPolicy, error) {
	result, err := s.org.ChangeLabelPolicy(ctx, labelPolicyRequestToModel(policy))
	if err != nil {
		return nil, err
	}
	return labelPolicyFromModel(result), nil
}

func (s *Server) RemoveLabelPolicy(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	err := s.org.RemoveLabelPolicy(ctx)
	return &empty.Empty{}, err
}
