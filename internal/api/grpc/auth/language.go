package auth

import (
	"context"

	"github.com/caos/zitadel/internal/api/grpc/text"
	auth_pb "github.com/caos/zitadel/pkg/grpc/auth"
)

func (s *Server) GetMyLanguages(ctx context.Context, req *auth_pb.GetMyLanguagesRequest) (*auth_pb.GetMyLanguagesResponse, error) {
	langs, err := s.repo.Languages(ctx)
	if err != nil {
		return nil, err
	}
	return &auth_pb.GetMyLanguagesResponse{Languages: text.LanguageTagsToStrings(langs)}, nil
}
