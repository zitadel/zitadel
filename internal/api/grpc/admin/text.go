package admin

import (
	"context"

	"github.com/caos/zitadel/pkg/grpc/admin"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) GetDefaultMailTexts(ctx context.Context, _ *empty.Empty) (*admin.DefaultMailTextsView, error) {
	result, err := s.iam.GetDefaultMailTexts(ctx)
	if err != nil {
		return nil, err
	}
	return textsViewFromModel(result), nil
}

func (s *Server) GetDefaultMailText(ctx context.Context, textType string, language string) (*admin.DefaultMailTextView, error) {
	result, err := s.iam.GetDefaultMailText(ctx, textType, language)
	if err != nil {
		return nil, err
	}
	return textViewFromModel(result), nil
}

func (s *Server) UpdateDefaultMailText(ctx context.Context, text *admin.DefaultMailTextUpdate) (*admin.DefaultMailText, error) {
	result, err := s.command.ChangeDefaultMailText(ctx, textToDomain(text))
	if err != nil {
		return nil, err
	}
	return textFromDomain(result), nil
}
