package auth

import (
	"context"

	auth_pb "github.com/caos/zitadel/pkg/grpc/auth"
	"github.com/caos/zitadel/v2/internal/api/grpc/text"
)

func (s *Server) GetSupportedLanguages(ctx context.Context, req *auth_pb.GetSupportedLanguagesRequest) (*auth_pb.GetSupportedLanguagesResponse, error) {
	langs, err := s.query.Languages(ctx)
	if err != nil {
		return nil, err
	}
	return &auth_pb.GetSupportedLanguagesResponse{Languages: text.LanguageTagsToStrings(langs)}, nil
}
