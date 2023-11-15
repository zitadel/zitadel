package auth

import (
	"context"

	"github.com/zitadel/zitadel/v2/internal/api/grpc/text"
	auth_pb "github.com/zitadel/zitadel/v2/pkg/grpc/auth"
)

func (s *Server) GetSupportedLanguages(ctx context.Context, req *auth_pb.GetSupportedLanguagesRequest) (*auth_pb.GetSupportedLanguagesResponse, error) {
	langs, err := s.query.Languages(ctx)
	if err != nil {
		return nil, err
	}
	return &auth_pb.GetSupportedLanguagesResponse{Languages: text.LanguageTagsToStrings(langs)}, nil
}
