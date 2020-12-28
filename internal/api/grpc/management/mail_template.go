package management

import (
	"context"

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
	result, err := s.org.AddMailTemplate(ctx, mailTemplateRequestToModel(template))
	if err != nil {
		return nil, err
	}
	return mailTemplateFromModel(result), nil
}

func (s *Server) UpdateMailTemplate(ctx context.Context, template *management.MailTemplateUpdate) (*management.MailTemplate, error) {
	result, err := s.org.ChangeMailTemplate(ctx, mailTemplateRequestToModel(template))
	if err != nil {
		return nil, err
	}
	return mailTemplateFromModel(result), nil
}

func (s *Server) RemoveMailTemplate(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	err := s.org.RemoveMailTemplate(ctx)
	return &empty.Empty{}, err
}
