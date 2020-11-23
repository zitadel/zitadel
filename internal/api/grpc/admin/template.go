package admin

import (
	"context"

	"github.com/caos/zitadel/pkg/grpc/admin"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) GetDefaultMailTemplate(ctx context.Context, _ *empty.Empty) (*admin.DefaultMailTemplateView, error) {
	result, err := s.iam.GetDefaultMailTemplate(ctx)
	if err != nil {
		return nil, err
	}
	return templateViewFromModel(result), nil
}

func (s *Server) UpdateDefaultMailTemplate(ctx context.Context, policy *admin.DefaultMailTemplateUpdate) (*admin.DefaultMailTemplate, error) {
	result, err := s.iam.ChangeDefaultMailTemplate(ctx, templateToModel(policy))
	if err != nil {
		return nil, err
	}
	return templateFromModel(result), nil
}
