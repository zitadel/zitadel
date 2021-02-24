package management

import (
	"context"
	"github.com/caos/zitadel/internal/api/authz"

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

func (s *Server) CreateMailText(ctx context.Context, mailText *management.MailTextUpdate) (*management.MailText, error) {
	result, err := s.command.AddMailText(ctx, authz.GetCtxData(ctx).OrgID, mailTextRequestToDomain(mailText))
	if err != nil {
		return nil, err
	}
	return mailTextFromDoamin(result), nil
}

func (s *Server) UpdateMailText(ctx context.Context, mailText *management.MailTextUpdate) (*management.MailText, error) {
	result, err := s.command.ChangeMailText(ctx, authz.GetCtxData(ctx).OrgID, mailTextRequestToDomain(mailText))
	if err != nil {
		return nil, err
	}
	return mailTextFromDoamin(result), nil
}

func (s *Server) RemoveMailText(ctx context.Context, mailText *management.MailTextRemove) (*empty.Empty, error) {
	err := s.command.RemoveMailText(ctx, authz.GetCtxData(ctx).OrgID, mailText.MailTextType, mailText.Language)
	return &empty.Empty{}, err
}
