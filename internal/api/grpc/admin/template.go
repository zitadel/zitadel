package admin

import (
	"context"

	"github.com/caos/zitadel/pkg/grpc/admin"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) GetDefaultTemplate(ctx context.Context, _ *empty.Empty) (*admin.DefaultTemplateView, error) {
	result, err := s.iam.GetDefaultMailTemplate(ctx)
	if err != nil {
		return nil, err
	}
	return templateViewFromModel(result), nil
}

func (s *Server) UpdateDefaultTemplate(ctx context.Context, policy *admin.DefaultTemplateUpdate) (*admin.DefaultTemplate, error) {
	result, err := s.iam.ChangeDefaultMailTemplate(ctx, templateToModel(policy))
	if err != nil {
		return nil, err
	}
	return templateFromModel(result), nil
}
