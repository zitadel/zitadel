package management

import (
	"context"

	"github.com/caos/zitadel/pkg/grpc/management"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) GetMailTexts(ctx context.Context, _ *empty.Empty) (*management.MailTextsView, error) {
	result, err := s.org.GetMailTexts(ctx)
	if err != nil {
		return nil, err
	}
	return mailTextsViewFromModel(result.Texts), nil
}

func (s *Server) GetDefaultMailTexts(ctx context.Context, _ *empty.Empty) (*management.MailTextsView, error) {
	result, err := s.org.GetDefaultMailTexts(ctx)
	if err != nil {
		return nil, err
	}
	return mailTextsViewFromModel(result.Texts), nil
}

func (s *Server) CreateMailText(ctx context.Context, template *management.MailTextUpdate) (*management.MailText, error) {
	result, err := s.org.AddMailText(ctx, mailTextRequestToModel(template))
	if err != nil {
		return nil, err
	}
	return mailTextFromModel(result), nil
}

func (s *Server) UpdateMailText(ctx context.Context, template *management.MailTextUpdate) (*management.MailText, error) {
	result, err := s.org.ChangeMailText(ctx, mailTextRequestToModel(template))
	if err != nil {
		return nil, err
	}
	return mailTextFromModel(result), nil
}

func (s *Server) RemoveMailText(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	err := s.org.RemoveMailText(ctx)
	return &empty.Empty{}, err
}
