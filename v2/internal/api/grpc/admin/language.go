package admin

import (
	"context"

	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
	"github.com/caos/zitadel/v2/internal/api/grpc/text"
)

func (s *Server) GetSupportedLanguages(ctx context.Context, req *admin_pb.GetSupportedLanguagesRequest) (*admin_pb.GetSupportedLanguagesResponse, error) {
	langs, err := s.query.Languages(ctx)
	if err != nil {
		return nil, err
	}
	return &admin_pb.GetSupportedLanguagesResponse{Languages: text.LanguageTagsToStrings(langs)}, nil
}
