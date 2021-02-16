package management

import (
	"context"
	"github.com/caos/zitadel/internal/api/authz"

	"github.com/caos/zitadel/pkg/grpc/management"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) GetMailTemplate(ctx context.Context, _ *empty.Empty) (*management.MailTemplateView, error) {
	result, err := s.org.GetMailTemplate(ctx)
	if err != nil {
		return nil, err
	}
	return mailTemplateViewFromModel(result), nil
}

func (s *Server) GetDefaultMailTemplate(ctx context.Context, _ *empty.Empty) (*management.MailTemplateView, error) {
	result, err := s.org.GetDefaultMailTemplate(ctx)
	if err != nil {
		return nil, err
	}
	return mailTemplateViewFromModel(result), nil
}

func (s *Server) CreateMailTemplate(ctx context.Context, template *management.MailTemplateUpdate) (*management.MailTemplate, error) {
	result, err := s.command.AddMailTemplate(ctx, authz.GetCtxData(ctx).OrgID, mailTemplateRequestToDomain(template))
	if err != nil {
		return nil, err
	}
	return mailTemplateFromDomain(result), nil
}

func (s *Server) UpdateMailTemplate(ctx context.Context, template *management.MailTemplateUpdate) (*management.MailTemplate, error) {
	result, err := s.command.ChangeMailTemplate(ctx, authz.GetCtxData(ctx).OrgID, mailTemplateRequestToDomain(template))
	if err != nil {
		return nil, err
	}
	return mailTemplateFromDomain(result), nil
}

func (s *Server) RemoveMailTemplate(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	err := s.command.RemoveMailTemplate(ctx, authz.GetCtxData(ctx).OrgID)
	return &empty.Empty{}, err
}
