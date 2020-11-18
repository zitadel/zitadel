package admin

import (
	"context"

	"github.com/caos/zitadel/pkg/grpc/admin"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) GetDefaultTexts(ctx context.Context, _ *empty.Empty) (*admin.DefaultTextsView, error) {
	result, err := s.iam.GetDefaultMailTexts(ctx)
	if err != nil {
		return nil, err
	}
	return textsViewFromModel(result), nil
}

func (s *Server) GetDefaultText(ctx context.Context, _ *empty.Empty) (*admin.DefaultTextView, error) {
	result, err := s.iam.GetDefaultMailText(ctx, "type", "language") // todo wag
	if err != nil {
		return nil, err
	}
	return textViewFromModel(result), nil
}

func (s *Server) UpdateDefaultText(ctx context.Context, text *admin.DefaultTextUpdate) (*admin.DefaultText, error) {
	result, err := s.iam.ChangeDefaultMailText(ctx, textToModel(text))
	if err != nil {
		return nil, err
	}
	return textFromModel(result), nil
}
